package handlers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// PrettyJSON returns indented JSON as HTML <pre> for debugging (optional, can be replaced with normal JSON)
func PrettyJSON(c *fiber.Ctx, v interface{}) error {
	pretty, _ := json.MarshalIndent(v, "", "  ")
	return c.Type("json").Send(pretty)
}

func escapeRediSearchValue(val string) string {
	// Escape double quotes
	val = strings.ReplaceAll(val, `"`, `\"`)
	// Escape other RediSearch special characters
	specialChars := []string{"-", "[", "]", "{", "}", "(", ")", "<", ">", ":", "~", "*", "?", "|", "&", "'", "!", "@", "#", "$", "%", "^", "="}
	for _, ch := range specialChars {
		val = strings.ReplaceAll(val, ch, `\`+ch)
	}
	return val
}

// BuildRediSearchQuery constructs a RediSearch query string from a map of identifiers
func BuildRediSearchQuery(identifiers map[string]string) string {
	var parts []string
	for k, v := range identifiers {
		escaped := escapeRediSearchValue(v)
		parts = append(parts, fmt.Sprintf("@%s:\"%s\"", k, escaped))
	}
	if len(parts) == 0 {
		return "*"
	}
	return strings.Join(parts, " ")
}
