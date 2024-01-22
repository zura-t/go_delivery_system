package httpclient

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func NewHttpRequest(req any, method string, url string) (*http.Request, error) {
	var body io.Reader
	if req != nil {
		data, err := json.Marshal(req)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(data)
	} else {
		body = nil
	}

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	return request, nil
}
