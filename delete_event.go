package main

import (
	"context"

	"github.com/nbd-wtf/go-nostr"
)

var (
	ctx = context.Background()
)

func DeleteEventFor0xchat() {
	logger.Printf("begin to run delete event")
	filters := nostr.Filter{
		Kinds:   []int{30078},
		Authors: []string{"093dff31a87bbf838c54fd39ff755e72b38bd6b7975c670c0f2633fa7c54ddd0"},
	}

	ch, err := db.QueryEvents(ctx, filters)
	if err != nil {
		logger.Printf("query event failed:%v\n", err)
		return
	}

	eventsCh := make(chan *nostr.Event)
	go func() {
		defer close(eventsCh)
		for event := range ch {
			db.DeleteEvent(ctx, event)
			logger.Printf("Delete Event,EVENT_ID:%s\n", event.ID)
		}
	}()
}
