package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"nutshell/nutlet/config"
	"nutshell/nutlet/proxy"
	"time"
)

var configPath = flag.String("config", "../_example/config.yaml", "path to service definition config")
var procPath = flag.String("proc", "../_example/procfile", "path to service definition config")

func main() {
	flag.Parse()

	//read env info
	envCfg := config.ReadFromSysEnv()
	proxy.ServeDns()

	grpc.EnableTracing = true
	appRouter := proxy.ServeAppProxy(6701, 6702)
	envRouter := proxy.ServeEnvProxy(80)

	f := func() error {
		svcDefCfg, err := config.LoadAppDefinitionsCfg(*configPath)
		if err != nil {
			return err
		}

		apps, err := config.LoadProcfile(*procPath)
		if err != nil {
			return err
		}
		// apply services into inbound router
		appRouter.Reset(svcDefCfg, apps)
		return nil
	}

	config.WatchFile(*configPath,  f)
	config.WatchFile(*procPath, f)
	envRouter.AddEnv("local", fmt.Sprintf("%s:%s", envCfg.Ip, envCfg.GrpcPort))
	config.DiscoveryEnv(envCfg, func(s string, config *config.RuntimeEnvCfg) {
		// add
		envRouter.AddEnv(config.Env, fmt.Sprintf("%s:%s", config.Ip, config.GrpcPort))
	}, func(s string, config *config.RuntimeEnvCfg) {
		// del
		envRouter.DelEnv(config.Env)
	})


	waitToQuit()
}


func waitToQuit() {
	timer := time.NewTimer(time.Second * 3)
	for {
		select {
		case <-timer.C:
			log.Info("nutlet is running...\n")
		}
		timer.Reset(time.Second * 3)
	}
	log.Info("nutlet exist...\n")
}

