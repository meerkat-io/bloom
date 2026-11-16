package strcase

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// SnakeCase : convert argument to snake_case style string.
func SnakeCase(src string) string {
	fields := splitToLowerFields(src)
	return strings.Join(fields, "_")
}

// ChainCase : convert argument to chain-case style string.
func ChainCase(src string) string {
	fields := splitToLowerFields(src)
	return strings.Join(fields, "-")
}

// CamelCase : convert argument to camelCase style string
func CamelCase(src string) string {
	fields := splitToLowerFields(src)
	for i, f := range fields {
		if i != 0 {
			fields[i] = capitalize(f)
		}
	}
	return strings.Join(fields, "")
}

// PascalCase : convert argument to PascalCase style string
func PascalCase(src string) string {
	fields := splitToLowerFields(src)
	for i, f := range fields {
		fields[i] = capitalize(f)
	}
	return strings.Join(fields, "")
}

func splitToLowerFields(src string) []string {
	var fields []string
	for _, word := range strings.FieldsFunc(src, isSeparator) {
		for _, sc := range splitCamelCase(word) {
			fields = append(fields, strings.ToLower(sc))
		}
	}
	return fields
}

func isSeparator(r rune) bool {
	return r == '_' || r == '-' || unicode.IsSpace(r)
}

// capitalize upper-cases the first rune (replacement for deprecated strings.Title).
func capitalize(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	if size == 0 {
		return s
	}
	return string(unicode.ToUpper(r)) + s[size:]
}

func splitCamelCase(src string) []string {
	if !utf8.ValidString(src) {
		return []string{src}
	}
	var entries []string
	var runes [][]rune
	lastClass := 0
	class := 0
	for _, r := range src {
		switch {
		case unicode.IsLower(r):
			class = 1
		case unicode.IsUpper(r):
			class = 2
		case unicode.IsDigit(r):
			class = 3
		default:
			class = 4
		}
		if class == lastClass {
			runes[len(runes)-1] = append(runes[len(runes)-1], r)
		} else {
			runes = append(runes, []rune{r})
		}
		lastClass = class
	}
	// handle upper case -> lower case sequences, e.g.
	// "PDFL", "oader" -> "PDF", "Loader"
	for i := 0; i < len(runes)-1; i++ {
		if unicode.IsUpper(runes[i][0]) && unicode.IsLower(runes[i+1][0]) {
			runes[i+1] = append([]rune{runes[i][len(runes[i])-1]}, runes[i+1]...)
			runes[i] = runes[i][:len(runes[i])-1]
		}
	}
	// construct []string from results
	for _, s := range runes {
		if len(s) > 0 {
			entries = append(entries, string(s))
		}
	}
	return entries
}
