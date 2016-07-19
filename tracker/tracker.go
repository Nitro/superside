package tracker

import (
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/newrelic/sidecar/catalog"
	"github.com/newrelic/sidecar/service"
	"github.com/nitro/superside/circular"
	"github.com/nitro/superside/notification"
	"github.com/satori/go.uuid"
)

const (
	INITIAL_RING_SIZE   = 500
	CHANNEL_BUFFER_SIZE = 25
	INITIAL_DEPLOYMENT_SIZE
	DEPLOYMENT_CUTOFF = 10*time.Minute
)

type Deployment struct {
	ID        string
	Name      string
	StartTime time.Time
	EndTime   time.Time
	Version   string
	Image     string
}

func (d *Deployment) Matches(other *Deployment) bool {
	// If a deployment was more than DEPLOYMENT_CUTOFF after the last event
	// for the same matching deployment, we'll call it a new deployment.
	timeDiff := other.StartTime.Sub(d.StartTime)

	return other.Version == d.Version &&
		other.Image == d.Image &&
		other.Name == d.Name &&
		timeDiff < DEPLOYMENT_CUTOFF
}

type Tracker struct {
	svcEvents           *circular.Buffer
	svcEventsChan       chan catalog.StateChangedEvent
	svcEventsListeners  []chan *notification.Notification
	deploymentListeners []chan *Deployment
	listenLock          sync.Mutex
	deployments         map[string][]*Deployment
}

func NewTracker(svcEventsRingSize int) *Tracker {
	return &Tracker{
		svcEventsChan: make(chan catalog.StateChangedEvent, CHANNEL_BUFFER_SIZE),
		svcEvents:     circular.NewBuffer(svcEventsRingSize),
		deployments:   make(map[string][]*Deployment, INITIAL_DEPLOYMENT_SIZE),
	}
}

// Enqueue an update to the channel. Rely on channel buffer. We block if channel is full
func (t *Tracker) EnqueueUpdate(evt catalog.StateChangedEvent) {
	t.svcEventsChan <- evt
}

// Subscribe a service events listener, returns a listening channel
func (t *Tracker) GetSvcEventsListener() chan *notification.Notification {
	listenChan := make(chan *notification.Notification, 100)

	t.listenLock.Lock()
	t.svcEventsListeners = append(t.svcEventsListeners, listenChan)
	t.listenLock.Unlock()

	return listenChan
}

// Subscribe a deployment events listener, returns a listening channel
func (t *Tracker) GetDeploymentListener() chan *Deployment {
	listenChan := make(chan *Deployment, 100)

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
		case listener <- notification.FromEvent(evt):
		default:
		}
	}
}

// Announce changes to all deployment listeners
func (t *Tracker) tellDeploymentListeners(deploy *Deployment) {
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

// Compare some stuff and decide if this notification looks like it's
// a deployment event.
func looksLikeDeployment(notice *notification.Notification) bool {
	evt := notice.Event
	svc := evt.Service

	return (svc.Status == service.ALIVE || svc.Status == service.UNHEALTHY) &&
		(evt.PreviousStatus == service.UNKNOWN)
}

// Construct a deployment object from a notification object
func deploymentFromNotification(notice *notification.Notification) *Deployment {
	evt := notice.Event
	svc := evt.Service

	return &Deployment{
		ID:        uuid.NewV4().String(),
		Name:      svc.Name,
		StartTime: evt.Time,
		EndTime:   evt.Time,
		Version:   strings.Split(evt.Service.Image, ":")[1],
		Image:     evt.Service.Image,
	}
}

// Handle processing a single notification
func (t *Tracker) processOneDeployment(notice *notification.Notification) {
	evt := notice.Event
	svc := evt.Service

	thisDeploy := deploymentFromNotification(notice)

	if looksLikeDeployment(notice) {
		deploys := t.deployments[svc.Name]

		// We don't have any deployments for that service so let's add it
		if deploys == nil {
			t.insertDeployment(thisDeploy)
			log.Debug("Inserting deployment: ", thisDeploy)
			return
		}

		// We have some and the last one matches
		lastDeploy := deploys[len(deploys)-1]

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

func (t *Tracker) RemoveSvcEventsListener(victim chan *notification.Notification) {
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

func (t *Tracker) RemoveDeploymentListener(victim chan *Deployment) {
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
func (t *Tracker) insertDeployment(deploy *Deployment) {
	t.deployments[deploy.Name] = append(t.deployments[deploy.Name], deploy)
	t.tellDeploymentListeners(deploy)
}

func (t *Tracker) GetSvcEventsList() []notification.Notification {
	return t.svcEvents.List()
}

func (t *Tracker) GetDeployments() map[string][]*Deployment {
	return t.deployments
}

// Linearize the updates coming in from the async HTTP handler
func (t *Tracker) ProcessUpdates() {
	go t.processDeployments()

	for evt := range t.svcEventsChan {
		t.svcEvents.Insert(evt)
		t.tellSvcEventListeners(&evt)
	}
}
