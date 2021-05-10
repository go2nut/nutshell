package proxy

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/url"
	"nutshell/nutlet/config"
	"sync"
)

type EnvRouter struct {
	sync.RWMutex
	EnvRules map[string]string // env->protocol->address
}

func NewEnvRouter() *EnvRouter {
	return &EnvRouter{
		EnvRules: make(map[string]string, 0),
	}
}

func (router *EnvRouter) AddEnv(env, address string) {
	router.EnvRules[env] = address
	log.Infof("env router add rule, env:%s, address:%s", env, address)
}

func (router *EnvRouter) DelEnv(env string) {
	delete(router.EnvRules, env)
	log.Infof("env router del rule, env:%s", env)
}

func  (router *EnvRouter) Route(ctx context.Context, originAddr string, fullMethodName string) (*url.URL, error) {
	log.Infof("outbound router originAddr:%s fullMethodName:%s", originAddr, fullMethodName)

	env := ""
	if _, env1, err := config.ParseHost(originAddr); err != nil {
		return nil, fmt.Errorf("InvalidOriginAddr %s, err:%v", originAddr, err)
	} else {
		env = env1
	}
	target, exist := router.EnvRules[env]
	if exist == false {
		return nil, fmt.Errorf("InvalidOriginAddr %s", originAddr)
	}
	targetURL, err := url.Parse(fmt.Sprintf("http://%s%s", target, fullMethodName))

	log.Printf("Origin=%s Target=%s\n", originAddr, target)
	if err != nil {
		return nil, fmt.Errorf("InvalidTarget %s", target)
	}
	return targetURL, nil
}