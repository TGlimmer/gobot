package models

// PollAnswer https://core.telegram.org/bots/api#pollanswer
type PollAnswer struct {
	PollID              string   `json:"poll_id"`
	VoterChat           *Chat    `json:"voter_chat,omitempty"`
	User                *User    `json:"user"`
	OptionIDs           []int    `json:"option_ids,omitempty"`
	OptionPersistentIDs []string `json:"option_persistent_ids,omitempty"` // Bot API 9.6
}

// InputPollOption https://core.telegram.org/bots/api#inputpolloption
type InputPollOption struct {
	Text          string                `json:"text"`
	TextParseMode ParseMode             `json:"text_parse_mode,omitempty"`
	TextEntities  []MessageEntity       `json:"text_entities,omitempty"`
	Media         *InputPollOptionMedia `json:"media,omitempty"` // Bot API 10.0
}

// PollOption https://core.telegram.org/bots/api#polloption
type PollOption struct {
	Text         string          `json:"text"`
	TextEntities []MessageEntity `json:"text_entities,omitempty"`
	VoterCount   int             `json:"voter_count"`
	PersistentID string          `json:"persistent_id,omitempty"` // Bot API 9.6
	AddedByUser  *User           `json:"added_by_user,omitempty"`  // Bot API 9.6
	AddedByChat  *Chat           `json:"added_by_chat,omitempty"`  // Bot API 9.6
	AdditionDate int             `json:"addition_date,omitempty"`  // Bot API 9.6
	Media        *PollMedia      `json:"media,omitempty"`          // Bot API 10.0
}

// Poll https://core.telegram.org/bots/api#poll
type Poll struct {
	ID                     string          `json:"id"`
	Question               string          `json:"question"`
	QuestionEntities       []MessageEntity `json:"question_entities,omitempty"`
	Options                []PollOption    `json:"options"`
	TotalVoterCount        int             `json:"total_voter_count"`
	IsClosed               bool            `json:"is_closed"`
	IsAnonymous            bool            `json:"is_anonymous"`
	Type                   string          `json:"type"`
	AllowsMultipleAnswers  bool            `json:"allows_multiple_answers"`
	AllowsRevoting         bool            `json:"allows_revoting,omitempty"`           // Bot API 9.6
	ShuffleOptions         bool            `json:"shuffle_options,omitempty"`           // Bot API 9.6
	AllowAddingOptions     bool            `json:"allow_adding_options,omitempty"`      // Bot API 9.6
	HideResultsUntilCloses bool            `json:"hide_results_until_closes,omitempty"` // Bot API 9.6
	CorrectOptionID        int             `json:"correct_option_id,omitempty"`         // Deprecated by 9.6, kept for backwards compatibility
	CorrectOptionIDs       []int           `json:"correct_option_ids,omitempty"`        // Bot API 9.6 (replaces correct_option_id; quizzes may have multiple)
	Description            string          `json:"description,omitempty"`               // Bot API 9.6
	DescriptionEntities    []MessageEntity `json:"description_entities,omitempty"`      // Bot API 9.6
	Explanation            string          `json:"explanation,omitempty"`
	ExplanationEntities    []MessageEntity `json:"explanation_entities,omitempty"`
	ExplanationMedia       *PollMedia      `json:"explanation_media,omitempty"` // Bot API 10.0
	Media                  *PollMedia      `json:"media,omitempty"`             // Bot API 10.0
	MembersOnly            bool            `json:"members_only,omitempty"`      // Bot API 10.0
	CountryCodes           []string        `json:"country_codes,omitempty"`     // Bot API 9.6 / 10.0
	OpenPeriod             int             `json:"open_period,omitempty"`
	CloseDate              int             `json:"close_date,omitempty"`
}

// PollOptionAdded https://core.telegram.org/bots/api#polloptionadded — Bot API 9.6
type PollOptionAdded struct {
	PollID string     `json:"poll_id"`
	Option PollOption `json:"option"`
}

// PollOptionDeleted https://core.telegram.org/bots/api#polloptiondeleted — Bot API 9.6
type PollOptionDeleted struct {
	PollID             string `json:"poll_id"`
	OptionPersistentID string `json:"option_persistent_id"`
}
