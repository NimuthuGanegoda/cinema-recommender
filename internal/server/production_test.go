package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestBoundaries(t *testing.T) {
	app := setupTestServer()
	cases := []struct {
		name, path, body, content string
		status                    int
	}{
		{"negative budget", "/api/v1/recommend", `{"budget_lkr":-1}`, "application/json", 400},
		{"large party", "/api/v1/recommend", `{"party_size":51}`, "application/json", 400},
		{"unknown field", "/api/v1/recommend", `{"budegt":100}`, "application/json", 400},
		{"multiple objects", "/api/v1/recommend", `{} {}`, "application/json", 400},
		{"null", "/api/v1/recommend", `null`, "application/json", 400},
		{"unsupported dietary promise", "/api/v1/recommend", `{"dietary_preferences":["nut-free"]}`, "application/json", 400},
		{"unsupported tier", "/api/v1/recommend", `{"privilege_tier":"Admin"}`, "application/json", 400},
		{"wrong media type", "/api/v1/recommend", `{}`, "text/plain", 415},
		{"too large", "/api/v1/recommend", `{"city":"` + strings.Repeat("x", 70<<10) + `"}`, "application/json", 413},
		{"quantity bomb", "/api/v1/cart/evaluate", `{"items":[{"item_id":"KND-01","quantity":2147483647}]}`, "application/json", 400},
		{"negative quantity", "/api/v1/cart/evaluate", `{"items":[{"item_id":"KND-01","quantity":-1}]}`, "application/json", 400},
		{"unknown item", "/api/v1/cart/evaluate", `{"items":[{"item_id":"FAKE","quantity":1}]}`, "application/json", 400},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", tc.path, strings.NewReader(tc.body))
			req.Header.Set("Content-Type", tc.content)
			res := httptest.NewRecorder()
			app.GetHandler().ServeHTTP(res, req)
			if res.Code != tc.status {
				t.Fatalf("got %d, want %d: %s", res.Code, tc.status, res.Body.String())
			}
		})
	}
}
func TestOriginAndSecurityHeaders(t *testing.T) {
	app := setupTestServer()
	for _, origin := range []string{"https://untrusted.example", "http://example.com", ""} {
		req := httptest.NewRequest(http.MethodGet, "http://example.com/api/v1/health", nil)
		req.Header.Set("Origin", origin)
		res := httptest.NewRecorder()
		app.GetHandler().ServeHTTP(res, req)
		want := 200
		if origin == "https://untrusted.example" {
			want = 403
		}
		if res.Code != want {
			t.Fatalf("origin %q: %d", origin, res.Code)
		}
		if res.Header().Get("X-Content-Type-Options") != "nosniff" || res.Header().Get("Content-Security-Policy") == "" {
			t.Fatal("missing security headers")
		}
		if res.Header().Get("Access-Control-Allow-Origin") == "*" {
			t.Fatal("wildcard origin")
		}
	}
}

func TestStaticDirectoriesAndMissingRoutes(t *testing.T) {
	app := setupTestServer()
	for _, path := range []string{"/assets/", "/assets/js/", "/data/", "/embed.go", "/api/v1/missing", "/api/v1/cinemas/UNKNOWN/movies"} {
		res := httptest.NewRecorder()
		app.GetHandler().ServeHTTP(res, httptest.NewRequest("GET", path, nil))
		if res.Code != 404 {
			t.Errorf("%s: got %d, want 404", path, res.Code)
		}
	}
}
