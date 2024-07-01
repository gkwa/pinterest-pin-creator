package pinterest

import (
	"context"
	"fmt"
)

type ownerRequestBody struct {
	Username string `json:"username"`
}

type boardRequestBody struct {
	Id          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Owner       ownerRequestBody `json:"owner"`
	Privacy     string           `json:"privacy"`
}

func (c *Client) CreateBoard(ctx context.Context, boardData BoardData) error {
	createBoardRequestBody := BoardData{
		Name:        boardData.Name,
		Description: boardData.Description,
		Privacy:     boardData.Privacy,
	}

	return c.doCreateBoard(ctx, createBoardRequestBody)
}

func (c *Client) doCreateBoard(ctx context.Context, body BoardData) error {
	url := fmt.Sprintf("%s%s", c.baseUrl, "boards")

	req, err := c.createRequest("POST", url, body)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	_, err = c.executeRequest(ctx, req, 201)
	if err != nil {
		return fmt.Errorf("error executing request: %v", err)
	}

	return nil
}
