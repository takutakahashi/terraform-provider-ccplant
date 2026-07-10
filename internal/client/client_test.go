package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetJSONAddsBearerToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Header.Get("Authorization"), "Bearer test-key"; got != want {
			t.Fatalf("Authorization header = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	c, err := New(server.URL, "test-key")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	body, err := c.GetJSON(context.Background(), "/user/info")
	if err != nil {
		t.Fatalf("GetJSON() error = %v", err)
	}
	if got, want := string(body), `{"ok":true}`; got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}
}

func TestGetJSONRejectsInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`not json`))
	}))
	defer server.Close()

	c, err := New(server.URL, "")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if _, err := c.GetJSON(context.Background(), "/bad"); err == nil {
		t.Fatal("GetJSON() error = nil, want error")
	}
}

func TestIsNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer server.Close()

	c, err := New(server.URL, "")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	_, err = c.GetJSON(context.Background(), "/missing")
	if !IsNotFound(err) {
		t.Fatalf("IsNotFound(%v) = false, want true", err)
	}
}
