package stacapi

const APIVersion = "0.1.0"
const STACVersion = "1.1.0"

const (
	STACAPICoreConformanceURI        = "https://api.stacspec.org/v1.0.0/core"
	STACAPICollectionsConformanceURI = "https://api.stacspec.org/v1.0.0/collections"
)
const (
	RootPath            = "/"
	ServiceDescPath     = "/api"
	ServiceDocPath      = "/api.html"
	ListCollectionsPath = "/collections"
	GetCollectionPath   = "/collections/:collection_id"
)
const (
	SelfRel        = "self"
	RootRel        = "root"
	ParentRel      = "parent"
	ServiceDescRel = "service-desc"
	ServiceDocRel  = "service-doc"
	DataRel        = "data"
)

const (
	ApplicationJSONType = "application/json"
	OpenAPIYAMLType     = "application/x-yaml"
	HTMLType            = "text/html"
)
