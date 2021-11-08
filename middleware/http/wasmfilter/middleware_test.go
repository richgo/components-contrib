package wasmfilter

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"

	"github.com/dapr/kit/logger"
)

func TestNewWASMHandlerFuncRequests(t *testing.T) {
	testLogger := logger.NewLogger("test")
	tests := []struct {
		name                string
		inputRequestFactory func() *http.Request
		evaluateFactory     func(t *testing.T) func(ctx *fasthttp.RequestCtx)
	}{
		{
			"Get method is handled",
			func() *http.Request {
				return httptest.NewRequest("GET", "http://localhost:8080/test", nil)
			},
			func(t *testing.T) func(ctx *fasthttp.RequestCtx) {
				return func(ctx *fasthttp.RequestCtx) {
					result := string(ctx.Request.Body())

					fmt.Println(result) // 42!

					assert.Equal(t, "43", result)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.inputRequestFactory()
			handler := NewWASMHandlerFunc(testLogger, tt.evaluateFactory(t))

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)
		})
	}
}
