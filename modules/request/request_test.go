package request

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHttpGet(t *testing.T) {
	Convey("HTTP GET", t, func() {

		body := HTTPRequest{
			URL: "https://cdn.nba.com/static/json/liveData/boxscore/boxscore_0022000180.json",
		}
		resp, err := body.HTTPGet()
		So(err, ShouldBeNil)
		t.Log(resp)
	})
}
