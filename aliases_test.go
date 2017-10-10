package goprsc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestAlias_List(t *testing.T) {
	setup()
	defer shutdown()

	mux.HandleFunc("/api/v1/domains/example.com/aliases", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[
			{
				"id": 1,
				"alias": "contact",
				"email": "info@example.com",
				"created": "2016-07-11T08:42:31+0000",
				"updated": "2016-07-21T14:57:07+0000",
				"enabled": true
			}
		]`)
	})

	_, err := client.Aliases.List("example.com")
	if err != nil {
		t.Fatal("Failed to list aliases for domain 'example.com'")
	}
}

func TestAlias_Get(t *testing.T) {
	setup()
	defer shutdown()

	mux.HandleFunc("/api/v1/domains/example.com/aliases/contact", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{
			"id": 1,
			"alias": "contact",
			"email": "info@example.com",
			"created": "2016-07-11T08:42:31+0000",
			"updated": "2016-07-21T14:57:07+0000",
			"enabled": true
		}`)
	})

	alias, err := client.Aliases.Get("example.com", "contact")
	if err != nil {
		t.Fatal(err)
	}

	if alias.Email != "info@example.com" {
		t.Fatalf("expected: %v, got: %v", "test", alias.Name)
	}
}

func TestAlias_Create(t *testing.T) {
	setup()
	defer shutdown()

	mux.HandleFunc("/api/v1/domains/example.com/aliases", func(w http.ResponseWriter, r *http.Request) {
		var v AliasUpdateRequest

		if r.Method != http.MethodPost {
			t.Fatalf("expected method: %v, got: %v", http.MethodPost, r.Method)
		}

		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			t.Fatalf("unable to create alias: %v", err)
		}

		if v.Name != "contact" {
			t.Fatalf("expected: %v, got: %v", "contact", v.Name)
		}

		if v.Email != "info@example.com" {
			t.Fatalf("expected: %v, got: %v", "info@example.com", v.Email)
		}
	})

	if err := client.Aliases.Create("example.com", "contact", "info@example.com"); err != nil {
		t.Fatal(err)
	}
}

func TestAlias_Update(t *testing.T) {
	setup()
	defer shutdown()

	req := &AliasUpdateRequest{
		Name:  "contact",
		Email: "info@example.com",
	}

	mux.HandleFunc("/api/v1/domains/example.com/aliases/contact", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("Expected method: %v, got: %v", http.MethodPut, r.Method)
		}

		var v AliasUpdateRequest

		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if v.Name != "contact" {
			t.Fatalf("expected: %v, got: %v", "contact", v.Name)
		}
	})

	if err := client.Aliases.Update("example.com", "contact", req); err != nil {
		t.Fatal(err)
	}
}

func TestAlias_Delete(t *testing.T) {
	setup()
	defer shutdown()

	mux.HandleFunc("/api/v1/domains/example.com/aliases/contact", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("expected method: %v, got: %v", http.MethodDelete, r.Method)
		}
	})

	if err := client.Aliases.Delete("example.com", "contact"); err != nil {
		t.Fatal(err)
	}
}
