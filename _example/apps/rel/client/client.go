package client

import (
	"context"
	"flag"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"nutshell/_example/apps/shard"
	"time"
)

var relGrpcAddr = flag.String("rel_grpc_addr", "", "app to run")

var Client shard.RelSvcClient

func Init() error {
	opts := []grpc.DialOption{ grpc.WithInsecure() , grpc.WithBlock(), grpc.WithTimeout(time.Second*30)}
	if rpcConn, err := grpc.Dial(*relGrpcAddr, opts...); err != nil {
		log.Errorf("fail connect rel server:%s err:%v", *relGrpcAddr, err)
		return err
	} else {
		Client = shard.NewRelSvcClient(rpcConn)
		log.Infof("success connect rel server:%s", *relGrpcAddr)
		go func() {
			time.Sleep(time.Second)
			d , err := Client.IsFriend(context.Background(), &shard.UserPairRequest{
				Uid1:                 100,
				Uid2:                 101,
			})
			log.Infof("is friend:%#v err:%v", d, err)
		}()
		return nil
	}
}


