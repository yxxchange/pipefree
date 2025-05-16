package str

import "strings"

func Replace(template, holder, value string) string {
	return strings.Replace(template, holder, value, 1)
}
