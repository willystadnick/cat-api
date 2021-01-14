package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup() Api {
    var api Api
    api.Database = setupDatabase()
    api.Router = setupRouter(api)

	return api
}

func TestPingRoute(t *testing.T) {
	api := setup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	api.Router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestBreedsRouteWithoutName(t *testing.T) {
	api := setup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/breeds", nil)
	api.Router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "invalid breed name", w.Body.String())
}

func TestBreedsRouteWithoutRecords(t *testing.T) {
	api := setup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/breeds?name=foo", nil)
	api.Router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "no records", w.Body.String())
}

func TestBreedsRouteWithNameSib(t *testing.T) {
	api := setup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/breeds?name=sib", nil)
	api.Router.ServeHTTP(w, req)

    var cats []Cat
	json.Unmarshal(w.Body.Bytes(), &cats)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, 1, len(cats))
}

func TestLoginRouteSuccess(t *testing.T) {
	api := setup()

	body, _ := json.Marshal(map[string]string{
		"username": "admin",
		"password": "@#$RF@!718",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	api.Router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "login ok", w.Body.String())

}

func TestLoginRouteInvalidCredentials(t *testing.T) {
	api := setup()

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
