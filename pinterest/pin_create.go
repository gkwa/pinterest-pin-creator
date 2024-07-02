package pinterest

import (
	"context"
	"encoding/json"
	"fmt"

	"pin-creator/internal/logger"
)

type mediaSourceRequestBody struct {
	SourceType  string `json:"source_type"`
	ContentType string `json:"content_type"`
	Data        string `json:"data"`
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

func (c *Client) CreatePin(ctx context.Context, pinData PinData) error {
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

	return c.doCreatePin(ctx, createPinRequestBody)
}

func (c *Client) doCreatePin(ctx context.Context, body createPinRequestBody) error {
	log := logger.FromContext(ctx)
	url := fmt.Sprintf("%s%s", c.baseUrl, "pins")

	req, err := c.createRequest("POST", url, body)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	responseBody, err := c.executeRequest(ctx, req, 201)
	if err != nil {
		return fmt.Errorf("error executing request: %v", err)
	}

	var response interface{}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return fmt.Errorf("unable to decode response body: %v", err)
	}

	log.V(2).Info(fmt.Sprintf("Pin created successfully. Response: %s", string(responseBody)))
	return nil
}
