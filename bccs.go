package goprsc

import (
	"fmt"
	"net/http"
)

const (
	inputBccType  = "incomming"
	outputBccType = "outgoing"
)

type BccService interface {
	// Get makes a GET request and fetches the specified BCC.
	Get(domain, account string) (*Bcc, error)

	// Create makes a POST request to create a new BCC.
	Create(domain, account, email string) error

	// Update makes a PUT request to update the specified BCC.
	Update(domain, account string, ur *BccUpdateRequest) error

	// Delete removes a BCC.
	Delete(domain, account string) error
}

type bccServiceImpl struct {
	client  *Client
	bccType string
}

// IncommingBccService handles communication with the incomming BCC APIs in the Postfix REST Server.
type IncommingBccService struct {
	*bccServiceImpl
}

// NewIncommingBccService creates a new IncommingBccService instance.
func NewIncommingBccService(c *Client) *IncommingBccService {
	return &IncommingBccService{
		bccServiceImpl: &bccServiceImpl{
			client:  c,
			bccType: inputBccType,
		},
	}
}

// OutgoingBccService handles communication with the outgoing BCC APIs in the Postfix REST Server.
type OutgoingBccService struct {
	*bccServiceImpl
}

// NewOutgoingBccService creates a new OutgoingBccService instance.
func NewOutgoingBccService(c *Client) *OutgoingBccService {
	return &OutgoingBccService{
		bccServiceImpl: &bccServiceImpl{
			client:  c,
			bccType: outputBccType,
		},
	}
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

// Get makes a GET request and fetches the specified BCC.
func (s *bccServiceImpl) Get(domain, account string) (*Bcc, error) {
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
func (s *bccServiceImpl) Create(domain, account, email string) error {
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
func (s *bccServiceImpl) Update(domain, account string, ur *BccUpdateRequest) error {
	req, err := s.client.NewRequest(http.MethodPut, s.getBccsURL(domain, account), ur)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return err
}

// Delete removes a BCC.
func (s *bccServiceImpl) Delete(domain, account string) error {
	req, err := s.client.NewRequest(http.MethodDelete, s.getBccsURL(domain, account), nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return err
}

func (s *bccServiceImpl) getBccsURL(domain, username string) string {
	return fmt.Sprintf("%s/%s/accounts/%s/bccs/%s", domainsURL, domain, username, s.bccType)
}
