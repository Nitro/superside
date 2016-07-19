package tracker

import (
	"sync"

	"github.com/newrelic/sidecar/catalog"
	"github.com/nitro/superside/circular"
	"github.com/nitro/superside/notification"
)

const (
	INITIAL_RING_SIZE   = 500
	CHANNEL_BUFFER_SIZE = 25
)

type Tracker struct {
	changes     *circular.Buffer
	changesChan chan catalog.StateChangedEvent
	listeners   []chan notification.Notification
	listenLock  sync.Mutex
}

func NewTracker(changesRingSize int) *Tracker {
	return &Tracker{
		changesChan: make(chan catalog.StateChangedEvent, CHANNEL_BUFFER_SIZE),
		changes:     circular.NewBuffer(INITIAL_RING_SIZE),
	}
}

// Enqueue an update to the channel. Rely on channel buffer. We block if channel is full
func (t *Tracker) EnqueueUpdate(evt catalog.StateChangedEvent) {
	t.changesChan <- evt
}

// Subscribe a listener, returns a listening channel
func (t *Tracker) GetListener() chan notification.Notification {
	listenChan := make(chan notification.Notification, 100)

	t.listenLock.Lock()
	t.listeners = append(t.listeners, listenChan)
	t.listenLock.Unlock()

	return listenChan
}

// Announce changes to all listeners
func (t *Tracker) tellListeners(evt *catalog.StateChangedEvent) {
	t.listenLock.Lock()
	defer t.listenLock.Unlock()

	// Try to tell the listener about the change but use a select
	// to protect us from any blocking readers.
	for _, listener := range t.listeners {
		select {
		case listener <- *notification.FromEvent(evt):
		default:
		}
	}
}

func (t *Tracker) GetChangesList() []notification.Notification {
	return t.changes.List()
}

// Linearize the updates coming in from the async HTTP handler
func (t *Tracker) ProcessUpdates() {
	for evt := range t.changesChan {
		t.changes.Insert(evt)
		t.tellListeners(&evt)
	}
}
