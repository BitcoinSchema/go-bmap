package test

import (
	"os"
	"strings"
)

// GetTestHex gets the hex of the given test file
func GetTestHex(fileName string) string {
	fileData, err := os.ReadFile(fileName) //nolint:gosec // only used in testing
	if err != nil {
		return ""
	}

	return strings.Trim(string(fileData), "\n")
}
