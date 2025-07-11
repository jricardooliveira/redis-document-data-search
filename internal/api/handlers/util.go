package handlers

import (
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// PrettyJSON returns indented JSON as HTML <pre> for debugging (optional, can be replaced with normal JSON)
func PrettyJSON(c *fiber.Ctx, v interface{}) error {
	pretty, _ := json.MarshalIndent(v, "", "  ")
	return c.Type("html", "utf-8").SendString("<pre>" + string(pretty) + "</pre>")
}

// BuildRediSearchQuery constructs a RediSearch query string from a map of identifiers
func BuildRediSearchQuery(identifiers map[string]string) string {
	var parts []string
	for k, v := range identifiers {
		parts = append(parts, "@"+k+":"+v)
	}
	if len(parts) == 0 {
		return "*"
	}
	return strings.Join(parts, " ")
}
