package internal

import (
	"net/http"
	"net/url"

	"github.com/planetlabs/go-stac"
)

func RootLink(request *http.Request) stac.Link {
	root_href := constructHREF(request, RootPath)
	return stac.Link{
		Href:  root_href,
		Rel:   RootRel,
		Type:  ApplicationJSONType,
		Title: "Root",
	}
}

func SelfLink(request *http.Request, href, type_ string) stac.Link {
	href = constructHREF(request, href)
	return stac.Link{
		Href:  href,
		Rel:   SelfRel,
		Type:  type_,
		Title: "This Page",
	}
}

func ServiceDescLink(request *http.Request) stac.Link {
	href := constructHREF(request, ServiceDescPath)
	return stac.Link{
		Href:  href,
		Rel:   ServiceDescRel,
		Type:  OpenAPIYAMLType,
		Title: "OpenAPI YAML",
	}
}

func ServiceDocLink(request *http.Request) stac.Link {
	href := constructHREF(request, ServiceDocPath)
	return stac.Link{
		Href:  href,
		Rel:   ServiceDocRel,
		Type:  HTMLType,
		Title: "OpenAPI Docs",
	}
}

func constructHREF(request *http.Request, path string) string {
	scheme := "http"
	if request.TLS != nil {
		scheme = "https"
	}
	new_url := url.URL{
		Scheme: scheme,
		Host:   request.Host,
		Path:   path,
	}

	return new_url.String()
}
