package context

type Cron struct {
	Name string
	Cron string
	Fun  func()
}
