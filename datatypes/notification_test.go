package datatypes

import (
	"testing"

	"github.com/newrelic/sidecar/catalog"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_NotificationFromEvent(t *testing.T) {
	Convey("NotificationFromEvent() constructs a proper notification", t, func() {
		change := catalog.ChangeEvent{}

		evt := &catalog.StateChangedEvent{
			ChangeEvent: change,
			State:       catalog.ServicesState{ClusterName: "awesome-cluster"},
		}

		notice := NotificationFromEvent(evt)

		So(notice.Event, ShouldResemble, &change)
		So(notice.ClusterName, ShouldEqual, evt.State.ClusterName)
	})
}
