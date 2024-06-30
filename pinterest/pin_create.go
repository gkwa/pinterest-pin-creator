package pinterest

import (
   "bytes"
   "encoding/json"
   "errors"
   "fmt"
   "net/http"
)

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
   	return fmt.Errorf("statuscode not 201 while doCreatePin. ErrorCode: %d ErrorMessage: %s", errorResponse.Code, errorResponse.Message)
   }

   var response Pin

   err = json.NewDecoder(res.Body).Decode(&response)
   if err != nil {
   	return fmt.Errorf("unable to decode response body: %v", err)
   }

   prettyJSON, err := json.MarshalIndent(response, "", "    ")
   if err != nil {
   	return fmt.Errorf("unable to format JSON: %v", err)
   }

   fmt.Println(string(prettyJSON))
   fmt.Printf("status: %d\n", res.StatusCode)

   return nil
}

