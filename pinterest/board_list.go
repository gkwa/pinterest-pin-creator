package pinterest

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-logr/logr"
)

type listBoardResponseBody struct {
	Items    []boardRequestBody `json:"items"`
	Bookmark string             `json:"bookmark"`
}

func (c *Client) ListBoards(ctx context.Context) ([]BoardInfo, error) {
	listBoardResponseBody, err := c.doListBoards(ctx)
	if err != nil {
		return nil, err
	}

	boardInfos := make([]BoardInfo, 0, len(listBoardResponseBody.Items))
	for _, item := range listBoardResponseBody.Items {
		boardInfos = append(boardInfos, BoardInfo{
			Id:   item.Id,
			Name: item.Name,
		})
	}

	return boardInfos, nil
}

func (c *Client) doListBoards(ctx context.Context) (listBoardResponseBody, error) {
	url := c.buildBoardsURL()
	req, err := c.createRequest("GET", url, nil)
	if err != nil {
		return listBoardResponseBody{}, fmt.Errorf("error creating request: %v", err)
	}

	responseBody, err := c.executeRequest(ctx, req, 200)
	if err != nil {
		return listBoardResponseBody{}, fmt.Errorf("error executing request: %v", err)
	}

	var listBoardResponseBody listBoardResponseBody
	if err := json.Unmarshal(responseBody, &listBoardResponseBody); err != nil {
		return listBoardResponseBody, fmt.Errorf("unable to unmarshal response body: %v", err)
	}

	return listBoardResponseBody, nil
}

func BoardIdByName(boards []BoardInfo, boardName string) (string, error) {
	for _, board := range boards {
		if board.Name == boardName {
			return board.Id, nil
		}
	}

	return "", fmt.Errorf("board %s not found", boardName)
}

func findBoard(ctx context.Context, client ClientInterface, log logr.Logger, boardName string) (string, error) {
	log.V(2).Info("Attempting to list boards")
	boards, err := client.ListBoards(ctx)
	if err != nil {
		return "", fmt.Errorf("error listing boards: %w", err)
	}

	for _, board := range boards {
		if board.Name == boardName {
			log.V(2).Info("board found", "boardName", boardName, "boardId", board.Id)
			return board.Id, nil
		}
	}

	return "", ErrBoardNotFound{BoardName: boardName}
}
