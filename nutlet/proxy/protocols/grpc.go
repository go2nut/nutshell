package protocols

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
	"net"
	"net/url"
	"sync"
)

type GrpcRouter func(ctx context.Context, targetAddress string, fullnameMethod string) (*url.URL, error)

type GrpcProxyServer struct {
	Router GrpcRouter
	clientStreamDescForProxying *grpc.StreamDesc
	rpcConnPool map[string]*grpc.ClientConn
	lock  *sync.RWMutex
}

func NewGrpcProxyServer(Router GrpcRouter) *GrpcProxyServer{
	return &GrpcProxyServer{Router: Router,
		clientStreamDescForProxying: &grpc.StreamDesc{
			ServerStreams: true,
			ClientStreams: true,
		},
		rpcConnPool: map[string]*grpc.ClientConn{},
		lock:  &sync.RWMutex{},
	}
}

func (p *GrpcProxyServer) ServeGrpc(port int) {
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Printf("can not listen tcp connections, Err=%s\n", err)
		return
	}
	log.Printf("serving grpc proxy on address: %s\n", listen.Addr().String())

	server := grpc.NewServer(
		grpc.CustomCodec(Codec()),
		grpc.UnknownServiceHandler(p.StreamHandler),
	)
	log.Printf("serviceInfo:%#v", server.GetServiceInfo())

	log.Println(fmt.Sprintf("serving grpc proxy on:%d", port))

	if err := server.Serve(listen); err != nil {
		log.Printf("can not listen tcp connections %v \n", err)
		return
	}

	log.Println(fmt.Sprintf("closing grpc proxy on:%d", port))
}

func (p *GrpcProxyServer) getDialConn(host string) (conn *grpc.ClientConn, err error) {
	p.lock.RLock()
	var exist bool
	if conn, exist = p.rpcConnPool[host]; exist {
		if conn.GetState() != connectivity.Shutdown {
			p.lock.RUnlock()
			return
		}
		delete(p.rpcConnPool, host)
	}
	p.lock.RUnlock()
	p.lock.Lock()
	if conn, exist = p.rpcConnPool[host]; !exist {
		if conn, err = grpc.Dial(host, grpc.WithInsecure(), grpc.WithCodec(Codec())); err == nil {
			log.Printf("rcp New Target Ip: %s\n", host)
			p.rpcConnPool[host] = conn
		}
	}
	p.lock.Unlock()
	return
}

// StreamHandler 核心代码:
// 接受 Stream 帧数据; 并传递给目标服务; 如果为确立链接, 则确立;
func (p *GrpcProxyServer) StreamHandler(srv interface{}, serverStream grpc.ServerStream) error {
	println(">>>>>>> StreamHandler")
	fullMethodName, ok := grpc.MethodFromServerStream(serverStream)
	if !ok {
		return grpc.Errorf(codes.Internal, "lowLevelServerStream not exists in context")
	}

	ctx := serverStream.Context()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return fmt.Errorf("no metadata found for grpc request")
	}

	auths, found := md[":authority"]
	if !found || len(auths) == 0 || len(auths[0]) == 0 {
		return fmt.Errorf("no origin address found in metadata")
	}

	originAddr := auths[0]

	target, err := p.Router(ctx, originAddr, fullMethodName)
	if err != nil {
		return err
	}

	backendConn, err := p.getDialConn(target.Host)
	if err != nil {
		return err
	}

	fullMethodName = target.Path

	// TODO: 增加 `forwarded` 头信息到 Metadata, https://en.wikipedia.org/wiki/X-Forwarded-For.
	clientCtx, clientCancel := context.WithCancel(ctx)
	clientStream, err := grpc.NewClientStream(
		clientCtx,
		p.clientStreamDescForProxying,
		backendConn,
		fullMethodName,
	)
	if err != nil {
		return err
	}

	// Explicitly *do not close* s2cErrChan and c2sErrChan, otherwise the select below will not terminate.
	// Channels do not have to be closed, it is just a control flow mechanism, see
	// https://groups.google.com/forum/#!msg/golang-nuts/pZwdYRGxCIk/qpbHxRRPJdUJ
	s2cErrChan := p.forwardServerToClient(serverStream, clientStream)
	c2sErrChan := p.forwardClientToServer(clientStream, serverStream)

	// We don't know which side is going to stop sending first, so we need a select between the two.
	for i := 0; i < 2; i++ {
		select {
		case s2cErr := <-s2cErrChan:
			if s2cErr == io.EOF {
				// this is the happy case where the sender has encountered io.EOF, and won't be sending anymore./
				// the clientStream>serverStream may continue pumping though.
				clientStream.CloseSend()
				break
			} else {
				// however, we may have gotten a receive error (stream disconnected, a read error etc) in which case we need
				// to cancel the clientStream to the backend, let all of its goroutines be freed up by the CancelFunc and
				// exit with an error to the stack
				clientCancel()
				return grpc.Errorf(codes.Internal, "failed proxying s2c: %v", s2cErr)
			}
		case c2sErr := <-c2sErrChan:
			// This happens when the clientStream has nothing else to offer (io.EOF), returned a gRPC error. In those two
			// cases we may have received Trailers as part of the call. In case of other errors (stream closed) the trailers
			// will be nil.
			serverStream.SetTrailer(clientStream.Trailer())
			// c2sErr will contain RPC error from client code. If not io.EOF return the RPC error as server stream error.
			if c2sErr != io.EOF {
				return c2sErr
			}
			return nil
		}
	}

	return grpc.Errorf(codes.Internal, "gRPC proxying should never reach this stage.")
}

func (p *GrpcProxyServer) forwardClientToServer(src grpc.ClientStream, dst grpc.ServerStream) chan error {
	ret := make(chan error, 1)
	go func() {
		f := &frame{}
		for i := 0; ; i++ {
			if err := src.RecvMsg(f); err != nil {
				ret <- err // this can be io.EOF which is happy case
				break
			}
			if i == 0 {
				// This is a bit of a hack, but client to server headers are only readable after first client msg is
				// received but must be written to server stream before the first msg is flushed.
				// This is the only place to do it nicely.
				md, err := src.Header()
				if err != nil {
					ret <- err
					break
				}
				if err := dst.SendHeader(md); err != nil {
					ret <- err
					break
				}
			}
			if err := dst.SendMsg(f); err != nil {
				ret <- err
				break
			}
		}
	}()
	return ret
}

func (p *GrpcProxyServer) forwardServerToClient(src grpc.ServerStream, dst grpc.ClientStream) chan error {
	ret := make(chan error, 1)
	go func() {
		f := &frame{}
		for i := 0; ; i++ {
			if err := src.RecvMsg(f); err != nil {
				ret <- err // this can be io.EOF which is happy case
				break
			}
			if err := dst.SendMsg(f); err != nil {
				ret <- err
				break
			}
		}
	}()
	return ret
}
