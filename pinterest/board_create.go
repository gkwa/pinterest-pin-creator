package pinterest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (c *Client) CreateBoard(ctx context.Context, boardData BoardData) error {
	createBoardRequestBody := BoardData{
		Name:        boardData.Name,
		Description: boardData.Description,
		Privacy:     boardData.Privacy,
	}

	return c.doCreateBoard(createBoardRequestBody)
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
