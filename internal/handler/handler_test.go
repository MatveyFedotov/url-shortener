package handler

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"url-shortener/internal/service"
	"url-shortener/internal/storage/memory"
)

func TestCreateShorten(t *testing.T) {
	gin.SetMode(gin.TestMode)

	store := memory.New()
	svc := service.New(store)
	h := New(svc)

	router := gin.Default()
	router.POST("/shorten", h.Create)

	body := `{"url":"https://google.com"}`
	req := httptest.NewRequest("POST", "/shorten", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Error("Expected 200, got", w.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal("Failed to decode response:", err)
	}

	code, ok := resp["code"].(string)
	if !ok {
		t.Fatal("Code not found in response")
	}
	if len(code) != 10 {
		t.Error("Code should be 10 chars, got", len(code))
	}
}

func TestGetLink(t *testing.T) {
	gin.SetMode(gin.TestMode)

	store := memory.New()
	svc := service.New(store)
	h := New(svc)

	router := gin.Default()
	router.POST("/shorten", h.Create)
	router.GET("/:code", h.Get)

	// create
	createReq := httptest.NewRequest("POST", "/shorten", bytes.NewBufferString(`{"url":"https://google.com"}`))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)

	var createResp map[string]interface{}
	if err := json.NewDecoder(createW.Body).Decode(&createResp); err != nil {
		t.Fatal("Failed to decode create response:", err)
	}

	code, ok := createResp["code"].(string)
	if !ok {
		t.Fatal("Code not found in create response")
	}

	// get
	getReq := httptest.NewRequest("GET", "/"+code, nil)
	getW := httptest.NewRecorder()
	router.ServeHTTP(getW, getReq)

	if getW.Code != 200 {
		t.Error("Expected 200, got", getW.Code)
	}

	var getResp map[string]interface{}
	if err := json.NewDecoder(getW.Body).Decode(&getResp); err != nil {
		t.Fatal("Failed to decode get response:", err)
	}

	url, ok := getResp["url"].(string)
	if !ok {
		t.Fatal("URL not found in response")
	}
	if url != "https://google.com" {
		t.Error("Expected https://google.com, got", url)
	}
}

func TestGetNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	store := memory.New()
	svc := service.New(store)
	h := New(svc)

	router := gin.Default()
	router.GET("/:code", h.Get)

	req := httptest.NewRequest("GET", "/notexist", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Error("Expected 404, got", w.Code)
	}
}
