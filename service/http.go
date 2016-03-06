package service

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

var (
	httpClient *http.Client
)

func init() {
	defTransport := http.DefaultTransport.(*http.Transport)
	transport := *defTransport
	transport.TLSNextProto = make(map[string]func(authority string, c *tls.Conn) http.RoundTripper)
	httpClient = &http.Client{
		Transport: &transport,
	}
}

// jsonRequest encodes the given data as JSON object and performs an HTTP request with it.
func jsonRequest(method, url string, data interface{}) (*http.Response, error) {
	json, err := json.Marshal(data)
	if err != nil {
		return nil, maskAny(err)
	}
	resp, err := request(method, url, "application/json", json)
	if err != nil {
		return nil, maskAny(err)
	}
	return resp, nil
}

// request creates and performs an HTTP request.
func request(method, url, contentType string, data []byte) (*http.Response, error) {
	var reader io.Reader
	var contentLength string
	if data != nil {
		reader = bytes.NewReader(data)
		contentLength = strconv.Itoa(len(data))
	}
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, maskAny(err)
	}
	if contentLength != "" {
		req.Header.Set("Content-Length", contentLength)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, maskAny(err)
	}
	return resp, err
}

func parseResponse(resp *http.Response, data interface{}) error {
	defer resp.Body.Close()
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return maskAny(err)
	}

	if !isStatusOK(resp.StatusCode) {
		var er ErrorResponse
		if err := json.Unmarshal(raw, &er); err == nil {
			return er
		}

		return maskAny(fmt.Errorf("invalid status code=%d body='%s'", resp.StatusCode, string(raw)))
	}

	if err := json.Unmarshal(raw, data); err != nil {
		return maskAny(err)
	}
	return nil
}

func isStatusOK(code int) bool {
	return code >= 200 && code < 300
}
