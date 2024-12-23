package models

// CrackPayload - describes a possible payload User can use to start
// cracking process.
type CrackPayload struct {
	Name        string `json:"name"`
	Keyword     string `json:"keyword"`
	PasslistUrl string `json:"passlistUrl"`
}
