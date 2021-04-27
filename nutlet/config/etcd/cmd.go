package etcd

import (
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"log"
	"time"
)

func SetupRegister(endpoint string, key string, value string) {
	ser, err := NewServiceRegister([]string{endpoint}, fmt.Sprintf("/nutshell/envs/%s", key), value, 5)
	if err != nil {
		log.Fatalln(err)
	}
	//监听续租相应chan
	go ser.ListenLeaseRespChan()
}

func SetupDiscovery(endpoint string, load func(map[string]string), watch func(action string, key string, value string)) {
	ser := NewServiceDiscovery([]string{endpoint})
	defer ser.Close()
	ser.WatchService("/nutshell/envs", func(resp *clientv3.GetResponse) {
		m := make(map[string]string, 0)
		for _, kv := range resp.Kvs {
			m[string(kv.Key)] = string(kv.Value)
		}
		load(m)
	}, func(event *clientv3.Event) {
		switch event.Type {
		case mvccpb.PUT:
			watch("put", string(event.Kv.Key), string(event.Kv.Value))
		case mvccpb.DELETE:
			watch("delete", string(event.Kv.Key), "")
		default:
			log.Println("unknown watch event type:", event.Type)
		}
	})
	for {
		select {
		case <-time.Tick(10 * time.Second):
			log.Println(ser.GetServices())
		}
	}
}
