package generator

import (
	"fmt"
)

const (
	Vowels     = "aeiou"
	Consonants = "bcdfghjklmnpqrstvwxyz"
)

// TokenGeneratorV1 - generates btcrecover-compliant tokens.
type TokenGeneratorV1 struct {
	ruleset *TokenRulesetV1
}

// NewTokenGeneratorV1 - returns new TokenGeneratorV1 instance.
func NewTokenGeneratorV1() *TokenGeneratorV1 {
	return &TokenGeneratorV1{
		ruleset: &TokenRulesetV1{
			BaseTokenRuleset: &BaseTokenRuleset{
				prefixes: []*BaseTokenRule{
					{charset: "o0", min: 0, max: 2},
				},
				suffixes: []*BaseTokenRule{
					{charset: "123", min: 0, max: 3},
					{charset: "69", min: 0, max: 2},
				},
				letterReplacements: map[byte]*BaseTokenRule{
					'a': {
						charset: "a4",
						min:     1,
						max:     3,
					},
					'e': {
						charset: "e3",
						min:     1,
						max:     3,
					},
					'o': {
						charset: "o0",
						min:     1,
						max:     3,
					},
				},
			},
		},
	}
}

// Generate - generates list of tokens for btcrecover to use. At least one token is created each time.
func (tg *TokenGeneratorV1) Generate(keyword string) string {
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

	return currentToken
}

func (tg *TokenGeneratorV1) getWildCardFromRule(rule *BaseTokenRule) string {
	rangePart := ""

	if rule.min == rule.max || rule.max == 0 {
		rangePart = fmt.Sprintf("%d", rule.min)
	} else {
		rangePart = fmt.Sprintf("%d,%d", rule.min, rule.max)
	}

	wildcard := fmt.Sprintf("%%%s[%s]", rangePart, rule.charset)

	return wildcard
}
