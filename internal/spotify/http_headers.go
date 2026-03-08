package spotify

import "net/http"

const defaultSpotifyAppPlatform = "WebPlayer"

type requestHeaders struct {
	AccessToken   string
	ClientToken   string
	ClientVersion string
	Accept        string
	ContentType   string
	Language      string
	AppPlatform   string
	ConnectionID  string
}

func applyRequestHeaders(req *http.Request, headers requestHeaders) {
	if req == nil {
		return
	}
	req.Header.Set("User-Agent", defaultUserAgent())
	if headers.AccessToken != "" {
		req.Header.Set("Authorization", "Bearer "+headers.AccessToken)
	}
	if headers.ClientToken != "" {
		req.Header.Set("Client-Token", headers.ClientToken)
	}
	if headers.ClientVersion != "" {
		req.Header.Set("Spotify-App-Version", headers.ClientVersion)
	}
	if headers.Accept != "" {
		req.Header.Set("Accept", headers.Accept)
	}
	if headers.ContentType != "" {
		req.Header.Set("Content-Type", headers.ContentType)
	}
	if headers.Language != "" {
		req.Header.Set("Accept-Language", headers.Language)
	}
	if headers.AppPlatform != "" {
		req.Header.Set("app-platform", headers.AppPlatform)
	}
	if headers.ConnectionID != "" {
		req.Header.Set("x-spotify-connection-id", headers.ConnectionID)
	}
}
