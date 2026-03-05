//go:build !debug

package spotify

func debugDumpQueue(_ map[string]any) {}
func debugDumpTracksRaw(_ []byte)     {}
func debugEnrichError(_ error)        {}
