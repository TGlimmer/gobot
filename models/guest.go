package models

// SentGuestMessage https://core.telegram.org/bots/api#sentguestmessage — Bot API 10.0
//
// Result returned by answerGuestQuery: the bare identifiers needed to refer
// back to a message sent in response to a guest query.
type SentGuestMessage struct {
	GuestQueryID string `json:"guest_query_id"`
	MessageID    int    `json:"message_id"`
}
