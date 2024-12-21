package generator

// GeneratorResult - result of generating a passlist file for given combination of
// base password and rules.
type GeneratorResult struct {
	Keyword     string
	Rules       []string
	PasslistUrl string
}

// TokenRule - describe a single aspect of token creation.
type TokenRule struct {
	charset string
	min     uint8
	max     uint8
	// If true, entire charset will be used as a value (resulting in new token entry in tokenlist).
	useFullCharset bool
}

// TokenRuleset - a versioned description of how token should be created from a base keyword.
type TokenRuleset struct {
	version            uint8
	includeUpperCase   bool
	prefixes           []*TokenRule
	suffixes           []*TokenRule
	letterReplacements map[byte]*TokenRule
}
