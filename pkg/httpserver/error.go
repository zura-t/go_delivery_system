package httpserver

import (
	"encoding/json"
	"io"
)

func HttpErrorResponse(body io.ReadCloser) (any, error) {
	content, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	var errorMessage any
	err = json.Unmarshal(content, &errorMessage)
	if err != nil {
		return nil, err
	}

	return errorMessage, nil
}