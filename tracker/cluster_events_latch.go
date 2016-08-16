package tracker

import (
	"time"

	"github.com/newrelic/sidecar/catalog"
)

const (
	LATCH_INTERVAL = 5 * time.Second
)

// We will receive many duplicate events from hosts in the same cluster. One way
// to de-dupe is to latch on to one host at a time for a short time interval. If
// the hostname matches, or if the interval has expired, then we update with the
// new hostname and keep rolling along.
//
// Pros: Really simple, not much state
// Cons: Up to LATCH_INTERVAL missed events. Needs to run serialized or locked.

// Maps clustername => hostname and expiry time for this latch
type LatchEntry struct {
	Hostname string
	Expiry   time.Time
}
type ClusterEventsLatch struct {
	Latches map[string]*LatchEntry
}

func NewClusterEventsLatch() *ClusterEventsLatch {
	return &ClusterEventsLatch{
		Latches: make(map[string]*LatchEntry, 5),
	}
}

func (l *ClusterEventsLatch) ShouldAccept(event *catalog.StateChangedEvent) bool {
	// No *Name, no shoes, no service
	if event.State.ClusterName == "" || event.State.Hostname == "" {
		return false
	}

	latch, ok := l.Latches[event.State.ClusterName]

	// Don't have an entry for this cluster, so add it and allow
	if !ok {
		l.Latches[event.State.ClusterName] = &LatchEntry{
			Hostname: event.State.Hostname,
			Expiry:   time.Now().UTC().Add(LATCH_INTERVAL),
		}
		return true
	}

	// We have this cluster, so validate the hostname or expiry
	if latch.Hostname == event.State.Hostname {
		// Reset the expiry and approve
		latch.Expiry = time.Now().UTC().Add(LATCH_INTERVAL)
		return true
	}

	// Now we have an entry but the hostname doesn't match

	// If the latch is expired then reset it to now and this host
	if latch.Expiry.Before(time.Now().UTC()) {
		latch.Expiry = time.Now().UTC().Add(LATCH_INTERVAL)
		latch.Hostname = event.State.Hostname
		return true
	}

	// Otherwise it's junk
	return false
}
