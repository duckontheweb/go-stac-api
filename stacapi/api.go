package stacapi

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oaswrap/spec/adapter/ginopenapi"
	"github.com/oaswrap/spec/option"
	"github.com/planetlabs/go-stac"
	_ "gopkg.in/yaml.v3"
)

const APIVersion = "0.1.0"

type Api struct {
	config         ApiConfig
	backend_client BackendClient
	conformances   []*IStacConformance
}

func NewApi(config ApiConfig) Api {
	backend_name := config.Backend.Type

	backend, exists := backends[backend_name]
	if !exists {
		log.Fatalf("Backend with name %s not installed", backend_name)
	}
	backend_client, err := backend.GetClient(config.Backend)
	if err != nil {
		log.Fatalf("Unable to create client for %s backend: %s", backend_name, err)
	}

	api := Api{config: config, backend_client: backend_client}

	core_conformance := StacCoreConformance{}
	core_conformance.AttachTo(&api)

	var ok bool

	var collections_backend ICollectionsBackend
	collections_backend, ok = backend_client.(ICollectionsBackend)
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
		conforms_to = append(conforms_to, (*stac_router).ConformanceClasses()...)
	}

	return conforms_to
}

func (api *Api) LandingPageLinks(request *http.Request) []*stac.Link {
	links := make([]*stac.Link, 0)

	for _, stac_router := range api.conformances {
		links = append(links, (*stac_router).LandingPageLinks(request)...)
	}

	return links
}

func (api Api) Shutdown() error {
	return api.backend_client.Close()
}

type IStacConformance interface {
	AttachHandlers(router ginopenapi.Generator)
	ConformanceClasses() []string
	LandingPageLinks(request *http.Request) []*stac.Link
}
