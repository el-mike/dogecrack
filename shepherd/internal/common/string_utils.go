package common

func RemoveStringDuplicates(items []string) []string {
	keys := make(map[string]bool)
	result := []string{}

	for _, item := range items {
		if _, value := keys[item]; !value {
			keys[item] = true
			result = append(result, item)
		}
	}

	return result
}
