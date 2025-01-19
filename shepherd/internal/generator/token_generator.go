package generator

// TokenGenerator - a token generator interface.
type TokenGenerator interface {
	// Generate - generates tokenlist for btcrecover to use. At least one token is created each time.
	Generate(keyword string) string
}
