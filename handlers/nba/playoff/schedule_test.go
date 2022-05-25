package playoff

import (
	"testing"

	"github.com/lzzzzl/basketball-go/modules/prepare"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPlayOffSchedulePrinter(t *testing.T) {
	prepare.InitTestConfig()
	Convey("PlayOffSchedule", t, func() {
		p := &PlayOffSchedule{Year: 2022}
		err := p.PlayOffSchedulePrinter()
		So(err, ShouldBeNil)
	})
}
