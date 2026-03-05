package spotify

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

func (c *ConnectClient) playback(ctx context.Context) (PlaybackStatus, error) {
	return withConnectState(ctx, c, func(state connectState) (PlaybackStatus, error) {
		status := mapPlaybackStatus(state)
		if status.Item != nil && status.Item.Type == "track" && status.Item.ID != "" {
			if track, err := c.trackInfo(ctx, status.Item.ID); err == nil {
				status.Item.Artists = track.Artists
				if status.Item.Album == "" {
					status.Item.Album = track.Album
				}
				if status.Item.DurationMS == 0 {
					status.Item.DurationMS = track.DurationMS
				}
			}
		}
		return status, nil
	})
}

func (c *ConnectClient) devices(ctx context.Context) ([]Device, error) {
	return withConnectState(ctx, c, func(state connectState) ([]Device, error) {
		return mapDevices(state), nil
	})
}

func (c *ConnectClient) transfer(ctx context.Context, deviceID string) error {
	return withConnectStateErr(ctx, c, func(state connectState) error {
		fromID := connectTransferSourceID(state)
		if fromID == "" {
			return c.transferViaWebAPI(ctx, deviceID)
		}
		return c.sendConnectCommand(ctx, fmt.Sprintf("%s/connect/transfer/from/%s/to/%s", connectStateBase, fromID, deviceID), map[string]any{
			"transfer_options": map[string]any{
				"restore_paused": "resume",
			},
			"command_id": randomHex(32),
		})
	})
}

func (c *ConnectClient) transferViaWebAPI(ctx context.Context, deviceID string) error {
	return withWebFallback(c, func(web *Client) error {
		return web.Transfer(ctx, deviceID)
	})
}

func (c *ConnectClient) play(ctx context.Context, uri string) error {
	return withConnectStateErr(ctx, c, func(state connectState) error {
		if uri == "" {
			return c.sendPlayerCommand(ctx, state, "resume", nil)
		}
		return c.sendPlayerCommand(ctx, state, "play", playCommandPayload(uri))
	})
}

func (c *ConnectClient) pause(ctx context.Context) error {
	return c.sendStateCommand(ctx, "pause", nil)
}

func (c *ConnectClient) next(ctx context.Context) error {
	return c.sendStateCommand(ctx, "skip_next", nil)
}

func (c *ConnectClient) previous(ctx context.Context) error {
	return c.sendStateCommand(ctx, "skip_prev", nil)
}

func (c *ConnectClient) seek(ctx context.Context, positionMS int) error {
	if positionMS < 0 {
		positionMS = 0
	}
	return c.sendStateCommand(ctx, "seek_to", map[string]any{
		"command": map[string]any{
			"endpoint": "seek_to",
			"value":    positionMS,
			"logging_params": map[string]any{
				"command_id": randomHex(32),
			},
		},
	})
}

func (c *ConnectClient) volume(ctx context.Context, volume int) error {
	volume = clampVolume(volume)
	return withConnectStateErr(ctx, c, func(state connectState) error {
		fromID := connectTransferSourceID(state)
		if fromID == "" || state.activeDeviceID == "" {
			return errors.New("missing device id")
		}
		return c.sendConnectCommand(ctx, fmt.Sprintf("%s/connect/volume/from/%s/to/%s", connectStateBase, fromID, state.activeDeviceID), map[string]any{
			"volume": int(float64(volume) / 100 * 65535),
		})
	})
}

func (c *ConnectClient) shuffle(ctx context.Context, enabled bool) error {
	return c.sendStateCommand(ctx, "set_shuffling_context", map[string]any{
		"command": map[string]any{
			"endpoint": "set_shuffling_context",
			"value":    enabled,
			"logging_params": map[string]any{
				"command_id": randomHex(32),
			},
		},
	})
}

func (c *ConnectClient) repeat(ctx context.Context, mode string) error {
	command := map[string]any{
		"endpoint": "set_options",
		"logging_params": map[string]any{
			"command_id": randomHex(32),
		},
	}
	repeatingTrack, repeatingContext := repeatFlags(mode)
	command["repeating_track"] = repeatingTrack
	command["repeating_context"] = repeatingContext
	return c.sendStateCommand(ctx, "set_options", map[string]any{"command": command})
}

