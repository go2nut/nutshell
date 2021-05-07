package proxy

import (
	"nutshell/nutlet/proxy/protocols"
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
	protocols.ServeDns(53, map[string]string{"*.nutshell": "127.0.0.1"})
}