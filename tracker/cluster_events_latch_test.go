package tracker

import (
	"testing"
	"time"

	"github.com/newrelic/sidecar/catalog"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_ClusterEventsLatch(t *testing.T) {
	Convey("ClusterEventsLatch()", t, func() {
		latch := NewClusterEventsLatch()

		evt1 := &catalog.StateChangedEvent{
			State: catalog.ServicesState{
				ClusterName: "france",
				Hostname: "joffre",
			},
		}

		evt2 := &catalog.StateChangedEvent{
			State: catalog.ServicesState{
				ClusterName: "france",
				Hostname: "foch",
			},
		}

		Convey("Accepts an event from any new cluster", func() {
			So(latch.ShouldAccept(evt1), ShouldBeTrue)
		})

		Convey("Accepts an event from the latched host", func() {
			latch.ShouldAccept(evt1)
			So(latch.ShouldAccept(evt1), ShouldBeTrue)
		})

		Convey("Does not accept events for the same cluster from another host", func() {
			latch.ShouldAccept(evt1)
			So(latch.ShouldAccept(evt2), ShouldBeFalse)
		})

		Convey("Accepts events for the cluster from another host if entry expired", func() {
			latch.ShouldAccept(evt1)
			latch.Latches[evt1.State.ClusterName].Expiry = time.Unix(0,0)
			So(latch.ShouldAccept(evt2), ShouldBeTrue)
		})
	})
}
