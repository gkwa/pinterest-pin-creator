package pinterest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"pin-creator/internal/logger"
)

func (c *Client) buildBoardsURL() string {
	return fmt.Sprintf("%s%s", c.baseUrl, "boards")
}

func (c *Client) createBoardsRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create new http request while doListBoards")
	}
	return req, nil
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
	log.V(1).Info(fmt.Sprintf("Request: %s", string(reqDumpJSON)))
}

func (c *Client) sendRequest(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}

func (c *Client) checkResponseStatus(res *http.Response) error {
	if res.StatusCode != 200 {
		errorResponse, err := handleWrongStatuscode(res)
		if err != nil {
			return err
		}
		return fmt.Errorf("statuscode not 200 while doListBoards. ErrorCode: %d ErrorMessage: %s", errorResponse.Code, errorResponse.Message)
	}
	return nil
}

func readResponseBody(res *http.Response) ([]byte, error) {
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %v", err)
	}
	return bodyBytes, nil
}

func unmarshalResponse(bodyBytes []byte) (listBoardResponseBody, error) {
	var listBoardResponseBody listBoardResponseBody
	err := json.Unmarshal(bodyBytes, &listBoardResponseBody)
	if err != nil {
		return listBoardResponseBody, fmt.Errorf("unable to unmarshal response body while doListBoards: %v", err)
	}
	return listBoardResponseBody, nil
}

func checkEmptyResponse(bodyBytes []byte) error {
	if len(bodyBytes) == 0 {
		return fmt.Errorf("empty JSON input")
	}
	return nil
}

func (c *Client) logResponse(ctx context.Context, response interface{}) {
	log := logger.FromContext(ctx)
	prettyJSON, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		log.Error(err, "unable to format JSON")
		return
	}
	log.V(1).Info(string(prettyJSON))
}

func (c *Client) prettyPrintJSON(ctx context.Context, data interface{}) {
	log := logger.FromContext(ctx)
	prettyJSON, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Error(err, "error marshaling JSON with indentation")
		return
	}
	log.V(1).Info(string(prettyJSON))
}
