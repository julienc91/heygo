// Some useful functions
package tools

import (
	"os"
	"unicode"
)

// Check if e is in a
func InArray(a []string, e string) bool {
	for _, x := range a {
		if x == e {
			return true
		}
	}
	return false
}

// Check if a file exists
func CheckFilePath(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func UnderscoreToCamelCase(s string) string {
	var res []rune
	var upperNextOne bool = false

	for i, c := range s {
		if i == 0 {
			res = append(res, unicode.ToUpper(c))
		} else if c == '_' {
			upperNextOne = true
		} else if upperNextOne {
			res = append(res, unicode.ToUpper(c))
			upperNextOne = false
		} else {
			res = append(res, c)
		}
	}
	return string(res)
}
