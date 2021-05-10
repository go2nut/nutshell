package rel

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	proto "nutshell/_example/apps/shard"
)
type RelationServer struct {
}

func (s *RelationServer) IsFriend(ctx context.Context, req *proto.UserPairRequest) (*proto.IsFiendResp, error) {
	if req.Uid2 == req.Uid1 {
		return &proto.IsFiendResp{IsFriend: false}, nil
	}

	k := fmt.Sprintf("%d:%d", req.Uid1, req.Uid2)
	if req.Uid2 < req.Uid1 {
		k = fmt.Sprintf("%d:%d", req.Uid2, req.Uid1)
	}

	_, exist := db.Friends[k]
	return &proto.IsFiendResp{
	    IsFriend:             exist,
	}, nil
}


func Run(httpPort int, grpcPort int) {
	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	server := grpc.NewServer(opts...)

	proto.RegisterRelSvcServer(server, &RelationServer{})
	log.Infof("serving serve rel grpc at %s", lis.Addr().String())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("fail serve rel grpc at %d", grpcPort)
	} else {
		log.Infof("fail serve rel grpc at %d", grpcPort)
	}
}
