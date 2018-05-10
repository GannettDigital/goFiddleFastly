package fiddle

import (
	"net/http"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
)

const BaseAddress = "https://fiddle.fastlydemo.net"

// Client contains information about our goFiddleFastly client.
type Client struct {
	// Address is the location of Fastly fiddle. Defaults to the Fastly Fiddle demo address.
	Address string

	// HTTPClient is the HTTP client used. Defaults to the go-cleanhttp default client.
	HTTPClient *http.Client
}

// DefaultClient creates a default goFiddleFastly client.
func DefaultClient() (*Client, error) {
	return &Client{
		Address:    BaseAddress,
		HTTPClient: cleanhttp.DefaultClient(),
	}, nil
}

// Post executes a POST request.
func (c *Client) Post(r *RequestInput) (*http.Response, error) {
	return c.Request("POST", r)
}

// Put executes a PUT request.
func (c *Client) Put(r *RequestInput) (*http.Response, error) {
	return c.Request("PUT", r)
}
