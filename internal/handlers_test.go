package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"
)

func TestLandingPageRoute(t *testing.T) {
	ja := jsonassert.New(t)
	w := httptest.NewRecorder()

	r := gin.Default()
	r.GET("/", HandleLandingPage)

	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	body := w.Body.String()

	ja.Assertf(body, `
	{
		"id": "<<PRESENCE>>",
		"type": "Catalog",
		"stac_version": "1.1.0",
		"description": "<<PRESENCE>>",
		"conformsTo": "<<PRESENCE>>",
		"links": "<<PRESENCE>>"
	}`)
}
