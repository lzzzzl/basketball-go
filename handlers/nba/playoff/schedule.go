package playoff

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/lzzzzl/basketball-go/modules/constant"
	"github.com/lzzzzl/basketball-go/modules/log"
	"github.com/lzzzzl/basketball-go/modules/print"

	"github.com/PuerkitoBio/goquery"
)

// PlayOffSchedule playoff shedule
type PlayOffSchedule struct {
	Year int
}

// PlayOffSchedulePrinter playoff shedule printer
func (p *PlayOffSchedule) PlayOffSchedulePrinter() error {
	url := fmt.Sprintf(constant.PLAYOFF_SCHEDULE_URL, p.Year)
	log.Println("Playoff Schedule URL: ", url)
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
	doc.Find("div[class*='article']").Each(func(i int, s *goquery.Selection) {
		content := s.Text()
		cnts := strings.Split(content, "\n")
		for _, s := range cnts[8:] {
			if strings.Contains(s, "Conference finals:") || strings.Contains(s, "Conference semifinals:") ||
				strings.Contains(s, "First-round series:") {
				print.FigurePrint(s, "", "green")
				fmt.Println()
			} else {
				fmt.Println(s)
			}
		}
	})
	return nil
}
