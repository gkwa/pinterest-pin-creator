package pinterest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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

func (c *Client) ListBoards() ([]BoardInfo, error) {
	boardInfos := make([]BoardInfo, 0)

	listBoardResponseBody, err := c.doListBoards()
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

func (c *Client) CreateBoard(boardData BoardData) error {
	createBoardRequestBody := BoardData{
		Name:        boardData.Name,
		Description: boardData.Description,
		Privacy:     boardData.Privacy,
	}

	return c.doCreateBoard(createBoardRequestBody)
}

func (c *Client) deleteBoard(boardId string) error {
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

func (c *Client) doListBoards() (listBoardResponseBody, error) {
	listBoardResponseBody := listBoardResponseBody{}

	url := fmt.Sprintf("%s%s", c.baseUrl, "boards")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return listBoardResponseBody, errors.New("unable to create new http request while doListBoards")
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	req.Header.Add("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return listBoardResponseBody, errors.New("unable to send request while doListBoards")
	}

	if res.StatusCode != 200 {
		errorResponse, err := handleWrongStatuscode(res)
		if err != nil {
			return listBoardResponseBody, err
		}
		return listBoardResponseBody, fmt.Errorf("statuscode not 200 while doListBoards. ErrorCode: %d ErrorMessage: %s", errorResponse.Code, errorResponse.Message)
	}

	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return listBoardResponseBody, errors.New("unable to read response body while doListBoards")
	}

	err = json.Unmarshal(bytes, &listBoardResponseBody)
	if err != nil {
		return listBoardResponseBody, errors.New("unable to unmarshal response body while doListBoards")
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
