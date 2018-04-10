package context

type IService interface {
	Load(fun func(c Config))
}
