package config

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"nutshell/nutlet/config/etcd"
	"os"
)

func ReadFromSysEnv() *RuntimeEnvCfg {
	cfg := &RuntimeEnvCfg{
		Env:          os.Getenv("nutshell_env"),
		Ip:           os.Getenv("nutshell_ip"),
		HttpPort:     os.Getenv("nutshell_http_port"),
		GrpcPort:     os.Getenv("nutshell_grpc_port"),
		EtcdEndpoint: os.Getenv("nutshell_etcd"),
		WorkSpace: 	  os.Getenv("nutshell_ws"),
	}
	if cfg.Env == "" {
		log.Println("nutshell_env is empty")
		cfg.Env = "test1"
	}
	if cfg.Ip == "" {
		log.Println("nutshell_host is empty")
		cfg.Ip = "127.0.0.1"
	}

	if cfg.HttpPort == "" {
		log.Println("nutshell_http_port is empty")
		cfg.HttpPort = "6701"
	}

	if cfg.GrpcPort == "" {
		log.Println("nutshell_grpc_port is empty")
		cfg.GrpcPort = "6702"
	}

	if cfg.EtcdEndpoint == "" {
		cfg.EtcdEndpoint = "127.0.0.1:2379"
	}

	if cfg.WorkSpace == "" {
		cfg.WorkSpace = "../_example"
	}
	return cfg
}

func (envCfg *RuntimeEnvCfg) WorkSpaceInfo() (configPath string, procPath string) {
	configPath = fmt.Sprintf("%s/.nutshell/apps.yaml", envCfg.WorkSpace)
	procPath = fmt.Sprintf("%s/.nutshell/apps.Procfile", envCfg.WorkSpace)
	return configPath, procPath
}


func EnvConfigFromString(env string, str string) *RuntimeEnvCfg {
	cfg := new(RuntimeEnvCfg)
	err := json.Unmarshal([]byte(str), cfg)
	if err != nil {
		log.Errorf("fail parse etcd config env:%s value:%s",  env, str)
	}
	return cfg
}


func DiscoveryEnv(localEnvCfg *RuntimeEnvCfg, add func(string, *RuntimeEnvCfg), del func(string, *RuntimeEnvCfg)) {
	// register env at etcd and register
	envCfgJson, _ := json.Marshal(localEnvCfg)
	etcd.SetupRegister(localEnvCfg.EtcdEndpoint, localEnvCfg.Env, string(envCfgJson))
	etcd.SetupDiscovery(localEnvCfg.EtcdEndpoint, func(m map[string]string) {
		for envName, envBody := range m {
			add(envName, EnvConfigFromString(envName, envBody))
		}
	}, func(action string, key string, value string) {
		switch action {
		case "put":
			add(key, EnvConfigFromString(key, value))
		case "delete":
			del(key, EnvConfigFromString(key, value))
		}
	})
}
