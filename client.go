package goprsc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	defaultProtocol = "http"
	defaultHost     = "localhost"
	defaultPort     = "8080"
	defaultAPIPath  = "api/v1/"
	libraryVersion  = "0.1.0"
	userAgent       = "goprsc/" + libraryVersion
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

	// Domains is the service used for communication with the domains API.
	Domains *DomainService

	// Accounts is the service used for communication with the accounts API.
	Accounts *AccountService

	// Aliases is the service used for communication with the aliases API.
	Aliases *AliasService

	// OutputBccs is the service used for communication with the output BCC API.
	OutputBccs *OutgoingBccService

	// InputBccs is the service used for communication with the input BCC API.
	InputBccs *IncommingBccService
}

type service struct {
	client *Client
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
	c.Domains = (*DomainService)(&s)
	c.Accounts = (*AccountService)(&s)
	c.Aliases = (*AliasService)(&s)
	// Allocate separate structs for the bcc services as they have specific state
	c.OutputBccs = NewOutgoingBccService(s.client)
	c.InputBccs = NewIncommingBccService(s.client)

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

// NewRequest creates an API request. An URL relative to the API version path must be provided in urlStr.
func (c Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
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
	req.Header.Add("User-Agent", userAgent)

	return req, nil
}

// Do sends a request and returns an API response. The respose is JSON decoded and stored in the value
// pointed to by v.
func (c Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
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

func checkResponse(response *http.Response) error {
	if sc := response.StatusCode; sc >= 200 && sc <= 299 {
		return nil
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
	return fmt.Sprintf("%v %v %v", e.Method, e.Message, e.Path)
}
