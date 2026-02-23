package stacapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oaswrap/spec/adapter/ginopenapi"
	"github.com/oaswrap/spec/option"
	"github.com/planetlabs/go-stac"
	_ "gopkg.in/yaml.v3"
)

const APIVersion = "0.1.0"

type Api struct {
	config       ApiConfig
	backend      Backend
	conformances []*IStacConformance
	gin_router   *gin.IRouter
}

func NewApi(config ApiConfig, backend Backend) Api {
	api := Api{config: config, backend: backend}

	core_conformance := StacCoreConformance{}
	core_conformance.AttachTo(&api)

	var ok bool

	var collections_backend ICollectionsBackend
	collections_backend, ok = backend.(ICollectionsBackend)
	if ok {
		collections_conformance := NewStacCollectionsConformance(collections_backend)
		collections_conformance.AttachTo(&api)
	}

	return api
}

func (api *Api) AddConformance(stac_router IStacConformance) {
	api.conformances = append(api.conformances, &stac_router)
}

func (api *Api) AddToRouter(gr gin.IRouter) {
	openapi := ginopenapi.NewRouter(gr,
		option.WithTitle("STAC API - Go"),
		option.WithVersion(APIVersion),
		option.WithDescription("STAC API implementation in Go."),
		option.WithOpenAPIVersion("3.0.3"),
		option.WithServer("http://localhost:8080"),
		option.WithSpecPath(ServiceDescPath),
		option.WithDocsPath(ServiceDocPath),
	)

	for _, conformance := range api.conformances {
		(*conformance).AttachHandlers(openapi)
	}
}

func (api *Api) ConformsTo() []string {
	conforms_to := make([]string, 0)

	for _, stac_router := range api.conformances {
		for _, c := range (*stac_router).ConformanceClasses() {
			conforms_to = append(conforms_to, c)
		}
	}

	return conforms_to
}

func (api *Api) LandingPageLinks(request *http.Request) []*stac.Link {
	links := make([]*stac.Link, 0)

	for _, stac_router := range api.conformances {
		for _, l := range (*stac_router).LandingPageLinks(request) {
			links = append(links, l)
		}
	}

	return links
}

type ApiConfig struct {
	BackendConfig `yaml:"backend"`
}

type Backend interface{}

type BackendConfig struct {
	Type             string `yaml:"type"`
	ConnectionString string `yaml:"connection_string"`
}

type IStacConformance interface {
	AttachHandlers(router ginopenapi.Generator)
	ConformanceClasses() []string
	LandingPageLinks(request *http.Request) []*stac.Link
}
