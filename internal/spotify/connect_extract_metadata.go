package spotify

func extractArtistNames(value any) []string {
	artists := []string{}
	walkMap(value, func(m map[string]any) {
		if list, ok := m["artists"].([]any); ok {
			for _, entry := range list {
				if name := findFirstName(entry); name != "" {
					artists = append(artists, name)
				}
			}
		}
	})
	if len(artists) == 0 {
		if m, ok := value.(map[string]any); ok {
			if name := getString(m, "artistName"); name != "" {
				artists = append(artists, name)
			}
		}
	}
	return dedupeStrings(artists)
}

func extractAlbumName(value any) string {
	var album string
	walkMap(value, func(m map[string]any) {
		if album != "" {
			return
		}
		if inner, ok := m["album"].(map[string]any); ok {
			if name := getString(inner, "name"); name != "" {
				album = name
			}
		}
		if inner, ok := m["albumOfTrack"].(map[string]any); ok {
			if name := getString(inner, "name"); name != "" {
				album = name
			}
		}
	})
	return album
}

func extractOwnerName(value any) string {
	var owner string
	walkMap(value, func(m map[string]any) {
		if owner != "" {
			return
		}
		if inner, ok := m["owner"].(map[string]any); ok {
			if name := getString(inner, "name"); name != "" {
				owner = name
			}
		}
		if inner, ok := m["user"].(map[string]any); ok {
			if name := getString(inner, "name"); name != "" {
				owner = name
			}
		}
	})
	return owner
}

func walkMap(value any, fn func(map[string]any)) {
	switch typed := value.(type) {
	case map[string]any:
		fn(typed)
		for _, child := range typed {
			walkMap(child, fn)
		}
	case []any:
		for _, child := range typed {
			walkMap(child, fn)
		}
	}
}
