package context

import "github.com/SunMaybo/jewel-inject/inject"

type Plugin interface {
	Open(injector *inject.Injector) error
	Health() error
	Close()
	Interface() (string, interface{})
}
