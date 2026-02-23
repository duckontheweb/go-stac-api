package duckdb

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/duckdb/duckdb-go/v2"
	"github.com/duckontheweb/go-stac-api/stacapi"
	"github.com/go-viper/mapstructure/v2"
	"github.com/jmoiron/sqlx"
)

const BackendType = "duckdb"

type Backend struct {
	stacapi.BackendConfig
}

func NewBackend(config stacapi.BackendConfig) (Backend, error) {
	if config.Type != BackendType {
		return Backend{}, errors.New(fmt.Sprintf("Backend type must be '%s', found '%s'", BackendType, config.Type))
	}

	backend := Backend{
		BackendConfig: config,
	}

	return backend, nil
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
