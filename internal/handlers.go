package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/planetlabs/go-stac"
)

func HandleLandingPage(c *gin.Context) {
	root_link := RootLink(c.Request)
	self_link := SelfLink(c.Request, RootPath, ApplicationJSONType)
	service_desc_link := ServiceDescLink(c.Request)
	service_doc_link := ServiceDocLink(c.Request)

	landing_page := stac.Catalog{
		Version:     STACVersion,
		Id:          "go-stac-api",
		ConformsTo:  []string{STACAPICoreConformanceURI},
		Description: "Implementation of the STAC API Spec in Go",
		Links:       []*stac.Link{&root_link, &self_link, &service_desc_link, &service_doc_link},
	}
	c.JSON(http.StatusOK, landing_page)
}
