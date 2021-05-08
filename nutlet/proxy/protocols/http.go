package protocols

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type HttpProxyServer struct {
	t          http.RoundTripper
	httpRouter Router
	wsRouter   Router
}

func NewHttpProxy(httpR Router, wsR Router) *HttpProxyServer {
	return &HttpProxyServer{
		t:          http.DefaultTransport,
		httpRouter: httpR,
		wsRouter:   wsR,
	}
}

type Router func(ctx context.Context, req *http.Request) (*url.URL, bool)

func (s *HttpProxyServer) ServeHTTP(port int) {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Printf("can not listen tcp connections, Err=%s\n", err)
		return
	}
	log.Println("serving http proxy on address: ", l.Addr().String())

	h := http.NewServeMux()
	h.HandleFunc("/", s.handleRequest)
	srv := &http.Server{Handler: h}
	err = srv.Serve(l)
	if err != nil {
		log.Fatal("failed to serve http")
	}
}

func (s *HttpProxyServer) handleRequest(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	if target, exist := s.httpRouter(ctx, req); exist {

		log.Printf("schema:%s origin=%s target=%s\n", req.URL.Scheme, req.Host, target)

		if target.Host == "" {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte("502 - Bad Gateway"))
			return
		}
		s.serveProxy(target, w, req)
		return
	} else if target, exist := s.wsRouter(ctx, req); exist {
		NewProxy(target).ServeHTTP(w, req)
		return
	} else {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("502 - Bad Gateway"))
		return
	}
}

func (s *HttpProxyServer) serveProxy(target *url.URL, w http.ResponseWriter, req *http.Request) {
	log.Printf("serve proxy request req.url=%s req.host=%s\n", req.URL, req.Host)

	p := httputil.NewSingleHostReverseProxy(target)
	p.Director = func(req *http.Request) {
		req.URL = target
		req.Host = target.Host
		req.RequestURI = target.Path
	}
	p.Transport = s.t
	p.ServeHTTP(w, req)
}

//
//type httpTransport struct {
//}
//
//func (t *httpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
//	resp, err := http.DefaultTransport.RoundTrip(req)
//	if err != nil {
//		log.Printf("RoundTripErr:failed to send http request; addr=%s url=%s err=%s\n", req.RemoteAddr, req.RequestURI, err)
//		return nil, err
//	}
//
//	return resp, nil
//}
