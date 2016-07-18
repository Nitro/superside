package notification

import (
	"github.com/newrelic/sidecar/catalog"
)

type Notification struct {
	Event       *catalog.ChangeEvent
	ClusterName string
}

func FromEvent(evt *catalog.StateChangedEvent) *Notification {
	return &Notification{
		Event: &evt.ChangeEvent,
		ClusterName: evt.State.ClusterName,
	}
}
