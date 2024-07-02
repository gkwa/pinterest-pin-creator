package pinterest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) createRequest(method, url string, body interface{}) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal body: %v", err)
		}
		bodyReader = bytes.NewBuffer(bodyBytes)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("unable to create new http request: %v", err)
	}

	c.addRequestHeaders(req)
	return req, nil
}

func (c *Client) executeRequest(ctx context.Context, req *http.Request, expectedStatus int) ([]byte, error) {
	c.logRequestDetails(ctx, req)

	res, err := c.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("unable to send request: %v", err)
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %v", err)
	}

	if res.StatusCode != expectedStatus {
		errorResponse, err := handleWrongStatuscode(res)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("unexpected status code %d. ErrorCode: %d ErrorMessage: %s", res.StatusCode, errorResponse.Code, errorResponse.Message)
	}

	c.logResponse(ctx, bodyBytes)
	return bodyBytes, nil
}
