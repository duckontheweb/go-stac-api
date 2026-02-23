package duckdb

import (
	"database/sql"
	"errors"
	"log"
	"os"

	_ "github.com/duckdb/duckdb-go/v2"
	"github.com/duckontheweb/go-stac-api/stacapi"
	"github.com/go-viper/mapstructure/v2"
	"github.com/jmoiron/sqlx"
	"gopkg.in/yaml.v3"
)

type BackendConfig struct {
	ConnectionString string `yaml:"connection_string"`
}

type Backend struct {
	BackendConfig
}

func NewBackend(config_path string) Backend {
	var config BackendConfig

	contents, err := os.ReadFile(config_path)
	if err != nil {
		log.Fatalf("Could not read config file at path %s", config_path)
	}

	err = yaml.Unmarshal(contents, &config)
	if err != nil {
		log.Fatalf("Could not parse config file at path %s as config YAML", config_path)
	}

	backend := Backend{
		BackendConfig: config,
	}

	return backend
}

func (b Backend) ListCollections() []map[string]any {
	db, err := sqlx.Connect("duckdb", b.ConnectionString)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	collections := []*Collection{}
	err = db.Select(&collections, `SELECT content FROM collections;`)
	if err != nil {
		log.Fatal(err)
	}

	collections_list := make([]map[string]any, len(collections))

	for i, collection := range collections {
		m := map[string]any{}
		mapstructure.Decode(collection, &m)
		collections_list[i] = m
	}

	return collections_list
}

func (b Backend) GetCollection(id string) (map[string]any, error) {
	db, err := sqlx.Connect("duckdb", b.ConnectionString)
	if err != nil {
		return map[string]any{}, err
	}
	defer db.Close()

	collection := map[string]any{}
	err = db.Get(&collection, `SELECT content FROM collections WHERE id = $1;`, id)
	if err != nil {
		if err != sql.ErrNoRows {
			return map[string]any{}, err
		}
		return map[string]any{}, stacapi.CollectionNotFoundError{Id: id}
	}

	return collection, nil
}

type Collection struct {
	Id               string         `mapstructure:"id"`
	Links            []*Link        `mapstructure:"links"`
	AdditionalFields map[string]any `mapstructure:",remain"`
}

type Link struct {
	HREF  string
	Rel   string
	Type  string `json:"type,omitempty"`
	Title string `json:"title,omitempty"`
}

func (c *Collection) Scan(src interface{}) error {
	switch src.(type) {
	default:
		return errors.New("Failed to parse value.")
	case map[string]interface{}:
		mapstructure.Decode(src, c)
		return nil
	}
}
