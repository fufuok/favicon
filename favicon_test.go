package favicon

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type header struct {
	Key   string
	Value string
}

func performRequest(r http.Handler, method, path string, headers ...header) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

type mockFS struct{}

func (m mockFS) Open(name string) (http.File, error) {
	if name == "/" {
		name = "."
	} else {
		name = strings.TrimPrefix(name, "/")
	}
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func TestFavicon(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(New())

	r.Handle("GET", "/")
	w := performRequest(r, "GET", "/")
	if w.Code != http.StatusOK {
		t.Fatalf("expected: %d, actual: %d", http.StatusOK, w.Code)
	}

	r.GET("/favicon.ico")
	w = performRequest(r, "GET", "/favicon.ico")
	if w.Code != http.StatusNoContent {
		t.Fatalf("status code expected: %d, actual: %d", http.StatusNoContent, w.Code)
	}

	w = performRequest(r, "OPTIONS", "/favicon.ico")
	if w.Code != http.StatusOK {
		t.Fatalf("status code expected: %d, actual: %d", http.StatusOK, w.Code)
	}

	w = performRequest(r, "POST", "/favicon.ico")
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status code expected: %d, actual: %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestFaviconFile(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(New(Config{
		File: "example/favicon.ico",
	}))

	r.GET("/favicon.ico")
	resp := performRequest(r, "GET", "/favicon.ico").Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status code expected: %d, actual: %d", http.StatusOK, resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != "image/x-icon" {
		t.Fatalf("Content-Type expected: image/x-icon, actual: %s", resp.Header.Get("Content-Type"))
	}
	if resp.Header.Get("Cache-Control") != "public, max-age=31536000" {
		t.Fatalf("Content-Type expected: public, max-age=31536000, actual: %s", resp.Header.Get("Cache-Control"))
	}
}

func TestFaviconFileSystem(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(New(Config{
		File:       "example/favicon.ico",
		FileSystem: mockFS{},
	}))

	r.GET("/favicon.ico")
	resp := performRequest(r, "GET", "/favicon.ico").Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status code expected: %d, actual: %d", http.StatusOK, resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != "image/x-icon" {
		t.Fatalf("Content-Type expected: image/x-icon, actual: %s", resp.Header.Get("Content-Type"))
	}
	if resp.Header.Get("Cache-Control") != "public, max-age=31536000" {
		t.Fatalf("Content-Type expected: public, max-age=31536000, actual: %s", resp.Header.Get("Cache-Control"))
	}
}

func TestFaviconFileData(t *testing.T) {
	testData := "mock icon data"
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(New(Config{
		FileData: []byte(testData),
	}))

	r.GET("/favicon.ico")
	resp := performRequest(r, "GET", "/favicon.ico").Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status code expected: %d, actual: %d", http.StatusOK, resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != "image/x-icon" {
		t.Fatalf("Content-Type expected: image/x-icon, actual: %s", resp.Header.Get("Content-Type"))
	}
	if resp.Header.Get("Cache-Control") != "public, max-age=31536000" {
		t.Fatalf("Content-Type expected: public, max-age=31536000, actual: %s", resp.Header.Get("Cache-Control"))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read resp.Body err: %s", err)
	}
	if string(data) != testData {
		t.Fatalf("favicon data expected: %s, actual: %s", testData, string(data))
	}
}

func TestFaviconCacheControl(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(New(Config{
		File:         "example/favicon.ico",
		CacheControl: "no-cache",
	}))

	r.GET("/favicon.ico")
	resp := performRequest(r, "GET", "/favicon.ico").Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status code expected: %d, actual: %d", http.StatusOK, resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != "image/x-icon" {
		t.Fatalf("Content-Type expected: image/x-icon, actual: %s", resp.Header.Get("Content-Type"))
	}
	if resp.Header.Get("Cache-Control") != "no-cache" {
		t.Fatalf("Content-Type expected: public, max-age=31536000, actual: %s", resp.Header.Get("Cache-Control"))
	}
}
