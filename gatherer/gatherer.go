package gatherer

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	// APIUserAgent identifies this library with the Kraken API
	APIUserAgent = "GO API Agent (https://github.com/baunes/api-gatherer)"
)

// Client allows make HTTP Calls
type Client interface {
	Get(string) (*Response, error)
	Query(string, string) (*Response, error)
}

type client struct {
}

// Response holds the response from the server
type Response struct {
	StatusCode int                     // Status code from the server
	Body       *map[string]interface{} // Content of the body
	Headers    map[string][]string     // HTTP Headers from the server
}

// NewClient creates a new client
func NewClient() Client {
	return &client{}
}

func (client *client) client() *http.Client {
	return http.DefaultClient
}

// Get makes an HTTP GET call an return de response (Body)
func (client *client) Get(urlToQuery string) (*Response, error) {
	return client.Query("GET", urlToQuery)
}

// Query makes an HTTP call an return de response (Body)
func (client *client) Query(method string, urlToQuery string) (*Response, error) {
	return client.doRequest(method, urlToQuery, nil)
}

// doRequest executes a HTTP Request to the Kraken API and returns the result
func (client *client) doRequest(method string, reqURL string, headers map[string]string) (*Response, error) {
	// Create request
	req, err := http.NewRequest(method, reqURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", APIUserAgent)

	// Exeute request
	resp, err := client.client().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse request
	var jsonData map[string]interface{}

	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		return nil, err
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       &jsonData,
	}, nil
}
