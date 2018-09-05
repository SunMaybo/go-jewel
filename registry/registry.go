package registry

type Registry interface {
	refresh(id interface{}) error
	register() error
	Up() (interface{}, error)
	Down() error
}
