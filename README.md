go-stac-api
===========

A [STAC API](https://github.com/radiantearth/stac-api-spec) server written in Go using the [Gin
Framework](https://gin-gonic.com/) inspired by projects like
[`stac-fastapi`](https://stac-utils.github.io/stac-fastapi/) and [rustac](https://github.com/stac-utils/rustac).

The aspiration is to provide a fully compliant STAC API implementation supporting various backends (e.g. PgSTAC, STAC
Geoparquet, DuckDB) that can be used either as a configuration-driven command line tool or a Go module integrated into
other Gin-based services. This project is very young and we're still a ways from that goal...

## Run the Server

Clone this repo and use [`air`](https://github.com/air-verse/air) to run the server with hot reloading:

```console
$ git clone git@github.com:duckontheweb/go-stac-api.git
$ cd go-stac-api
$ air
```

By default, this will use the [example DuckDB config](./example/duckdb/config.yaml) which uses a [DuckDB
backend](./backend/duckdb/) connected to the example data at `./example/duckdb/data.duckdb`. If you want to

This will serve the application on http://localhost:8080 using the configuration at
[`./backend/duckdb/example/config.yaml`](./backend/duckdb/example/config.yaml). This uses the DuckDB backend and connects to the file at
`./backend/duckdb/example/data.duckdb`. See [Backends](#backends) for details on how to connect to other backends.

## Use as a Library

Running a server involves the following steps:

1. Import your desired backend.

    In this example, we import the [DuckDB backend](#duckdb)

    ```go
    import _ "github.com/duckontheweb/go-stac-api/backend/duckdb"
    ```

2. Create a new `stacapi.ApiConfig` instance

    You can either instantiate the struct directly or use `stacapi.ReadConfig` to read it from a YAML file.

    **Instantiate Directly**
    ```go
    import "github.com/duckontheweb/go-stac-api/stacapi"

    var config stacapi.ApiConfig
    // Instantiate config directly
    config = stacapi.ApiConfig{
        Backend: stacapi.BackendConfig{
            Type: "duckdb",
            ConnectionString: "./example/duckdb/data.duckdb"
        },
    }
    ```

    **Read from YAML**

    ```yaml
    # ./example/duckdb/config.yaml
    backend:
      type: duckdb
      connection_string: ./example/duckdb/data.duckdb
    ```

    ```go
    config = stacapi.ReadConfig("./example/duckdb/config.yaml")
    ```

3. Create a `stacapi.StacApi` instance from the config

    ```go
    import "log"

    stac_api := stacapi.NewApi(config)
    // Ensure the backend is closed before exiting
	defer func() {
		err := stac_api.Shutdown()
		if err != nil {
			log.Fatalf("Unable to close client for %s backend", err)
		}
	}()
    ```

4. Attach the `stacapi.StacApi` instance to a `gin` router/engine and run the `gin` app

    ```go
    import "github.com/gin-gonic/gin"

    router := gin.Default()

    stac_api.AddToRouter(router)

	err := router.Run()
	if err != nil {
		log.Fatal(err)
	}
    ```

See the [`stac-server` executable](./cmd/stac-server/main.go) for a full example.
Note that we use Go build tags in that package to allow us to build separate executables for each backend,
so the DuckDB backend import actually happens in a [separate file](./cmd/stac-server/duckdb.go).

## STAC API Spec Compliance

See [Backends](#backends) for details on the specific backend implementations.

| **Capability** | **DuckDB Backend** | **PgSTAC Backend** |
| -- | -- | -- |
| [STAC API - Core](https://github.com/radiantearth/stac-api-spec/blob/release/v1.0.0/core) | ✅ | ✅ |
| [STAC API - Collections](https://github.com/radiantearth/stac-api-spec/blob/release/v1.0.0/ogcapi-features/README.md#stac-api---collections) | ✅ | ✅ |
| [STAC API - Features](https://github.com/radiantearth/stac-api-spec/blob/release/v1.0.0/ogcapi-features) | ✖️ | ✖️ |
| [STAC API - Item Search](https://github.com/radiantearth/stac-api-spec/blob/release/v1.0.0/item-search) | ✖️ | ✖️ |
| [Aggregation extension](https://github.com/stac-api-extensions/aggregation) | ✖️ | ✖️ |
| [Browseable extension](https://github.com/stac-api-extensions/browseable) | ✖️ | ✖️ |
| [Children extension](https://github.com/stac-api-extensions/children) | ✖️ | ✖️ |
| [Collection search extension](https://github.com/stac-api-extensions/collection-search) | ✖️ | ✖️ |
| [Collection transaction extension](https://github.com/stac-api-extensions/collection-transaction) | ✖️ | ✖️ |
| [Fields extension](https://github.com/stac-api-extensions/fields) | ✖️ | ✖️ |
| [Filter extension](https://github.com/stac-api-extensions/filter) | ✖️ | ✖️ |
| [Free-text search extension](https://github.com/stac-api-extensions/freetext-search) | ✖️ | ✖️ |
| [Language (I18N) extension](https://github.com/stac-api-extensions/language) | ✖️ | ✖️ |
| [Query extension](https://github.com/stac-api-extensions/query) | ✖️ | ✖️ |
| [Sort extension](https://github.com/stac-api-extensions/sort) | ✖️ | ✖️ |
| [Transaction extension](https://github.com/stac-api-extensions/transaction) | ✖️ | ✖️ |

## Backends

This library takes a similar approach to [`stac-fastapi`](https://stac-utils.github.io/stac-fastapi/) by defining
API-level data structures in the `stacapi` package along with abstract backend interfaces that the API uses to read
and write STAC objects. Concrete backends must implement these interfaces to be used by the API.

### `duckdb`

**Go Package:** `github.com/duckontheweb/go-stac-api/backend/duckdb`
**Example:** [./example/duckdb](./example/duckdb/)

Uses the [`duckdb` driver](https://duckdb.org/docs/stable/clients/go) for [`sqlx`](https://jmoiron.github.io/sqlx/)
to read STAC objects from a [persistent DuckDB
database](https://duckdb.org/docs/stable/connect/overview#persistent-database). This persistent backend must have a
`collections` tables with the following columns:

- `id TEXT`: Unique Collection ID
- `content JSON`: Valid JSON-encoded STAC Collection

To use this backend, include the following `backend` config in your YAML config file:

```yaml
backend:
    type: duckdb
    connection_string: /path/to/some/database.duckdb
```

and build the executable with `-tags duckdb`.

### `pgstac`

**Go Package:** `github.com/duckontheweb/go-stac-api/backend/pgstac`
**Example:** [./example/pgstac](./example/pgstac/)

Uses the [`lib/pq`](https://github.com/lib/pq) driver for [`sqlx`](https://jmoiron.github.io/sqlx/) to read STAC
objects from a PostgreSQL database with [PgSTAC](https://stac-utils.github.io/pgstac/pgstac/) installed. You must run
the [PgSTAC migrations](https://stac-utils.github.io/pgstac/pypgstac/#running-migrations) prior to starting the API.
See [`./example/pgstac/load-data.sh`](./example/pgstac/load-data.sh) for an example of running migrations and bulk
loading data.

To use this backend, include the following `backend` config in your YAML config file:

```yaml
backend:
    type: pgstac
    connection_string: username:password@hostname:5432/dbname
```

and build the executable with `-tags pgstac`.

## Development

### Run Locally
The project is configured to use [`air`](https://github.com/air-verse/air) for hot reloading within containers managed
by Docker Compose.

To run the server with the default DuckDB backend:

```console
$ docker compose -f ./docker/docker-compose.yaml up
```

To use the PgSTAC backend instead, copy the [`docker/.env.example`](./docker/.env.example) file to `docker/.env` and
uncomment the lines defining the `BACKEND` and `POSTGRES_PASSWORD` variables, then run the same command as above.

### Debug

To run the server using Go's [Delve Debugger](https://github.com/go-delve/delve), set `DEBUG=TRUE` in either the
`docker/.env` or in you command:

```console
$ DEBUG=TRUE docker compose -f ./docker/docker-compose.yaml up
```

This will start the debug server listening on port `2345` (you can change this port by setting the `DEBUG_PORT`
enviromment variable). The API will _not start until a debugger is attached_.

To attach a debugger using VSCode:

1. Add the following configuration to your `launch.json` file

    ```json
    {
        "version": "0.2.0",
        "configurations": [
            {
                "name": "Connect to server",
                "type": "go",
                "debugAdapter": "dlv-dap",
                "request": "attach",
                "mode": "remote",
                "substitutePath": [{
                    "from": "${workspaceFolder}",
                    "to": "/mnt"
                }],
                "port": 2345,
            }
        ]
    }
    ```

2. Start the server in debug mode

    ```console
    $ DEBUG=TRUE docker compose -f ./docker/docker-compose.yaml up
    ```

3. Connect the debugger using the "Connect to server" option in the Run and Debug panel

   ![Connect Debugger](https://github.com/duckontheweb/go-stac-api/raw/refs/heads/pgstac-backend/img/connect_debugger.mp4)


### Code Quality
We also define pre-commit hooks using Python's `pre-commit` library. To install the pre-commit hooks you will need
to have [`uv` installed](https://docs.astral.sh/uv/getting-started/installation/). Then run:

```console
$ uv sync --all-groups
$ uv run pre-commit install
```

To validate the STAC API using [stac-api-validator](https://github.com/stac-utils/stac-api-validator) (after you have run `uv sync`):

Start the server:

```console
$ air
# or
$ stac-server
```

Run the validator:

```console
$ uv run stac-api-validator \
    --root-url http://localhost:8080 \
    --conformance core --conformance collections \
    --collection naip
```

Pre-commit checks and STAC API validation are also run in CI.
