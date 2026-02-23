package main

import (
	"flag"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-yaml"

	"github.com/duckontheweb/go-stac-api/backend/duckdb"
	"github.com/duckontheweb/go-stac-api/stacapi"
)

func main() {

	config_path := flag.String("config_path", "./stac-api-config.yaml", "Path to the config file.")
	config := readConfig(*config_path)

	var backend stacapi.Backend
	var err error
	switch config.BackendConfig.Type {
	default:
		log.Fatalf("Unsupported backend type: '%s'", config.BackendConfig.Type)
	case duckdb.BackendType:
		backend, err = duckdb.NewBackend(config.BackendConfig)
		if err != nil {
			log.Fatal(err)
		}
	}

	router := gin.Default()
	stac_api := stacapi.NewApi(config, backend)
	stac_api.AddToRouter(router)

	router.Run()
}

func readConfig(config_path string) stacapi.ApiConfig {

	if !filepath.IsAbs(config_path) {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal("Could not get current working directory, please provide an absolute path to the config file.")
		}
		config_path, err = filepath.Abs(path.Join(cwd, config_path))
	}
	contents, err := os.ReadFile(config_path)
	if err != nil {
		log.Fatalf("Could not read config file at path %s", config_path)
	}

	config := stacapi.ApiConfig{}
	err = yaml.Unmarshal(contents, &config)
	if err != nil {
		log.Fatalf("Could not parse config file at path %s as config YAML", config_path)
	}

	return config
}
