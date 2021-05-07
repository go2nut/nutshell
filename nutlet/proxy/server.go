package proxy

import (
	log "github.com/sirupsen/logrus"
	"nutshell/nutlet/proxy/protocols"
	"os/exec"
)

func ServeEnvProxy(grpcPort int) *EnvRouter {
	router := NewEnvRouter()
	go protocols.NewGrpcProxyServer(router.Route).ServeGrpc(grpcPort)
	return router
}

func ServeAppProxy(httpPort, grpcPort int) *AppRouter {
	appRouter := NewAppRouter()
	go protocols.NewHttpProxy(appRouter.RouteHttp).ServeHTTP(httpPort)
	go protocols.NewGrpcProxyServer(appRouter.RouteGrpc).ServeGrpc(grpcPort)
	return appRouter
}

func ServeDns() {
	//p := exec.Command("sed", "-i", "1s/^/nameserver 127.0.0.1\n", "/etc/resolv.conf")
	p := exec.Command("echo 'nameserver 127.0.0.1' >> /etc/resolv.conf2 && cat /etc/resolv.conf >> /etc/resolv.conf2 && mv /etc/resolv.conf2 /etc/resolv.conf")
	err := p.Run()
	log.Infof("edit nameserver command, cmd:%s err:%v", p.String(), err)
	protocols.ServeDns(53, map[string]string{"*.nutshell": "127.0.0.1"})
}