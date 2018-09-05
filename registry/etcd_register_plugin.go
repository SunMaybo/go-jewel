package registry

import (
	"github.com/SunMaybo/jewel-inject/inject"
	"net"
	"github.com/SunMaybo/go-jewel/context"
	"github.com/cihub/seelog"
)

type EtcRegisterPlugin struct {
	client *EtcRegistry
}

func (plugin EtcRegisterPlugin) Open(injector *inject.Injector) error {
	p := injector.Service(&JewelRegisterProperties{}).(JewelRegisterProperties)
	etcPlugin := p.Registry.JewelPlugin.EtcdPlugin
	if etcPlugin.Enabled != nil && !*etcPlugin.Enabled {
		return nil
	}
	jewel := injector.Service(&context.JewelProperties{}).(context.JewelProperties)
	etcPlugin.Name = jewel.Jewel.Name
	etcPlugin.Port = jewel.Jewel.Port

	if etcPlugin.IsRefresh == nil {
		etcPlugin.IsRefresh = new(int32)
	}
	if  etcPlugin.Address == nil {
		ip := getLocalIp()
		etcPlugin.Address = &ip
	}
	plugin.client=etcPlugin
	return etcPlugin.register()
}
func (plugin EtcRegisterPlugin) Health() error {
	return nil
}
func (plugin EtcRegisterPlugin) Close() {
	seelog.Error("close etcd service register")
	plugin.client.Down()
	plugin.client.client.Close()
}
func (plugin EtcRegisterPlugin) Interface() (string, interface{}) {
	return "etcd_register", plugin.client
}

func getLocalIp() (IpAddr string) {
	addrSlice, err := net.InterfaceAddrs()
	if nil != err {
		IpAddr = "localhost"
		return IpAddr
	}
	for _, addr := range addrSlice {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if nil != ipnet.IP.To4() {
				IpAddr = ipnet.IP.String()
				return
			}
		}
	}
	IpAddr = "localhost"
	return IpAddr
}
