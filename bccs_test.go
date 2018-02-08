package goprsc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestBcc_Get(t *testing.T) {
	setup()
	defer shutdown()

	mux.HandleFunc("/api/v1/domains/example.com/accounts/test/bccs/incomming", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id":1, "accountId": 1, "email": "bcc@example.com","created":"2016-07-18T14:16:25+0000","updated":"2016-07-20T07:45:16+0000","enabled":true}`)
	})

	bcc, err := client.InputBccs.Get("example.com", "test")
	if err != nil {
		t.Fatal(err)
	}

	if bcc.Email != "bcc@example.com" {
		t.Fatalf("expected: %v, got: %v", "bcc.example.com", bcc.Email)
	}
}

func TestBcc_Create(t *testing.T) {
	setup()
	defer shutdown()

	mux.HandleFunc("/api/v1/domains/example.com/accounts/test/bccs/incomming", func(w http.ResponseWriter, r *http.Request) {
		var ur BccUpdateRequest

		if r.Method != http.MethodPost {
			t.Fatalf("expected method: %v, got: %v", http.MethodPost, r.Method)
		}

		if err := json.NewDecoder(r.Body).Decode(&ur); err != nil {
			t.Fatalf("unable to create bcc: %v", err)
		}

		if ur.Email != "bcc@example.com" {
			t.Fatalf("expected: %v, got: %v", "bcc@example.com", ur.Email)
		}
	})

	if err := client.InputBccs.Create("example.com", "test", "bcc@example.com"); err != nil {
		t.Fatal(err)
	}
}

func TestBcc_Update(t *testing.T) {
	setup()
	defer shutdown()

	mux.HandleFunc("/api/v1/domains/example.com/accounts/test/bccs/incomming", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("Expected method: %v, got: %v", http.MethodPut, r.Method)
		}

		var v BccUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if v.Email != "bcc2@example.com" {
			t.Fatalf("expected: %v, got: %v", "bcc2@example.com", v.Email)
		}
	})

	ur := &BccUpdateRequest{Email: "bcc2@example.com"}
	if err := client.InputBccs.Update("example.com", "test", ur); err != nil {
		t.Fatal(err)
	}
}

func TestBcc_Delete(t *testing.T) {
	setup()
	defer shutdown()

	mux.HandleFunc("/api/v1/domains/example.com/accounts/test/bccs/incomming", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("expected method: %v, got: %v", http.MethodDelete, r.Method)
		}
	})

	if err := client.InputBccs.Delete("example.com", "test"); err != nil {
		t.Fatal(err)
	}
}
