package datatypes

import (
	"time"
)

const (
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
