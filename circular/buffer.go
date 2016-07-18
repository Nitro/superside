package circular

import (
	"container/ring"

	"github.com/newrelic/sidecar/catalog"
	"github.com/nitro/superside/notification"
)

const (
	INITIAL_RING_SIZE = 20
)

type Buffer struct {
	changes        *ring.Ring
	ringSize       int
}

func (b *Buffer) List() []notification.Notification {
	var changeHistory []notification.Notification
	b.changes.Do(func(evt interface{}) {
		if evt != nil {
			event := evt.(catalog.StateChangedEvent)
			changeHistory = append(changeHistory, *notification.FromEvent(&event))
		}
	})

	return changeHistory
}

func (b *Buffer) Insert(evt catalog.StateChangedEvent) {
	newEntry := &ring.Ring{Value: evt}

	if b.ringSize == 0 {
		b.changes = newEntry
		b.ringSize += 1
	} else if b.ringSize < INITIAL_RING_SIZE {
		b.changes.Prev().Link(newEntry)
		b.ringSize += 1
	} else {
		b.changes = b.changes.Prev()
		b.changes.Unlink(1)
		b.changes = b.changes.Next()
		b.changes.Prev().Link(newEntry)
	}
}
