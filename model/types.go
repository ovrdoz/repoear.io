package model

type Config struct {
	Host            string         `yaml:"host"`
	Port            string         `yaml:"port"`
	RefreshInterval int            `yaml:"refresh_interval"`
	Repositories    []Repositories `yaml:"repositories"`
}

type Repositories struct {
	Name     string `yaml:"name"`
	URL      string `yaml:"url"`
	AutoSync bool   `yaml:"auto_sync" `
	Script   string `yaml:"script"`
}
