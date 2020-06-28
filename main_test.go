package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestingIndex(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatalf("TestIndex: could not create HTTP GET request: %v", err)
	}
	rec := httptest.NewRecorder()
	index().ServeHTTP(rec, req)
	res := rec.Result()
	defer func() {
		err := res.Body.Close()
		if err != nil {
			t.Fatalf("TestIndex: could not close response body: %v", err)
		}
	}()

	body, err := ioutil.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		t.Errorf("TestIndex: could not read response body: %v", err)
	}
	if len(string(body)) == 0 {
		t.Errorf("TestIndex: unexpected empty response body")
	}
}
