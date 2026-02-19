package internal

import "github.com/planetlabs/go-stac"

func RootLink() stac.Link {
	return stac.Link{
		Href:  RootHREF,
		Rel:   RootRel,
		Type:  ApplicationJSONType,
		Title: "Root",
	}
}

func SelfLink(href, type_ string) stac.Link {
	return stac.Link{
		Href:  href,
		Rel:   SelfRel,
		Type:  type_,
		Title: "This Page",
	}
}
