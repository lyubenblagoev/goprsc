package goprsc

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

func TestPaginationURL(t *testing.T) {
	got, err := paginationURL("domains?sort=name", []PaginationOptions{{Page: 2, PageSize: 50}})
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := url.Parse(got)
	if err != nil {
		t.Fatal(err)
	}
	query := parsed.Query()
	if query.Get("page") != "2" || query.Get("pageSize") != "50" || query.Get("sort") != "name" {
		t.Fatalf("unexpected query parameters: %v", query)
	}

	got, err = paginationURL("domains", nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != "domains" {
		t.Fatalf("expected omitted pagination defaults, got: %s", got)
	}
}

func TestPaginatedCollectionRequests(t *testing.T) {
	setup()
	defer shutdown()

	assertPagination := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page") != "2" {
			t.Errorf("expected page=2, got: %q", r.URL.Query().Get("page"))
		}
		if r.URL.Query().Get("pageSize") != "50" {
			t.Errorf("expected pageSize=50, got: %q", r.URL.Query().Get("pageSize"))
		}
		fmt.Fprint(w, `{"items":[],"page":2,"pageSize":50,"total":137}`)
	}
	mux.HandleFunc("/api/v1/domains", assertPagination)
	mux.HandleFunc("/api/v1/domains/example.com/accounts", assertPagination)
	mux.HandleFunc("/api/v1/domains/example.com/aliases", assertPagination)
	mux.HandleFunc("/api/v1/domains/example.com/aliases/contact", assertPagination)

	options := PaginationOptions{Page: 2, PageSize: 50}
	if _, err := client.Domains.List(options); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Accounts.List("example.com", options); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Aliases.List("example.com", options); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Aliases.Get("example.com", "contact", options); err != nil {
		t.Fatal(err)
	}
}

func TestPaginatedCollectionReturnsBadRequest(t *testing.T) {
	setup()
	defer shutdown()

	mux.HandleFunc("/api/v1/domains", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"message":"page size exceeds the configured maximum"}`)
	})

	page, err := client.Domains.List(PaginationOptions{PageSize: 50})
	if err == nil {
		t.Fatal("expected the API's HTTP 400 response to be returned as an error")
	}
	if page != nil {
		t.Fatalf("expected no page on error, got: %+v", page)
	}
	apiError, ok := err.(*ErrorResponse)
	if !ok || apiError.Response.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected HTTP 400 ErrorResponse, got: %#v", err)
	}
}
