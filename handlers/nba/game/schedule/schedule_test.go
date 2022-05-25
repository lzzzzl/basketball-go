package schedule

import (
	"testing"

	"github.com/lzzzzl/basketball-go/modules/prepare"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSummaryPrinter(t *testing.T) {
	prepare.InitTestConfig()
	Convey("SummaryPrinter", t, func() {
		schedule := &Schedule{
			GameStartTime: "20220511",
			GameEndTime:   "20220511",
		}
		bList, err := schedule.ScheduleBoardPrinter()
		So(err, ShouldBeNil)
		t.Log(bList)
	})
}
