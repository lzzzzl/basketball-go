package schedule

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/lzzzzl/basketball-go/handlers/nba/game"
	"github.com/lzzzzl/basketball-go/modules/constant"
	"github.com/lzzzzl/basketball-go/modules/log"
	"github.com/lzzzzl/basketball-go/modules/print"
	"github.com/lzzzzl/basketball-go/modules/request"
	"github.com/lzzzzl/basketball-go/modules/time"
	"github.com/tidwall/gjson"
)

// Schedule playoff shedule
type Schedule struct {
	GameStartTime string
	GameEndTime   string
}

// ScheduleBoardPrinter schedule board printer
func (schedule *Schedule) ScheduleBoardPrinter() ([]*game.GameBoard, error) {
	bList := make([]*game.GameBoard, 0)
	if schedule.GameStartTime == "" || schedule.GameEndTime == "" {
		return bList, errors.New("not found time")
	}
	url := constant.CURRENT_SEASON_SCHEDULE_URL
	log.Println("Scheduler Board URL: ", url)
	request := &request.HTTPRequest{
		URL:     url,
		Headers: constant.LIVE_HEADER,
		TimeOut: 5,
	}
	result, err := request.HTTPGet()
	if err != nil {
		return nil, err
	}
	bMap, bList, err := schedule.ResultParser(result)
	if err != nil {
		return bList, err
	}
	err = schedule.PrintTable(bMap)
	if err != nil {
		return bList, err
	}
	return bList, nil
}

// ResultParser parse http json data
func (schedule *Schedule) ResultParser(result string) (bMap map[string][]*game.GameBoard, bList []*game.GameBoard, err error) {
	bMap = make(map[string][]*game.GameBoard)
	bList = make([]*game.GameBoard, 0)

	gameDates := gjson.Get(result, "leagueSchedule.gameDates")
	gameDates.ForEach(func(key, value gjson.Result) bool {
		// parse est time to local time
		s := strings.Split(value.Get("gameDate").String(), " ")
		if len(s) == 1 {
			log.Error("gameDates wrong")
			return false
		}
		estTime := fmt.Sprintf("%s %s", s[0], "12:00:00")
		localTime, err := time.EstToLocal(estTime, "1/2/2006 15:04:05", constant.EST_ZONE)
		if err != nil {
			log.Error(fmt.Sprintf("time format wrong, %s", estTime))
			return false
		}
		s = strings.Split(localTime, " ")
		if len(s) == 1 {
			log.Error("gameDates wrong")
			return false
		}
		t := strings.ReplaceAll(s[0], "-", "")
		if schedule.GameStartTime <= t && schedule.GameEndTime >= t {
			bMap[t] = make([]*game.GameBoard, 0)
			games := value.Get("games")
			games.ForEach(func(m, n gjson.Result) bool {
				gameTime, err := time.EstToLocal(n.Get("gameDateTimeEst").String(),
					"2006-01-02T15:04:05Z",
					constant.EST_ZONE)
				if err != nil {
					log.Error("gameDates wrong")
					return false
				}
				board := &game.GameBoard{
					GameID:        n.Get("gameId").String(),
					GameStartTime: gameTime,
					GameStatus:    n.Get("gameStatusText").String(),
					SeriesText:    n.Get("seriesText").String(),
					HomeTeamName:  n.Get("homeTeam.teamName").String(),
					HomeTeamCity:  n.Get("homeTeam.teamCity").String(),
					HomeTeamCode:  n.Get("homeTeam.teamTricode").String(),
					HomeTeamScore: n.Get("homeTeam.score").String(),
					AwayTeamName:  n.Get("awayTeam.teamName").String(),
					AwayTeamCity:  n.Get("awayTeam.teamCity").String(),
					AwayTeamCode:  n.Get("awayTeam.teamTricode").String(),
					AwayTeamScore: n.Get("awayTeam.score").String(),
				}
				bMap[t] = append(bMap[t], board)
				bList = append(bList, board)
				return true
			})
		}
		return true
	})
	return bMap, bList, err
}

// PrintTable print table
func (schedule *Schedule) PrintTable(bMap map[string][]*game.GameBoard) error {
	var keys []string
	for k := range bMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		print.FigurePrint(k, "starwars", "green")
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.SetColumnConfigs([]table.ColumnConfig{
			{Name: game.ColTitleTeams, AlignHeader: text.AlignCenter},
			{Name: game.ColTitleScores, AlignHeader: text.AlignCenter},
			{Name: game.ColTitleStatus, Align: text.AlignCenter, AlignHeader: text.AlignCenter},
			{Name: game.ColTitleSeries, Align: text.AlignCenter, AlignHeader: text.AlignCenter},
		})
		t.AppendHeader(table.Row{game.ColTitleIndex, game.ColTitleTeams, game.ColTitleScores,
			game.ColTitleStatus, game.ColTitleSeries})
		for i, b := range bMap[k] {
			homeTeamScore, awayTeamScore, err := print.MaxScoreColor(b.HomeTeamScore, b.AwayTeamScore)
			if err != nil {
				return err
			}
			s := b.SeriesText
			if s == "" {
				s = "/"
			}
			t.AppendRow([]interface{}{
				i + 1,
				fmt.Sprintf("%12s  🆚  %-12s", b.HomeTeamName, b.AwayTeamName),
				fmt.Sprintf("%s : %s", homeTeamScore, awayTeamScore),
				b.GameStatus,
				s})
		}
		t.SetStyle(table.StyleBold)
		t.Style().Options.SeparateRows = true
		t.Render()
	}
	return nil
}
