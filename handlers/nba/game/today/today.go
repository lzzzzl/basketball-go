package today

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/lzzzzl/basketball-go/modules/constant"
	"github.com/lzzzzl/basketball-go/modules/log"
	"github.com/lzzzzl/basketball-go/modules/print"
	"github.com/lzzzzl/basketball-go/modules/request"
	"github.com/lzzzzl/basketball-go/modules/time"

	"github.com/lzzzzl/basketball-go/handlers/nba/game"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/tidwall/gjson"
)

type Today struct {
}

func (today *Today) TodayBoardPrinter() ([]*game.GameBoard, error) {
	url := constant.TODAY_SCORE_BOARD_URL
	log.Println("Today Score Board URL: ", url)
	request := &request.HttpRequest{
		Url:     url,
		Headers: constant.LIVE_HEADER,
		TimeOut: 30,
	}
	result, err := request.HttpGet()
	if err != nil {
		return nil, err
	}
	boards, err := today.ResultParser(result)
	if err != nil {
		return boards, err
	}
	err = PrintTime(result)
	if err != nil {
		return boards, err
	}
	today.PrintTable(boards)

	return boards, nil
}

func (today *Today) ResultParser(result string) ([]*game.GameBoard, error) {
	var boards []*game.GameBoard
	if gjson.Get(result, "meta.code").Int() != 200 {
		return boards, errors.New("error response")
	}
	games := gjson.Get(result, "scoreboard.games")
	games.ForEach(func(key, value gjson.Result) bool {
		utcTime, err := time.UtcToLocal(value.Get("gameTimeUTC").String(),
			"2006-01-02T15:04:05Z",
			constant.EST_ZONE)
		if err != nil {
			log.Error(err)
			return false
		}
		boards = append(boards, &game.GameBoard{
			GameStartTime: utcTime,
			GameId:        value.Get("gameId").String(),
			GameStatus:    value.Get("gameStatusText").String(),
			SeriesText:    value.Get("seriesText").String(),
			HomeTeamName:  value.Get("homeTeam.teamName").String(),
			HomeTeamCity:  value.Get("homeTeam.teamCity").String(),
			HomeTeamCode:  value.Get("homeTeam.teamTricode").String(),
			HomeTeamScore: value.Get("homeTeam.score").String(),
			AwayTeamName:  value.Get("awayTeam.teamName").String(),
			AwayTeamCity:  value.Get("awayTeam.teamCity").String(),
			AwayTeamCode:  value.Get("awayTeam.teamTricode").String(),
			AwayTeamScore: value.Get("awayTeam.score").String(),
		})
		return true
	})

	return boards, nil
}

func PrintTime(json string) error {
	if gjson.Get(json, "meta.code").Int() != 200 {
		return errors.New("error response")
	}
	gameDate := gjson.Get(json, "scoreboard.gameDate").String()
	gameDate = strings.ReplaceAll(gameDate, "-", "")
	print.FigurePrint(gameDate, "starwars", "green")
	fmt.Println("")
	return nil
}

func (today *Today) PrintTable(boards []*game.GameBoard) error {
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
	for i, v := range boards {
		homeTeamScore, awayTeamScore, err := print.MaxScoreColor(v.HomeTeamScore, v.AwayTeamScore)
		if err != nil {
			return err
		}
		t.AppendRow([]interface{}{
			i + 1,
			fmt.Sprintf("%12s  🆚  %-12s", v.HomeTeamName, v.AwayTeamName),
			fmt.Sprintf("%s : %s", homeTeamScore, awayTeamScore),
			v.GameStatus,
			v.SeriesText})
	}
	t.SetStyle(table.StyleBold)
	t.Style().Options.SeparateRows = true
	t.Render()

	return nil
}
