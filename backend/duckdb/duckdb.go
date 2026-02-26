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

func init() {
	stacapi.RegisterBackend(BackendType, Backend{})
}

const BackendType = "duckdb"

type Backend struct{}

type BackendClient struct {
	db *sqlx.DB
}

func (b BackendClient) Close() error {
	return b.db.Close()
}

func (b Backend) GetClient(config stacapi.BackendConfig) (stacapi.BackendClient, error) {
	if config.Type != BackendType {
		return BackendClient{}, fmt.Errorf("Backend type must be '%s', found '%s'", BackendType, config.Type)
	}

	connection_string := config.ConnectionString
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
		log.Fatal("Unable to connect to database.")
	}

	return BackendClient{db: db}, nil
}

func (b BackendClient) ListCollections() []map[string]any {
	collections := []map[string]any{}
	err := b.db.Select(&collections, `SELECT content FROM collections;`)
	if err != nil {
		log.Fatal(err)
	}

	return collections
}

func (b BackendClient) GetCollection(id string) (map[string]any, error) {
	collection := map[string]any{}
	err := b.db.Get(&collection, `SELECT content FROM collections WHERE id = $1;`, id)
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
