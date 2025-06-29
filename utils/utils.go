package utils

func Deduplicate(list []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0)

	for _, item := range list {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

func StringSliceToMap(list []string) map[string]struct{} {
	result := make(map[string]struct{})
	for _, item := range list {
		result[item] = struct{}{}
	}
	return result
}
