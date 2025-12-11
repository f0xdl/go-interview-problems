package matchstatsprocessor

import (
	"context"
	"sync"
)

const workers = 4

type MatchStatsProcessor interface {
	Process(ctx context.Context, events <-chan Event)
}

type Event struct {
	MatchID string
	Type    string
	Team    string
}

type MatchStats struct {
	events map[Event]int
	mu     sync.Mutex
}

func NewMatchStats() *MatchStats {
	return &MatchStats{
		events: map[Event]int{},
		mu:     sync.Mutex{},
	}
}

func (ms *MatchStats) Process(ctx context.Context, events <-chan Event) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg := &sync.WaitGroup{}
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case event, ok := <-events:
					if !ok {
						return
					}
					ms.mu.Lock()
					ms.events[event]++
					ms.mu.Unlock()
				}
			}

		}()
	}
	wg.Wait()
}

func (ms *MatchStats) GetMatchEvents(MatchID string) []Event {
	var events []Event
	for k, v := range ms.events {
		if k.MatchID == MatchID {
			for range v {
				events = append(events, k)
			}
		}
	}
	return events
}

func (ms *MatchStats) GetMap() map[Event]int {
	return ms.events
}
