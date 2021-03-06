package goprsc

import (
	"fmt"
	"net/http"
)

// AliasService handles communication with the alias APIs in the Postfix REST Server.
type AliasService service

// Alias is an email address alias.
type Alias struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	Enabled bool     `json:"enabled"`
	Created DateTime `json:"created"`
	Updated DateTime `json:"updated"`
}

// AliasUpdateRequest carries alias update information.
type AliasUpdateRequest struct {
	Name    string `json:"name,omitempty"`
	Email   string `json:"email,omitempty"`
	Enabled bool   `json:"enabled"`
}

// List makes a GET request for all aliases for the given domain.
func (s *AliasService) List(domain string) ([]Alias, error) {
	req, err := s.client.NewRequest(http.MethodGet, getAliasesURL(domain), nil)
	if err != nil {
		return nil, err
	}

	var aliases []Alias
	_, err = s.client.Do(req, &aliases)
	if err != nil {
		return nil, err
	}

	return aliases, err
}

// Get retrieves information for an alias.
func (s *AliasService) Get(domain, alias string) ([]Alias, error) {
	req, err := s.client.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", getAliasesURL(domain), alias), nil)
	if err != nil {
		return nil, err
	}

	var aliases []Alias
	_, err = s.client.Do(req, &aliases)
	if err != nil {
		return nil, err
	}

	return aliases, err
}

// GetForEmail retrieves an alias for specific account and target email.
func (s *AliasService) GetForEmail(domain, alias, email string) (*Alias, error) {
	req, err := s.client.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s/%s", getAliasesURL(domain), alias, email), nil)
	if err != nil {
		return nil, err
	}

	var a Alias
	_, err = s.client.Do(req, &a)
	if err != nil {
		return nil, err
	}

	return &a, err
}

// Create makes a POST request to create a new alias.
func (s *AliasService) Create(domain, alias, email string) error {
	ur := &AliasUpdateRequest{
		Name:    alias,
		Email:   email,
		Enabled: true,
	}

	req, err := s.client.NewRequest(http.MethodPost, getAliasesURL(domain), ur)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return err
}

// Update makes a PUT request and updates the specified alias.
func (s *AliasService) Update(domain, alias, email string, ur *AliasUpdateRequest) error {
	req, err := s.client.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s/%s", getAliasesURL(domain), alias, email), ur)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return err
}

// Delete removes an alias.
func (s *AliasService) Delete(domain, alias, email string) error {
	req, err := s.client.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s/%s", getAliasesURL(domain), alias, email), nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return err
}

func getAliasesURL(domain string) string {
	return fmt.Sprintf("%s/%s/aliases", domainsURL, domain)
}
