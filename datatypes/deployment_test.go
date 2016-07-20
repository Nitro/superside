package datatypes

import (
	"testing"
	"time"

    . "github.com/smartystreets/goconvey/convey"
	"github.com/newrelic/sidecar/catalog"
	"github.com/newrelic/sidecar/service"
)

func Test_Matches(t *testing.T) {
	Convey("Matches()", t, func() {
		baseTime := time.Now().UTC()

		thisDeploy := &Deployment{
			ID:          "marechalfoch",
			Name:        "awesome-svc",
			StartTime:   baseTime.Add(-1*time.Second),
			EndTime:     baseTime,
			Version:     "0.2",
			Image:       "awesome-svc:0.2",
			ClusterName: "awesome-cluster",
		}

		thatDeploy := *thisDeploy

		Convey("Identifies when deploys don't match", func() {
			thatDeploy.Version = "0.1"
			So(thatDeploy.Matches(thisDeploy), ShouldBeFalse)
		})

		Convey("Identifies when deploys match", func() {
			thatDeploy.StartTime = baseTime.Add(-5*time.Second)
			So(thatDeploy.Matches(thisDeploy), ShouldBeTrue)
		})

		Convey("Properly separates deployments by time threshold", func() {
			thatDeploy.StartTime = baseTime.Add(-12*time.Minute)
			So(thatDeploy.Matches(thisDeploy), ShouldBeFalse)
		})
	})
}

func Test_DeploymentFromNotification(t *testing.T) {
	Convey("Generates a properly formed Deployment", t, func() {
		notice := &Notification{
			Event: &catalog.ChangeEvent{
				Service: service.Service{Image: "awesome-svc:0.1"},
			},
			ClusterName: "awesome-cluster",
		}

		deploy := DeploymentFromNotification(notice)

		So(deploy, ShouldNotBeNil)
		So(deploy.ID, ShouldNotBeEmpty)
		So(deploy.ClusterName, ShouldEqual, notice.ClusterName)
		So(deploy.Version, ShouldEqual, "0.1")
	})
}
