package circular

import (
	"container/ring"

	"github.com/Nitro/sidecar/catalog"
	"github.com/Nitro/superside/datatypes"
)

// A Ring buffer for SvcEvents
type SvcEventsBuffer struct {
	changes *ring.Ring
}

// Return a new, properly configured circular buffer
func NewSvcEventsBuffer(size int) *SvcEventsBuffer {
	newRing := ring.New(size)

	return &SvcEventsBuffer{changes: newRing}
}

// Get all the items from the buffer that have a value, return as linear slice
func (b *SvcEventsBuffer) All() []datatypes.Notification {
	var changeHistory []datatypes.Notification
	b.changes.Do(func(evt interface{}) { // Start from oldest node
		if evt == nil {
			return
		}

		event := evt.(catalog.StateChangedEvent)
		changeHistory = append(changeHistory, *datatypes.NotificationFromEvent(&event))
	})

	return changeHistory
}

func (b *SvcEventsBuffer) AllRaw() []catalog.StateChangedEvent {
	var changeHistory []catalog.StateChangedEvent
	b.changes.Do(func(evt interface{}) { // Start from oldest node
		if evt == nil {
			return
		}
		event := evt.(catalog.StateChangedEvent)
		changeHistory = append(changeHistory, event)
	})

	return changeHistory
}

func (b *SvcEventsBuffer) Insert(evt catalog.StateChangedEvent) {
	b.changes.Value = evt
	b.changes = b.changes.Next()
}

// A Ring buffer for Deployments
type DeploymentsBuffer struct {
	deploys *ring.Ring
}

// Return a new, properly configured circular buffer
func NewDeploymentsBuffer(size int) *DeploymentsBuffer {
	newRing := ring.New(size)

	return &DeploymentsBuffer{deploys: newRing}
}

// Get all the items from the buffer that have a value, return as linear slice
func (b *DeploymentsBuffer) All() []*datatypes.Deployment {
	var deploys []*datatypes.Deployment
	b.deploys.Do(func(item interface{}) { // Start from oldest node
		if item == nil {
			return
		}

		deploy := item.(*datatypes.Deployment)
		deploys = append(deploys, deploy)
	})

	return deploys
}

func (b *DeploymentsBuffer) Insert(deploy *datatypes.Deployment) {
	b.deploys.Value = deploy
	b.deploys = b.deploys.Next()
}

func (b *DeploymentsBuffer) GetLast() *datatypes.Deployment {
	deploy := b.deploys.Prev()
	value := deploy.Value.(*datatypes.Deployment)
	return value
}
