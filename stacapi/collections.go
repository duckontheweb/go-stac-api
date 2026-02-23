package stacapi

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/go-viper/mapstructure/v2"
	"github.com/oaswrap/spec/adapter/ginopenapi"
	"github.com/oaswrap/spec/option"
	"github.com/planetlabs/go-stac"
)

const StacCollectionsConformanceURI = "https://api.stacspec.org/v1.0.0/collections"

const (
	ListCollectionsPath = "/collections"
	GetCollectionPath   = "/collections/:id"
)

type StacCollectionsConformance struct {
	api     *Api
	backend ICollectionsBackend
}

func NewStacCollectionsConformance(backend ICollectionsBackend) StacCollectionsConformance {
	return StacCollectionsConformance{
		backend: backend,
	}
}

func (r *StacCollectionsConformance) AttachTo(api *Api) {
	api.AddConformance(r)
	r.api = api
}

func (r StacCollectionsConformance) ConformanceClasses() []string {
	return []string{StacCollectionsConformanceURI}
}

func (r StacCollectionsConformance) LandingPageLinks(request *http.Request) []*stac.Link {
	data_link := dataLink(request, ListCollectionsPath)
	return []*stac.Link{&data_link}
}

func (r StacCollectionsConformance) AttachHandlers(router ginopenapi.Generator) {
	router.GET(ListCollectionsPath, r.HandleListCollections).With(
		option.Summary("List Collections"),
		option.Response(http.StatusOK, new([]stac.Collection)),
	)
	router.GET(GetCollectionPath, r.HandleGetCollection).With(
		option.Summary("Get Collection"),
		option.Request(new(GetCollectionRequest)),
		option.Response(http.StatusOK, new(map[string]any)),
	)
}

func (r StacCollectionsConformance) HandleListCollections(c *gin.Context) {
	root_link := rootLink(c.Request)
	self_link := selfLink(c.Request, ApplicationJSONType)

	collections := r.backend.ListCollections()
	collections_list := map[string]any{
		"collections": collections,
		"links":       []*stac.Link{&root_link, &self_link},
	}
	c.JSON(http.StatusOK, collections_list)
}

func (r StacCollectionsConformance) HandleGetCollection(c *gin.Context) {
	var req GetCollectionRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request."})
	}
	root_link := rootLink(c.Request)
	self_link := selfLink(c.Request, ApplicationJSONType)
	parent_link := parentLink(c.Request)

	raw_collection, err := r.backend.GetCollection(req.Id)

	if err != nil {
		switch err.(type) {
		case CollectionNotFoundError:
			c.JSON(http.StatusNotFound, err.Error())
		default:
			c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error."})
		}
		return
	}
	raw_collection["version"] = raw_collection["stac_version"]
	collection := stac.Collection{}
	err = mapstructure.Decode(raw_collection, &collection)

	response_links := []*stac.Link{&root_link, &self_link, &parent_link}
	reserved_rels := []string{RootRel, SelfRel, ParentRel}
	for _, link := range collection.Links {
		if !slices.Contains(reserved_rels, link.Rel) {
			response_links = append(response_links, link)
		}
	}
	collection.Links = response_links

	c.JSON(http.StatusOK, collection)

}

type GetCollectionRequest struct {
	Id string `path:"id" uri:"id" required:"true"`
}
type ICollectionsBackend interface {
	ListCollections() []map[string]any
	GetCollection(id string) (map[string]any, error)
}
