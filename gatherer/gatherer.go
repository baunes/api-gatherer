package gatherer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
)

const (
	// APIUserAgent identifies this library with the Kraken API
	APIUserAgent = "GO API Agent (https://github.com/baunes/api-gatherer)"
)

// Client allows make HTTP Calls
type Client struct {
}

// Response holds the response from the server
type Response struct {
	StatusCode int                     // Status code from the server
	Body       *map[string]interface{} // Content of the body
	Headers    map[string][]string     // HTTP Headers from the server
}

func (client *Client) client() *http.Client {
	return http.DefaultClient
}

// Get makes an HTTP GET call an return de response (Body)
func (client *Client) Get(urlToQuery string) (*Response, error) {
	return client.Query("GET", urlToQuery)
}

// Query makes an HTTP call an return de response (Body)
func (client *Client) Query(method string, urlToQuery string) (*Response, error) {
	return client.doRequest(method, urlToQuery, nil)
}

// doRequest executes a HTTP Request to the Kraken API and returns the result
func (client *Client) doRequest(method string, reqURL string, headers map[string]string) (*Response, error) {
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

	// Check mime type of response
	mimeType, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil || mimeType != "application/json" {
		return nil, fmt.Errorf("heder Content-Type is '%s', but should be 'application/json'", mimeType)
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
