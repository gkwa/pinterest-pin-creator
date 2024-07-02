package pinterest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"pin-creator/internal/logger"
)

func (c *Client) buildBoardsURL() string {
	return fmt.Sprintf("%s%s", c.baseUrl, "boards")
}

func (c *Client) addRequestHeaders(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	req.Header.Add("Content-Type", "application/json")
}

func (c *Client) logRequestDetails(ctx context.Context, req *http.Request) {
	log := logger.FromContext(ctx)
	reqDump := struct {
		Method  string      `json:"method"`
		URL     string      `json:"url"`
		Headers http.Header `json:"headers"`
	}{
		Method:  req.Method,
		URL:     req.URL.String(),
		Headers: req.Header,
	}
	reqDumpJSON, err := json.MarshalIndent(reqDump, "", "  ")
	if err != nil {
		log.Error(err, "error dumping request")
		return
	}
	log.V(2).Info(fmt.Sprintf("Request: %s", string(reqDumpJSON)))
}

func (c *Client) logResponse(ctx context.Context, bodyBytes []byte) {
	log := logger.FromContext(ctx)
	var response interface{}
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		log.Error(err, "unable to unmarshal response body")
		return
	}
	c.prettyPrintJSON(ctx, response)
}

func (c *Client) prettyPrintJSON(ctx context.Context, data interface{}) {
	log := logger.FromContext(ctx)
	prettyJSON, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Error(err, "error marshaling JSON with indentation")
		return
	}
	log.V(2).Info(string(prettyJSON))
}
