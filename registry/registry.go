package registry

type Registry interface {
	refresh(id interface{})
	register() error
	Up() (interface{}, error)
	Down() error
	Services() ([]Server, error)
}
