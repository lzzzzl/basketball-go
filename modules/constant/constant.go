package constant

var LIVE_HEADER = map[string]string{
	"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	"Accept-Encoding": "",
	"Accept-Language": "en-US,en;q=0.9",
	"Cache-Control":   "max-age=0",
	"Connection":      "keep-alive",
	"Host":            "cdn.nba.com",
	"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36",
}

const (
	// BOARD
	GAME_DETAILS_URL            = "https://cdn.nba.com/static/json/liveData/boxscore/boxscore_%s.json"
	CURRENT_SEASON_SCHEDULE_URL = "https://cdn.nba.com/static/json/staticData/scheduleLeagueV2_1.json"
	TODAY_SCORE_BOARD_URL       = "https://cdn.nba.com/static/json/liveData/scoreboard/todaysScoreboard_00.json"

	// PLAYOFF
	PLAYOFF_BRACKET_URL  = "https://stats.nba.com/stats/playoffbracket?LeagueID=00&SeasonYear=%d&State=2"
	PLAYOFF_SCHEDULE_URL = "https://www.nba.com/news/%d-nba-playoffs-schedule"

	// TEAM
	TEAM_URL         = "https://www.nba.com/teams"
	TEAM_RANKING_URL = "https://www.espn.com/nba/standings/_/season/%s"

	// ZONE
	EST_ZONE = "America/New_York"
)
