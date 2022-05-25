package today

import (
	"testing"

	"github.com/lzzzzl/basketball-go/modules/prepare"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSummaryPrinter(t *testing.T) {
	prepare.InitTestConfig()
	Convey("SummaryPrinter", t, func() {
		today := &Today{}
		_, err := today.TodayBoardPrinter()
		So(err, ShouldBeNil)
	})
}
