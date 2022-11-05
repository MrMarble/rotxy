package utils

import (
	"bufio"
	"io"
)

// FilterBy removes elements from a slice that exist in another slice.
func FilterBy(base []string, filter []string) []string {
	newBase := []string{}
	for _, b := range base {
		if !contains(filter, b) {
			newBase = append(newBase, b)
		}
	}
	return newBase
}

// contains checks if a slice contains a string.
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// ReadLines reads a file and returns a slice of lines.
func ReadLines(r io.Reader) []string {
	scanner := bufio.NewScanner(r)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}
