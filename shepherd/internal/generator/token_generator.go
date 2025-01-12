package generator

// TokenGenerator - a token generator interface.
type TokenGenerator interface {
	// Generate - generates list of tokens for btcrecover to use. At least one token is created each time.
	Generate(keyword string) []string
}
