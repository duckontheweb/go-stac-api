package main

import (
	"github.com/gin-gonic/gin"

	"github.com/duckontheweb/go-stac-api/stacapi"
)

func main() {
	router := gin.Default()

	stacapi.AddSTACRoutes(router)

	router.Run()
}
