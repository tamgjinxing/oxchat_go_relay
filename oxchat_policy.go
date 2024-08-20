package main

import (
	"context"

	"github.com/fiatjaf/khatru"
	"github.com/nbd-wtf/go-nostr"
)

func PreventTimestampsFor1059(thresholdSeconds nostr.Timestamp) func(context.Context, *nostr.Event) (bool, string) {
	return func(ctx context.Context, event *nostr.Event) (reject bool, msg string) {
		if event.Kind == 1059 {
			if nostr.Now() < event.CreatedAt {
				return true, "event in the future"
			}
		} else {
			if nostr.Now()-event.CreatedAt > thresholdSeconds {
				return true, "event too old"
			}

			if event.CreatedAt-nostr.Now() > thresholdSeconds {
				return true, "event too much in the future"
			}
		}

		return false, ""
	}
}

func DBSaveEventLogic(ctx context.Context, evt *nostr.Event) error {
	defer func() {
		if r := recover(); r != nil {
			logger.Printf("Recovered from panic:%v\n", r)
			logger.Printf("Event Saved failed:%v\n", evt)
		}
	}()

	return db.SaveEvent(ctx, evt)
}

func PrintHeartEvent(ctx context.Context, event *nostr.Event) {
	if event.Kind == 22456 {
		logger.Printf("0xchat心跳事件kind=%d.EVENT为:%v\n", event.Kind, event)
	}
}

func OnConnectLogic(ctx context.Context) {
	logger.Printf("connected from:%s\n", khatru.GetIP(ctx))
}

func OnDisconnectLogic(ctx context.Context) {
	logger.Printf("disconnect from: %s\n", khatru.GetIP(ctx))
}
