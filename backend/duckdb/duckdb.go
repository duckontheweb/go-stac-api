package duckdb

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	_ "github.com/duckdb/duckdb-go/v2"
	"github.com/duckontheweb/go-stac-api/stacapi"
	"github.com/jmoiron/sqlx"
)

const BackendType = "duckdb"

type Backend struct {
	Config stacapi.BackendConfig
}

func (b Backend) GetConnection() (*sqlx.DB, error) {
	connection_string := b.Config.ConnectionString
	if !filepath.IsAbs(connection_string) {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal("Could not get current working directory, please provide an absolute path to the config file.")
		}
		var abs_connection_string string
		abs_connection_string, err = filepath.Abs(path.Join(cwd, connection_string))
		if err != nil {
			log.Fatalf("Could not construct absolute path from relative path: %s", connection_string)
		}
		connection_string = abs_connection_string
	}
	db, err := sqlx.Connect("duckdb", connection_string)
	if err != nil {
		return new(sqlx.DB), err
	}

	return db, nil
}

func NewBackend(config stacapi.BackendConfig) (Backend, error) {
	if config.Type != BackendType {
		return Backend{}, fmt.Errorf("Backend type must be '%s', found '%s'", BackendType, config.Type)
	}

	backend := Backend{
		Config: config,
	}

	return backend, nil
}

func (b Backend) ListCollections() []map[string]any {
	db, err := sqlx.Connect("duckdb", b.Config.ConnectionString)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		err = db.Close()
	}()
	if err != nil {
		log.Print("Unable to close backend connection")
	}

	collections := []map[string]any{}
	err = db.Select(&collections, `SELECT content FROM collections;`)
	if err != nil {
		log.Fatal(err)
	}

	return collections
}

func (b Backend) GetCollection(id string) (map[string]any, error) {
	db, err := b.GetConnection()
	if err != nil {
		return map[string]any{}, err
	}
	defer func() {
		err = db.Close()
	}()
	if err != nil {
		log.Print("Unable to close backend connection")
	}

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
