package gmidware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	gero "github.com/DiscoFighter47/gEro"
	gmidware "github.com/DiscoFighter47/gMidware"
)

func testHandler() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.String() {
		case "/apierr":
			panic(gero.NewAPIerror("API Error", http.StatusInternalServerError, fmt.Errorf("api error")))
		case "/err":
			panic(fmt.Errorf("error"))
		case "/panic":
			panic("panic")
		}
	}
	return http.HandlerFunc(fn)
}

func TestRecoverer(t *testing.T) {
	svr := gmidware.Recoverer(testHandler())

	testData := []struct {
		des  string
		url  string
		code int
		body string
	}{
		{
			des:  "api error",
			url:  "/apierr",
			body: `{"error":{"title":"API Error","detail":"api error"}}`,
		},
		{
			des:  "error",
			url:  "/err",
			body: `{"error":{"title":"Internal Server Error","detail":"error"}}`,
		},
		{
			des:  "panic",
			url:  "/panic",
			body: `{"error":{"title":"Internal Server Error","detail":"panic"}}`,
		},
	}

	for _, td := range testData {
		t.Run(td.des, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, td.url, nil)
			res := httptest.NewRecorder()
			svr.ServeHTTP(res, req)
			assert.Equal(t, http.StatusInternalServerError, res.Code)
			assert.JSONEq(t, td.body, res.Body.String())
		})
	}
}
