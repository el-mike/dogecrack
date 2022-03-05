package generator

// PasswordGenerator - entity responsible for handling passlist files generation.
type PasswordGenerator struct {
}

// NewPasswordGenerator - returns new PasswordGenerator instance.
func NewPasswordGenerator() *PasswordGenerator {
	return &PasswordGenerator{}
}

// Generate - creates new password file based on given basePassword and rules.
func (pg *PasswordGenerator) Generate(keyword string, rules []string) (*GeneratorResult, error) {
	result := &GeneratorResult{
		Keyword: keyword,
		Rules:   rules,
		// Big file
		// PasslistUrl: "https://drive.google.com/uc?id=1hztPGWlG4bfjLXjIG80-7si5M3PqoNzx",
		// Small file, win:
		// PasslistUrl: "https://drive.google.com/uc?id=12ULZTz8X5tIZ4243DI_-pHwe65PgQ33E",
		// Medium file, win:
		// PasslistUrl: "https://drive.google.com/uc?id=15Ao40uZK44whbS6BlM8RSo7LewIymKwq",
		// Medium file, loose:
		PasslistUrl: "https://drive.google.com/uc?id=1g8TfnlFYBW77Dh2lXPr7iqPFaCbwYxj6",
	}

	return result, nil
}
