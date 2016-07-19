package circular

import (
	"container/ring"

	"github.com/newrelic/sidecar/catalog"
	"github.com/nitro/superside/notification"
)

type Buffer struct {
	changes     *ring.Ring
}

// Return a new, properly configured circular buffer
func NewBuffer(size int) *Buffer {
	newRing := ring.New(size)

	return &Buffer{changes: newRing}
}

// Get all the items from the buffer that have a value, return as linear slice
func (b *Buffer) List() []notification.Notification {
	var changeHistory []notification.Notification
	b.changes.Do(func(evt interface{}) { // Start from oldest node
		if evt == nil {
			return
		}

		event := evt.(catalog.StateChangedEvent)
		changeHistory = append(changeHistory, *notification.FromEvent(&event))
	})

	return changeHistory
}

func (b *Buffer) Insert(evt catalog.StateChangedEvent) {
	b.changes.Value = evt
	b.changes = b.changes.Next()
}
