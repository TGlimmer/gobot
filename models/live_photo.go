package models

import (
	"encoding/json"
	"io"
)

// LivePhoto https://core.telegram.org/bots/api#livephoto — Bot API 10.0
type LivePhoto struct {
	Photo        []PhotoSize `json:"photo,omitempty"`
	FileID       string      `json:"file_id"`
	FileUniqueID string      `json:"file_unique_id"`
	Width        int         `json:"width"`
	Height       int         `json:"height"`
	Duration     int         `json:"duration"`
	MimeType     string      `json:"mime_type,omitempty"`
	FileSize     int64       `json:"file_size,omitempty"`
}

// InputMediaLivePhoto https://core.telegram.org/bots/api#inputmedialivephoto — Bot API 10.0
type InputMediaLivePhoto struct {
	Media                 string          `json:"media"`
	Thumbnail             InputFile       `json:"thumbnail,omitempty"`
	Caption               string          `json:"caption,omitempty"`
	ParseMode             ParseMode       `json:"parse_mode,omitempty"`
	CaptionEntities       []MessageEntity `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia bool            `json:"show_caption_above_media,omitempty"`
	Width                 int             `json:"width,omitempty"`
	Height                int             `json:"height,omitempty"`
	Duration              int             `json:"duration,omitempty"`
	HasSpoiler            bool            `json:"has_spoiler,omitempty"`

	MediaAttachment io.Reader `json:"-"`
}

func (m *InputMediaLivePhoto) Attachment() io.Reader { return m.MediaAttachment }
func (m *InputMediaLivePhoto) GetMedia() string      { return m.Media }
func (m InputMediaLivePhoto) MarshalInputMedia() ([]byte, error) {
	return json.Marshal(&struct {
		Type string `json:"type"`
		InputMediaLivePhoto
	}{
		Type:                "live_photo",
		InputMediaLivePhoto: m,
	})
}
func (InputMediaLivePhoto) inputMediaTag() {}

// PaidMediaLivePhoto https://core.telegram.org/bots/api#paidmedialivephoto — Bot API 10.0
type PaidMediaLivePhoto struct {
	Type PaidMediaType

	LivePhoto LivePhoto `json:"live_photo"`
}

// InputPaidMediaLivePhoto https://core.telegram.org/bots/api#inputpaidmedialivephoto — Bot API 10.0
type InputPaidMediaLivePhoto struct {
	Media     string    `json:"media"`
	Thumbnail InputFile `json:"thumbnail,omitempty"`
	Width     int       `json:"width,omitempty"`
	Height    int       `json:"height,omitempty"`
	Duration  int       `json:"duration,omitempty"`

	MediaAttachment io.Reader `json:"-"`
}

func (m *InputPaidMediaLivePhoto) inputPaidMediaTag() {}
func (m *InputPaidMediaLivePhoto) Attachment() io.Reader {
	return m.MediaAttachment
}
func (m *InputPaidMediaLivePhoto) GetMedia() string { return m.Media }
func (m *InputPaidMediaLivePhoto) MarshalInputMedia() ([]byte, error) {
	ret := struct {
		Type string `json:"type"`
		*InputPaidMediaLivePhoto
	}{
		Type:                    "live_photo",
		InputPaidMediaLivePhoto: m,
	}
	return json.Marshal(&ret)
}
