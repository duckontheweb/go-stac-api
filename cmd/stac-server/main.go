package main

import (
	"github.com/gin-gonic/gin"

	"github.com/duckontheweb/go-stac-api/stacapi"
)

func main() {
	router := gin.Default()

	stac_api := stacapi.STACApi{}

	core_router := stacapi.STACCoreRouter{}
	collections_router := stacapi.STACCollectionsRouter{}

	core_router.AttachTo(&stac_api)
	collections_router.AttachTo(&stac_api)

	stac_api.AttachHandlers(router)

	router.Run()
}
