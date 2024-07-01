package pinterest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"pin-creator/internal/logger"
)

type BoardInfo struct {
	Id   string
	Name string
}

type BoardData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Privacy     string `json:"privacy"`
}

type listBoardResponseBody struct {
	Items    []boardRequestBody `json:"items"`
	Bookmark string             `json:"bookmark"`
}

type boardRequestBody struct {
	Id          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Owner       ownerRequestBody `json:"owner"`
	Privacy     string           `json:"privacy"`
}

type ownerRequestBody struct {
	Username string `json:"username"`
}

func (c *Client) ListBoards(ctx context.Context) ([]BoardInfo, error) {
	boardInfos := make([]BoardInfo, 0)

	listBoardResponseBody, err := c.doListBoards(ctx)
	if err != nil {
		return nil, err
	}

	for _, item := range listBoardResponseBody.Items {
		boardInfos = append(boardInfos, BoardInfo{
			Id:   item.Id,
			Name: item.Name,
		})
		// log.Printf("Board found: %s", item.Name)
	}

	return boardInfos, nil
}

func (c *Client) CreateBoard(ctx context.Context, boardData BoardData) error {
	createBoardRequestBody := BoardData{
		Name:        boardData.Name,
		Description: boardData.Description,
		Privacy:     boardData.Privacy,
	}

	return c.doCreateBoard(createBoardRequestBody)
}

func (c *Client) doDeleteBoard(boardId string) error {
	url := fmt.Sprintf("%s%s/%s", c.baseUrl, "boards", boardId)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return errors.New("unable to create new http request while deleteBoard")
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	req.Header.Add("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return errors.New("unable to send request while deleteBoard")
	}

	if res.StatusCode != 204 {
		errorResponse, err := handleWrongStatuscode(res)
		if err != nil {
			return err
		}
		return fmt.Errorf("statuscode not 204 while deleteBoard. ErrorCode: %d ErrorMessage: %s", errorResponse.Code, errorResponse.Message)
	}

	return nil
}

func (c *Client) doListBoards(ctx context.Context) (listBoardResponseBody, error) {
	log := logger.FromContext(ctx)

	listBoardResponseBody := listBoardResponseBody{}

	url := fmt.Sprintf("%s%s", c.baseUrl, "boards")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return listBoardResponseBody, fmt.Errorf("unable to create new http request while doListBoards")
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	req.Header.Add("Content-Type", "application/json")

	// Marshal and indent request
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
		return listBoardResponseBody, fmt.Errorf("error dumping request: %v", err)
	}
	log.V(1).Info(fmt.Sprintf("Request: %s", string(reqDumpJSON)))

	res, err := c.httpClient.Do(req)
	if err != nil {
		return listBoardResponseBody, errors.New("unable to send request while doListBoards")
	}
	defer res.Body.Close()

	// https://developers.pinterest.com/docs/api/v5/boards-list
	if res.StatusCode != 200 {
		errorResponse, err := handleWrongStatuscode(res)
		if err != nil {
			return listBoardResponseBody, err
		}
		return listBoardResponseBody, fmt.Errorf("statuscode not 200 while doListBoards. ErrorCode: %d ErrorMessage: %s", errorResponse.Code, errorResponse.Message)
	}

	var response interface{}

	// Read the entire body
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return listBoardResponseBody, fmt.Errorf("unable to read response body: %v", err)
	}

	// Always close the original body when you're done with it
	defer res.Body.Close()

	// First unmarshaling
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return listBoardResponseBody, fmt.Errorf("unable to decode response body: %v", err)
	}

	prettyJSON, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		return listBoardResponseBody, fmt.Errorf("unable to format JSON: %v", err)
	}

	log.V(1).Info(string(prettyJSON))
	log.V(1).Info(fmt.Sprintf("status: %d", res.StatusCode))

	// Check if the JSON is empty
	if len(bodyBytes) == 0 {
		return listBoardResponseBody, fmt.Errorf("empty JSON input")
	}

	// Unmarshal the JSON data into a map
	var data map[string]interface{}
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		return listBoardResponseBody, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	// Marshal the data with indentation
	prettyJSON, err = json.MarshalIndent(data, "", "    ")
	if err != nil {
		return listBoardResponseBody, fmt.Errorf("error marshaling JSON with indentation: %v", err)
	}

	// Print the indented JSON
	log.V(1).Info(string(prettyJSON))

	err = json.Unmarshal(bodyBytes, &listBoardResponseBody)
	if err != nil {
		return listBoardResponseBody, fmt.Errorf("unable to unmarshal response body while doListBoards: %v", err)
	}

	return listBoardResponseBody, nil
}

func (c *Client) doCreateBoard(body BoardData) error {
	url := fmt.Sprintf("%s%s", c.baseUrl, "boards")

	bodyBytes, err := json.MarshalIndent(body, "", "  ")
	if err != nil {
		return errors.New("unable to marshal body while doCreateBoard")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return errors.New("unable to create new http request while doCreateBoard")
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	req.Header.Add("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return errors.New("unable to send request while doCreateBoard")
	}

	if res.StatusCode != 201 {
		errorResponse, err := handleWrongStatuscode(res)
		if err != nil {
			return err
		}
		return fmt.Errorf("statuscode not 201 while doCreateBoard. ErrorCode: %d ErrorMessage: %s", errorResponse.Code, errorResponse.Message)
	}

	return nil
}
