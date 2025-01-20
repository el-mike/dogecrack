package pitbull_client

import "fmt"

// BuildRunCommand - a helper for building arguments for Pitbull's "run" command.
func BuildRunCommand(walletString, passlistUrl, tokenlist string, skipCount, minlength, maxLength int64) string {
	baseCommand := "pitbull run"

	// Single quotes are needed, as it may contain "?" characters that bash can interpret as a part of regex.
	baseCommand += fmt.Sprintf(" -w '%s'", walletString)

	if passlistUrl != "" {
		// Single quotes are needed, as it may contain "?" characters that bash can interpret as a part of regex.
		baseCommand += fmt.Sprintf(" -u '%s'", passlistUrl)
	}

	if tokenlist != "" {
		// Single quotes are needed, as it may contain whitespaces and characters like "?".
		baseCommand += fmt.Sprintf(" -t '%s'", tokenlist)
	}

	if skipCount > 0 {
		baseCommand += fmt.Sprintf(" --skip %d", skipCount)
	}

	if minlength > 0 {
		baseCommand += fmt.Sprintf(" --length-min %d", minlength)
	}

	if maxLength > 0 {
		baseCommand += fmt.Sprintf(" --length-max %d", maxLength)
	}

	return baseCommand
}
