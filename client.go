package goprsc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
)

const (
	defaultProtocol = "http"
	defaultHost     = "localhost"
	defaultPort     = "8080"
	defaultAPIPath  = "api/v1/"
	libraryVersion  = "0.1.0"
	userAgent       = "goprsc/" + libraryVersion + " (" + runtime.GOOS + " " + runtime.GOARCH + ")"
	mediaType       = "application/json"
)

// Client manages communication with the Postfix REST Server API.
type Client struct {
	client *http.Client

	// The protocol used for API requests (defaults to http).
	Protocol string

	// The host to connect to (defaults to localhost).
	Host string

	// The port on which the client should connect to the server (defaults to 8080).
	Port string

	// UserAgent is the client user agent
	UserAgent string

	// Login represents the username of the logged in user
	Login string

	// The authentication token
	AuthToken string

	// The refresh token for retrieving new authentication token
	RefreshToken string

	// Auth is the service used for communication with the authentication API.
	Auth *AuthService

	// Domains is the service used for communication with the domains API.
	Domains *DomainService

	// Accounts is the service used for communication with the accounts API.
	Accounts *AccountService

	// Aliases is the service used for communication with the aliases API.
	Aliases *AliasService

	// OutputBccs is the service used for communication with the output BCC API.
	OutputBccs *OutgoingBccService

	// InputBccs is the service used for communication with the input BCC API.
	InputBccs *IncomingBccService
}

type service struct {
	client *Client
}

// RefreshTokenRequest represents a request for retrieving a new access token
type RefreshTokenRequest struct {
	Login        string `json:"login"`
	RefreshToken string `json:"refreshToken"`
}

// DefaultClient is the default Client that works with the default HTTP client and
// connects to port 8080 on localhost. Use this one if you don't need specific http.Client
// implementation, protocol, host/IP or port number.
var DefaultClient = NewClient(nil)

// NewClient returns a new Postfix REST Server API client.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	c := &Client{
		client:    httpClient,
		Protocol:  defaultProtocol,
		Host:      defaultHost,
		Port:      defaultPort,
		UserAgent: userAgent,
	}
	s := service{client: c} // Reuse a single struct instead of allocating one for each service
	c.Auth = (*AuthService)(&s)
	c.Domains = (*DomainService)(&s)
	c.Accounts = (*AccountService)(&s)
	c.Aliases = (*AliasService)(&s)
	// Allocate separate structs for the bcc services as they have specific state
	c.OutputBccs = NewOutgoingBccService(s.client)
	c.InputBccs = NewIncomingBccService(s.client)

	return c
}

// NewClientWithOptions returns a new client with the given ClientOptions applied.
func NewClientWithOptions(httpClient *http.Client, options ...ClientOption) (*Client, error) {
	client := NewClient(httpClient)
	for _, option := range options {
		if err := option(client); err != nil {
			return nil, err
		}
	}
	return client, nil
}

// ClientOption is an option to the client that can be passed to NewClientWithOptions().
type ClientOption func(*Client) error

// HTTPSProtocolOption is a client option for using the HTTPS protocol.
func HTTPSProtocolOption() ClientOption {
	return func(c *Client) error {
		c.Protocol = "https"
		return nil
	}
}

// HostOption is a client option for setting the hostname or IP address of the server.
func HostOption(host string) ClientOption {
	return func(c *Client) error {
		_, err := url.Parse(fmt.Sprintf("%s://%s:%s", c.Protocol, host, c.Port))
		if err != nil {
			return err
		}

		c.Host = host

		return nil
	}
}

// PortOption is a client option for setting the port number on which the server is listening.
func PortOption(port string) ClientOption {
	return func(c *Client) error {
		_, err := url.Parse(fmt.Sprintf("%s://%s:%s", c.Protocol, c.Host, port))
		if err != nil {
			return err
		}

		c.Port = port

		return nil
	}
}

// UserAgentOption is a client option for setting the user agent.
func UserAgentOption(userAgent string) ClientOption {
	return func(c *Client) error {
		c.UserAgent = fmt.Sprintf("%s+%s", userAgent, c.UserAgent)
		return nil
	}
}

// AuthOption is a client option for setting the authentication tokens.
func AuthOption(login, authToken, refreshToken string) ClientOption {
	return func(c *Client) error {
		c.Login = login
		c.AuthToken = authToken
		c.RefreshToken = refreshToken
		return nil
	}
}

// NewRequest creates an API request. An URL relative to the API version path must be provided in urlStr.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rurl, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	baseURLStr := fmt.Sprintf("%s://%s:%s/%s", c.Protocol, c.Host, c.Port, defaultAPIPath)

	baseURL, err := url.Parse(baseURLStr)
	if err != nil {
		return nil, err
	}
	url := baseURL.ResolveReference(rurl)

	buf := new(bytes.Buffer)
	if body != nil {
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", c.UserAgent)
	if len(c.AuthToken) > 0 {
		req.Header.Add("Authorization", "Bearer "+c.AuthToken)
	}

	return req, nil
}

// Do sends a request and returns an API response. The respose is JSON decoded and stored in the value
// pointed to by v.
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusUnauthorized && len(req.Header.Get("X-GOPRSC-Refresh")) == 0 && len(c.RefreshToken) > 0 {
		authResponse, err := c.refreshTokens()
		if err != nil {
			return nil, err
		}
		// Resend the original request using the new authentication token
		req.Header.Set("Authorization", "Bearer "+authResponse.AuthToken)
		resp, err = c.client.Do(req)
	}

	if err := checkResponse(resp); err != nil {
		return resp, err
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return nil, err
		}
	}

	return resp, err
}

func (c *Client) refreshTokens() (*AuthResponse, error) {
	c.AuthToken = ""
	rr := &RefreshTokenRequest{
		Login:        c.Login,
		RefreshToken: c.RefreshToken,
	}
	refreshRequest, err := c.NewRequest(http.MethodPost, "auth/refresh-token", rr)
	if err != nil {
		return nil, err
	}
	refreshRequest.Header.Add("X-GOPRSC-Refresh", "1")
	authResponse := &AuthResponse{}
	_, err = c.Do(refreshRequest, &authResponse)
	if err != nil {
		return nil, err
	}
	c.AuthToken = authResponse.AuthToken
	c.RefreshToken = authResponse.RefreshToken
	return authResponse, nil
}

func checkResponse(response *http.Response) error {
	if sc := response.StatusCode; sc >= 200 && sc <= 299 {
		return nil
	}

	if sc := response.StatusCode; sc == http.StatusUnauthorized || sc == http.StatusForbidden {
		unauthorized := &ErrorResponse{
			Response: response,
			Message:  "Unauthorized. Please log in",
		}
		return unauthorized
	}

	errResponse := &ErrorResponse{Response: response}
	data, err := ioutil.ReadAll(response.Body)
	if err == nil && len(data) > 0 {
		if err := json.Unmarshal(data, errResponse); err != nil {
			return err
		}
	}

	return errResponse
}

// ErrorResponse represents an error caused by an API request.
type ErrorResponse struct {
	// The HTTP response
	Response *http.Response

	// Message is the error message received as response to the API request.
	Message string `json:"message"`

	// Path is the URL of the request.
	Path string `json:"path"`

	// Method is the HTTP method of the request.
	Method string `json:"method"`
}

func (e ErrorResponse) Error() string {
	if len(e.Method) > 0 {
		return fmt.Sprintf("%v %v %v", e.Method, e.Message, e.Path)
	} else if len(e.Message) > 0 {
		return fmt.Sprintf("%v", e.Message)
	} else {
		return fmt.Sprintln("Unknown error")
	}
}
