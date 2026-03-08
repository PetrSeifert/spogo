package cli

import (
	"context"
	"errors"
	"fmt"

	"github.com/steipete/spogo/internal/app"
	"github.com/steipete/spogo/internal/spotify"
)

type itemLookup func(context.Context, spotify.API, string) (spotify.Item, error)

func spotifyClient(ctx *app.Context) (spotify.API, context.Context, error) {
	if ctx == nil {
		return nil, context.Background(), errors.New("nil context")
	}
	client, err := ctx.Spotify()
	if err != nil {
		return nil, ctx.CommandContext(), err
	}
	return client, ctx.CommandContext(), nil
}

func runInfoLookup(ctx *app.Context, input string, kind string, lookup itemLookup) error {
	client, cmdCtx, err := spotifyClient(ctx)
	if err != nil {
		return err
	}
	res, err := spotify.ParseTypedID(input, kind)
	if err != nil {
		return err
	}
	item, err := lookup(cmdCtx, client, res.ID)
	if err != nil {
		return err
	}
	return emitItem(ctx, item)
}

func emitItem(ctx *app.Context, item spotify.Item) error {
	return ctx.Output.Emit(item, []string{itemPlain(item)}, []string{itemHuman(ctx.Output, item)})
}

func emitItems(ctx *app.Context, items []spotify.Item, total int, extras map[string]any) error {
	plain, human := renderItems(ctx.Output, items)
	payload := map[string]any{
		"total": total,
		"items": items,
	}
	for key, value := range extras {
		payload[key] = value
	}
	return ctx.Output.Emit(payload, plain, human)
}

func emitOK(ctx *app.Context, payload map[string]any, human string) error {
	if payload == nil {
		payload = map[string]any{"status": "ok"}
	}
	return ctx.Output.Emit(payload, []string{"ok"}, []string{human})
}

func emitCountStatus(ctx *app.Context, count int, label string) error {
	return emitOK(ctx, map[string]any{"status": "ok", "count": count}, fmt.Sprintf("%s %d items", label, count))
}
