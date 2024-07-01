package pinterest

import (
	"context"
	"net/http"
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
