package textcase

import "strings"

func ToSnakeCase(input string) string {
	var builder strings.Builder
	for i, r := range input {
		if i > 0 && r >= 'A' && r <= 'Z' {
			builder.WriteByte('_')
		}
		if r >= 'A' && r <= 'Z' {
			builder.WriteRune(r + 32)
			continue
		}
		builder.WriteRune(r)
	}
	return builder.String()
}
