package generator

import (
	"fmt"
	"testing"
)

func TestTokenGeneratorV1(t *testing.T) {
	generator := NewTokenGeneratorV1()

	token := generator.Generate("test")
	fmt.Println(token)
}
