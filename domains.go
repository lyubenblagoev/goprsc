package goprsc

import "fmt"
import "net/http"

// domainsURL is the base address for all domain-related urls.
const domainsURL = "domains"

// DomainService is an interface for managing domains with the Postfix REST Server API.
type DomainService interface {
	List() ([]Domain, error)
	Get(string) (*Domain, error)
	Create(string) error
	Update(string, *DomainUpdateRequest) error
	Delete(string) error
}

// DomainServiceImpl handles communication with the domain related API.
type DomainServiceImpl struct {
	client *Client
}

// Domain represents a domain (e.g. example.com)
type Domain struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Enabled bool     `json:"enabled"`
	Created DateTime `json:"created"`
	Updated DateTime `json:"updated"`
}

// DomainUpdateRequest represents a request for domain update
type DomainUpdateRequest struct {
	Name    string `json:"name,omitempty"`
	Enabled bool   `json:"enabled"`
}

// List makes a GET request for all registered domains.
func (s DomainServiceImpl) List() ([]Domain, error) {
	req, err := s.client.NewRequest(http.MethodGet, domainsURL, nil)
	if err != nil {
		return nil, err
	}

	var domains []Domain
	if _, err := s.client.Do(req, &domains); err != nil {
		return nil, err
	}

	return domains, err
}

// Get makes a GET request for a specific domain specified with the domain parameter.
func (s DomainServiceImpl) Get(domain string) (*Domain, error) {
	req, err := s.client.NewRequest(http.MethodGet, fmt.Sprintf("%v/%v", domainsURL, domain), nil)
	if err != nil {
		return nil, err
	}

	d := new(Domain)
	if _, err := s.client.Do(req, d); err != nil {
		return nil, err
	}
	return d, err
}

// Create makes a POST request to the API to create a new domain.
func (s DomainServiceImpl) Create(domain string) error {
	ur := &DomainUpdateRequest{
		Name:    domain,
		Enabled: true,
	}

	req, err := s.client.NewRequest(http.MethodPost, domainsURL, ur)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return err
}

// Update makes a PUT request to update domain parameters
func (s DomainServiceImpl) Update(name string, updateRequest *DomainUpdateRequest) error {
	req, err := s.client.NewRequest(http.MethodPut, fmt.Sprintf("%v/%v", domainsURL, name), updateRequest)
	if err != nil {
		return nil
	}
	_, err = s.client.Do(req, nil)
	return err
}

// Delete makes a DELETE request to the API to delete the specified domain
func (s DomainServiceImpl) Delete(name string) error {
	req, err := s.client.NewRequest(http.MethodDelete, fmt.Sprintf("%v/%v", domainsURL, name), nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}
