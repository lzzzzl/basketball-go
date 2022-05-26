package board

import (
	"testing"

	"github.com/lzzzzl/basketball-go/modules/prepare"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPlayerBoardPrinter(t *testing.T) {
	prepare.InitTestConfig()
	Convey("PlayerBoardPrinter", t, func() {
		board := &PlayerBoard{GameID: "0042100102"}
		err := board.PlayerBoardPrinter()
		So(err, ShouldBeNil)
	})
}
