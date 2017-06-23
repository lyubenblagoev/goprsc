package goprsc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestDomain_List(t *testing.T) {
	setup()
	defer shutdown()

	mux.HandleFunc("/api/v1/domains", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"id":1,"created":"2016-07-18T14:16:25+0000","updated":"2016-07-18T14:16:25+0000","enabled":true,"name":"example.com"},{"id":2,"created":"2016-07-18T14:16:25+0000","updated":"2016-07-18T14:16:25+0000","enabled":true,"name":"example2.com"}]`)
	})

	_, err := client.Domains.List()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDomain_Get(t *testing.T) {
	setup()
	defer shutdown()

	mux.HandleFunc("/api/v1/domains/example.com", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id":1,"created":"2016-07-18T14:16:25+0000","updated":"2016-07-18T14:16:25+0000","enabled":true,"name":"example.com"}`)
	})

	_, err := client.Domains.Get("example.com")
	if err != nil {
		t.Fatal(err)
	}

}

func TestDomain_Create(t *testing.T) {
	setup()
	defer shutdown()

	mux.HandleFunc("/api/v1/domains", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected method: %v, got: %v", http.MethodPost, r.Method)
		}
		fmt.Fprint(w, `{"id":1, "enabled":true,"name":"example.net"}`)
	})

	if err := client.Domains.Create("example.net"); err != nil {
		t.Fatal(err)
	}
}

func TestDomain_Update(t *testing.T) {
	setup()
	defer shutdown()

	v := &DomainUpdateRequest{
		Name:    "example.org",
		Enabled: true,
	}

	mux.HandleFunc("/api/v1/domains/example.net", func(w http.ResponseWriter, r *http.Request) {
		var v DomainUpdateRequest

		if r.Method != http.MethodPut {
			t.Fatalf("expected method: %v, got: %v", http.MethodPut, r.Method)
		}

		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if v.Name != "example.org" {
			t.Fatalf("expected: %v, got: %v", "example.org", v.Name)
		}
	})

	if err := client.Domains.Update("example.net", v); err != nil {
		t.Fatal("failed to update domain", err)
	}
}

func TestDomain_Delete(t *testing.T) {
	setup()
	defer shutdown()

	mux.HandleFunc("/api/v1/domains/example.org", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("expected method: %v, got: %v", http.MethodDelete, r.Method)
		}
	})

	if err := client.Domains.Delete("example.org"); err != nil {
		t.Fatal("Failed to delete domain example.org")
	}
}
