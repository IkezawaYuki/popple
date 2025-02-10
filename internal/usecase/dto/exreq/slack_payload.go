package exreq

type SlackPayload struct {
	IconEmoji string `json:"icon_emoji"`
	Text      string `json:"text"`
	Username  string `json:"username"`
}
