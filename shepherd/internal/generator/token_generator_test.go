package generator

import (
	"fmt"
	"testing"
)

func TestTokenGenerator(t *testing.T) {
	generator := NewTokenGenerator(tokenRulesetOne)

	token := generator.Generate("test")
	fmt.Println(token)
}
