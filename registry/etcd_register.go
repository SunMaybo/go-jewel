package registry

import (
	"time"
	"strings"
	"crypto/tls"
	"github.com/cihub/seelog"
	"context"
	"go.etcd.io/etcd/clientv3"
	"strconv"
	"encoding/json"
	"errors"
	"sync/atomic"
)

type JewelRegisterProperties struct {
	Registry *RegisterProperties `yaml:"jewel" xml:"jewel" json:"jewel"`
}

type RegisterProperties struct {
	JewelPlugin *JewelPlugin `yaml:"register" xml:"register" json:"register"`
}

type JewelPlugin struct {
	EtcdPlugin *EtcRegistry `yaml:"etcd_plugin" xml:"etcd_plugin" json:"etcd_plugin"`
}

type EtcRegistry struct {
	Urls    *string `yaml:"urls" xml:"urls" json:"urls"`
	Enabled *bool   `yaml:"enabled" xml:"enabled" json:"enabled"`

	AutoSyncInterval *int64 `json:"auto-sync-interval" xml:"auto-sync-interval" yaml:"auto-sync-interval"`

	// DialTimeout is the timeout for failing to establish a connection.
	DialTimeout *int64 `json:"dial-timeout" xml:"dial-timeout" yaml:"dial-timeout"`

	// DialKeepAliveTime is the time after which client pings the server to see if
	// transport is alive.
	DialKeepAliveTime *int64 `json:"dial-keep-alive-time" xml:"dial-keep-alive-time" yaml:"dial-keep-alive-time"`

	// DialKeepAliveTimeout is the time that the client waits for a response for the
	// keep-alive probe. If the response is not received in this time, the connection is closed.
	DialKeepAliveTimeout *int64 `json:"dial-keep-alive-timeout" xml:"dial-keep-alive-timeout" yaml:"dial-keep-alive-timeout"`

	// TLS holds the client secure credentials, if any.
	InsecureSkipVerify *bool `json:"insecure_skip_verify" xml:"insecure_skip_verify" yaml:"insecure_skip_verify"`

	// Username is a user name for authentication.
	Username *string `json:"username" xml:"username" yaml:"username"`

	// Password is a password for authentication.
	Password *string `json:"password" xml:"password" yaml:"password"`

	// RejectOldCluster when set will refuse to create a client against an outdated cluster.
	RejectOldCluster *bool  `json:"reject-old-cluster" xml:"reject-old-cluster" yaml:"reject-old-cluster"`
	RefreshTimeOut   *int64 `json:"refresh_timeout" xml:"refresh_timeout" yaml:"refresh_timeout"`
	IsRefresh        *int32
	Server
	client           *clientv3.Client
}

func (etcPlugin EtcRegistry) refresh(id interface{}) {
	if atomic.LoadInt32(etcPlugin.IsRefresh) != 0 {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	leaveKeepAliveResp, err := etcPlugin.client.KeepAlive(ctx, id.(clientv3.LeaseID))
	if err != nil {
		seelog.Errorf("lease_id:%x,renewed failed :%s", id, err.Error())
	}
	for {
		if resp, ok := <-leaveKeepAliveResp; ok {
			seelog.Infof("lease_id:%x,renewed success...", resp.ID)
		}
		time.Sleep(time.Duration(*etcPlugin.RefreshTimeOut * 1000000))
	}

}
func (etcPlugin EtcRegistry) register() error {
	cfg := clientv3.Config{}
	if etcPlugin.Urls != nil {
		cfg.Endpoints = strings.Split(*etcPlugin.Urls, ",")
	} else {
		return errors.New("urls is  required")
	}

	if etcPlugin.Username != nil {
		cfg.Username = *etcPlugin.Username
	}
	if etcPlugin.Password != nil {
		cfg.Password = *etcPlugin.Password
	}
	if etcPlugin.DialTimeout != nil {
		cfg.DialTimeout = time.Duration(*etcPlugin.DialTimeout * 1000000)
	}
	if etcPlugin.AutoSyncInterval != nil {
		cfg.AutoSyncInterval = time.Duration(*etcPlugin.AutoSyncInterval * 1000000)
	}
	if etcPlugin.DialKeepAliveTime != nil {
		cfg.DialKeepAliveTime = time.Duration(*etcPlugin.DialKeepAliveTime * 1000000)
	}
	if etcPlugin.InsecureSkipVerify != nil && *etcPlugin.InsecureSkipVerify {
		cfg.TLS = &tls.Config{InsecureSkipVerify: true}
	}
	if etcPlugin.RejectOldCluster != nil {
		cfg.RejectOldCluster = *etcPlugin.RejectOldCluster
	}
	client, err := clientv3.New(cfg)
	if err != nil {
		return err
	}
	etcPlugin.client = client
	id, err := etcPlugin.Up()
	if err != nil {
		return err
	}
	go func() {
		etcPlugin.refresh(id)
	}()
	return nil
}
func (etcPlugin EtcRegistry) Up() (interface{}, error) {
	//register
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var key string
	if (etcPlugin.UseAddress == nil || !*etcPlugin.UseAddress) && etcPlugin.Host != nil {
		key = etcPlugin.Name + "/" + *etcPlugin.Host + ":" + strconv.FormatInt(int64(etcPlugin.Port), 10)
	} else {
		key = etcPlugin.Name + "/" + *etcPlugin.Address + ":" + strconv.FormatInt(int64(etcPlugin.Port), 10)
	}
	serverBuff, err := json.Marshal(etcPlugin.Server)
	if err != nil {
		return nil, err
	}
	if etcPlugin.RefreshTimeOut == nil {
		etcPlugin.RefreshTimeOut = new(int64)
		*etcPlugin.RefreshTimeOut = 30000
	}
	leaseResp, err := etcPlugin.client.Lease.Grant(ctx, *etcPlugin.RefreshTimeOut/1000*3)
	if err != nil {
		return nil, err
	}
	_, err = etcPlugin.client.Put(ctx, key, string(serverBuff), clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return nil, err
	}
	atomic.CompareAndSwapInt32(etcPlugin.IsRefresh, 1, 0)
	seelog.Infof("successful registration to the %s", *etcPlugin.Urls)
	return leaseResp.ID, nil
}
func (etcPlugin EtcRegistry) Down() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var key string
	if (etcPlugin.UseAddress == nil || !*etcPlugin.UseAddress) && etcPlugin.Host != nil {
		key = etcPlugin.Name + "/" + *etcPlugin.Host + ":" + strconv.FormatInt(int64(etcPlugin.Port), 10)
	} else {
		key = etcPlugin.Name + "/" + *etcPlugin.Address + ":" + strconv.FormatInt(int64(etcPlugin.Port), 10)
	}

	_, err := etcPlugin.client.Delete(ctx, key)
	if err != nil {
		seelog.Error(err)
	}
	atomic.CompareAndSwapInt32(etcPlugin.IsRefresh, 0, 1)
	seelog.Infof("the service was  logged from the %s", *etcPlugin.Urls)
	return nil
}