func (c *ConnectClient) queueAdd(ctx context.Context, uri string) error {
	return c.sendStateCommand(ctx, "add_to_queue", map[string]any{
		"command": map[string]any{
			"endpoint": "add_to_queue",
			"track": map[string]any{
				"uri": uri,
			},
			"logging_params": map[string]any{
				"command_id": randomHex(32),
			},
		},
	})
}

func (c *ConnectClient) queue(ctx context.Context, limit int) (Queue, error) {
	return withConnectState(ctx, c, func(state connectState) (Queue, error) {
		debugDumpQueue(state.raw)
		queue := mapQueue(state)
		if limit > 0 && len(queue.Queue) > limit {
			queue.Queue = queue.Queue[:limit]
		}
		debugEnrichError(c.enrichQueueNames(ctx, queue.Queue))
		return queue, nil
	})
}

func (c *ConnectClient) enrichQueueNames(ctx context.Context, items []Item) error {
	type slot struct {
		idx int
		id  string
	}
	slots := make([]slot, 0, len(items))
	for i, item := range items {
		if item.Name == "" && item.Type == "track" && item.ID != "" {
			slots = append(slots, slot{idx: i, id: item.ID})
		}
	}
	if len(slots) == 0 {
		return nil
	}

	type result struct {
		idx  int
		item Item
		err  error
	}
	results := make(chan result, len(slots))

	const maxConcurrent = 5
	sem := make(chan struct{}, maxConcurrent)

	for _, current := range slots {
		current := current
		sem <- struct{}{}
		go func() {
			defer func() { <-sem }()
			track, err := c.trackInfo(ctx, current.id)
			results <- result{idx: current.idx, item: track, err: err}
		}()
	}

	var firstErr error
	for range slots {
		result := <-results
		if result.err != nil {
			if firstErr == nil {
				firstErr = result.err
			}
			continue
		}
		items[result.idx].Name = result.item.Name
		items[result.idx].Artists = result.item.Artists
		items[result.idx].Album = result.item.Album
		items[result.idx].DurationMS = result.item.DurationMS
	}
	return firstErr
}

func (c *ConnectClient) sendStateCommand(ctx context.Context, endpoint string, payload map[string]any) error {
	return withConnectStateErr(ctx, c, func(state connectState) error {
		return c.sendPlayerCommand(ctx, state, endpoint, payload)
	})
}

func withConnectState[T any](ctx context.Context, c *ConnectClient, fn func(connectState) (T, error)) (T, error) {
	state, err := c.connectState(ctx)
	if err != nil {
		var zero T
		return zero, err
	}
	return fn(state)
}

func withConnectStateErr(ctx context.Context, c *ConnectClient, fn func(connectState) error) error {
	_, err := withConnectState(ctx, c, func(state connectState) (struct{}, error) {
		return struct{}{}, fn(state)
	})
	return err
}

func connectTransferSourceID(state connectState) string {
	fromID := state.originDeviceID
	if fromID == "" {
		fromID = state.activeDeviceID
	}
	return fromID
}

func playCommandPayload(uri string) map[string]any {
	command := map[string]any{
		"endpoint": "play",
		"logging_params": map[string]any{
			"command_id": randomHex(32),
		},
	}
	command["context"] = map[string]any{"uri": uri, "url": "context://" + uri}
	if !isContextURI(uri) {
		command["options"] = map[string]any{
			"skip_to": map[string]any{"track_uri": uri},
		}
	}
	return map[string]any{"command": command}
}

func clampVolume(volume int) int {
	if volume < 0 {
		return 0
	}
	if volume > 100 {
		return 100
	}
	return volume
}

func repeatFlags(mode string) (bool, bool) {
	switch strings.ToLower(mode) {
	case "track":
		return true, false
	case "context":
		return false, true
	default:
		return false, false
	}
}
