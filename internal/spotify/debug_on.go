//go:build debug

package spotify

import (
	"encoding/json"
	"fmt"
	"os"
)

func debugDumpQueue(raw map[string]any) {
	if data, err := json.MarshalIndent(raw, "", "  "); err == nil {
		_ = os.WriteFile("queue_debug.json", data, 0o644)
	}
}

func debugDumpStatus(raw map[string]any) {
	if data, err := json.MarshalIndent(raw, "", "  "); err == nil {
		_ = os.WriteFile("status_debug.json", data, 0o644)
	}
}

func debugDumpStatusRaw(data []byte) {
	var raw any
	if err := json.Unmarshal(data, &raw); err == nil {
		if pretty, err := json.MarshalIndent(raw, "", "  "); err == nil {
			_ = os.WriteFile("status_debug.json", pretty, 0o644)
			return
		}
	}
	_ = os.WriteFile("status_debug.json", data, 0o644)
}

func debugDumpTracksRaw(data []byte) {
	_ = os.WriteFile("tracks_debug.json", data, 0o644)
}

func debugEnrichError(err error) {
	if err == nil {
		_ = os.WriteFile("enrich_debug.txt", []byte("ok\n"), 0o644)
		return
	}
	_ = os.WriteFile("enrich_debug.txt", []byte(fmt.Sprintf("error: %v\n", err)), 0o644)
}
