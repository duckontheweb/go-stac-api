go-stac-api
===========

A [STAC API](https://github.com/radiantearth/stac-api-spec) server written in Go using the [Gin
Framework](https://gin-gonic.com/) inspired by [`stac-fastapi`](https://stac-utils.github.io/stac-fastapi/).

The aspiration is to provide a fully compliant STAC API implementation supporting various backends (e.g. PgSTAC, STAC
Geoparquet, DuckDB) that can be used either as a configuration-driven command line tool or a Go module integrated into
other Gin-based services.

## Install

Install with `go`:

```console
$ go install github.com/duckontheweb/go-stac-api/cmd/stac-server
```

Add the `go` installation directory to your `PATH` if you have not done so already. On a Mac, this is usually `~/go/bin`.

## Quick Start

Run installed binary:

```console
$ stac-server
```

This will serve the application on http://localhost:8080 using the configuration at
[`./data/example-config.yaml](./data/example-config.yaml). This uses the DuckDB backend and connects to the file at
`./data/example.duckdb`.

## STAC API Spec Compliance

See [Backends](#backends) for details on the specific backend implementations.

| **Capability** | **DuckDB Backend** |
| -- | -- |
| [STAC API - Core](https://github.com/radiantearth/stac-api-spec/blob/release/v1.0.0/core) | ✅ |
| [STAC API - Collections](https://github.com/radiantearth/stac-api-spec/blob/release/v1.0.0/ogcapi-features/README.md#stac-api---collections) | ✅ |
| [STAC API - Features](https://github.com/radiantearth/stac-api-spec/blob/release/v1.0.0/ogcapi-features) | ✖️ |
| [STAC API - Item Search](https://github.com/radiantearth/stac-api-spec/blob/release/v1.0.0/item-search) | ✖️ |
| [Aggregation extension](https://github.com/stac-api-extensions/aggregation) | ✖️ |
| [Browseable extension](https://github.com/stac-api-extensions/browseable) | ✖️ |
| [Children extension](https://github.com/stac-api-extensions/children) | ✖️ |
| [Collection search extension](https://github.com/stac-api-extensions/collection-search) | ✖️ |
| [Collection transaction extension](https://github.com/stac-api-extensions/collection-transaction) | ✖️ |
| [Fields extension](https://github.com/stac-api-extensions/fields) | ✖️ |
| [Filter extension](https://github.com/stac-api-extensions/filter) | ✖️ |
| [Free-text search extension](https://github.com/stac-api-extensions/freetext-search) | ✖️ |
| [Language (I18N) extension](https://github.com/stac-api-extensions/language) | ✖️ |
| [Query extension](https://github.com/stac-api-extensions/query) | ✖️ |
| [Sort extension](https://github.com/stac-api-extensions/sort) | ✖️ |
| [Transaction extension](https://github.com/stac-api-extensions/transaction) | ✖️ |

## Backends

This library takes a similar approach to [`stac-fastapi`](https://stac-utils.github.io/stac-fastapi/) by defining
backend-agnostic API data structures in the `stacapi` package along with abstract interfaces that the API structures
will interact with to read and write STAC objects. Concrete backends must implement these interfaces to be used by
the API.

Only the `duckdb` backend is currently supported.

### `duckdb`

This backend uses [DuckDB](https://duckdb.org/) to read STAC objects from a [persistent
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

## Development

The project is configured to use [`air`](https://github.com/air-verse/air) for live reloading:

```console
$ air
```

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
