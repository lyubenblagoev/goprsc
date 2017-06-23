package goprsc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestAccount_List(t *testing.T) {
	setup()
	defer shutdown()

	mux.HandleFunc("/api/v1/domains/example.com/accounts", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"id":1,"username":"test","domain":"example.com","domainId":1,"created":"2016-07-18T14:16:25+0000","updated":"2016-07-20T07:45:16+0000","enabled":true}]`)
	})

	_, err := client.Accounts.List("example.com")
	if err != nil {
		t.Fatal("Failed to list accounts for domain 'example.com'")
	}
}

func TestAccount_Get(t *testing.T) {
	setup()
	defer shutdown()

	mux.HandleFunc("/api/v1/domains/example.com/accounts/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id":1,"username":"test","domain":"example.com","domainId":1,"created":"2016-07-18T14:16:25+0000","updated":"2016-07-20T07:45:16+0000","enabled":true}`)
	})

	acc, err := client.Accounts.Get("example.com", "test")
	if err != nil {
		t.Fatal(err)
	}

	if acc.Username != "test" {
		t.Fatalf("expected: %v, got: %v", "test", acc.Username)
	}
}

func TestAccount_Create(t *testing.T) {
	setup()
	defer shutdown()

	mux.HandleFunc("/api/v1/domains/example.com/accounts", func(w http.ResponseWriter, r *http.Request) {
		var v AccountUpdateRequest

		if r.Method != http.MethodPost {
			t.Fatalf("expected method: %v, got: %v", http.MethodPost, r.Method)
		}

		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			t.Fatalf("unable to create account: %v", err)
		}

		if v.Username != "test" {
			t.Fatalf("expected: %v, got: %v", "test", v.Username)
		}

		if v.Password != "testpass" {
			t.Fatalf("expected: %v, got: %v", "testpass", v.Password)
		}
	})

	if err := client.Accounts.Create("example.com", "test", "testpass"); err != nil {
		t.Fatal(err)
	}
}

func TestAccount_Update(t *testing.T) {
	setup()
	defer shutdown()

	req := &AccountUpdateRequest{
		Username: "test2",
	}

	mux.HandleFunc("/api/v1/domains/example.com/accounts/test", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("Expected method: %v, got: %v", http.MethodPut, r.Method)
		}

		var v AccountUpdateRequest

		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if v.Username != "test2" {
			t.Fatalf("expected: %v, got: %v", "test2", v.Username)
		}
	})

	if err := client.Accounts.Update("example.com", "test", req); err != nil {
		t.Fatal(err)
	}
}

func TestAccount_Delete(t *testing.T) {
	setup()
	defer shutdown()

	mux.HandleFunc("/api/v1/domains/example.com/accounts/test2", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("expected method: %v, got: %v", http.MethodDelete, r.Method)
		}
	})

	if err := client.Accounts.Delete("example.com", "test2"); err != nil {
		t.Fatal(err)
	}
}
