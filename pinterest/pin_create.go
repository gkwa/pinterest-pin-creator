package pinterest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"pin-creator/internal/logger"
)

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
		return fmt.Errorf("statuscode not 201 while doCreatePin. ErrorCode: %d ErrorMessage: %s", errorResponse.Code, errorResponse.Message)
	}

	var response interface{}

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return fmt.Errorf("unable to decode response body: %v", err)
	}

	prettyJSON, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		return fmt.Errorf("unable to format JSON: %v", err)
	}

	log.V(1).Info(string(prettyJSON))
	log.V(1).Info("status: %d", res.StatusCode)

	return nil
}
