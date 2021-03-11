package repository

type SlackMessage struct {
	IconEmoji string `json:"icon_emoji,omitempty"`
	Text      string `json:"text,omitempty"`
}
