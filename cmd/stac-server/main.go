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

	config_path := flag.String("config-path", "./stac-api-config.yaml", "Path to the config file.")

	flag.Parse()

	config := readConfig(*config_path)

	var backend stacapi.Backend
	var err error
	switch config.Backend.Type {
	default:
		log.Fatalf("Unsupported backend type: '%s'", config.Backend.Type)
	case duckdb.BackendType:
		backend, err = duckdb.NewBackend(config.Backend)
		if err != nil {
			log.Fatal(err)
		}
	}

	router := gin.Default()
	stac_api := stacapi.NewApi(config, backend)
	stac_api.AddToRouter(router)

	err = router.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func readConfig(config_path string) stacapi.ApiConfig {

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

	config := stacapi.ApiConfig{}
	err = yaml.Unmarshal(contents, &config)
	if err != nil {
		log.Fatalf("Could not parse config file at path %s as config YAML", config_path)
	}

	return config
}
