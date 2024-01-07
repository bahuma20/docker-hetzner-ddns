package request

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type loggingTransport struct{}

func (s *loggingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	bytes, _ := httputil.DumpRequestOut(r, true)

	resp, err := http.DefaultTransport.RoundTrip(r)
	// err is returned after dumping the response

	respBytes, _ := httputil.DumpResponse(resp, true)
	bytes = append(bytes, respBytes...)

	fmt.Printf("%s\n", bytes)

	return resp, err
}

func Request(httpMethod string, url url.URL, headers map[string]string, body []byte) ([]byte, error) {
	// Create client
	client := &http.Client{
		Transport: &loggingTransport{},
	}

	// Create request
	req, err := http.NewRequest(httpMethod, url.String(), bytes.NewBuffer(body))

	if err != nil {
		fmt.Println("Failure : ", err)
		return []byte{}, err
	}

	// Headers
	for key, element := range headers {
		req.Header.Add(key, element)
	}

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Failure : ", err)
		return []byte{}, err
	}

	// Read Response Body
	respBody, _ := ioutil.ReadAll(resp.Body)

	return respBody, nil
}
