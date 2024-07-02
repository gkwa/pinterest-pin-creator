package pinterest

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/go-logr/logr"
)

const (
	defaultRetryTimeout = 60 * time.Second
	maxRetries          = 5
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

func CreateOrFindBoard(ctx context.Context, client ClientInterface, log logr.Logger, boardName string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRetryTimeout)
	defer cancel()

	var boardID string
	operation := func() error {
		var err error
		boardID, err = findOrCreateBoard(ctx, client, log, boardName)
		if err != nil {
			log.Error(err, "Failed to find or create board", "boardName", boardName)
		}
		return err
	}

	backOff := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), maxRetries)

	err := backoff.Retry(operation, backoff.WithContext(backOff, ctx))
	if err != nil {
		if err == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout occurred while trying to create or find board: %w", err)
		}
		return "", fmt.Errorf("failed to create or find board after retries: %w", err)
	}

	return boardID, nil
}
