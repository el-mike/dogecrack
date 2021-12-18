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
		// Small file
		FileUrl: "https://drive.google.com/file/d/1a9EHdCaTWq_btpJw9DkQLgrUHg74233S/view?usp=sharing",
	}

	return result, nil
}
