package thebus

import (
	"sync/atomic"
	"time"
)

// Counters holds numeric counters for published, delivered, etc.
type Counters struct {
	Published uint64
	Delivered uint64
	Failed    uint64
	Dropped   uint64
}

// StatsResults represents aggregated statistics of the bus.
type StatsResults struct {
	StartedAt   time.Time
	Open        bool
	Topics      int
	Subscribers int
	Totals      Counters
	PerTopic    map[string]TopicStats
}

// TopicStats represents statistics for a single topic.
type TopicStats struct {
	Subscribers int
	Buffered    int
	Counters
}

type atomicCounters struct {
	Published atomic.Uint64
	Delivered atomic.Uint64
	Failed    atomic.Uint64
	Dropped   atomic.Uint64
}
