//go:build !debug

package spotify

func debugDumpQueue(_ map[string]any)  {}
func debugDumpStatus(_ map[string]any) {}
func debugDumpStatusRaw(_ []byte)      {}
func debugDumpTracksRaw(_ []byte)      {}
func debugEnrichError(_ error)         {}
