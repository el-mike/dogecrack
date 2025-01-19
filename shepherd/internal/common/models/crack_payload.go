package models

// CrackPayload - describes a possible payload User can use to start
// cracking process.
type CrackPayload struct {
	Name                  string                    `json:"name"`
	Keywords              []string                  `json:"keywords"`
	Tokenlist             string                    `json:"tokenlist"`
	PasslistUrl           string                    `json:"passlistUrl"`
	TokenGeneratorVersion TokenGeneratorVersionEnum `json:"tokenGeneratorVersion"`
}
