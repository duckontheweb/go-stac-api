package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/duckontheweb/go-stac-api/stacapi"
)

func main() {

	config_path := flag.String("config-path", "./stac-api-config.yaml", "Path to the config file.")

	flag.Parse()

	config := stacapi.ReadConfig(*config_path)

	router := gin.Default()
	stac_api := stacapi.NewApi(config)
	defer func() {
		err := stac_api.Shutdown()
		if err != nil {
			log.Fatalf("Unable to close client for %s backend", err)
		}
	}()
	stac_api.AddToRouter(router)

	err := router.Run()
	if err != nil {
		log.Fatal(err)
	}
}
