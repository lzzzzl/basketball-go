package team

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/lzzzzl/basketball-go/modules/constant"
	"github.com/lzzzzl/basketball-go/modules/log"
)

type Team struct{}

func (t *Team) TeamListPrinter() error {
	url := constant.TEAM_URL
	log.Println("Team URL: ", url)
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return errors.New(fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status))
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	tMap := make(map[int][]string)
	headerRow := table.Row{}
	tb := table.NewWriter()
	tb.SetOutputMirror(os.Stdout)
	doc.Find("div[class*='TeamDivisions_wrapper']").Each(func(i int, s *goquery.Selection) {
		s.Find("div[class*='TeamDivisions_division']").Each(func(i int, m *goquery.Selection) {
			distinct := m.Find("div[class*='TeamDivisions_divisionName']").Text()
			if distinct != "" {
				headerRow = append(headerRow, text.FgGreen.Sprintf("%10s", strings.ToUpper(distinct)))
			}
			m.Find("div[class*='TeamDivisions_divisionTeams']").Each(func(i int, n *goquery.Selection) {
				n.Find("a[class*='TeamFigure_tfMainLink']").Each(func(i int, k *goquery.Selection) {
					if k.Text() != "" {
						if _, ok := tMap[i]; !ok {
							tMap[i] = make([]string, 0)
						}
						tMap[i] = append(tMap[i], k.Text())
					}
				})
			})
		})
	})
	tb.AppendRow(headerRow)
	for _, m := range tMap {
		row := table.Row{}
		for _, n := range m {
			row = append(row, n)
		}
		tb.AppendRow(row)
	}
	var configRows []table.ColumnConfig
	for i := 1; i <= len(headerRow); i++ {
		configRows = append(configRows, table.ColumnConfig{Number: i, Align: text.AlignCenter})
	}
	tb.SetColumnConfigs(configRows)
	tb.SetStyle(table.StyleBold)
	tb.Style().Options.SeparateRows = true
	tb.Render()

	return nil
}
