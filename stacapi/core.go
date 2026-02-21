package stacapi

import (
	"net/http"

	"github.com/duckontheweb/go-stac-api/internal"
	"github.com/gin-gonic/gin"
	"github.com/oaswrap/spec/adapter/ginopenapi"
	"github.com/oaswrap/spec/option"
	"github.com/planetlabs/go-stac"
)

type STACCoreRouter struct {
	api *STACApi
}

func (r *STACCoreRouter) AttachTo(api *STACApi) {
	api.AddSTACRouter(r)
	r.api = api
}

func (r STACCoreRouter) ConformanceClasses() []string {
	return []string{internal.STACAPICoreConformanceURI}
}

func (r STACCoreRouter) LandingPageLinks(request *http.Request) []*stac.Link {
	root_link := internal.RootLink(request)
	self_link := internal.SelfLink(request, internal.ApplicationJSONType)
	service_desc_link := internal.ServiceDescLink(request)
	service_doc_link := internal.ServiceDocLink(request)
	return []*stac.Link{
		&root_link,
		&self_link,
		&service_desc_link,
		&service_doc_link,
	}
}

func (r STACCoreRouter) AttachHandlers(gr gin.IRouter) {
	api := ginopenapi.NewRouter(gr,
		option.WithTitle("STAC API - Go"),
		option.WithVersion(internal.APIVersion),
		option.WithDescription("STAC API implementation in Go."),
		option.WithOpenAPIVersion("3.0.3"),
		option.WithServer("http://localhost:8080"),
		option.WithSpecPath(internal.ServiceDescPath),
		option.WithDocsPath(internal.ServiceDocPath),
	)

	api.GET(internal.RootPath, r.HandleLandingPage).With(
		option.Summary("Landing Page"),
		option.Response(http.StatusOK, new(stac.Catalog)),
	)
}

func (r STACCoreRouter) HandleLandingPage(c *gin.Context) {

	landing_page := stac.Catalog{
		Version:     internal.STACVersion,
		Id:          "go-stac-api",
		ConformsTo:  r.api.ConformsTo(),
		Description: "Implementation of the STAC API Spec in Go",
		Links:       r.api.LandingPageLinks(c.Request),
	}
	c.JSON(http.StatusOK, landing_page)
}
