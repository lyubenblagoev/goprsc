package goprsc

import (
	"fmt"
	"net/http"
)

const (
	inputBccType  = "incomming"
	outputBccType = "outgoing"
)

// BccService is an interface for managing BCCs with the Postfix REST Server API.
type BccService interface {
	Get(domain, account string) (*Bcc, error)
	Create(domain, account, email string) error
	Update(domain, account string, ur *BccUpdateRequest) error
	Delete(domain, account string) error
}

// Bcc is a blind carbon copy for specific account
type Bcc struct {
	ID        int      `json:"id"`
	AccountID int      `json:"accountId"`
	Email     string   `json:"email"`
	Enabled   bool     `json:"enabled"`
	Created   DateTime `json:"created"`
	Updated   DateTime `json:"updated"`
}

// BccUpdateRequest carrires BCC update information.
type BccUpdateRequest struct {
	Email   string `json:"email,omitempty"`
	Enabled bool   `json:"enabled"`
}

// bccServiceImpl is an internal implementation that is wrapped by more
// concrete BccService implementations.
type bccServiceImpl struct {
	client  *Client
	bccType string
}

// OutputBccServiceImpl handles communication with the output BCC API.
type OutputBccServiceImpl struct {
	*bccServiceImpl
}

// NewOutputBccService creates a new OutputBccServiceImpl instance.
func NewOutputBccService(client *Client) *OutputBccServiceImpl {
	return &OutputBccServiceImpl{
		bccServiceImpl: &bccServiceImpl{
			client:  client,
			bccType: outputBccType,
		},
	}
}

// InputBccServiceImpl handles communication with the input BCC API.
type InputBccServiceImpl struct {
	*bccServiceImpl
}

// NewInputBccService creates a new InputBccServiceImpl instance.
func NewInputBccService(client *Client) *InputBccServiceImpl {
	return &InputBccServiceImpl{
		bccServiceImpl: &bccServiceImpl{
			client:  client,
			bccType: inputBccType,
		},
	}
}

// Get makes a GET request and fetches the specified BCC.
func (s bccServiceImpl) Get(domain, account string) (*Bcc, error) {
	req, err := s.client.NewRequest(http.MethodGet, s.getBccsURL(domain, account), nil)
	if err != nil {
		return nil, err
	}

	var bcc Bcc
	_, err = s.client.Do(req, &bcc)
	if err != nil {
		return nil, err
	}

	return &bcc, nil
}

// Create makes a POST request to create a new BCC.
func (s bccServiceImpl) Create(domain, account, email string) error {
	ur := &BccUpdateRequest{
		Email:   email,
		Enabled: true,
	}

	req, err := s.client.NewRequest(http.MethodPost, s.getBccsURL(domain, account), ur)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return err
}

// Update makes a PUT request to update the specified BCC.
func (s bccServiceImpl) Update(domain, account string, ur *BccUpdateRequest) error {
	req, err := s.client.NewRequest(http.MethodPut, s.getBccsURL(domain, account), ur)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return err
}

// Delete removes a BCC.
func (s bccServiceImpl) Delete(domain, account string) error {
	req, err := s.client.NewRequest(http.MethodDelete, s.getBccsURL(domain, account), nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return err
}

func (s bccServiceImpl) getBccsURL(domain, username string) string {
	return fmt.Sprintf("%s/%s/accounts/%s/bccs/%s", domainsURL, domain, username, s.bccType)
}
