package generator

// BaseTokenRule - describe a single aspect of token creation.
// This base type should be extended by version-specific Generator rules.
type BaseTokenRule struct {
	charset string
	min     uint8
	max     uint8
	// If true, entire charset will be used as a value (resulting in new token entry in tokenlist).
	useFullCharset bool
}

// BaseTokenRuleset - a versioned description of how token should be created from a base keyword.
// This base type should be extended by version-specific Generator ruleset.
type BaseTokenRuleset struct {
	prefixes           []*BaseTokenRule
	suffixes           []*BaseTokenRule
	letterReplacements map[byte]*BaseTokenRule
}

type TokenRuleV1 struct {
	*BaseTokenRule
}

type TokenRulesetV1 struct {
	*BaseTokenRuleset
}
