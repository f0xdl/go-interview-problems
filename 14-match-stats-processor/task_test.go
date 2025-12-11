package matchstatsprocessor_test

import (
	"context"
	"testing"

	. "github.com/blindlobstar/go-interview-problems/14-match-stats-processor"
	"github.com/stretchr/testify/assert"
)

func Test_Process_Empty(t *testing.T) {
	events := make(chan Event)
	stats := NewMatchStats()
	close(events)
	stats.Process(context.Background(), events)
	assert.Empty(t, stats.GetMap(), "check empty")
}

func Test_Process_1Match(t *testing.T) {
	events := make(chan Event, 10)
	go func() {
		events <- Event{MatchID: "1", Type: "goal", Team: "guess"}
		events <- Event{MatchID: "1", Type: "goal", Team: "home"}
		events <- Event{MatchID: "1", Type: "goal", Team: "guess"}
		close(events)

	}()

	stats := NewMatchStats()
	stats.Process(context.Background(), events)
	matchEvents := stats.GetMatchEvents("1")

	assert.Equal(t, 3, len(matchEvents), "count event")
	for id, event := range matchEvents {
		assert.Equal(t, "goal", event.Type, "event %d: type: %s team: %s", id, event.Type, event.Team)
	}
}
