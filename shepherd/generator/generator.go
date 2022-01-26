package generator

// PasswordGenerator - entity responsible for handling passlist files generation.
type PasswordGenerator struct {
}

// NewPasswordGenerator - returns new PasswordGenerator instance.
func NewPasswordGenerator() *PasswordGenerator {
	return &PasswordGenerator{}
}

// Generate - creates new password file based on given basePassword and rules.
func (pg *PasswordGenerator) Generate(basePassword string, rules []string) (*GeneratorResult, error) {
	result := &GeneratorResult{
		BasePassword: basePassword,
		Rules:        rules,
		// Big file
		// fileUrl: "https://drive.google.com/file/d/1hztPGWlG4bfjLXjIG80-7si5M3PqoNzx/view?usp=sharing",
		// Small file, win:
		FileUrl: "https://drive.google.com/uc?id=12ULZTz8X5tIZ4243DI_-pHwe65PgQ33E",
		// Medium file, win:
		// FileUrl: "https://drive.google.com/uc?id=15Ao40uZK44whbS6BlM8RSo7LewIymKwq",
	}

	return result, nil
}
