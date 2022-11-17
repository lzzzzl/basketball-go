package playerranking

import (
	"testing"

	"github.com/lzzzzl/basketball-go/modules/prepare"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPlayerRankingPrinter(t *testing.T) {
	prepare.InitTestConfig()
	Convey("Test Player Ranking Printer", t, func() {
		pRanking := &PlayerRanking{}
		err := pRanking.PlayerRankingPrinter("2022", "3")
		So(err, ShouldBeNil)
	})
}
