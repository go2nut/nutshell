package user

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	"nutshell/_example/apps/shard"
	"strconv"
	"strings"
)

type InfoGrpcService struct {
}

func (svc *InfoGrpcService) UserInfo(ctx context.Context, req *shard.UserReq) (*shard.User, error) {
	if user, exist := userIdIdx[req.ReqUser]; exist {
		return &shard.User{
			UserId:               user.Id,
			Name:                 user.NickName,
			Birthday:             user.Birthday,
			Gender:               string(user.Gender),
		}, nil
	}
	return nil, errors.New(fmt.Sprintf("uid:%d not exist", req.ReqUser))
}

func (svc *InfoGrpcService) UserByToken(ctx context.Context, req *shard.TokenReq) (*shard.User, error) {
	uid, _ := strconv.ParseInt(strings.ReplaceAll(req.Token, "token_", ""), 10, 64)
	if user, exist := userIdIdx[uid]; exist {
		return &shard.User{
			UserId:               user.Id,
			Name:                 user.NickName,
			Birthday:             user.Birthday,
			Gender:               string(user.Gender),
		}, nil
	}
	return nil, errors.New(fmt.Sprintf("uid:%d not exist", req.Token))
}


func handleHttpLogin(c *gin.Context) {
	type EmailPasswd struct {
		Email string
		Passwd string
	}
	ep := new(EmailPasswd)
	if err := c.BindJSON(ep); err != nil {
		c.String(400, "error params")
	} else {
		log.Printf("login: input:%v\n", ep)
		user := userIdx[ep.Email]
		if user.Passwd == ep.Passwd {
			c.JSON(200, map[string]interface{}{"user": user, "token": user.Token()})
		} else {
			c.JSON(400, "unauthorized")
		}
	}
}



func Run(httpPort int, grpcPort int) {

	go func() {
		r := gin.Default()
		r.POST("/login", handleHttpLogin)
		r.Run(fmt.Sprintf("0.0.0.0:%d", httpPort)) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	shard.RegisterUserSvcServer(grpcServer, &InfoGrpcService{})
	log.Infof("serving user grpc at %s", lis.Addr().String())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("fail serve user grpc at %d, err:%v", grpcPort, err)
	} else {
		log.Infof("close serve user grpc at %d", grpcPort)
	}
}
