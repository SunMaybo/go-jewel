package registry

type Server struct {
	Name       string
	Host       *string `yaml:"host" xml:"host" json:"host"`
	Port       int
	Address    *string `yaml:"address" xml:"address" json:"address"`
	UseAddress *bool   `yaml:"use_address" xml:"use_address" json:"use_address"`
}
