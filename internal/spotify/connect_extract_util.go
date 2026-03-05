package spotify

import (
	"strconv"
	"strings"
)

func getMap(value any, path ...string) (map[string]any, bool) {
	current := value
	for _, key := range path {
		m, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}
		next, ok := m[key]
		if !ok {
			return nil, false
		}
		current = next
	}
	m, ok := current.(map[string]any)
	return m, ok
}

func getString(m map[string]any, key string) string {
	if m == nil {
		return ""
	}
	if value, ok := m[key].(string); ok {
		return value
	}
	return ""
}

func getInt(m map[string]any, key string) int {
	if m == nil {
		return 0
	}
	switch value := m[key].(type) {
	case int:
		return value
	case int64:
		return int(value)
	case float64:
		return int(value)
	case string:
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return 0
}

func getInt64(m map[string]any, key string) int64 {
	if m == nil {
		return 0
	}
	switch value := m[key].(type) {
	case int:
		return int64(value)
	case int64:
		return value
	case float64:
		return int64(value)
	case string:
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
			return parsed
		}
	}
	return 0
}

func getBool(m map[string]any, key string) bool {
	if m == nil {
		return false
	}
	if value, ok := m[key].(bool); ok {
		return value
	}
	return false
}

func dedupeStrings(values []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}
