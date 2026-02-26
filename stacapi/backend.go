package stacapi

import (
	"log"
)

var backends = make(map[string]Backend)

func RegisterBackend(name string, backend Backend) {
	if backend == nil {
		log.Fatal("sql: Register driver is nil")
	}
	if _, exists := backends[name]; exists {
		log.Fatal("sql: Register called twice for driver " + name)
	}
	backends[name] = backend
}

type BackendConfig struct {
	Type             string `yaml:"type"`
	ConnectionString string `yaml:"connection_string"`
}

type Backend interface {
	GetClient(config BackendConfig) (BackendClient, error)
}

type BackendClient interface {
	Close() error
}
