package team

import (
	"testing"

	"github.com/lzzzzl/basketball-go/modules/prepare"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTeamListPrinter(t *testing.T) {
	prepare.InitTestConfig()
	Convey("PlayOffSchedule", t, func() {
		t := &Team{}
		err := t.TeamListPrinter()
		So(err, ShouldBeNil)
	})
}
