package cmd

import (
	"fmt"

	"github.com/lzzzzl/basketball-go/handlers/nba/board"
	"github.com/lzzzzl/basketball-go/handlers/nba/game"
	"github.com/lzzzzl/basketball-go/handlers/nba/game/schedule"
	"github.com/lzzzzl/basketball-go/handlers/nba/game/today"
	"github.com/lzzzzl/basketball-go/handlers/nba/playoff"
	playerranking "github.com/lzzzzl/basketball-go/handlers/nba/ranking/player_ranking"
	teamranking "github.com/lzzzzl/basketball-go/handlers/nba/ranking/team_ranking"
	"github.com/lzzzzl/basketball-go/handlers/nba/team"
	"github.com/lzzzzl/basketball-go/modules/log"
	"github.com/lzzzzl/basketball-go/modules/time"
	"github.com/spf13/cobra"
)

var (
	gameDay             string
	isYesterday         bool
	isToday             bool
	isTomorrow          bool
	beforeDays          int
	nextDays            int
	playoffYear         int
	playoffScheduleYear int
	isTeamList          bool
)

var rootCmd = &cobra.Command{
	Use:   "basketball-go",
	Short: "A Fast and Flexible 🏀 CLI",
	// Long:  `A Fast and Flexible 🏀 CLI`,
}

// PickUpGameDetails pickup game details
func PickUpGameDetails(boards []*game.GameBoard) {
	if len(boards) <= 0 {
		fmt.Println("no game exists...")
		return
	}
	var index int
	fmt.Println("chose game id # e.g. 1、2、3...")
	fmt.Scanf("%d", &index)
	if index <= 0 || index > len(boards) {
		fmt.Println("something wrong...")
		return
	}
	var choice string
	fmt.Println("chose game summary or players stats you want to see # e.g. g/p")
	fmt.Scanf("%s\n", &choice)
	switch choice {
	case "g":
		s := board.SummaryBoard{GameID: boards[index-1].GameID}
		err := s.SummaryBoardPrinter()
		if err != nil {
			fmt.Println(err)
			return
		}
	case "p":
		p := board.PlayerBoard{GameID: boards[index-1].GameID}
		err := p.PlayerBoardPrinter()
		if err != nil {
			fmt.Println(err)
			return
		}
	default:
		fmt.Println("something wrong...")
		return
	}
}

var scheduleCmd = &cobra.Command{
	Use: "game",
	Run: func(cmd *cobra.Command, args []string) {
		switch {
		case isToday:
			today := &today.Today{}
			boards, err := today.TodayBoardPrinter()
			if err != nil {
				fmt.Println("something went wrong...")
				log.Error(err)
				return
			}
			PickUpGameDetails(boards)
		case isYesterday:
			schedule := &schedule.Schedule{
				GameStartTime: time.GetPlusDate("20060102", -1),
				GameEndTime:   time.GetPlusDate("20060102", -1),
			}
			boards, err := schedule.ScheduleBoardPrinter()
			if err != nil {
				fmt.Println("something went wrong...")
				log.Error(err)
				return
			}
			PickUpGameDetails(boards)
		case isTomorrow:
			schedule := &schedule.Schedule{
				GameStartTime: time.GetPlusDate("20060102", 1),
				GameEndTime:   time.GetPlusDate("20060102", 1),
			}
			_, err := schedule.ScheduleBoardPrinter()
			if err != nil {
				fmt.Println("something went wrong...")
				log.Error(err)
				return
			}
		case gameDay != "":
			t, err := time.Str2Time(gameDay, "2006/1/2", "20060102")
			if err != nil {
				fmt.Println("time format is wrong")
				return
			}
			schedule := &schedule.Schedule{
				GameStartTime: t,
				GameEndTime:   t,
			}
			boards, err := schedule.ScheduleBoardPrinter()
			if err != nil {
				fmt.Println("something went wrong...")
				log.Error(err)
				return
			}
			PickUpGameDetails(boards)
		case beforeDays > 0 && beforeDays <= 30:
			schedule := &schedule.Schedule{
				GameStartTime: time.GetPlusDate("20060102", -beforeDays),
				GameEndTime:   time.GetPlusDate("20060102", -1),
			}
			_, err := schedule.ScheduleBoardPrinter()
			if err != nil {
				fmt.Println("something went wrong...")
				log.Error(err)
				return
			}
		case nextDays > 0 && nextDays <= 30:
			schedule := &schedule.Schedule{
				GameStartTime: time.GetPlusDate("20060102", 1),
				GameEndTime:   time.GetPlusDate("20060102", nextDays),
			}
			_, err := schedule.ScheduleBoardPrinter()
			if err != nil {
				fmt.Println("something went wrong...")
				log.Error(err)
				return
			}
		default:
			fmt.Println("nothing happen...")
		}
	},
}

