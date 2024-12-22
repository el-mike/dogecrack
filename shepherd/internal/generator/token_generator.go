package generator

import (
	"fmt"
)

const (
	Vowels     = "aeiou"
	Consonants = "bcdfghjklmnpqrstvwxyz"
)

// TokenGenerator - generates btcrecover-compliant tokens.
type TokenGenerator struct {
	ruleset *TokenRuleset
}

// NewTokenGenerator - returns new TokenGenerator instance.
func NewTokenGenerator(ruleset *TokenRuleset) *TokenGenerator {
	return &TokenGenerator{
		ruleset: ruleset,
	}
}

// Generate - generates list of tokens for btcrecover to use. At least one token is created each time.
func (tg *TokenGenerator) Generate(keyword string) []string {
	currentToken := ""

	for _, prefix := range tg.ruleset.prefixes {
		wildcard := tg.getWildCardFromRule(prefix)
		currentToken += wildcard
	}

	for _, c := range keyword {
		if replacementRule, ok := tg.ruleset.letterReplacements[byte(c)]; ok {
			wildcard := tg.getWildCardFromRule(replacementRule)
			currentToken += wildcard
		} else {
			currentToken += string(c)
		}
	}

	for _, suffix := range tg.ruleset.suffixes {
		wildcard := tg.getWildCardFromRule(suffix)
		currentToken += wildcard
	}

	return []string{currentToken}
}

func (tg *TokenGenerator) getWildCardFromRule(rule *TokenRule) string {
	rangePart := ""

	if rule.min == rule.max || rule.max == 0 {
		rangePart = fmt.Sprintf("%d", rule.min)
	} else {
		rangePart = fmt.Sprintf("%d,%d", rule.min, rule.max)
	}

	wildcard := fmt.Sprintf("%%%s[%s]", rangePart, rule.charset)

	return wildcard
}
