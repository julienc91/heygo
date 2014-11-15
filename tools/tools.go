// Some useful functions
package tools

import (
	"os"
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
