package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingRoute(t *testing.T) {
	router := setupRouter()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/actor/1", nil)
	router.ServeHTTP(recorder, request)
	if recorder.Code != 200 {
		t.Errorf("%d != 200\n", recorder.Code)
	}
	if recorder.Body.String() != `{"Id":1,"FirstName":"Penelope","LastName":"Guiness"}` {
		t.Errorf("Bad response: %s", recorder.Body.String())
	}
}
