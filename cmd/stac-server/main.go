package main

import (
	"github.com/gin-gonic/gin"

	"github.com/duckontheweb/go-stac-api/backend/duckdb"
	"github.com/duckontheweb/go-stac-api/stacapi"
)

func main() {
	router := gin.Default()

	stac_api := stacapi.STACApi{}

	core_router := stacapi.CoreRouter{}
	core_router.AttachTo(&stac_api)

	backend := duckdb.NewBackend("./stac-api-config.yaml")
	collections_router := stacapi.NewCollectionsRouter(&backend)
	collections_router.AttachTo(&stac_api)

	stac_api.AttachHandlers(router)

	router.Run()
}
