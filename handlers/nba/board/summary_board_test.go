package board

import (
	"testing"

	"github.com/lzzzzl/basketball-go/modules/prepare"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSummaryPrinter(t *testing.T) {
	prepare.InitTestConfig()
	Convey("SummaryPrinter", t, func() {
		board := &SummaryBoard{GameId: "0042100102"}
		err := board.SummaryBoardPrinter()
		So(err, ShouldBeNil)
	})
}
