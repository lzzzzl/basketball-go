package teamranking

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lzzzzl/basketball-go/modules/constant"
	"github.com/lzzzzl/basketball-go/modules/errno"
	"github.com/lzzzzl/basketball-go/modules/log"
	"github.com/lzzzzl/basketball-go/modules/print"
)

// HEADER ...
var HEADER = table.Row{"TEAM", "W", "L", "HOME", "AWAY", "PPG", "OPP PPG", "DIFF", "STRK", "L10"}

// TeamRanking ...
type TeamRanking struct {
}

// Details ...
type Details struct {
	Name   string
	Win    string
	Loss   string
	Home   string
	Away   string
	PPG    string
	OPPPPG string
	Diff   string
	STRK   string
	L10    string
}

// TeamRankingPrinter ...
func (t *TeamRanking) TeamRankingPrinter(year string) error {
	url := fmt.Sprintf(constant.TEAM_RANKING_URL, year)
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
	details, err := t.ResultParser(res)
	if err != nil {
		return err
	}
	err = t.PrintTable(details)
	if err != nil {
		return err
	}

	return nil
}

// ResultParser ...
func (t *TeamRanking) ResultParser(result *http.Response) (map[string][]*Details, error) {
	trds := make(map[string][]*Details)

	doc, err := goquery.NewDocumentFromReader(result.Body)
	if err != nil {
		return trds, err
	}

	doc.Find("div[class*='ResponsiveTable']").Each(func(_ int, s *goquery.Selection) {
		teams := make([]string, 0)
		details := make([]*Details, 0)

		s.Find("div[class*='flex']").Each(func(_ int, m *goquery.Selection) {
			m.Find("tbody[class*='Table__TBODY']").Each(func(_ int, n *goquery.Selection) {
				n.Find("abbr").Each(func(_ int, t *goquery.Selection) {
					teams = append(teams, t.Text())
				})
			})

			m.Find("div[class*='Table__Scroller']").Each(func(_ int, n *goquery.Selection) {
				n.Find("tr[class*='Table__TR']").Each(func(i int, t *goquery.Selection) {
					if i > 0 {
						data := make([]string, 0)
						t.Find("span").Each(func(i int, q *goquery.Selection) {
							data = append(data, q.Text())
						})
						details = append(details, &Details{
							Win:    data[0],
							Loss:   data[1],
							Home:   data[4],
							Away:   data[5],
							PPG:    data[8],
							OPPPPG: data[9],
							Diff:   data[10],
							STRK:   data[11],
							L10:    data[12],
						})
					}
				})
			})
		})

		conference := s.Find("div[class*='Table__Title']").Text()
		for i, t := range teams {
			details[i].Name = t
			trds[conference] = append(trds[conference], details[i])
		}
	})

	return trds, nil
}

// PrintTable ...
func (t *TeamRanking) PrintTable(result map[string][]*Details) error {
	for conf, details := range result {
		fmt.Println("\n" + strings.ToUpper(conf))

		tb := table.NewWriter()
		tb.SetOutputMirror(os.Stdout)
		tb.AppendHeader(HEADER)

		for i, d := range details {
			var diff, strk string
			if strings.HasPrefix(d.Diff, "+") {
				diff = print.RedString(d.Diff)
			} else {
				diff = print.GreenString(d.Diff)
			}
			if strings.HasPrefix(d.STRK, "W") {
				strk = print.RedString(d.STRK)
			} else {
				strk = print.GreenString(d.STRK)
			}
			row := table.Row{
				fmt.Sprintf("%d %s", i+1, d.Name),
				d.Win,
				d.Loss,
				d.Home,
				d.Away,
				d.PPG,
				d.OPPPPG,
				diff,
				strk,
				d.L10,
			}
			tb.AppendRow(row)
		}
		tb.SetStyle(table.StyleBold)
		tb.Style().Options.SeparateRows = true
		tb.Render()

		fmt.Println()
	}

	return nil
}
