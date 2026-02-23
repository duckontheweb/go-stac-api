package stacapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oaswrap/spec/adapter/ginopenapi"
	"github.com/oaswrap/spec/option"
	"github.com/planetlabs/go-stac"
)

type CoreRouter struct {
	api *STACApi
}

func (r *CoreRouter) AttachTo(api *STACApi) {
	api.AddSTACRouter(r)
	r.api = api
}

func (r CoreRouter) ConformanceClasses() []string {
	return []string{STACAPICoreConformanceURI}
}

func (r CoreRouter) LandingPageLinks(request *http.Request) []*stac.Link {
	root_link := RootLink(request)
	self_link := SelfLink(request, ApplicationJSONType)
	service_desc_link := ServiceDescLink(request)
	service_doc_link := ServiceDocLink(request)
	return []*stac.Link{
		&root_link,
		&self_link,
		&service_desc_link,
		&service_doc_link,
	}
}

func (r CoreRouter) AttachHandlers(gr gin.IRouter) {
	api := ginopenapi.NewRouter(gr,
		option.WithTitle("STAC API - Go"),
		option.WithVersion(APIVersion),
		option.WithDescription("STAC API implementation in Go."),
		option.WithOpenAPIVersion("3.0.3"),
		option.WithServer("http://localhost:8080"),
		option.WithSpecPath(ServiceDescPath),
		option.WithDocsPath(ServiceDocPath),
	)

	api.GET(RootPath, r.HandleLandingPage).With(
		option.Summary("Landing Page"),
		option.Response(http.StatusOK, new(stac.Catalog)),
	)
}

func (r CoreRouter) HandleLandingPage(c *gin.Context) {

	landing_page := stac.Catalog{
		Version:     STACVersion,
		Id:          "go-stac-api",
		ConformsTo:  r.api.ConformsTo(),
		Description: "Implementation of the STAC API Spec in Go",
		Links:       r.api.LandingPageLinks(c.Request),
	}
	c.JSON(http.StatusOK, landing_page)
}
