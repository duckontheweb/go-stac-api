package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/planetlabs/go-stac"
)

func HandleLandingPage(c *gin.Context) {
	root_link := RootLink()
	self_link := SelfLink(RootHREF, ApplicationJSONType)
	service_desc_link := stac.Link{
		Href:  ServiceDescHREF,
		Rel:   ServiceDescRel,
		Type:  OpenAPIYAMLType,
		Title: "OpenAPI YAML",
	}
	service_doc_link := stac.Link{
		Href:  ServiceDocHREF,
		Rel:   ServiceDocRel,
		Type:  HTMLType,
		Title: "OpenAPI Docs",
	}

	landing_page := stac.Catalog{
		Version:     "1.0.0",
		Id:          "go-stac-api",
		Description: "Golang implementation of the STAC API Spec",
		Links:       []*stac.Link{&root_link, &self_link, &service_desc_link, &service_doc_link},
	}
	c.JSON(http.StatusOK, landing_page)
}
