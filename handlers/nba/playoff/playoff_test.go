package playoff

import (
	"testing"

	"github.com/lzzzzl/basketball-go/modules/prepare"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPlayoffBracketPrinter(t *testing.T) {
	prepare.InitTestConfig()
	Convey("SummaryPrinter", t, func() {
		playoff := Playoff{Year: 2020}
		err := playoff.PlayoffBracketPrinter()
		So(err, ShouldBeNil)
	})
}
