package main

import (
	"encoding/json"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"nutshell/nutlet/config"
	"nutshell/nutlet/proxy"
	"os/exec"
	"time"
)

func main() {
	flag.Parse()

	//read env info
	envCfg := config.ReadFromSysEnv()
	proxy.ServeDns()
	if envCfg.EtcdEndpoint == "127.0.0.1" || envCfg.EtcdEndpoint == "localhost"{
		go func() {
			proc := exec.Command("etcd")
			proc.Run()
		}()
	}

	grpc.EnableTracing = true
	appRouter := proxy.ServeAppProxy(6701, 6702)
	envRouter := proxy.ServeEnvProxy(80)

	configPath, procPath := envCfg.WorkSpaceInfo()

	f := func() (err error) {
		defer func() {
			log.Infof("reload cfg file configPath:%s procPath:%s err:%v", configPath, procPath, err)
		}()
		svcDefCfg, err := config.LoadAppDefinitionsCfg(configPath)
		if err != nil {
			return err
		}

		apps, err := config.LoadProcfile(procPath)
		if err != nil {
			return err
		}

		svcDefCfgJ, _ := json.MarshalIndent(svcDefCfg, ">", "\t")
		println(fmt.Sprintf( "reload cfg file svc:%v apps:%v", string(svcDefCfgJ), apps))
		// apply services into inbound router
		appRouter.Reset(svcDefCfg, apps)
		return nil
	}

	config.WatchFile(configPath,  f)
	config.WatchFile(procPath, f)
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

