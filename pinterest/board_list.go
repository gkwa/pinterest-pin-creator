package pinterest

import (
	"context"
	"encoding/json"
	"fmt"
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
