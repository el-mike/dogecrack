package models

// CrackPayload - describes a possible payload User can use to start
// cracking process.
type CrackPayload struct {
	Name        string   `json:"name"`
	Keyword     string   `json:"keyword"`
	Tokens      []string `json:"tokens"`
	PasslistUrl string   `json:"passlistUrl"`
}
