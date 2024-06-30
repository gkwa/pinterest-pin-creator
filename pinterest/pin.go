package pinterest

import (
	"time"
)

type Pin struct {
	ID              string      `json:"id"`
	CreatedAt       CustomTime  `json:"created_at"`
	Link            string      `json:"link"`
	Title           string      `json:"title"`
	Description     string      `json:"description"`
	DominantColor   string      `json:"dominant_color"`
	AltText         string      `json:"alt_text"`
	CreativeType    string      `json:"creative_type"`
	BoardID         string      `json:"board_id"`
	BoardSectionID  *string     `json:"board_section_id"`
	BoardOwner      BoardOwner  `json:"board_owner"`
	IsOwner         bool        `json:"is_owner"`
	Media           Media       `json:"media"`
	ParentPinID     *string     `json:"parent_pin_id"`
	IsStandard      bool        `json:"is_standard"`
	HasBeenPromoted bool        `json:"has_been_promoted"`
	Note            string      `json:"note"`
	PinMetrics      *PinMetrics `json:"pin_metrics"`
	ProductTags     []string    `json:"product_tags"`
}

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	t, err := time.Parse(`"2006-01-02T15:04:05"`, s)
	if err != nil {
		return err
	}
	*ct = CustomTime{t}
	return nil
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(ct.Time.Format(`"2006-01-02T15:04:05"`)), nil
}

type BoardOwner struct {
	Username string `json:"username"`
}

type Media struct {
	MediaType string           `json:"media_type"`
	Images    map[string]Image `json:"images"`
}

type Image struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
}

type PinMetrics struct {
	PinMetrics []PinMetricsData `json:"pin_metrics"`
}

type PinMetricsData struct {
	NinetyD MetricsData `json:"90d"`
	AllTime MetricsData `json:"all_time"`
}

type MetricsData struct {
	PinClick     int `json:"pin_click"`
	Impression   int `json:"impression"`
	Clickthrough int `json:"clickthrough"`
	Reaction     int `json:"reaction,omitempty"`
	Comment      int `json:"comment,omitempty"`
}
