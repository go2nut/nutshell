package proxy

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"net/http"
	"net/url"
	"nutshell/nutlet/config"
	"strings"
	"sync"
)

type PrefixRule struct {
	Prefix string
	Target string
}

type AppRouter struct {
	sync.RWMutex
	HttpRules []*PrefixRule
	WebsocketRules []*PrefixRule
	GrpcRules []*PrefixRule
}

func NewAppRouter() *AppRouter {
	return &AppRouter{
		HttpRules: make([]*PrefixRule, 0),
		GrpcRules: make([]*PrefixRule, 0),
	}
}

func (ir *AppRouter) Reset(appDefCfg *config.NutshellCfg, localApps []string) {
	ir.Lock()
	defer ir.Unlock()

	for _, app := range appDefCfg.Apps {
		for _, expose := range app.Exposes {
			for _, path := range expose.Paths {
				switch config.Protocol(strings.ToLower(strings.Trim(expose.Protocol, " \n\t")) ){
				case config.ProtocolHttp:
					if in(app.Name, localApps...) {
						ir.HttpRules = append(ir.HttpRules, &PrefixRule{path, config.BuildHost(app.Name, config.EnvLocal, expose.Port)})
					} else {
						if defaultAddr, defaultAddrExist := appDefCfg.FindAddr(config.EnvDefault, app.Name, config.ProtocolHttp);  defaultAddrExist {
							ir.HttpRules = append(ir.HttpRules, &PrefixRule{path, defaultAddr})
						} else {
							log.Warnf("not found target addr for http path:%s", path)
						}
					}
				case config.ProtocolWebsocket:
					log.Warnf("app:%s expose:%s port:%d path:%s", app, expose.Protocol, expose.Port, path)
					if in(app.Name, localApps...) {
						ir.WebsocketRules = append(ir.WebsocketRules, &PrefixRule{path, config.BuildHost(app.Name, config.EnvLocal, expose.Port)})
					} else {
						if defaultAddr, defaultAddrExist := appDefCfg.FindAddr(config.EnvDefault, app.Name, config.ProtocolWebsocket);  defaultAddrExist {
							ir.WebsocketRules = append(ir.WebsocketRules, &PrefixRule{path, defaultAddr})
						} else {
							log.Warnf("not found target addr for http path:%s", path)
						}
					}
				case config.ProtocolGrpc:
					if in(app.Name, localApps...) {
						ir.GrpcRules = append(ir.GrpcRules, &PrefixRule{path, config.BuildHost(app.Name, config.EnvLocal, expose.Port)})
					} else {
						if defaultAddr, defaultAddrExist := appDefCfg.FindAddr(config.EnvDefault, app.Name, config.ProtocolGrpc);  defaultAddrExist {
							ir.GrpcRules = append(ir.GrpcRules, &PrefixRule{path, defaultAddr})
						} else {
							log.Warnf("not found target addr for grpc path:%s", path)
						}
					}
				}
			}
		}
	}
	d, _ := json.MarshalIndent(ir, ">", "\t")
	println(fmt.Sprintf("app router info: %s", string(d)))
}

func (ir *AppRouter) RouteHttp(ctx context.Context, req *http.Request) (*url.URL, bool) {
	reqURL := req.URL.RequestURI()

	target := ""
	for _, rule:= range ir.HttpRules{
		if strings.HasPrefix(reqURL, rule.Prefix) {
			target = rule.Target
			break
		}
	}
	if target == "" {
		log.Printf("prefix:%s not matched\n", reqURL)
		return nil, false
	} else {
		targetURL, err := url.Parse(fmt.Sprintf("http://%s%s", target, reqURL))
		if err != nil {
			log.Printf("skip routing, target of test can not be parsed as url.URL, test ski    pped; Target=%s Err=%s\n", target, err)
			targetURL = req.URL
			return targetURL, false
		} else {
			log.Debugf("inbound route req:%s host:%s target:%s", req.URL, req.Host, targetURL)
			return targetURL, true
		}
	}
}

func (ir *AppRouter) RouteWebsocket(ctx context.Context, req *http.Request) (*url.URL, bool) {
	reqURL := req.URL.RequestURI()

	target := ""
	for _, rule:= range ir.WebsocketRules{
		if strings.HasPrefix(reqURL, rule.Prefix) {
			target = rule.Target
			break
		}
	}
	if target == "" {
		log.Printf("prefix:%s not matched\n", reqURL)
		return nil, false
	} else {
	    targetURL, err := url.Parse(fmt.Sprintf("ws://%s%s", target, reqURL))
	    if err != nil {
	    	log.Printf("skip routing, target of test can not be parsed as url.URL, test ski    pped; Target=%s Err=%s\n", target, err)
	    	targetURL = req.URL
	    	return targetURL, false
	    } else {
	    	log.Debugf("inbound route req:%s host:%s target:%s", req.URL, req.Host, targetURL)
	    	return targetURL, true
	    }
	}
}

func (ir *AppRouter) RouteGrpc(ctx context.Context, originAddr string, fullMethodName string) (*url.URL, error) {

	target := ""
	for _, rule:= range ir.GrpcRules{
		if strings.HasPrefix(fullMethodName, rule.Prefix) {
			target = rule.Target
			break
		}
	}
	targetURL, err := url.Parse(fmt.Sprintf("http://%s%s", target, fullMethodName))

	log.Printf("Origin=%s Target=%s\n", originAddr, target)
	if err != nil {
		return nil, fmt.Errorf("InvalidTarget %s", target)
	}

	return targetURL, nil
}


func in(i string, l ...string) bool {
	for _, l0 := range l {
		if i == l0 {
			return true
		}
	}
	return false
}
