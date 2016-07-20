package tracker

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/newrelic/sidecar/catalog"
	"github.com/newrelic/sidecar/service"
	"github.com/nitro/superside/circular"
	"github.com/nitro/superside/datatypes"
)

const (
	INITIAL_RING_SIZE       = 500 // We'll track 500 service events globally
	CHANNEL_BUFFER_SIZE     = 25
	INITIAL_DEPLOYMENT_SIZE = 20
)

type Tracker struct {
	svcEvents           *circular.SvcEventsBuffer
	svcEventsChan       chan catalog.StateChangedEvent
	svcEventsListeners  []chan *datatypes.Notification
	deploymentListeners []chan *datatypes.Deployment
	listenLock          sync.Mutex
	deployments         map[string]*circular.DeploymentsBuffer
}

func NewTracker(svcEventsRingSize int) *Tracker {
	return &Tracker{
		svcEventsChan: make(chan catalog.StateChangedEvent, CHANNEL_BUFFER_SIZE),
		svcEvents:     circular.NewSvcEventsBuffer(svcEventsRingSize),
		deployments:   make(map[string]*circular.DeploymentsBuffer, INITIAL_DEPLOYMENT_SIZE),
	}
}

// Enqueue an update to the channel. Rely on channel buffer. We block if channel is full
func (t *Tracker) EnqueueUpdate(evt catalog.StateChangedEvent) {
	t.svcEventsChan <- evt
}

// Subscribe a service events listener, returns a listening channel
func (t *Tracker) GetSvcEventsListener() chan *datatypes.Notification {
	listenChan := make(chan *datatypes.Notification, 100)

	t.listenLock.Lock()
	t.svcEventsListeners = append(t.svcEventsListeners, listenChan)
	t.listenLock.Unlock()

	return listenChan
}

// Subscribe a deployment events listener, returns a listening channel
func (t *Tracker) GetDeploymentListener() chan *datatypes.Deployment {
	listenChan := make(chan *datatypes.Deployment, 100)

	t.listenLock.Lock()
	t.deploymentListeners = append(t.deploymentListeners, listenChan)
	t.listenLock.Unlock()

	return listenChan
}

// Announce changes to all service event listeners
func (t *Tracker) tellSvcEventListeners(evt *catalog.StateChangedEvent) {
	t.listenLock.Lock()
	defer t.listenLock.Unlock()

	// Try to tell the listener about the change but use a select
	// to protect us from any blocking readers.
	for _, listener := range t.svcEventsListeners {
		select {
		case listener <- datatypes.NotificationFromEvent(evt):
		default:
		}
	}
}

// Announce changes to all deployment listeners
func (t *Tracker) tellDeploymentListeners(deploy *datatypes.Deployment) {
	t.listenLock.Lock()
	defer t.listenLock.Unlock()

	// Try to tell the listener about the change but use a select
	// to protect us from any blocking readers.
	for _, listener := range t.deploymentListeners {
		select {
		case listener <- deploy:
		default:
		}
	}
}

// Compare some stuff and decide if this datatypes looks like it's
// a deployment event.
func looksLikeDeployment(notice *datatypes.Notification) bool {
	evt := notice.Event
	svc := evt.Service

	return (svc.Status == service.ALIVE || svc.Status == service.UNHEALTHY) &&
		(evt.PreviousStatus == service.UNKNOWN)
}

// Handle processing a single datatypes
func (t *Tracker) processOneDeployment(notice *datatypes.Notification) {
	evt := notice.Event
	svc := evt.Service

	thisDeploy := datatypes.DeploymentFromNotification(notice)

	if looksLikeDeployment(notice) {
		deploys := t.deployments[svc.Name]

		// We don't have any deployments for that service so let's add it
		if deploys == nil {
			t.insertDeployment(thisDeploy)
			log.Debug("Inserting deployment: ", thisDeploy)
			return
		}

		// We have some and the last one matches
		lastDeploy := deploys.GetLast()

		if lastDeploy.Matches(thisDeploy) {
			log.Debug("Found matching deployment: ", lastDeploy)

			if thisDeploy.StartTime.Before(lastDeploy.StartTime) {
				lastDeploy.StartTime = thisDeploy.StartTime
			}

			if lastDeploy.EndTime.Before(thisDeploy.EndTime) {
				lastDeploy.EndTime = thisDeploy.EndTime
			}

			t.tellDeploymentListeners(lastDeploy) // Send the updated original
			return
		}

		log.Debug("Found no matching deployments, inserting: ", thisDeploy)
		// We have some but they don't match
		t.insertDeployment(thisDeploy)
	}
}

func (t *Tracker) RemoveSvcEventsListener(victim chan *datatypes.Notification) {
	t.listenLock.Lock()
	defer t.listenLock.Unlock()

	for i, listener := range t.svcEventsListeners {
		if listener == victim {
			// Delete the item from the list
			t.svcEventsListeners = append(t.svcEventsListeners[:i], t.svcEventsListeners[i+1:]...)
			close(listener)
			return
		}
	}
}

func (t *Tracker) RemoveDeploymentListener(victim chan *datatypes.Deployment) {
	t.listenLock.Lock()
	defer t.listenLock.Unlock()

	for i, listener := range t.deploymentListeners {
		if listener == victim {
			// Delete the item from the list
			t.deploymentListeners = append(t.deploymentListeners[:i], t.deploymentListeners[i+1:]...)
			close(listener)
			return
		}
	}

}

// Try to extrapolate when a deployment started and stopped for each service
func (t *Tracker) processDeployments() {
	notifyChan := t.GetSvcEventsListener()
	defer close(notifyChan)

	for notice := range notifyChan {
		t.processOneDeployment(notice)
	}
}

// Add a new deployment, also announce it to listeners
func (t *Tracker) insertDeployment(deploy *datatypes.Deployment) {
	if t.deployments[deploy.Name] == nil {
		t.deployments[deploy.Name] = circular.NewDeploymentsBuffer(INITIAL_DEPLOYMENT_SIZE)
	}

	t.deployments[deploy.Name].Insert(deploy)
	t.tellDeploymentListeners(deploy)
}

func (t *Tracker) GetSvcEventsList() []datatypes.Notification {
	return t.svcEvents.All()
}

func (t *Tracker) GetDeployments() map[string][]*datatypes.Deployment {
	allDeploys := make(map[string][]*datatypes.Deployment, len(t.deployments))
	for name, ring := range t.deployments {
		allDeploys[name] = ring.All()
	}
	return allDeploys
}

// Linearize the updates coming in from the async HTTP handler
func (t *Tracker) ProcessUpdates() {
	go t.processDeployments()

	for evt := range t.svcEventsChan {
		t.svcEvents.Insert(evt)
		t.tellSvcEventListeners(&evt)
	}
}
