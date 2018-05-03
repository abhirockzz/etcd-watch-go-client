package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	etcdHost := flag.String("etcdHost", "192.168.99.100:2379", "etcd host")
	etcdWatchKey := flag.String("etcdWatchKey", "foo", "etcd key to watch")

	flag.Parse()

	fmt.Println("connecting to etcd - " + *etcdHost)

	etcd, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://" + *etcdHost},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("connected to etcd - " + *etcdHost)

	defer etcd.Close()

	rch := etcd.Watch(context.Background(), *etcdWatchKey)
	fmt.Println("set WATCH on " + *etcdWatchKey)

	go func() {
		fmt.Println("started goroutine for PUT...")
		for {
			etcd.Put(context.Background(), *etcdWatchKey, time.Now().String())
			fmt.Println("populated " + *etcdWatchKey + " with a value..")
			time.Sleep(2 * time.Second)
		}

	}()

	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("Event received! %s executed on %q with value %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
}
