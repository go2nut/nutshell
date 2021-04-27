package etcd

import (
	"context"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"log"
	"sync"
	"time"

	"go.etcd.io/etcd/clientv3"
)

//EnvironmentDiscovery 服务发现
type EnvironmentDiscovery struct {
	cli             *clientv3.Client  //etcd client
	environmentList map[string]string //服务列表
	lock            sync.Mutex
}

//NewServiceDiscovery  新建发现服务
func NewServiceDiscovery(endpoints []string) *EnvironmentDiscovery {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	return &EnvironmentDiscovery{
		cli:             cli,
		environmentList: make(map[string]string),
	}
}

//WatchService 初始化服务列表和监视
func (s *EnvironmentDiscovery) WatchService(prefix string, initHook func(resp *clientv3.GetResponse), watchHook func(event *clientv3.Event)) error {
	//根据前缀获取现有的key
	resp, err := s.cli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, ev := range resp.Kvs {
		s.SetServiceList(string(ev.Key), string(ev.Value))
	}
	initHook(resp)

	//监视前缀，修改变更的server
	go s.watcher(prefix, watchHook)
	return nil
}

//watcher 监听前缀
func (s *EnvironmentDiscovery) watcher(prefix string, watchHook func(event *clientv3.Event)) {
	rch := s.cli.Watch(context.Background(), prefix, clientv3.WithPrefix())
	log.Printf("watching prefix:%s now...", prefix)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			watchHook(ev)
			switch ev.Type {
			case mvccpb.PUT: //修改或者新增
				s.SetServiceList(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE: //删除
				s.DelServiceList(string(ev.Kv.Key))
			}
		}
	}
}

//SetServiceList 新增服务地址
func (s *EnvironmentDiscovery) SetServiceList(key, val string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.environmentList[key] = string(val)
	log.Println("put key :", key, "val:", val)
}

//DelServiceList 删除服务地址
func (s *EnvironmentDiscovery) DelServiceList(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.environmentList, key)
	log.Println("del key:", key)
}

//GetServices 获取服务地址
func (s *EnvironmentDiscovery) GetServices() []string {
	s.lock.Lock()
	defer s.lock.Unlock()
	addrs := make([]string, 0)

	for _, v := range s.environmentList {
		addrs = append(addrs, v)
	}
	return addrs
}

//Close 关闭服务
func (s *EnvironmentDiscovery) Close() error {
	return s.cli.Close()
}
