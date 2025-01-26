package models

type CheckedKeyword struct {
	Keyword           string                      `bson:"keyword" json:"keyword"`
	GeneratorVersions []TokenGeneratorVersionEnum `bson:"generatorVersions" json:"generatorVersions"`
	Tokenlists        []string                    `bson:"tokenlists" json:"tokenlists"`
	RunsCount         int                         `bson:"runsCount" json:"runsCount"`
}

type CheckedPasslist struct {
	PasslistUrl string `bson:"passlistUrl" json:"passlistUrl"`
	Name        string `bson:"name" json:"name"`
}

type CheckedIdeas struct {
	CheckedKeywords  []*CheckedKeyword  `json:"checkedKeywords"`
	CheckedPasslists []*CheckedPasslist `json:"checkedPasslists"`
}
