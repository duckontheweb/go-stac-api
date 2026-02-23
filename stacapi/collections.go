package stacapi

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/go-viper/mapstructure/v2"
	"github.com/planetlabs/go-stac"
)

type CollectionsRouter struct {
	api     *STACApi
	backend ICollectionsBackend
}

func NewCollectionsRouter(backend ICollectionsBackend) CollectionsRouter {
	return CollectionsRouter{
		backend: backend,
	}
}

func (r *CollectionsRouter) AttachTo(api *STACApi) {
	api.AddSTACRouter(r)
	r.api = api
}

func (r CollectionsRouter) ConformanceClasses() []string {
	return []string{STACAPICollectionsConformanceURI}
}

func (r CollectionsRouter) LandingPageLinks(request *http.Request) []*stac.Link {
	data_link := DataLink(request, ListCollectionsPath)
	return []*stac.Link{&data_link}
}

func (r CollectionsRouter) AttachHandlers(gr gin.IRouter) {
	gr.GET(ListCollectionsPath, r.HandleListCollections)
	gr.GET(GetCollectionPath, r.HandleGetCollection)
}

func (r CollectionsRouter) HandleListCollections(c *gin.Context) {
	root_link := RootLink(c.Request)
	self_link := SelfLink(c.Request, ApplicationJSONType)

	collections := r.backend.ListCollections()
	collections_list := map[string]any{
		"collections": collections,
		"links":       []*stac.Link{&root_link, &self_link},
	}
	c.JSON(http.StatusOK, collections_list)
}

func (r CollectionsRouter) HandleGetCollection(c *gin.Context) {
	root_link := RootLink(c.Request)
	self_link := SelfLink(c.Request, ApplicationJSONType)
	parent_link := ParentLink(c.Request)

	collection_id := c.Param("collection_id")
	raw_collection, err := r.backend.GetCollection(collection_id)

	if err != nil {
		switch err.(type) {
		case CollectionNotFoundError:
			c.JSON(http.StatusNotFound, err.Error())
		default:
			c.JSON(http.StatusInternalServerError, "Internal server error.")
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

type ICollectionsBackend interface {
	ListCollections() []map[string]any
	GetCollection(id string) (map[string]any, error)
}
