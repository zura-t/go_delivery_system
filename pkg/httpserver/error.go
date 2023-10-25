package httpserver

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

func HttpErrorResponse(body io.ReadCloser) (map[string]any, error) {
	content, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	var errorMessage map[string]any
	err = json.Unmarshal(content, &errorMessage)
	if err != nil {
		return nil, err
	}

	return errorMessage, nil
}