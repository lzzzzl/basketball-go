package playerranking

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lzzzzl/basketball-go/modules/constant"
	"github.com/lzzzzl/basketball-go/modules/errno"
	"github.com/lzzzzl/basketball-go/modules/print"
)

// PlayerRanking ...
type PlayerRanking struct {
}

// PlayerRankingData ...
type PlayerRankingData struct {
	Offense []*Statistics
	Defense []*Statistics
}

// Statistics ...
type Statistics struct {
	Type      string
	TypeShort string
	Details   []*Detail
}

// Details ...
type Detail struct {
	Index int
	Team  string
	Name  string
	Value string
}

// PlayerRankingPrinter ...
func (p *PlayerRanking) PlayerRankingPrinter(year, seasonType string) error {
	url := fmt.Sprintf(constant.PLAYER_RANKING_URL, year, seasonType)
	fmt.Println(url)
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
	result, err := p.ResultParser(res)
	if err != nil {
		return err
	}
	if err := p.PrintTable(*result); err != nil {
		return err
	}

	return nil
}

// ResultParser ...
func (p *PlayerRanking) ResultParser(result *http.Response) (*PlayerRankingData, error) {
	prds := &PlayerRankingData{}

	doc, err := goquery.NewDocumentFromReader(result.Body)
	if err != nil {
		return prds, err
	}

	doc.Find("div[class*='mb1']").Each(func(_ int, s *goquery.Selection) {
		title := s.Find("div[class*='Table__Title']").Text()
		stats := make([]*Statistics, 0)

		s.Find("table").Each(func(_ int, m *goquery.Selection) {
			details := make([]*Detail, 0)

			m.Find("tbody[class*='Table__TBODY']").Each(func(_ int, n *goquery.Selection) {
				n.Find("tr[class*='Table__TR']").Each(func(i int, t *goquery.Selection) {
					team := t.Find("span[class*='team']").Text()
					name, _ := t.Find("img").Attr("title")
					val := t.Find("td[class*='Table__TD']").Last().Text()

					if name != "" {
						details = append(details, &Detail{
							Index: i + 1,
							Team:  team,
							Name:  name,
							Value: val,
						})
					}
				})
			})

			tp := m.Find("th[class*='Table__TH']").First().Text()
			tp_short := m.Find("th[class*='Table__TH']").Last().Text()

			if strings.Contains(tp, "3-Pointers") {
				tp = "3PT Made"
			}

			stat := &Statistics{
				Type:      tp,
				TypeShort: tp_short,
				Details:   details,
			}
			stats = append(stats, stat)

		})

		if strings.Contains(title, "Offensive") {
			prds.Offense = stats
		} else {
			prds.Defense = stats
		}

	})

	return prds, nil
}

func (p *PlayerRanking) PrintTable(result PlayerRankingData) error {
	if err := printer("Offensive Leaders", result.Offense); err != nil {
		return err
	}
	if err := printer("Defensive Leaders", result.Defense); err != nil {
		return err
	}
	return nil
}

func printer(title string, statistics []*Statistics) error {
	fmt.Println(print.RedString(title))

	tb := table.NewWriter()
	tb.SetOutputMirror(os.Stdout)
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}

	for _, s := range statistics {
		tb.AppendRow(table.Row{"", s.Type, s.Type, s.TypeShort}, rowConfigAutoMerge)
		for _, d := range s.Details {
			tb.AppendRow(table.Row{d.Index, d.Name, d.Team, d.Value})
		}
	}

	tb.SetStyle(table.StyleBold)
	tb.Style().Options.SeparateRows = true
	tb.Render()

	fmt.Println()

	return nil
}
