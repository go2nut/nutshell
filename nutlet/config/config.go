package config

import (
	"errors"
	"fmt"
	"strings"
)

type Protocol string

const ProtocolHttp Protocol  = "http"
const ProtocolGrpc Protocol  = "grpc"
const ProtocolWebsocket Protocol = "ws"

const EnvDefault = "default"
const EnvLocal = "local"

type AppDefinition struct {
	Name string `json:"name"`
	Exposes []PortExpose `json:"exposes"`
}

type PortExpose struct {
	Protocol string `json:"protocol"`
	Port int `json:"port"`
	Paths []string `json:"paths"`
}

type EnvDefinition struct {
	Name string `json:"name"`
	Apps map[string]map[Protocol]string /*app, protocol, dest*/ `json:"apps"`
}

type NutshellCfg struct {
	Apps         []AppDefinition `json:"apps"`
	Environments []EnvDefinition `json:"environments"`
}

type RuntimeEnvCfg struct {
	Env          string `json:"env"`
	Ip           string `json:"ip"`
	HttpPort     string `json:"http_port"`
	GrpcPort     string `json:"grpc_port"`
	EtcdEndpoint string `json:"-"`
}

func (cfg *NutshellCfg) FindAddr(envName string, appName string, protocolName Protocol) (addr string, exist bool) {
	for _, env := range cfg.Environments {
		if env.Name == envName {
			if app, appExist := env.Apps[appName]; appExist {
				if addr, addrExist := app[protocolName]; addrExist {
					return addr, true
				}
			} else {
				return "", false
			}
		}
	}
	return "", false
}

func BuildHost(appName string, envName string, port int) string {
	return fmt.Sprintf("%s.%s.nutshell:%d", appName, envName, port)
}

func ParseHost(host string) (appName string, envName string, err error) {
	fields := strings.Split(host, ".")
	if len(fields) != 3 {
		return "", "", errors.New("not nutshell style host")
	} else if strings.Contains(fields[2], "nutshell") == false {
		return "", "", errors.New("not nutshell style host")
	} else {
		return fields[0], fields[1], nil
	}
}
