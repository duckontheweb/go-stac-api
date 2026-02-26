package stacapi

import (
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type ApiConfig struct {
	Backend BackendConfig `yaml:"backend"`
}

func ReadConfig(config_path string) ApiConfig {

	if !filepath.IsAbs(config_path) {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal("Could not get current working directory, please provide an absolute path to the config file.")
		}
		var abs_config_path string
		abs_config_path, err = filepath.Abs(path.Join(cwd, config_path))
		if err != nil {
			log.Fatalf("Could not construct absolute path from relative path: %s", config_path)
		}
		config_path = abs_config_path
	}
	contents, err := os.ReadFile(config_path)
	if err != nil {
		log.Fatalf("Could not read config file at path %s", config_path)
	}

	config := ApiConfig{}
	err = yaml.Unmarshal(contents, &config)
	if err != nil {
		log.Fatalf("Could not parse config file at path %s as config YAML", config_path)
	}

	return config
}
