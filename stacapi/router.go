package stacapi

import (
	"net/http"

	"github.com/duckontheweb/go-stac-api/internal"
	"github.com/gin-gonic/gin"
	"github.com/oaswrap/spec/adapter/ginopenapi"
	"github.com/oaswrap/spec/option"
	"github.com/planetlabs/go-stac"
)

func AddSTACRoutes(r gin.IRouter) {
	api := ginopenapi.NewRouter(r,
		option.WithTitle("STAC API - Go"),
		option.WithVersion(internal.APIVersion),
		option.WithDescription("STAC API implementation in Go."),
		option.WithOpenAPIVersion("3.0.3"),
		option.WithServer("http://localhost:8080"),
		option.WithSpecPath(internal.ServiceDescPath),
		option.WithDocsPath(internal.ServiceDocPath),
	)

	api.GET("/", internal.HandleLandingPage).With(
		option.Summary("Landing Page"),
		option.Response(http.StatusOK, new(stac.Catalog)),
	)

}
