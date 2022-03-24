package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingUnit(t *testing.T) {
	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	ping(rr, r)

	rs := rr.Result()

	if rs.StatusCode != http.StatusOK {
		t.Errorf("expected %d; got %d", http.StatusOK, rs.StatusCode)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("expected %s; got %s", "OK", string(body))
	}
}

func TestPingEndToEnd(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")
	if code != http.StatusOK {
		t.Errorf("expected %d; got %d", http.StatusOK, code)

	}
	if string(body) != "OK" {
		t.Errorf("expected %s; got %s", "OK", string(body))
	}
}
