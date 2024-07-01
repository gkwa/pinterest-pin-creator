package pinterest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"pin-creator/internal/logger"
)

func (c *Client) DeleteBoards(ctx context.Context, regex string) error {
	log := logger.FromContext(ctx)
	boards, err := c.ListBoards(ctx)
	if err != nil {
		return err
	}

	r := regexp.MustCompile(regex)
	for _, board := range boards {
		if r.MatchString(board.Name) {
			err := c.doDeleteBoard(board.Id)
			if err != nil {
				log.Error(err, fmt.Sprintf("error deleting board %s", board.Name))
			} else {
				log.Info(fmt.Sprintf("Deleted board: %s", board.Name))
			}
		}
	}

	return nil
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
