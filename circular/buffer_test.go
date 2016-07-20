package circular

import (
	"testing"

	"github.com/newrelic/sidecar/catalog"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/nitro/superside/datatypes"
)

func Test_SvcEventsBuffer(t *testing.T) {
	Convey("Working with SvcEventsBuffer", t, func() {
		buffer := NewSvcEventsBuffer(10)

		change := catalog.ChangeEvent{}

		evt := catalog.StateChangedEvent{
			ChangeEvent: change,
			State:       catalog.ServicesState{ClusterName: "awesome-cluster"},
		}

		Convey("Inserts new values", func() {
			for i := 0; i < 10; i++ {
				buffer.Insert(evt)
			}

			all := buffer.All()
			So(&all[0], ShouldResemble, datatypes.NotificationFromEvent(&evt))
		})

		Convey("Inserts more than the size", func() {
			for i := 0; i < 20; i++ {
				buffer.Insert(evt)
			}

			all := buffer.All()
			So(&all[0], ShouldResemble, datatypes.NotificationFromEvent(&evt))
		})
	})
}
