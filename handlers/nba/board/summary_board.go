package board

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/lzzzzl/basketball-go/handlers/nba/stats"
	"github.com/lzzzzl/basketball-go/modules/constant"
	"github.com/lzzzzl/basketball-go/modules/log"
	"github.com/lzzzzl/basketball-go/modules/print"
	"github.com/lzzzzl/basketball-go/modules/request"
	"github.com/lzzzzl/basketball-go/modules/time"
	"github.com/tidwall/gjson"
)

var statsRows = table.Row{" / ", "ast", "blk", "stl", "tor", "reb"}

type SummaryBoard struct {
	GameId         string
	GameTime       string
	AreanCity      string
	AreanName      string
	HomeTeamCode   string
	AwayTeamCode   string
	HomeTeamScores []int
	AwayTeamScores []int
	HomeTeamStats  stats.Stats
	AwayTeamStats  stats.Stats
}

func (board *SummaryBoard) SummaryBoardPrinter() error {
	url := fmt.Sprintf(constant.GAME_DETAILS_URL, board.GameId)
	log.Println("Game Summary Board URL: ", url)
	request := &request.HttpRequest{
		Url:     url,
		Headers: constant.LIVE_HEADER,
		TimeOut: 5,
	}
	result, err := request.HttpGet()
	if err != nil {
		return err
	}
	err = board.ResultParser(result)
	if err != nil {
		return err
	}
	err = board.PrintTable()
	if err != nil {
		return err
	}
	return nil
}

func (board *SummaryBoard) ResultParser(result string) error {
	if gjson.Get(result, "meta.code").Int() != 200 {
		return errors.New("response code wrong")
	}
	utcTime, err := time.UtcToLocal(gjson.Get(result, "game.gameTimeUTC").String(),
		"2006-01-02T15:04:05Z",
		constant.EST_ZONE)
	if err != nil {
		log.Error(err)
		return err
	}
	*board = SummaryBoard{
		GameId:         board.GameId,
		GameTime:       utcTime,
		AreanCity:      gjson.Get(result, "game.arena.arenaCity").String(),
		AreanName:      gjson.Get(result, "game.arena.arenaName").String(),
		HomeTeamCode:   gjson.Get(result, "game.homeTeam.teamTricode").String(),
		AwayTeamCode:   gjson.Get(result, "game.awayTeam.teamTricode").String(),
		HomeTeamScores: getPeriodsScore(gjson.Get(result, "game.homeTeam.periods")),
		AwayTeamScores: getPeriodsScore(gjson.Get(result, "game.awayTeam.periods")),
		HomeTeamStats: stats.Stats{
			Assits:    int(gjson.Get(result, "game.homeTeam.statistics.assists").Int()),
			Blocks:    int(gjson.Get(result, "game.homeTeam.statistics.blocks").Int()),
			Steals:    int(gjson.Get(result, "game.homeTeam.statistics.steals").Int()),
			Turnovers: int(gjson.Get(result, "game.homeTeam.statistics.turnovers").Int()),
			Rebounds:  int(gjson.Get(result, "game.homeTeam.statistics.reboundsTotal").Int()),
		},
		AwayTeamStats: stats.Stats{
			Assits:    int(gjson.Get(result, "game.awayTeam.statistics.assists").Int()),
			Blocks:    int(gjson.Get(result, "game.awayTeam.statistics.blocks").Int()),
			Steals:    int(gjson.Get(result, "game.awayTeam.statistics.steals").Int()),
			Turnovers: int(gjson.Get(result, "game.awayTeam.statistics.turnovers").Int()),
			Rebounds:  int(gjson.Get(result, "game.awayTeam.statistics.reboundsTotal").Int()),
		},
	}
	if len(board.HomeTeamScores) != len(board.AwayTeamScores) {
		return errors.New("score wrong")
	}
	return nil
}

func (board *SummaryBoard) PrintTable() error {
	fmt.Printf("\n🗓  %s | 🏟  %s %s\n", board.GameTime, board.AreanCity, board.AreanName)
	var headerRow, homeTeamRow, awayTeamRow table.Row
	headerRow = append(headerRow, "")
	homeTeamRow = append(homeTeamRow, board.HomeTeamCode)
	awayTeamRow = append(awayTeamRow, board.AwayTeamCode)
	var homePointTotal, awayPointTotal int
	for i := 0; i < len(board.HomeTeamScores); i++ {
		var quarter string
		if i+1 <= 4 {
			quarter = fmt.Sprintf("Q%d", i+1)
		}
		if i+1 > 4 {
			quarter = fmt.Sprintf("OT%d", i-4)
		}
		homePointTotal += board.HomeTeamScores[i]
		awayPointTotal += board.AwayTeamScores[i]
		headerRow = append(headerRow, quarter)
		s1, s2, err := print.MaxTextColor(strconv.Itoa(board.HomeTeamScores[i]), strconv.Itoa(board.AwayTeamScores[i]))
		if err != nil {
			return err
		}
		homeTeamRow = append(homeTeamRow, s1)
		awayTeamRow = append(awayTeamRow, s2)
	}
	headerRow = append(headerRow, "Total")
	headerRow = append(headerRow, statsRows...)
	hpt, apt, err := print.MaxTextColor(strconv.Itoa(homePointTotal), strconv.Itoa(awayPointTotal))
	if err != nil {
		return err
	}
	hts, ats, err := print.MaxTextColor(strconv.Itoa(board.HomeTeamStats.Assits),
		strconv.Itoa(board.AwayTeamStats.Assits))
	if err != nil {
		return err
	}
	htb, atb, err := print.MaxTextColor(strconv.Itoa(board.HomeTeamStats.Blocks),
		strconv.Itoa(board.AwayTeamStats.Blocks))
	if err != nil {
		return err
	}
	htss, atss, err := print.MaxTextColor(strconv.Itoa(board.HomeTeamStats.Steals),
		strconv.Itoa(board.AwayTeamStats.Steals))
	if err != nil {
		return err
	}
	htt, att, err := print.MaxTextColor(strconv.Itoa(board.HomeTeamStats.Turnovers),
		strconv.Itoa(board.AwayTeamStats.Turnovers))
	if err != nil {
		return err
	}
	htr, atr, err := print.MaxTextColor(strconv.Itoa(board.HomeTeamStats.Rebounds),
		strconv.Itoa(board.AwayTeamStats.Rebounds))
	if err != nil {
		return err
	}
	homeTeamRow = append(homeTeamRow, hpt, "/", hts, htb, htss, htt, htr)
	awayTeamRow = append(awayTeamRow, apt, "/", ats, atb, atss, att, atr)

	var configRows []table.ColumnConfig
	for i := 1; i <= len(headerRow); i++ {
		configRows = append(configRows, table.ColumnConfig{Number: i, Align: text.AlignCenter})
	}

	t := table.NewWriter()
	t.AppendHeader(headerRow)
	t.AppendRow(homeTeamRow)
	t.AppendRow(awayTeamRow)
	t.SetColumnConfigs(configRows)
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleBold)
	t.Style().Options.SeparateRows = true
	t.Render()
	return nil
}

func getPeriodsScore(result gjson.Result) []int {
	scores := make([]int, 0)
	result.ForEach(func(key, value gjson.Result) bool {
		scores = append(scores, int(value.Get("score").Int()))
		return true
	})
	return scores
}
