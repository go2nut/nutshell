package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"nutshell/nutlet/proxy/protocols"
)

var (
	flagBackend = flag.String("backend", "ws://127.0.0.1:8002", "Backend URL for proxying")
)

func main() {
	u, err := url.Parse(*flagBackend)
	if err != nil {
		log.Fatalln(err)
	}

	err = http.ListenAndServe(":8081", protocols.NewProxy(u))
	if err != nil {
		log.Fatalln(err)
	}
}
