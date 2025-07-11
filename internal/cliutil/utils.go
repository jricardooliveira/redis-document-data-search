package cliutil

import (
	"fmt"
	"strings"
)

func ParseInt(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}

func BuildRediSearchQuery(identifiers map[string]string) string {
	var parts []string
	for k, v := range identifiers {
		parts = append(parts, fmt.Sprintf("@%s:%s", k, v))
	}
	return strings.Join(parts, " ")
}
