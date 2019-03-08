package registry

import (
	"testing"
	"log"
	"context"
	"time"
	"fmt"
	"github.com/etcd-io/etcd/clientv3"
)

func TestRegister(t *testing.T) {

}
func TestLease(t *testing.T) {
	client, err := clientv3.NewFromURL("http://localhost:2379")
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	leaseResp, err := client.Lease.Grant(ctx, 90)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", leaseResp)
	putResp, err := client.Put(ctx, "test/128.0.0.1", "Hello World !", clientv3.WithLease(leaseResp.ID))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", putResp)
	cancel()
}
func TestReleaseGet(t *testing.T) {
	client, err := clientv3.NewFromURL("http://localhost:2379")
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	getResp, err := client.Get(ctx, "test", clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", getResp)
	cancel()
}
func TestLeaseKeepAlive(t *testing.T) {
	client, err := clientv3.NewFromURL("http://localhost:2379")
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	leaseResp, err := client.Grant(ctx, 3)
	client.Put(ctx, "test/va", "Hello", clientv3.WithLease(leaseResp.ID))
	resp,err:=client.KeepAlive(context.TODO(),leaseResp.ID)
	if err != nil {
		log.Fatal(err)
	}
	for  {
		if resp,ok:=<-resp;ok {
			fmt.Println(resp.ID)
		}
	}
}
