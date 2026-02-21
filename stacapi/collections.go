package stacapi

import (
	"net/http"

	"github.com/duckontheweb/go-stac-api/internal"
	"github.com/gin-gonic/gin"
	"github.com/planetlabs/go-stac"
)

type STACCollectionsRouter struct {
	api *STACApi
}

func (r *STACCollectionsRouter) AttachTo(api *STACApi) {
	api.AddSTACRouter(r)
	r.api = api
}

func (r STACCollectionsRouter) ConformanceClasses() []string {
	return []string{internal.STACAPICollectionsConformanceURI}
}

func (r STACCollectionsRouter) LandingPageLinks(request *http.Request) []*stac.Link {
	data_link := internal.DataLink(request, internal.ListCollectionsPath)
	return []*stac.Link{&data_link}
}

func (r STACCollectionsRouter) AttachHandlers(gr gin.IRouter) {
	gr.GET(internal.ListCollectionsPath, r.HandleListCollections)
	gr.GET(internal.GetCollectionPath, r.HandleGetCollection)
}

func (r STACCollectionsRouter) HandleListCollections(c *gin.Context) {
	root_link := internal.RootLink(c.Request)
	self_link := internal.SelfLink(c.Request, internal.ApplicationJSONType)
	collections := stac.CollectionsList{
		Collections: make([]*stac.Collection, 0),
		Links:       []*stac.Link{&root_link, &self_link},
	}
	c.JSON(http.StatusOK, collections)
}

func (r STACCollectionsRouter) HandleGetCollection(c *gin.Context) {
	root_link := internal.RootLink(c.Request)
	self_link := internal.SelfLink(c.Request, internal.ApplicationJSONType)
	parent_link := internal.ParentLink(c.Request)

	collection_id := c.Param("collection_id")

	collection := stac.Collection{
		Id:      collection_id,
		Version: internal.STACVersion,
		Links:   []*stac.Link{&root_link, &self_link, &parent_link},
	}
	c.JSON(http.StatusOK, collection)
}
