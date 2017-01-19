package datatypes

import (
	"github.com/Nitro/sidecar/catalog"
)

type Notification struct {
	Event       *catalog.ChangeEvent
	ClusterName string
}

func NotificationFromEvent(evt *catalog.StateChangedEvent) *Notification {
	return &Notification{
		Event: &evt.ChangeEvent,
		ClusterName: evt.State.ClusterName,
	}
}
