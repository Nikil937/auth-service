package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegisterInvalidEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authHandler := NewAuthHandler(nil)

	router := gin.New()
	router.POST("/auth/register", authHandler.Register)

	w := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/register",
		strings.NewReader(`{
			"email": "bad-email",
			"password": "password123"
		}`),
	)

	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			w.Code,
		)
	}
}

func TestRegisterShortPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authHandler := NewAuthHandler(nil)

	router := gin.New()
	router.POST("/auth/register", authHandler.Register)

	w := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/register",
		strings.NewReader(`{
			"email": "test@example.com",
			"password": "123"
		}`),
	)

	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			w.Code,
		)
	}
}

func TestRegisterInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authHandler := NewAuthHandler(nil)

	router := gin.New()
	router.POST("/auth/register", authHandler.Register)

	w := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/register",
		strings.NewReader(`{
			"email": "test@example.com",
			"password":
		}`),
	)

	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			w.Code,
		)
	}
}
