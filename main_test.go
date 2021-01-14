package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	api := setupApi()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	api.Router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestBreedsRouteWithoutJWT(t *testing.T) {
	api := setupApi()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/breeds", nil)
	api.Router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestBreedsRouteWithInvalidJWT(t *testing.T) {
	api := setupApi()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/breeds", nil)

	jwt := "foo"
	req.Header.Add("Authorization", "Bearer " + jwt)

	api.Router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestBreedsRouteWithoutName(t *testing.T) {
	api := setupApi()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/breeds", nil)

	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.t-IDcSemACt8x4iTMCda8Yhe3iZaWbvV5XKSTbuAn0M"
	req.Header.Add("Authorization", "Bearer " + jwt)

	api.Router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "invalid breed name", w.Body.String())
}

func TestBreedsRouteWithoutRecords(t *testing.T) {
	api := setupApi()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/breeds?name=foo", nil)

	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.t-IDcSemACt8x4iTMCda8Yhe3iZaWbvV5XKSTbuAn0M"
	req.Header.Add("Authorization", "Bearer " + jwt)

	api.Router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "no records", w.Body.String())
}

func TestBreedsRouteWithNameSib(t *testing.T) {
	api := setupApi()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/breeds?name=sib", nil)

	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.t-IDcSemACt8x4iTMCda8Yhe3iZaWbvV5XKSTbuAn0M"
	req.Header.Add("Authorization", "Bearer " + jwt)

	api.Router.ServeHTTP(w, req)

    var cats []Cat
	json.Unmarshal(w.Body.Bytes(), &cats)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, 1, len(cats))
}

func TestLoginRouteSuccess(t *testing.T) {
	api := setupApi()

	body, _ := json.Marshal(map[string]string{
		"username": "admin",
		"password": os.Getenv("ADMIN_PASS"),
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	api.Router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, true, strings.HasPrefix(w.Body.String(), "eyJ"))

}

func TestLoginRouteInvalidCredentials(t *testing.T) {
	api := setupApi()

	body, _ := json.Marshal(map[string]string{
		"username": "admin",
		"password": "admin",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	api.Router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "invalid credentials", w.Body.String())

}
