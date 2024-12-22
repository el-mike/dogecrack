package generator

var TokenRulesetOne = &TokenRuleset{
	version:          1,
	includeUpperCase: false,
	prefixes: []*TokenRule{
		{charset: "o0", min: 0, max: 2},
	},
	suffixes: []*TokenRule{
		{charset: "123", min: 0, max: 3},
		{charset: "69", min: 0, max: 2},
	},
	letterReplacements: map[byte]*TokenRule{
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
}
