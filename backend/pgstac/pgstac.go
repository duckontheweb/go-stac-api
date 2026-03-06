package pgstac

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/duckontheweb/go-stac-api/stacapi"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func init() {
	stacapi.RegisterBackend(BackendType, Backend{})
}

const BackendType = "pgstac"

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
	log.Printf("Connection string: %s", connection_string)
	db, err := sqlx.Connect("postgres", connection_string)
	if err != nil {
		log.Fatalf("Unable to connect to database: %s", err)
	}

	return BackendClient{db: db}, nil
}

func (b BackendClient) ListCollections() []map[string]any {
	collections := Collections{}
	err := b.db.Get(&collections, `SELECT all_collections AS "result" FROM all_collections();`)
	if err != nil {
		log.Fatal(err)
	}

	return collections
}

func (b BackendClient) GetCollection(id string) (map[string]any, error) {
	collection := Collection{}
	err := b.db.Get(&collection, `SELECT get_collection AS "result" FROM get_collection($1);`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return map[string]any{}, stacapi.CollectionNotFoundError{Id: id}
		}
		return map[string]any{}, stacapi.CollectionNotFoundError{Id: id}
	}
	return collection, nil
}

type Collection map[string]any

func (c *Collection) Scan(value interface{}) error {
	if value == nil {
		return sql.ErrNoRows
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("unable to convert value to bytes array")
	}
	err := json.Unmarshal(b, &c)
	if err != nil {
		return err
	}
	return nil
}

type Collections []map[string]any

func (c *Collections) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("unable to convert value to bytes array")
	}
	err := json.Unmarshal(b, &c)
	if err != nil {
		return err
	}
	return nil
}
