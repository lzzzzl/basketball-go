package time

import (
	"testing"

	"github.com/lzzzzl/basketball-go/modules/constant"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUtcToLocal(t *testing.T) {
	Convey("GameScore", t, func() {
		utcTime := "2022-04-13T23:00:00Z"
		localTime, err := UtcToLocal(utcTime, "2006-01-02T15:04:05Z", constant.EST_ZONE)
		So(err, ShouldBeNil)
		t.Log(localTime)
	})
}
