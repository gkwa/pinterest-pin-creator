package pinterest

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
