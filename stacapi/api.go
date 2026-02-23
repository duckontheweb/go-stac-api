package stacapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/planetlabs/go-stac"
)

type StacApi struct {
	routers []*IStacConformance
}

func (api *StacApi) AddStacRouter(stac_router IStacConformance) {
	api.routers = append(api.routers, &stac_router)
}

func (api *StacApi) AttachHandlers(gin_router gin.IRouter) {
	for _, stac_router := range api.routers {
		(*stac_router).AttachHandlers(gin_router)
	}
}

func (api *StacApi) ConformsTo() []string {
	conforms_to := make([]string, 0)

	for _, stac_router := range api.routers {
		for _, c := range (*stac_router).ConformanceClasses() {
			conforms_to = append(conforms_to, c)
		}
	}

	return conforms_to
}

func (api *StacApi) LandingPageLinks(request *http.Request) []*stac.Link {
	links := make([]*stac.Link, 0)

	for _, stac_router := range api.routers {
		for _, l := range (*stac_router).LandingPageLinks(request) {
			links = append(links, l)
		}
	}

	return links
}

type IStacConformance interface {
	AttachHandlers(gin_router gin.IRouter)
	ConformanceClasses() []string
	LandingPageLinks(request *http.Request) []*stac.Link
}
