package goprsc

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	expectedProtocol = "https"
	expectedPort     = "80"
	expectedHost     = "127.0.0.1"
)

var (
	client *Client

	mux    *http.ServeMux
	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	hostport := server.URL[len("http://"):]
	host, port, err := net.SplitHostPort(hostport)
	if err != nil {
		panic(err)
	}

	client, err = NewClientWithOptions(nil, HostOption(host), PortOption(port))
	if err != nil {
		panic(err)
	}
}

func shutdown() {
	server.Close()
}

func TestClient_NewClientWithOptions(t *testing.T) {
	client, err := NewClientWithOptions(nil, PortOption(expectedPort), HTTPSProtocolOption(), HostOption(expectedHost))
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		desc     string
		expected string
		actual   string
	}{
		{"Protocol", expectedProtocol, client.Protocol},
		{"Port", expectedPort, client.Port},
		{"Host", expectedHost, client.Host},
	}

	for _, tc := range testCases {
		if tc.expected != tc.actual {
			t.Errorf("%s: expected=%s, got=%s", tc.desc, tc.expected, tc.actual)
		}
	}
}
