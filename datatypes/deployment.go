package datatypes

import (
	"strings"
	"time"

	"github.com/satori/go.uuid"
)

const (
	DEPLOYMENT_CUTOFF = 10 * time.Minute
)

type Deployment struct {
	ID          string
	Name        string
	StartTime   time.Time
	EndTime     time.Time
	Version     string
	Image       string
	ClusterName string
	Hostnames   []string
}

func (d *Deployment) Matches(other *Deployment) bool {
	// If a deployment was more than DEPLOYMENT_CUTOFF after the last event
	// for the same matching deployment, we'll call it a new deployment.
	timeDiff := other.StartTime.Sub(d.StartTime)

	return other.Version == d.Version &&
		other.Image == d.Image &&
		other.Name == d.Name &&
		other.ClusterName == d.ClusterName &&
		timeDiff < DEPLOYMENT_CUTOFF
}

func (d *Deployment) Aggregate(other *Deployment) {
	if other.StartTime.Before(d.StartTime) {
		d.StartTime = other.StartTime
	}

	if d.EndTime.Before(other.EndTime) {
		d.EndTime = other.EndTime
	}

	// Dupes are desirable here... we might deploy more than once on
	// the same host.
	d.Hostnames = append(d.Hostnames, other.Hostnames...)
}

// Construct a deployment object from a datatypes object
func DeploymentFromNotification(notice *Notification) *Deployment {
	evt := notice.Event
	if evt == nil {
		return nil
	}

	svc := evt.Service

	return &Deployment{
		ID:          uuid.NewV4().String(),
		Name:        svc.Name,
		StartTime:   evt.Time,
		EndTime:     evt.Time,
		Version:     strings.Split(evt.Service.Image, ":")[1],
		Image:       evt.Service.Image,
		ClusterName: notice.ClusterName,
		Hostnames:   []string{evt.Service.Hostname},
	}
}
