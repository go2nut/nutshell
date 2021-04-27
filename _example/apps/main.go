package main

import (
	"flag"
	"nutshell/_example/apps/im"
	"nutshell/_example/apps/rel"
	relCli "nutshell/_example/apps/rel/client"
	"nutshell/_example/apps/user"
	userCli "nutshell/_example/apps/user/client"
)

var app = flag.String("app", "app", "app to run")
var httpPort = flag.Int("http_port", 8080, "http service address")
var grpcPort = flag.Int("grpc_port", 9090, "grpc service address")


func main() {
	flag.Parse()
	switch *app {
	case "user":
		user.Run(*httpPort, *grpcPort)
	case "rel":
		userCli.Init()
		rel.Run(*httpPort, *grpcPort)
	case "im":
		userCli.Init()
		relCli.Init()
		im.Run(*httpPort, *grpcPort)
	}
}
