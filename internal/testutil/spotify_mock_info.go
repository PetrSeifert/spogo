package testutil

import (
	"context"

	"github.com/steipete/spogo/internal/spotify"
)

func (m *SpotifyMock) Search(ctx context.Context, kind, query string, limit, offset int) (spotify.SearchResult, error) {
	if m.SearchFn == nil {
		return spotify.SearchResult{}, ErrNotImplemented
	}
	return m.SearchFn(ctx, kind, query, limit, offset)
}

func (m *SpotifyMock) GetTrack(ctx context.Context, id string) (spotify.Item, error) {
	if m.GetTrackFn == nil {
		return spotify.Item{}, ErrNotImplemented
	}
	return m.GetTrackFn(ctx, id)
}

func (m *SpotifyMock) GetAlbum(ctx context.Context, id string) (spotify.Item, error) {
	if m.GetAlbumFn == nil {
		return spotify.Item{}, ErrNotImplemented
	}
	return m.GetAlbumFn(ctx, id)
}

func (m *SpotifyMock) GetArtist(ctx context.Context, id string) (spotify.Item, error) {
	if m.GetArtistFn == nil {
		return spotify.Item{}, ErrNotImplemented
	}
	return m.GetArtistFn(ctx, id)
}

func (m *SpotifyMock) GetPlaylist(ctx context.Context, id string) (spotify.Item, error) {
	if m.GetPlaylistFn == nil {
		return spotify.Item{}, ErrNotImplemented
	}
	return m.GetPlaylistFn(ctx, id)
}

func (m *SpotifyMock) GetShow(ctx context.Context, id string) (spotify.Item, error) {
	if m.GetShowFn == nil {
		return spotify.Item{}, ErrNotImplemented
	}
	return m.GetShowFn(ctx, id)
}

func (m *SpotifyMock) GetEpisode(ctx context.Context, id string) (spotify.Item, error) {
	if m.GetEpisodeFn == nil {
		return spotify.Item{}, ErrNotImplemented
	}
	return m.GetEpisodeFn(ctx, id)
}

func (m *SpotifyMock) ArtistTopTracks(ctx context.Context, id string, limit int) ([]spotify.Item, error) {
	if m.ArtistTopTracksFn == nil {
		return nil, ErrNotImplemented
	}
	return m.ArtistTopTracksFn(ctx, id, limit)
}
