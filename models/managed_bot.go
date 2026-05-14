package models

// KeyboardButtonRequestManagedBot https://core.telegram.org/bots/api#keyboardbuttonrequestmanagedbot — Bot API 9.6
type KeyboardButtonRequestManagedBot struct {
	RequestID int32 `json:"request_id"`
}

// PreparedKeyboardButton https://core.telegram.org/bots/api#preparedkeyboardbutton — Bot API 9.6
//
// A keyboard button whose state can be saved by the bot via
// savePreparedKeyboardButton and later attached to a freshly built keyboard.
type PreparedKeyboardButton struct {
	ID     string         `json:"id"`
	Button KeyboardButton `json:"button"`
}

// ManagedBotCreated https://core.telegram.org/bots/api#managedbotcreated — Bot API 9.6
type ManagedBotCreated struct {
	Bot User `json:"bot"`
}

// ManagedBotUpdated https://core.telegram.org/bots/api#managedbotupdated — Bot API 9.6
type ManagedBotUpdated struct {
	User User `json:"user"`
	Bot  User `json:"bot"`
}

// BotAccessSettings https://core.telegram.org/bots/api#botaccesssettings — Bot API 10.0
//
// Controls how a managed bot can be accessed by other users / chats.
// Fields are optional and only those set by the caller are sent.
type BotAccessSettings struct {
	CanBeAddedToGroups   *bool `json:"can_be_added_to_groups,omitempty"`
	CanBeAddedToChannels *bool `json:"can_be_added_to_channels,omitempty"`
	CanReadAllMessages   *bool `json:"can_read_all_messages,omitempty"`
	CanRespondInGroups   *bool `json:"can_respond_in_groups,omitempty"`
	CanReceiveStars      *bool `json:"can_receive_stars,omitempty"`
}
