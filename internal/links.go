package internal

import (
	"net/http"
	"net/url"

	"github.com/planetlabs/go-stac"
)

func RootLink(request *http.Request) stac.Link {
	href := constructHREF(request, RootPath)
	return stac.Link{
		Href:  href,
		Rel:   RootRel,
		Type:  ApplicationJSONType,
		Title: "Root",
	}
}

func SelfLink(request *http.Request, type_ string) stac.Link {
	href := constructHREF(request, request.URL.Path)
	return stac.Link{
		Href:  href,
		Rel:   SelfRel,
		Type:  type_,
		Title: "This Page",
	}
}

func ParentLink(request *http.Request) stac.Link {
	href := constructHREF(request, RootPath)
	return stac.Link{
		Href:  href,
		Rel:   ParentRel,
		Type:  ApplicationJSONType,
		Title: "Root",
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

func DataLink(request *http.Request, path string) stac.Link {
	href := constructHREF(request, path)
	return stac.Link{
		Href:  href,
		Rel:   DataRel,
		Type:  ApplicationJSONType,
		Title: "Collections",
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
