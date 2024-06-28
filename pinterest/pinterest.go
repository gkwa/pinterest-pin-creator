package pinterest

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
)

const (
	baseUrl = "https://api-sandbox.pinterest.com/v5/"
)

type ClientInterface interface {
	CreatePin(pinData PinData) error
	ListBoards() ([]BoardInfo, error)
	CreateBoard(boardData BoardData) error
}

type Client struct {
	httpClient  *http.Client
	accessToken string
	baseUrl     string
}

type PinData struct {
	BoardId     string
	ImgPath     string
	Link        string
	Title       string
	Description string
	AltText     string
}

type BoardInfo struct {
	Id   string
	Name string
}

type BoardData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Privacy     string `json:"privacy"`
}

type createPinRequestBody struct {
	Link           string                 `json:"link"`
	Title          string                 `json:"title"`
	Description    string                 `json:"description"`
	AltText        string                 `json:"alt_text"`
	BoardId        string                 `json:"board_id"`
	BoardSectionId string                 `json:"board_section_id,omitempty"`
	MediaSource    mediaSourceRequestBody `json:"media_source"`
}

type mediaSourceRequestBody struct {
	SourceType  string `json:"source_type"`
	ContentType string `json:"content_type"`
	Data        string `json:"data"`
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

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewClient(accessToken string) *Client {
	return &Client{
		httpClient:  &http.Client{},
		accessToken: accessToken,
		baseUrl:     baseUrl,
	}
}

func (c *Client) CreatePin(pinData PinData) error {
	createPinRequestBody := createPinRequestBody{
		Link:        pinData.Link,
		Title:       pinData.Title,
		Description: pinData.Description,
		AltText:     pinData.AltText,
		BoardId:     pinData.BoardId,
		MediaSource: mediaSourceRequestBody{
			SourceType:  "image_base64",
			ContentType: "image/png",
			Data:        toBase64(pinData.ImgPath),
		},
	}

	err := c.doCreatePin(createPinRequestBody)
	if err != nil {
		return err
	}

	return nil
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
		log.Printf("Board found: %s", item.Name)
	}

	return boardInfos, nil
}

func (c *Client) CreateBoard(boardData BoardData) error {
	createBoardRequestBody := BoardData{
		Name:        boardData.Name,
		Description: boardData.Description,
		Privacy:     boardData.Privacy,
	}

	err := c.doCreateBoard(createBoardRequestBody)
	if err != nil {
		return err
	}

	return nil
}

func toBase64(imgPath string) string {
	bytes, err := ioutil.ReadFile(imgPath)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(bytes)
}

func (c *Client) doCreatePin(body createPinRequestBody) error {

	url := fmt.Sprintf("%s%s", c.baseUrl, "pins")

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return errors.New("unable to marshal body while doCreatePin")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return errors.New("unable to create new http request while doCreatePin")
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	req.Header.Add("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return errors.New("unable to send request while doCreatePin")
	}

	if res.StatusCode != 201 {
		errorResponse, err := handleWrongStatuscode(res)
		if err != nil {
			return err
		}
		return errors.New(fmt.Sprintf("statuscode not 201 while doCreatePin. ErrorCode: %d ErrorMessage: %s", errorResponse.Code, errorResponse.Message))
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
		return listBoardResponseBody, errors.New(fmt.Sprintf("statuscode not 201 while doListBoards. ErrorCode: %d ErrorMessage: %s", errorResponse.Code, errorResponse.Message))
	}

	defer res.Body.Close()

	bytes, err := ioutil.ReadAll(res.Body)
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

	// Print the entire HTTP request
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		return errors.New("unable to dump request while doCreateBoard")
	}
	fmt.Println("Request:")
	fmt.Println(string(requestDump))

	res, err := c.httpClient.Do(req)
	if err != nil {
		return errors.New("unable to send request while doCreateBoard")
	}

	if res.StatusCode != 201 {
		errorResponse, err := handleWrongStatuscode(res)
		if err != nil {
			return err
		}
		return errors.New(fmt.Sprintf("statuscode not 201 while doCreateBoard. ErrorCode: %d ErrorMessage: %s", errorResponse.Code, errorResponse.Message))
	}

	return nil
}

func handleWrongStatuscode(res *http.Response) (errorResponse, error) {
	errorResponse := errorResponse{}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errorResponse, errors.New("unable to read response body while handleWrongStatuscode")
	}

	err = json.Unmarshal(bytes, &errorResponse)
	if err != nil {
		return errorResponse, errors.New("unable to unmarshal response body while handleWrongStatuscode")
	}

	return errorResponse, nil
}
