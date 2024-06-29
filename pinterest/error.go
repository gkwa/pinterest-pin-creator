package pinterest

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func handleWrongStatuscode(res *http.Response) (errorResponse, error) {
	errorResponse := errorResponse{}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return errorResponse, errors.New("unable to read response body while handleWrongStatuscode")
	}

	err = json.Unmarshal(bytes, &errorResponse)
	if err != nil {
		return errorResponse, errors.New("unable to unmarshal response body while handleWrongStatuscode")
	}

	return errorResponse, nil
}
