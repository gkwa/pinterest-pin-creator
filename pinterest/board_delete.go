package pinterest

import (
	"context"
	"fmt"
	"regexp"

	"pin-creator/internal/logger"
)

func (client *Client) DeleteBoards(ctx context.Context, regex string) error {
	log := logger.FromContext(ctx)
	boards, err := client.ListBoards(ctx)
	if err != nil {
		return err
	}

	r := regexp.MustCompile(regex)
	for _, board := range boards {
		if r.MatchString(board.Name) {
			err := client.doDeleteBoard(ctx, board.Id)
			if err != nil {
				log.Error(err, fmt.Sprintf("error deleting board %s", board.Name))
			} else {
				log.Info(fmt.Sprintf("Deleted board: %s", board.Name))
			}
		}
	}

	return nil
}

func (client *Client) doDeleteBoard(ctx context.Context, boardId string) error {
	url := fmt.Sprintf("%s%s/%s", client.baseUrl, "boards", boardId)

	req, err := client.createRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	_, err = client.executeRequest(ctx, req, 204)
	if err != nil {
		return fmt.Errorf("error executing request: %v", err)
	}

	return nil
}
