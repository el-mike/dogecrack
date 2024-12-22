package generator

import (
	"fmt"
	"testing"
)

func TestTokenGenerator(t *testing.T) {
	generator := NewTokenGenerator(TokenRulesetOne)

	token := generator.Generate("test")
	fmt.Println(token)
}
