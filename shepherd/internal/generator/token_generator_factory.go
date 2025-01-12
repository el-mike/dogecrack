package generator

import (
	"fmt"
	"github.com/el-mike/dogecrack/shepherd/internal/common/models"
)

// TokenGeneratorFactory - factory creating TokenGenerator instances.
type TokenGeneratorFactory struct{}

// NewTokenGeneratorFactory - returns new TokenGeneratorFactory.
func NewTokenGeneratorFactory() *TokenGeneratorFactory {
	return &TokenGeneratorFactory{}
}

func (tf *TokenGeneratorFactory) CreateGenerator(version models.TokenGeneratorVersionEnum) (TokenGenerator, error) {
	switch version {
	case models.TokenGeneratorVersion.V1:
		return NewTokenGeneratorV1(), nil
	default:
		return nil, fmt.Errorf("unsupported token generator version: %s", version)
	}
}
