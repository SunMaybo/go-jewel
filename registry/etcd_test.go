package registry

import (
	"testing"
	"go.etcd.io/etcd/clientv3"
	"log"
	"context"
	"time"
)

func TestRegister(t *testing.T) {

}
func TestLease(t *testing.T) {
	client, err := clientv3.NewFromURL("http://localhost:2379")
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	leaseResp, err :=client.Lease.Grant(ctx,90)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", leaseResp)
    putResp,err:=client.Put(ctx,"test/128.0.0.1","Hello World !",clientv3.WithLease(leaseResp.ID))
    if err!=nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", putResp)
    cancel()
}
func TestReleaseGet(t *testing.T)  {
	client, err := clientv3.NewFromURL("http://localhost:2379")
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	getResp,err:=client.Get(ctx,"test",clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", getResp)
	cancel()
}
