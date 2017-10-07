package goprsc

import "fmt"
import "net/http"

// AccountService is an interface for managing accounts with the Postfix REST Server API.
type AccountService interface {
	List(domain string) ([]Account, error)
	Get(domain, username string) (*Account, error)
	Create(domain, username, password string) error
	Update(domain, username string, req *AccountUpdateRequest) error
	Delete(domain, username string) error
}

// AccountServiceImpl handles communication with the account related API.
type AccountServiceImpl struct {
	client *Client
}

// Account is an instance of a account (an email address)
type Account struct {
	ID       int      `json:"id"`
	Username string   `json:"username"`
	Domain   string   `json:"domain"`
	DomainID int      `json:"domainId"`
	Enabled  bool     `json:"enabled"`
	Created  DateTime `json:"created"`
	Updated  DateTime `json:"updated"`
}

// AccountUpdateRequest is a data structure that caries account update information
type AccountUpdateRequest struct {
	Username        string `json:"username,omitempty"`
	Password        string `json:"password,omitempty"`
	ConfirmPassword string `json:"confirmPassword,omitempty"`
	Enabled         bool   `json:"enabled"`
}

// List makes a GET request for all registered accounts in the specified domain.
func (s AccountServiceImpl) List(domain string) ([]Account, error) {
	req, err := s.client.NewRequest(http.MethodGet, getAccountsURL(domain), nil)
	if err != nil {
		return nil, err
	}

	var accounts []Account
	_, err = s.client.Do(req, &accounts)
	if err != nil {
		return nil, err
	}

	return accounts, err
}

// Get returns the account with the given username on the given domain.
func (s AccountServiceImpl) Get(domain, username string) (*Account, error) {
	req, err := s.client.NewRequest(http.MethodGet, fmt.Sprintf("%v/%v", getAccountsURL(domain), username), nil)
	if err != nil {
		return nil, err
	}

	var account Account
	_, err = s.client.Do(req, &account)
	if err != nil {
		return nil, err
	}

	return &account, err
}

// Create creates a new account with the given username in the given domain.
func (s AccountServiceImpl) Create(domain, username, password string) error {
	ur := &AccountUpdateRequest{
		Username:        username,
		Password:        password,
		ConfirmPassword: password,
		Enabled:         true,
	}

	req, err := s.client.NewRequest(http.MethodPost, getAccountsURL(domain), ur)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return err
}

// Update updates the specified account.
func (s AccountServiceImpl) Update(domain, username string, updateRequest *AccountUpdateRequest) error {
	req, err := s.client.NewRequest(http.MethodPut, fmt.Sprintf("%v/%v", getAccountsURL(domain), username), updateRequest)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return err
}

// Delete removes the account specified with the given domain and username
func (s AccountServiceImpl) Delete(domain, username string) error {
	req, err := s.client.NewRequest(http.MethodDelete, fmt.Sprintf("%v/%v", getAccountsURL(domain), username), nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return err
}

func getAccountsURL(domain string) string {
	return fmt.Sprintf("%s/%s/accounts", domainsURL, domain)
}
