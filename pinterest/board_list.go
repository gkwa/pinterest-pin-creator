package pinterest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

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
	}

	return boardInfos, nil
}

func (c *Client) doListBoards(ctx context.Context) (listBoardResponseBody, error) {
	url := c.buildBoardsURL()
	req, err := c.createBoardsRequest(url)
	if err != nil {
		return listBoardResponseBody{}, err
	}

	c.addRequestHeaders(req)
	c.logRequestDetails(ctx, req)

	res, err := c.sendRequest(req)
	if err != nil {
		return listBoardResponseBody{}, errors.New("unable to send request while doListBoards")
	}
	defer res.Body.Close()

	if err := c.checkResponseStatus(res); err != nil {
		return listBoardResponseBody{}, err
	}

	bodyBytes, err := readResponseBody(res)
	if err != nil {
		return listBoardResponseBody{}, err
	}

	if err := checkEmptyResponse(bodyBytes); err != nil {
		return listBoardResponseBody{}, err
	}

	var response interface{}
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return listBoardResponseBody{}, fmt.Errorf("unable to decode response body: %v", err)
	}

	c.logResponse(ctx, response)
	c.prettyPrintJSON(ctx, response)

	listBoardResponseBody, err := unmarshalResponse(bodyBytes)
	if err != nil {
		return listBoardResponseBody, err
	}

	return listBoardResponseBody, nil
}
