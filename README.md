go-stac-api
===========

A [STAC API](https://github.com/radiantearth/stac-api-spec) server written in Go.

## Installation

Install with `go`:

```console
$ go install github.com/duckontheweb/go-stac-api/cmd/stac-server
```

Add the `go` installation directory to your `PATH` if you have not done so already. On a Mac, this is usually `~/go/bin`.

## Usage

Run installed binary:

```console
$ stac-server
```

This will serve the application on http://localhost:8080.

## STAC API Spec Compliance

| **Capability** | **Supported** |
| -- | -- |
| [STAC API - Core](https://github.com/radiantearth/stac-api-spec/blob/release/v1.0.0/core) | ✅ |
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

## Development

The project is configured to use [`air`](https://github.com/air-verse/air) for live reloading:

```console
$ air ./cmd/stac-server
```
