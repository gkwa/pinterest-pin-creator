package pinterest

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type PinData struct {
	BoardId     string
	ImgPath     string
	Link        string
	Title       string
	Description string
	AltText     string
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

	return c.doCreatePin(createPinRequestBody)
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

func toBase64(imgPath string) string {
	bytes, err := ioutil.ReadFile(imgPath)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(bytes)
}