var playOffCmd = &cobra.Command{
	Use: "playoff",
	Run: func(cmd *cobra.Command, args []string) {
		switch {
		case playoffYear > 0:
			playoff := &playoff.Playoff{Year: playoffYear}
			if err := playoff.PlayoffBracketPrinter(); err != nil {
				fmt.Println("playoff went wrong...")
			}
		case playoffScheduleYear > 0:
			p := &playoff.PlayOffSchedule{Year: 2022}
			if err := p.PlayOffSchedulePrinter(); err != nil {
				fmt.Println("playoff went wrong")
			}
		}
	},
}

var teamCmd = &cobra.Command{
	Use: "team",
	Run: func(cmd *cobra.Command, args []string) {
		switch {
		case isTeamList:
			t := &team.Team{}
			if err := t.TeamListPrinter(); err != nil {
				fmt.Println("team went wrong")
			}
		}
	},
}

// ####################### rank #######################
var rankCmd = &cobra.Command{
	Use:   "rank",
	Short: "ranking team & player statictis",
}

var rankTeamCmd = &cobra.Command{
	Use: "team",
	Run: func(cmd *cobra.Command, args []string) {
		year := time.GetPlusYear(1)
		if len(args) > 0 {
			year = args[0]
		}
		r := &teamranking.TeamRanking{}
		if err := r.TeamRankingPrinter(year); err != nil {
			fmt.Println("team ranking wrong")
		}
	},
}

var rankPlayerCmd = &cobra.Command{
	Use: "player",
	Run: func(cmd *cobra.Command, args []string) {
		year := time.GetPlusYear(1)
		if len(args) >= 1 {
			year = args[0]
		}
		_type := "2"
		if len(args) >= 2 {
			_type = args[1]
		}
		r := &playerranking.PlayerRanking{}
		if err := r.PlayerRankingPrinter(year, _type); err != nil {
			fmt.Println("player ranking wrong")
		}
	},
}

func init() {
	scheduleCmd.PersistentFlags().BoolVarP(&isYesterday, "yesterday", "y", false, "yesterday game schedule")
	scheduleCmd.PersistentFlags().BoolVarP(&isToday, "today", "t", false, "today game schedule")
	scheduleCmd.PersistentFlags().BoolVarP(&isTomorrow, "tomorrow", "T", false, "tomorrow game schedule")
	scheduleCmd.PersistentFlags().StringVarP(&gameDay, "day", "d", "", "input game day eg. year/month/day")
	scheduleCmd.PersistentFlags().IntVarP(&beforeDays, "before", "b", 0, "before game day eg. 30")
	scheduleCmd.PersistentFlags().IntVarP(&nextDays, "next", "n", 0, "next game day eg. 30")

	playOffCmd.PersistentFlags().IntVarP(&playoffYear, "bracket", "b", 0, "input year playoff eg. 2022")
	playOffCmd.PersistentFlags().IntVarP(&playoffScheduleYear, "schedule", "s", 0, "input year playoff eg. 2022")

	teamCmd.PersistentFlags().BoolVarP(&isTeamList, "list", "l", false, "list teams")

	rankTeamCmd.PersistentFlags().StringP("year", "y", "", "input year eg. 2022")
	rankPlayerCmd.PersistentFlags().StringP("year", "y", "", "input year eg. 2022")
	rankPlayerCmd.PersistentFlags().StringP("type", "t", "", "input season type, eg. 2(regular season) / 3(postseasonn)")
}

// Execute add sub commands and execute root command
func Execute() error {
	rootCmd.AddCommand(
		scheduleCmd,
		playOffCmd,
		teamCmd,
		rankCmd,
	)
	rankCmd.AddCommand(rankTeamCmd, rankPlayerCmd)
	return rootCmd.Execute()
}
