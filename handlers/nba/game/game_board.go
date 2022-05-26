package game

//
var (
	ColTitleIndex  = "#"
	ColTitleTeams  = "Home  🏀  Away"
	ColTitleScores = "Score"
	ColTitleStatus = "Status"
	ColTitleSeries = "Series"
)

// GameBoard game board
type GameBoard struct {
	GameID        string
	GameStartTime string
	GameStatus    string
	SeriesText    string
	HomeTeamName  string
	HomeTeamCity  string
	HomeTeamCode  string
	HomeTeamScore string
	AwayTeamName  string
	AwayTeamCity  string
	AwayTeamCode  string
	AwayTeamScore string
}
