package pinterest

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	log "github.com/sirupsen/logrus"
)

const (
	baseUrl = "https://api-sandbox.pinterest.com/v5/"
)

type ClientInterface interface {
	CreatePin(ctx context.Context, pinData PinData) error
	ListBoards(ctx context.Context) ([]BoardInfo, error)
	CreateBoard(ctx context.Context, boardData BoardData) error
	DeleteBoards(ctx context.Context, regex string) error
}

type Client struct {
	httpClient  *http.Client
	accessToken string
	baseUrl     string
}

func NewClient(accessToken string) *Client {
	return &Client{
		httpClient:  &http.Client{},
		accessToken: accessToken,
		baseUrl:     baseUrl,
	}
}

func (c *Client) DeleteBoards(ctx context.Context, regex string) error {
	boards, err := c.ListBoards(ctx)
	if err != nil {
		return err
	}

	r := regexp.MustCompile(regex)
	for _, board := range boards {
		if r.MatchString(board.Name) {
			err := c.doDeleteBoard(board.Id)
			if err != nil {
				log.Errorf("Error deleting board %s: %v", board.Name, err)
			} else {
				log.Infof(fmt.Sprintf("Deleted board: %s", board.Name))
			}
		}
	}

	return nil
}
