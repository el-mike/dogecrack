package models

// KeywordSuggestionsPayload - a payload for retrieving keyword suggestions.
type KeywordSuggestionsPayload struct {
	TokenGeneratorVersion TokenGeneratorVersionEnum `json:"tokenGeneratorVersion"`
	PresetsOnly           bool                      `json:"presetsOnly"`
}
