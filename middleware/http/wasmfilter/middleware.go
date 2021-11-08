package wasmfilter

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/dapr/kit/logger"
	"github.com/valyala/fasthttp"
	wasmer "github.com/wasmerio/wasmer-go/wasmer"
)

// NewNetHTTPHandlerFunc wraps a fasthttp.RequestHandler in a http.HandlerFunc.
func NewWASMHandlerFunc(logger logger.Logger, h fasthttp.RequestHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := fasthttp.RequestCtx{}

		wasmBytes, _ := ioutil.ReadFile("simple.wasm")

		engine := wasmer.NewEngine()
		store := wasmer.NewStore(engine)

		// Compiles the module
		module, _ := wasmer.NewModule(store, wasmBytes)

		// Instantiates the module
		importObject := wasmer.NewImportObject()
		instance, _ := wasmer.NewInstance(module, importObject)

		// Gets the `sum` exported function from the WebAssembly instance.
		sum, _ := instance.Exports.GetFunction("sum")

		// Calls that exported function with Go standard values. The WebAssembly
		// types are inferred and values are casted automatically.
		result, _ := sum(5, 37)

		fmt.Println(result) // 42!
		resultInt := result.(int32)
		w.Write([]byte(strconv.Itoa(int(resultInt))))
		c.Response.BodyWriteTo(w)
	})
}
