package dns

import "strings"

func NormalizeSearchSuffixes(values []string) []string {
	seen := map[string]bool{}
	var result []string
	for _, value := range values {
		trimmed := strings.Trim(strings.TrimSpace(value), ".")
		if trimmed == "" || !validSuffix(trimmed) {
			continue
		}
		key := strings.ToLower(trimmed)
		if seen[key] {
			continue
		}
		seen[key] = true
		result = append(result, trimmed)
	}
	return result
}

func validSuffix(value string) bool {
	for _, char := range value {
		if char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char >= '0' && char <= '9' || char == '-' || char == '.' {
			continue
		}
		return false
	}
	return true
}
