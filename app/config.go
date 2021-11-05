package app

import (
	"fmt"
	"os"

	model "repoear.io/model"

	yaml "gopkg.in/yaml.v2"
)

func LoadConfiguration() (model.Config, error) {

	f, err := os.Open("config/config.yml")
	if err != nil {
		return model.Config{}, fmt.Errorf("os.Open() failed with %s", err)
	}
	defer f.Close()

	dec := yaml.NewDecoder(f)

	var config model.Config
	err = dec.Decode(&config)
	if err != nil {
		return model.Config{}, fmt.Errorf("dec.Decode() failed with %s", err)
	}

	return validator(config)
}

func validator(config model.Config) (model.Config, error) {

	if config.Refresh < 30 {
		return model.Config{}, fmt.Errorf("refresh_interval is mandatory, and cannot be less than 30 seconds")
	}
	if len(config.Repositories) > 0 {
		for _, element := range config.Repositories {
			if element.URL == "" {
				return model.Config{}, fmt.Errorf("the repository URL is riquired")
			}
		}
	} else {
		return model.Config{}, fmt.Errorf("the repositories list is riquired")
	}

	return config, nil
}
