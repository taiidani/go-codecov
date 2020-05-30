// Package codecov defines an API client for the CodeCov REST API.
//
// https://docs.codecov.io/reference
package codecov

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type (
	// Client represents the Codecov API client
	Client struct {
		endpoint url.URL // The endpoint to send requests to
		token    string  // Token contains the authentication token used to make requests with
	}

	// Response tracks the base information returned from most Codecov API responses
	Response struct {
		Meta Meta
	}

	// ResponseError tracks the base information returned from Codecov API responses that have errored
	ResponseError struct {
		Meta  Meta
		Error Error
	}

	// Meta may be returned as part of some API responses
	Meta struct {
		Status int
	}

	// Error may be returned as part of some API responses
	Error struct {
		Reason string
	}
)

// NewClient returns an instance of the Codecov API client
func NewClient(token string) *Client {
	const defaultClientScheme = "https"
	const defaultClientHost = "codecov.io"
	const defaultClientPath = "/api/"

	c := &Client{
		token: token,
	}

	c.SetEndpoint(url.URL{
		Scheme: defaultClientScheme,
		Host:   defaultClientHost,
		Path:   defaultClientPath,
	})

	return c
}

// SetEndpoint will override the default Codecov endpoint with a custom location
// This is effective for testing strategies using `httptest`.
func (c *Client) SetEndpoint(u url.URL) {
	c.endpoint = u
}

func (c *Client) doRequest(request *http.Request, response interface{}) error {
	// Prefix the URL
	request.URL.Scheme = c.endpoint.Scheme
	request.URL.Host = c.endpoint.Host
	request.URL.Path = "/" + strings.Trim(c.endpoint.Path, "/") + "/" + strings.TrimPrefix(request.URL.Path, "/")

	// Authenticate the request for sending
	c.addAuthorization(request)

	// Send the request
	responseData, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer responseData.Body.Close()

	switch responseData.StatusCode {
	case 200:
		// Success, parse the response successfully
		err = json.NewDecoder(responseData.Body).Decode(&response)
	default:
		// Error handling
		// Extract the error message from the response and return it
		responseError := ResponseError{}
		err = json.NewDecoder(responseData.Body).Decode(&responseError)
		if err == nil {
			err = fmt.Errorf("%d: %s", responseData.StatusCode, responseError.Error.Reason)
		}
	}

	return err
}
