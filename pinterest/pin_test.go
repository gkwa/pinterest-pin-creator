package pinterest

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPinUnmarshal(t *testing.T) {
	jsonData := `{
   	"alt_text": "WATCH IT NOW!",
   	"board_id": "1055883143825836483",
   	"board_owner": {
   		"username": "taylormonacelligkwa"
   	},
   	"board_section_id": null,
   	"created_at": "2024-06-30T04:54:04",
   	"creative_type": "REGULAR",
   	"description": "WATCH IT NOW!",
   	"dominant_color": "#f6f6f0",
   	"has_been_promoted": false,
   	"id": "1055883075131041293",
   	"is_owner": true,
   	"is_standard": true,
   	"link": "https://www.youtube.com/watch?v=e2fFMAPzZs4",
   	"media": {
   		"images": {
   			"1200x": {
   				"height": 411,
   				"url": "https://i.pinimg.com/1200x/91/01/76/9101764ceea582fc46316e1280f94484.jpg",
   				"width": 1200
   			},
   			"150x150": {
   				"height": 150,
   				"url": "https://i.pinimg.com/150x150/91/01/76/9101764ceea582fc46316e1280f94484.jpg",
   				"width": 150
   			},
   			"400x300": {
   				"height": 300,
   				"url": "https://i.pinimg.com/400x300/91/01/76/9101764ceea582fc46316e1280f94484.jpg",
   				"width": 400
   			},
   			"600x": {
   				"height": 193,
   				"url": "https://i.pinimg.com/564x/91/01/76/9101764ceea582fc46316e1280f94484.jpg",
   				"width": 564
   			}
   		},
   		"media_type": "image"
   	},
   	"note": "",
   	"parent_pin_id": null,
   	"pin_metrics": null,
   	"product_tags": [],
   	"title": "Second Video"
   }`

	var pin Pin
	err := json.Unmarshal([]byte(jsonData), &pin)
	assert.NoError(t, err)

	assert.Equal(t, "WATCH IT NOW!", pin.AltText)
	assert.Equal(t, "1055883143825836483", pin.BoardID)
	assert.Equal(t, "taylormonacelligkwa", pin.BoardOwner.Username)
	assert.Nil(t, pin.BoardSectionID)
	assert.Equal(t, CustomTime{time.Date(2024, 6, 30, 4, 54, 4, 0, time.UTC)}, pin.CreatedAt)
	assert.Equal(t, "REGULAR", pin.CreativeType)
	assert.Equal(t, "WATCH IT NOW!", pin.Description)
	assert.Equal(t, "#f6f6f0", pin.DominantColor)
	assert.False(t, pin.HasBeenPromoted)
	assert.Equal(t, "1055883075131041293", pin.ID)
	assert.True(t, pin.IsOwner)
	assert.True(t, pin.IsStandard)
	assert.Equal(t, "https://www.youtube.com/watch?v=e2fFMAPzZs4", pin.Link)
	assert.Equal(t, "image", pin.Media.MediaType)
	assert.Equal(t, 4, len(pin.Media.Images))
	assert.Equal(t, 1200, pin.Media.Images["1200x"].Width)
	assert.Equal(t, 411, pin.Media.Images["1200x"].Height)
	assert.Equal(t, "https://i.pinimg.com/1200x/91/01/76/9101764ceea582fc46316e1280f94484.jpg", pin.Media.Images["1200x"].URL)
	assert.Empty(t, pin.Note)
	assert.Nil(t, pin.ParentPinID)
	assert.Nil(t, pin.PinMetrics)
	assert.Empty(t, pin.ProductTags)
	assert.Equal(t, "Second Video", pin.Title)
}
