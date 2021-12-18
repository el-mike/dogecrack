package generator

// GeneratorResult - result of generating a passlist file for given combination of
// base password and rules.
type GeneratorResult struct {
	BasePassword string
	Rules        []string
	FileUrl      string
}
