package board

import (
	"errors"
	"fmt"
	"math"
	"os"
	"reflect"
	"regexp"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/lzzzzl/basketball-go/handlers/nba/stats"
	"github.com/lzzzzl/basketball-go/modules/constant"
	"github.com/lzzzzl/basketball-go/modules/log"
	"github.com/lzzzzl/basketball-go/modules/request"
	"github.com/tidwall/gjson"
)

var (
	statsList = []string{"player", "min", "fgm", "fga", "fg%", "3pm", "3pa", "3p%", "ftm",
		"fta", "ft%", "reb", "ast", "stl", "blk", "to", "pts", "+/-"}
	statsMap = map[string]string{
		"player": "Name",
		"min":    "Minutes",
		"fgm":    "FieldGoalsMade",
		"fga":    "FieldGoalsAttempted",
		"fg%":    "FieldGoalsPercentage",
		"3pm":    "ThreePointersMade",
		"3pa":    "ThreePointersAttempted",
		"3p%":    "ThreePointersPercentage",
		"ftm":    "FreeThrowsMade",
		"fta":    "FreeThrowsAttempted",
		"ft%":    "FreeThrowsPercentage",
		"reb":    "Rebounds",
		"ast":    "Assits",
		"stl":    "Steals",
		"blk":    "Blocks",
		"to":     "Turnovers",
		"pts":    "Points",
		"+/-":    "PlusMinusPoints",
	}
)

type PlayerBoard struct {
	GameId          string
	HomeTeamName    string
	HomeTeamPlayers []*stats.PlayerStats
	AwayTeamName    string
	AwayTeamPlayers []*stats.PlayerStats
}

func (board *PlayerBoard) PlayerBoardPrinter() error {
	url := fmt.Sprintf(constant.GAME_DETAILS_URL, board.GameId)
	log.Println("Player Board URL: ", url)
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

func (board *PlayerBoard) ResultParser(result string) error {
	if gjson.Get(result, "meta.code").Int() != 200 {
		return errors.New("response code wrong")
	}
	*board = PlayerBoard{
		GameId: board.GameId,
		HomeTeamName: fmt.Sprintf("%s %s", gjson.Get(result, "game.homeTeam.teamCity").String(),
			gjson.Get(result, "game.homeTeam.teamName").String()),
		AwayTeamName: fmt.Sprintf("%s %s", gjson.Get(result, "game.awayTeam.teamCity").String(),
			gjson.Get(result, "game.awayTeam.teamName").String()),
		HomeTeamPlayers: getPlayerStats(gjson.Get(result, "game.homeTeam.players")),
		AwayTeamPlayers: getPlayerStats(gjson.Get(result, "game.awayTeam.players")),
	}
	return nil
}

func (board *PlayerBoard) PrintTable() error {
	t1 := table.NewWriter()
	t2 := table.NewWriter()

	var tableRow table.Row
	var configRows []table.ColumnConfig
	for i, v := range statsList {
		tableRow = append(tableRow, v)
		configRows = append(configRows, table.ColumnConfig{Number: i + 1, Align: text.AlignCenter})
	}
	t1.AppendRow(tableRow)
	t2.AppendRow(tableRow)

	hr, err := getPlayersRows(board.HomeTeamPlayers)
	if err != nil {
		return err
	}
	ar, err := getPlayersRows(board.AwayTeamPlayers)
	if err != nil {
		return err
	}
	for _, v := range hr {
		t1.AppendRow(*v)
	}
	for _, v := range ar {
		t2.AppendRow(*v)
	}

	fmt.Println()
	fmt.Println(text.FgGreen.Sprint(board.HomeTeamName))
	t1.SetStyle(table.StyleBold)
	t1.SetColumnConfigs(configRows)
	t1.SetOutputMirror(os.Stdout)
	t1.Style().Options.SeparateRows = true
	t1.Render()

	fmt.Println(text.FgGreen.Sprint(board.AwayTeamName))
	t2.SetStyle(table.StyleBold)
	t2.SetColumnConfigs(configRows)
	t2.SetOutputMirror(os.Stdout)
	t2.Style().Options.SeparateRows = true
	t2.Render()

	return nil
}

func getPlayerStats(result gjson.Result) []*stats.PlayerStats {
	playerStats := make([]*stats.PlayerStats, 0)
	result.ForEach(func(key, value gjson.Result) bool {
		playerStats = append(playerStats, &stats.PlayerStats{
			Name: fmt.Sprintf("%s %s %s", value.Get("firstName").String(),
				value.Get("familyName").String(), value.Get("position")),
			Position:                value.Get("position").String(),
			Status:                  value.Get("status").String(),
			Minutes:                 value.Get("statistics.minutes").String(),
			IsStarter:               isStatrer(value.Get("starter").String()),
			JerseyNum:               value.Get("jerseyNum").String(),
			Assits:                  int(value.Get("statistics.assists").Int()),
			Points:                  int(value.Get("statistics.points").Int()),
			Blocks:                  int(value.Get("statistics.blocks").Int()),
			Steals:                  int(value.Get("statistics.steals").Int()),
			Turnovers:               int(value.Get("statistics.turnovers").Int()),
			Rebounds:                int(value.Get("statistics.reboundsTotal").Int()),
			PlusMinusPoints:         value.Get("statistics.plusMinusPoints").String(),
			FieldGoalsAttempted:     int(value.Get("statistics.fieldGoalsAttempted").Int()),
			FieldGoalsMade:          int(value.Get("statistics.fieldGoalsMade").Int()),
			FieldGoalsPercentage:    value.Get("statistics.fieldGoalsPercentage").Float(),
			ThreePointersAttempted:  int(value.Get("statistics.threePointersAttempted").Int()),
			ThreePointersMade:       int(value.Get("statistics.threePointersMade").Int()),
			ThreePointersPercentage: value.Get("statistics.threePointersPercentage").Float(),
			FreeThrowsAttempted:     int(value.Get("statistics.freeThrowsAttempted").Int()),
			FreeThrowsMade:          int(value.Get("statistics.freeThrowsMade").Int()),
			FreeThrowsPercentage:    value.Get("statistics.freeThrowsPercentage").Float(),
		})
		return true
	})
	return playerStats
}

func isStatrer(str string) bool {
	return str == "1"
}

func getPlayersRows(stats []*stats.PlayerStats) ([]*table.Row, error) {
	statsRows := make([]*table.Row, 0)
	for _, s := range stats {
		i := reflect.ValueOf(*s)
		var row table.Row
		for _, v := range statsList {
			switch v {
			case "min":
				min, err := parseMinutes(s.Minutes)
				if err != nil {
					return statsRows, err
				}
				row = append(row, min)
			case "player":
				name := fmt.Sprintf("%25s", s.Name)
				row = append(row, name)
			case "fg%":
				fgp := math.Floor(s.FieldGoalsPercentage * 100)
				row = append(row, fgp)
			case "3p%":
				tpg := math.Floor(s.ThreePointersPercentage * 100)
				row = append(row, tpg)
			case "ft%":
				ftp := math.Floor(s.FreeThrowsPercentage * 100)
				row = append(row, ftp)
			default:
				row = append(row, i.FieldByName(statsMap[v]))
			}
		}
		statsRows = append(statsRows, &row)
	}

	return statsRows, nil
}

func parseMinutes(str string) (min string, err error) {
	r := regexp.MustCompile("([0-9]+)")
	list := r.FindAllString(str, -1)
	if len(list) < 2 {
		return min, fmt.Errorf("minutes formation wrong")
	}
	min = fmt.Sprintf("%s:%s", list[0], list[1])

	return
}
