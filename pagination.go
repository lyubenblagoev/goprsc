package goprsc

import (
	"fmt"
	"net/url"
	"strconv"
)

// PaginationOptions configures a paginated collection request.
// Zero values leave the corresponding server default in effect.
type PaginationOptions struct {
	Page     int
	PageSize int
}

// PaginationMetadata describes the page returned by the API.
type PaginationMetadata struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

// PaginatedDomains is a page of domains returned by the API.
type PaginatedDomains struct {
	Items []Domain `json:"items"`
	PaginationMetadata
}

// PaginatedAccounts is a page of accounts returned by the API.
type PaginatedAccounts struct {
	Items []Account `json:"items"`
	PaginationMetadata
}

// PaginatedAliases is a page of aliases or recipients returned by the API.
type PaginatedAliases struct {
	Items []Alias `json:"items"`
	PaginationMetadata
}

func paginationURL(urlStr string, options []PaginationOptions) (string, error) {
	if len(options) == 0 {
		return urlStr, nil
	}
	if len(options) > 1 {
		return "", fmt.Errorf("at most one pagination options value may be provided")
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	query := u.Query()
	if options[0].Page != 0 {
		query.Set("page", strconv.Itoa(options[0].Page))
	}
	if options[0].PageSize != 0 {
		query.Set("pageSize", strconv.Itoa(options[0].PageSize))
	}
	u.RawQuery = query.Encode()
	return u.String(), nil
}
