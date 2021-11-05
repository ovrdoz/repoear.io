package model

type Config struct {
	Refresh      int            `yaml:"refresh"`
	Repositories []Repositories `yaml:"repositories"`
}

type Repositories struct {
	Name   string `yaml:"name"`
	URL    string `yaml:"url"`
	Sync   bool   `yaml:"sync" `
	Force  bool   `yaml:"force" `
	Script string `yaml:"script"`
}
