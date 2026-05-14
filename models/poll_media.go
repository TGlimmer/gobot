package models

import (
	"encoding/json"
	"io"
)

// PollMedia https://core.telegram.org/bots/api#pollmedia — Bot API 10.0
//
// Server-rendered media for polls/quizzes (poll descriptions, explanations,
// poll options). Exactly one of the optional fields is populated by the
// server.
type PollMedia struct {
	Animation *Animation `json:"animation,omitempty"`
	Audio     *Audio     `json:"audio,omitempty"`
	Document  *Document  `json:"document,omitempty"`
	LivePhoto *LivePhoto `json:"live_photo,omitempty"`
	Location  *Location  `json:"location,omitempty"`
	Photo     []PhotoSize `json:"photo,omitempty"`
	Sticker   *Sticker   `json:"sticker,omitempty"`
	Venue     *Venue     `json:"venue,omitempty"`
	Video     *Video     `json:"video,omitempty"`
}

// InputPollMedia https://core.telegram.org/bots/api#inputpollmedia — Bot API 10.0
//
// Used for poll description / quiz explanation. The concrete value must be
// one of the InputMedia* variants accepted in this context. Stored as the
// InputMedia interface so it integrates with the existing multipart upload
// pipeline.
type InputPollMedia struct {
	Media InputMedia
}

func (p InputPollMedia) MarshalJSON() ([]byte, error) {
	if p.Media == nil {
		return []byte("null"), nil
	}
	return p.Media.MarshalInputMedia()
}

// InputPollOptionMedia https://core.telegram.org/bots/api#inputpolloptionmedia — Bot API 10.0
//
// Media attached to an individual poll option (animation, live_photo,
// location, photo, sticker, venue, video).
type InputPollOptionMedia struct {
	Media InputMedia
}

func (p InputPollOptionMedia) MarshalJSON() ([]byte, error) {
	if p.Media == nil {
		return []byte("null"), nil
	}
	return p.Media.MarshalInputMedia()
}

// InputMediaSticker https://core.telegram.org/bots/api#inputmediasticker — Bot API 10.0
type InputMediaSticker struct {
	Media     string `json:"media"`
	EmojiList []string `json:"emoji_list,omitempty"`

	MediaAttachment io.Reader `json:"-"`
}

func (m *InputMediaSticker) Attachment() io.Reader { return m.MediaAttachment }
func (m *InputMediaSticker) GetMedia() string      { return m.Media }
func (m InputMediaSticker) MarshalInputMedia() ([]byte, error) {
	return json.Marshal(&struct {
		Type string `json:"type"`
		InputMediaSticker
	}{
		Type:              "sticker",
		InputMediaSticker: m,
	})
}
func (InputMediaSticker) inputMediaTag() {}

// InputMediaLocation https://core.telegram.org/bots/api#inputmedialocation — Bot API 10.0
type InputMediaLocation struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	HorizontalAccuracy   float64 `json:"horizontal_accuracy,omitempty"`
	LivePeriod           int     `json:"live_period,omitempty"`
	Heading              int     `json:"heading,omitempty"`
	ProximityAlertRadius int     `json:"proximity_alert_radius,omitempty"`
}

func (m *InputMediaLocation) Attachment() io.Reader { return nil }
func (m *InputMediaLocation) GetMedia() string      { return "" }
func (m InputMediaLocation) MarshalInputMedia() ([]byte, error) {
	return json.Marshal(&struct {
		Type string `json:"type"`
		InputMediaLocation
	}{
		Type:               "location",
		InputMediaLocation: m,
	})
}
func (InputMediaLocation) inputMediaTag() {}

// InputMediaVenue https://core.telegram.org/bots/api#inputmediavenue — Bot API 10.0
type InputMediaVenue struct {
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	Title           string  `json:"title"`
	Address         string  `json:"address"`
	FoursquareID    string  `json:"foursquare_id,omitempty"`
	FoursquareType  string  `json:"foursquare_type,omitempty"`
	GooglePlaceID   string  `json:"google_place_id,omitempty"`
	GooglePlaceType string  `json:"google_place_type,omitempty"`
}

func (m *InputMediaVenue) Attachment() io.Reader { return nil }
func (m *InputMediaVenue) GetMedia() string      { return "" }
func (m InputMediaVenue) MarshalInputMedia() ([]byte, error) {
	return json.Marshal(&struct {
		Type string `json:"type"`
		InputMediaVenue
	}{
		Type:            "venue",
		InputMediaVenue: m,
	})
}
func (InputMediaVenue) inputMediaTag() {}
