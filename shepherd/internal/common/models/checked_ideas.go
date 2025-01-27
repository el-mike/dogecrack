package models

type UsedKeyword struct {
	Keyword           string                      `bson:"keyword" json:"keyword"`
	GeneratorVersions []TokenGeneratorVersionEnum `bson:"generatorVersions" json:"generatorVersions"`
	Tokenlists        []string                    `bson:"tokenlists" json:"tokenlists"`
	RunsCount         int                         `bson:"runsCount" json:"runsCount"`
}

type UsedPasslist struct {
	PasslistUrl string `bson:"passlistUrl" json:"passlistUrl"`
	Name        string `bson:"name" json:"name"`
}

type UsedIdeas struct {
	CheckedKeywords  []*UsedKeyword  `json:"checkedKeywords"`
	CheckedPasslists []*UsedPasslist `json:"checkedPasslists"`
}
