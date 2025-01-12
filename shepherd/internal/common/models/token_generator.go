package models

// TokenGeneratorVersionEnum - describes the version of TokenGenerator.
type TokenGeneratorVersionEnum int8

var TokenGeneratorVersion = struct {
	V1 TokenGeneratorVersionEnum
}{
	V1: 1,
}

var LatestTokenGeneratorVersion = TokenGeneratorVersion.V1
