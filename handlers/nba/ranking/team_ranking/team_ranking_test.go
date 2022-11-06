package teamranking

import (
	"testing"

	"github.com/lzzzzl/basketball-go/modules/prepare"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTeamRankingPrinter(t *testing.T) {
	prepare.InitTestConfig()
	Convey("Test Team Ranking Printer", t, func() {
		tRanking := &TeamRanking{}
		err := tRanking.TeamRankingPrinter("2023")
		So(err, ShouldBeNil)
	})
}
