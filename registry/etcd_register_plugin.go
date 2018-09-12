package registry

import (
	"github.com/SunMaybo/jewel-inject/inject"
	"net"
	"github.com/SunMaybo/go-jewel/context"
	"github.com/cihub/seelog"
)

type EtcRegisterPlugin struct {
	Client *EtcRegistry
}

func (plugin *EtcRegisterPlugin) Open(injector *inject.Injector) error {
	p := injector.Service(&JewelRegisterProperties{}).(JewelRegisterProperties)
	jewelPlugin := p.Registry.JewelPlugin
	if jewelPlugin == nil || jewelPlugin.EtcdPlugin == nil || jewelPlugin.EtcdPlugin.Enabled == nil || (jewelPlugin.EtcdPlugin.Enabled != nil && !*jewelPlugin.EtcdPlugin.Enabled) {
		return nil
	}
	jewel := injector.Service(&context.JewelProperties{}).(context.JewelProperties)
	jewelPlugin.EtcdPlugin.Name = jewel.Jewel.Name
	jewelPlugin.EtcdPlugin.Port = int(*jewel.Jewel.Server.Port)

	if jewelPlugin.EtcdPlugin.IsRefresh == nil {
		jewelPlugin.EtcdPlugin.IsRefresh = new(int32)
	}
	if jewelPlugin.EtcdPlugin.Address == nil {
		ip := getLocalIp()
		jewelPlugin.EtcdPlugin.Address = &ip
	}
	plugin.Client = jewelPlugin.EtcdPlugin
	return jewelPlugin.EtcdPlugin.register()
}
func (plugin *EtcRegisterPlugin) Health() error {
	return nil
}
func (plugin *EtcRegisterPlugin) Close() {
	seelog.Error("close etcd service register")
	plugin.Client.Down()
	plugin.Client.client.Close()
}
func (plugin *EtcRegisterPlugin) Interface() (string, interface{}) {
	return "etcd_register", plugin.Client
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
