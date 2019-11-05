package gatherer

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// client is the CircleCI client being tested.
	client *Client

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
)

func setup() *url.URL {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	url, err := url.Parse(server.URL)
	if err != nil {
		panic(fmt.Sprintf("couldn't parse test server URL: %s", server.URL))
	}

	client = &Client{}
	return url
}

func teardown() {
	defer server.Close()
}

func TestQueryWithInvalidUrl(t *testing.T) {
	_, err := client.Query("", "69699ññ3ñ2ñ1ñ3ñdsqñdññ·$%&/()")
	if err == nil || !strings.Contains(err.Error(), "invalid URL") {
		t.Errorf(`Client.Query(...) must fail: %s`, err)
	}
}

func TestQueryWithEmptyUrl(t *testing.T) {
	_, err := client.Query("", "")
	if err == nil || !strings.Contains(err.Error(), "unsupported protocol scheme") {
		t.Errorf(`Client.Query(...) must fail: %s`, err)
	}
}

func TestGetOnlyAcceptsJSON(t *testing.T) {
	baseURL := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")

	})

	_, err := client.Get(baseURL.String())
	if err == nil || !strings.Contains(err.Error(), "should be 'application/json'") {
		t.Errorf(`Client.Query(...) must fail: %s`, err)
	}

}

func TestGetJSON(t *testing.T) {
	baseURL := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.Header.Get("User-Agent"), "GO API Agent")

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"key_string": "value 1",
			"key_number": 12,
			"key_float": 99.0,
			"key_boolean_true": true,
			"key_boolean_false": false,
			"key_array_string": ["a", "b", "c"],
			"key_array_int": [1, 2, 3, 4, 5],
			"key_array_float": [6.0, 7.0, 8.0, 9.0, 10.0],
			"key_array_boolean": [true, false, true],
			"key_null": null,
			"key_object": {
				"key_string": "value 2",
				"key_number": 112,
				"key_float": 199.0,
				"key_boolean_true": true,
				"key_boolean_false": false,
				"key_array_string": ["x", "y", "z"],
				"key_array_int": [10, 20, 30, 40, 50],
				"key_array_float": [60.0, 70.0, 80.0, 90.0, 100.00],
				"key_array_boolean": [false, true, false],
				"key_null": null
			}
		}`)
	})

	response, err := client.Get(baseURL.String())
	if err != nil {
		t.Errorf(`Client.Get(%s) errored with %s`, baseURL.String(), err)
	}

	assert.Equal(t, "value 1", (*response.Body)["key_string"])
	assert.Equal(t, 12.0, (*response.Body)["key_number"])
	assert.Equal(t, 99.0, (*response.Body)["key_float"])
	assert.Equal(t, true, (*response.Body)["key_boolean_true"])
	assert.Equal(t, false, (*response.Body)["key_boolean_false"])
	assert.ElementsMatch(t, []string{"a", "b", "c"}, (*response.Body)["key_array_string"])
	assert.ElementsMatch(t, []float64{1, 2, 3, 4, 5}, (*response.Body)["key_array_int"])
	assert.ElementsMatch(t, []float64{6.0, 7.0, 8.0, 9.0, 10.0}, (*response.Body)["key_array_float"])
	assert.ElementsMatch(t, []bool{true, false, true}, (*response.Body)["key_array_boolean"])
	assert.Nil(t, (*response.Body)["key_null"])
	keyObject, isKeyObject := (*response.Body)["key_object"].(map[string]interface{})
	assert.True(t, isKeyObject, "key_object must not be nil")
	assert.Equal(t, "value 2", keyObject["key_string"])
	assert.Equal(t, 112.0, keyObject["key_number"])
	assert.Equal(t, 199.0, keyObject["key_float"])
	assert.Equal(t, true, keyObject["key_boolean_true"])
	assert.Equal(t, false, keyObject["key_boolean_false"])
	assert.ElementsMatch(t, []string{"x", "y", "z"}, keyObject["key_array_string"])
	assert.ElementsMatch(t, []float64{10, 20, 30, 40, 50}, keyObject["key_array_int"])
	assert.ElementsMatch(t, []float64{60.0, 70.0, 80.0, 90.0, 100.0}, keyObject["key_array_float"])
	assert.ElementsMatch(t, []bool{false, true, false}, keyObject["key_array_boolean"])
	assert.Nil(t, keyObject["key_null"])
}
