package httpClient

import (
	"fmt"
	"io"
	"net/http"
)

type RequestBuildError struct {
	Url    string
	Method string
	Err    error
}

type RequestResponseError struct {
	statusCode int
	err        string
}

type RequestConfig struct {
	Url     string
	Params  map[string]string
	Headers map[string]string
	Body 	io.Reader
}

func (r *RequestBuildError) Error() string {
	return "Error making request to " + r.Url + " with method " + r.Method + " : " + r.Err.Error()
}

func (r *RequestResponseError) Error() string {
	return "Error making request: " + r.err + " with status code " + fmt.Sprint(r.statusCode)
}

// custom wrapper to make a get request
func Get(config RequestConfig) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, config.Url, config.Body)
	if err != nil {
		return nil, &RequestBuildError{Url: config.Url, Method: http.MethodGet, Err: err}
	}
	// set query params
	if config.Params != nil {
		queryParams := req.URL.Query()
		// add query params for each key value pair
		for key, value := range config.Params {
			queryParams.Add(key, value)
		}
		req.URL.RawQuery = queryParams.Encode()
	}

	// set headers if needed
	if config.Headers != nil {
		// add headers for each key value pair
		for key, value := range config.Headers {
			req.Header.Set(key, value)
		}
	}

	// make the request
	response, err := http.DefaultClient.Do(req)

	fmt.Println(req.URL)

	if err != nil && response == nil {
		return nil, &RequestResponseError{statusCode: 500, err: err.Error()}
	}
	if err != nil {
		return nil, &RequestResponseError{statusCode: response.StatusCode, err: err.Error()}
	}



 // return result
	return response, nil
}
