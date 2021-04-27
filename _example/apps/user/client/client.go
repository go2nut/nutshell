package client

import (
    "context"
    "flag"
    log "github.com/sirupsen/logrus"
    "google.golang.org/grpc"
    "nutshell/_example/apps/shard"
    "time"
)

var userGrpcAddr = flag.String("user_grpc_addr", "", "app to run")

var Client shard.UserSvcClient

func Init() error {
    opts := []grpc.DialOption{grpc.WithInsecure() , grpc.WithBlock(), grpc.WithTimeout(time.Second*30)}
    //opts := []grpc.DialOption{grpc.WithInsecure() , grpc.WithTimeout(time.Second*30)}
    if rpcConn, err := grpc.Dial(*userGrpcAddr, opts...); err != nil {
        log.Errorf("fail connect user server:%s, err:%v", *userGrpcAddr, err)
        return err
    } else {
        Client = shard.NewUserSvcClient(rpcConn)
        log.Infof("success connect user server:%s", *userGrpcAddr)
        go func() {
        	time.Sleep(time.Second)
            u, err := Client.UserInfo(context.Background(), &shard.UserReq{ ReqUser: 100}, )
            log.Infof("resp user:%v err:%v", u, err)
        }()
        return nil
    }
}
