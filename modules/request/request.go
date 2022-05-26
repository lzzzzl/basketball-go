package request

import (
	"io/ioutil"
	"net/http"
)

// HTTPRequest http request
type HTTPRequest struct {
	URL     string
	Params  string
	Headers map[string]string
	Proxy   string
	TimeOut int
}

// HTTPGet http get
func (request *HTTPRequest) HTTPGet() (str string, err error) {
	req, err := http.NewRequest("GET", request.URL, nil)
	for key, value := range request.Headers {
		req.Header.Add(key, value)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
}
