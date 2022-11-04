package teamranking

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/lzzzzl/basketball-go/modules/constant"
	"github.com/lzzzzl/basketball-go/modules/errno"
	"github.com/lzzzzl/basketball-go/modules/log"
)

type TeamRanking struct {
	Season string
}

type TeamRankingDetails struct {
	TeamName string
	Win      int32
	Loss     int32
}

func (t *TeamRanking) TeamRankingPrinter() error {
	url := fmt.Sprintf(constant.TEAM_RANKING_URL, t.Season)
	url = "https://www.espn.com/nba/standings"
	log.Println("Team Ranking URL: ", url)
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != errno.OK.Code {
		return &errno.Err{
			Code:    res.StatusCode,
			Message: res.Status,
		}
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}
	trds := make(map[string]map[string]*TeamRankingDetails)
	doc.Find("div[class*='ResponsiveTable']").Each(func(_ int, s *goquery.Selection) {
		conference := s.Find("div[class*='Table__Title']").Text()
		trds[conference] = make(map[string]*TeamRankingDetails)

		s.Find("div[class*='flex']").Each(func(_ int, m *goquery.Selection) {

			// team name
			m.Find("tbody[class*='Table__TBODY']").Each(func(_ int, n *goquery.Selection) {
				n.Find("abbr").Each(func(_ int, t *goquery.Selection) {
					team := t.Text()
					trds[conference][team] = &TeamRankingDetails{
						TeamName: team,
					}
				})
			})

			// ranking details
			m.Find("div[class*='Table__Scroller']").Each(func(_ int, n *goquery.Selection) {
				fmt.Println("=================")
				n.Find("tr[class*='Table__TR']").Each(func(_ int, t *goquery.Selection) {
					data := make([]string, 0)
					t.Find("span").Each(func(i int, q *goquery.Selection) {
						data = append(data, q.Text())
					})
					fmt.Println("data: ", data)
				})
			})
		})

	})
	return nil
}
