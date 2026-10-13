//go:build js && wasm

package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"syscall/js"

	"cinema-recommender/internal/server"
	"cinema-recommender/pkg/recommender"
	"cinema-recommender/pkg/scraper"
)

type browserResponse struct {
	headers http.Header
	status  int
	body    bytes.Buffer
}

func (w *browserResponse) Header() http.Header            { return w.headers }
func (w *browserResponse) WriteHeader(status int)         { w.status = status }
func (w *browserResponse) Write(data []byte) (int, error) { return w.body.Write(data) }

func main() {
	log.SetOutput(io.Discard)
	app := server.NewServer(server.Config{}, scraper.NewRegionalScraper(), recommender.NewEngine())
	app.StartBackgroundTasks()
	handler := app.GetHandler()
	invoke := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 3 {
			return `{"status":400,"body":{"error":"Invalid browser request"}}`
		}
		method, path, body := args[0].String(), args[1].String(), args[2].String()
		if !strings.HasPrefix(path, "/api/v1/") || (method != "GET" && method != "POST") {
			return `{"status":400,"body":{"error":"Unsupported browser request"}}`
		}
		request, err := http.NewRequest(method, path, strings.NewReader(body))
		if err != nil {
			return `{"status":400,"body":{"error":"Invalid request path"}}`
		}
		request.Header.Set("Content-Type", "application/json")
		response := &browserResponse{headers: make(http.Header), status: http.StatusOK}
		handler.ServeHTTP(response, request)
		data := json.RawMessage(response.body.Bytes())
		if !json.Valid(data) {
			data = json.RawMessage(`{"error":"This feature is unavailable"}`)
		}
		encoded, err := json.Marshal(struct {
			Status int             `json:"status"`
			Body   json.RawMessage `json:"body"`
		}{response.status, data})
		if err != nil {
			return `{"status":500,"body":{"error":"Unable to prepare a result"}}`
		}
		return string(encoded)
	})
	js.Global().Set("cinebiteInvoke", invoke)
	js.Global().Call("cinebiteReady")
	select {}
}
