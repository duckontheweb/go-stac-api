package stacapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oaswrap/spec/adapter/ginopenapi"
	"github.com/oaswrap/spec/option"
	"github.com/planetlabs/go-stac"
)

const StacCoreConformanceURI = "https://api.stacspec.org/v1.0.0/core"

const STACVersion = "1.1.0"

type StacCoreConformance struct {
	api *Api
}

func (r *StacCoreConformance) AttachTo(api *Api) {
	api.AddConformance(r)
	r.api = api
}

func (r StacCoreConformance) ConformanceClasses() []string {
	return []string{StacCoreConformanceURI}
}

func (r StacCoreConformance) LandingPageLinks(request *http.Request) []*stac.Link {
	root_link := rootLink(request)
	self_link := selfLink(request, ApplicationJSONType)
	service_desc_link := serviceDescLink(request)
	service_doc_link := serviceDocLink(request)
	return []*stac.Link{
		&root_link,
		&self_link,
		&service_desc_link,
		&service_doc_link,
	}
}

func (r StacCoreConformance) AttachHandlers(router ginopenapi.Generator) {
	router.GET(RootPath, r.HandleLandingPage).With(
		option.Summary("Landing Page"),
		option.Response(http.StatusOK, new(stac.Catalog)),
	)
}

func (r StacCoreConformance) HandleLandingPage(c *gin.Context) {
	landing_page := stac.Catalog{
		Version:     STACVersion,
		Id:          "go-stac-api",
		ConformsTo:  r.api.ConformsTo(),
		Description: "Implementation of the STAC API Spec in Go",
		Links:       r.api.LandingPageLinks(c.Request),
	}
	c.JSON(http.StatusOK, landing_page)
}
