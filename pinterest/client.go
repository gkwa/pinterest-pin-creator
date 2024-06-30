package pinterest

import (
	"net/http"
	"regexp"

	log "github.com/sirupsen/logrus"
)

const (
	baseUrl = "https://api-sandbox.pinterest.com/v5/"
)

type ClientInterface interface {
	CreatePin(pinData PinData) error
	ListBoards() ([]BoardInfo, error)
	CreateBoard(boardData BoardData) error
	DeleteBoards(regex string) error
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

func (c *Client) DeleteBoards(regex string) error {
	boards, err := c.ListBoards()
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
				log.Infof("Deleted board: %s", board.Name)
			}
		}
	}

	return nil
}
