package rel

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	proto "nutshell/_example/apps/shard"
)
type RelationServer struct {
}

func (s *RelationServer) Friends(ctx context.Context, req *proto.UserReq) (*proto.UserRsp, error) {
	uid := req.ReqUser
	if _, exist := db.Users[uid]; exist  {
		if frs, frsExist := db.Friends[uid]; frsExist {
			rsp := &proto.UserRsp{
				OtherUsers:     make([]*proto.User, 0),
			}
			for _, fr := range frs {
				if name, nameExist := db.Users[fr]; nameExist {
					rsp.OtherUsers = append(rsp.OtherUsers, &proto.User{
						UserId:               fr,
						Name:             name,
					})
				}
			}
			return rsp, nil
		} else {
			return &proto.UserRsp{
				OtherUsers:     make([]*proto.User, 0),
			}, nil
		}
	} else {
		return nil, errors.New(fmt.Sprintf("request user:%d not exist.", uid))
	}
}

func (s *RelationServer) IsFriend(ctx context.Context, req *proto.UserPairRequest) (*proto.IsFiendResp, error) {
	if _, name1Exist := db.Users[req.Uid1]; name1Exist {
		if _, name2Exist := db.Users[req.Uid2]; name2Exist {
			if frs, frsExist := db.Friends[req.Uid1]; frsExist {
				for _, fr := range frs {
					if fr == req.Uid2 {
						return &proto.IsFiendResp{
							IsFriend:             true,
						}, nil
					}
				}
			}
			return &proto.IsFiendResp{
				IsFriend:             false,
			}, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("user not exist"))
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
