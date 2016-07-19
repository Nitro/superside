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
)

const (
	INITIAL_RING_SIZE   = 500
	CHANNEL_BUFFER_SIZE = 25
	INITIAL_DEPLOYMENT_SIZE
)

type Deployment struct {
	Name            string
	StartTime       time.Time
	EndTime         time.Time
	Version         string
	Image           string
}

func (d *Deployment) Matches(other *Deployment) bool {
	return other.Version == d.Version &&
		other.Image == d.Image &&
		other.Name == d.Name
}

type Tracker struct {
	changes     *circular.Buffer
	changesChan chan catalog.StateChangedEvent
	listeners   []chan *notification.Notification
	listenLock  sync.Mutex
	deployments map[string][]*Deployment
}

func NewTracker(changesRingSize int) *Tracker {
	return &Tracker{
		changesChan: make(chan catalog.StateChangedEvent, CHANNEL_BUFFER_SIZE),
		changes:     circular.NewBuffer(INITIAL_RING_SIZE),
		deployments: make(map[string][]*Deployment, INITIAL_DEPLOYMENT_SIZE),
	}
}

// Enqueue an update to the channel. Rely on channel buffer. We block if channel is full
func (t *Tracker) EnqueueUpdate(evt catalog.StateChangedEvent) {
	t.changesChan <- evt
}

// Subscribe a listener, returns a listening channel
func (t *Tracker) GetListener() chan *notification.Notification {
	listenChan := make(chan *notification.Notification, 100)

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
		case listener <- notification.FromEvent(evt):
		default:
		}
	}
}

func looksLikeDeployment(notice *notification.Notification) bool {
	evt := notice.Event
	svc := evt.Service

	return (svc.Status == service.ALIVE || svc.Status == service.UNHEALTHY) &&
		(evt.PreviousStatus == service.UNKNOWN)
}

func deploymentFromNotification(notice *notification.Notification) *Deployment {
	evt := notice.Event
	svc := evt.Service

	return &Deployment{
		Name:            svc.Name,
		StartTime:       evt.Time,
		EndTime:         evt.Time,
		Version:         strings.Split(evt.Service.Image, ":")[1],
		Image:           evt.Service.Image,
	}
}

func (t *Tracker) processOneDeployment(notice *notification.Notification) {
	evt := notice.Event
	svc := evt.Service

	thisDeploy := deploymentFromNotification(notice)

	if looksLikeDeployment(notice) {
		deploys := t.deployments[svc.Name]

		// We don't have any deployments for that service so let's add it
		if deploys == nil {
			t.insertDeployment(thisDeploy)
			log.Debug("Inserting deployment: %#v", thisDeploy)
			return
		}

		// We have some and the last one matches
		lastDeploy := deploys[len(deploys)-1]

		if lastDeploy.Matches(thisDeploy) {
			log.Debug("Found matching deployment: %#v", lastDeploy)

			if thisDeploy.StartTime.Before(lastDeploy.StartTime) {
				lastDeploy.StartTime = thisDeploy.StartTime
			}

			if lastDeploy.EndTime.Before(thisDeploy.EndTime) {
				lastDeploy.EndTime = thisDeploy.EndTime
			}
			return
		}

		log.Debug("Found no matching deployments, inserting: %#v", thisDeploy)
		// We have some but they don't match
		t.insertDeployment(thisDeploy)
	}
}

// Try to extrapolate when a deployment started and stopped for each service
func (t *Tracker) processDeployments() {
	notifyChan := t.GetListener()
	defer close(notifyChan)

	for notice := range notifyChan {
		t.processOneDeployment(notice)
	}
}

func (t *Tracker) insertDeployment(deploy *Deployment) {
	t.deployments[deploy.Name] = append(t.deployments[deploy.Name], deploy)
}

func (t *Tracker) GetChangesList() []notification.Notification {
	return t.changes.List()
}

// Linearize the updates coming in from the async HTTP handler
func (t *Tracker) ProcessUpdates() {
	go t.processDeployments()

	for evt := range t.changesChan {
		t.changes.Insert(evt)
		t.tellListeners(&evt)
	}
}
