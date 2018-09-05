package registry

type Server struct {
	Name       string  `json:"name"`
	Host       *string `yaml:"host" xml:"host" json:"host"`
	Port       int     `json:"port"`
	Address    *string `yaml:"address" xml:"address" json:"address"`
	UseAddress *bool   `yaml:"use_address" xml:"use_address" json:"use_address"`
}
