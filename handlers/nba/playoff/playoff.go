package playoff

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/lzzzzl/basketball-go/modules/constant"
	"github.com/lzzzzl/basketball-go/modules/log"
	"github.com/lzzzzl/basketball-go/modules/print"
	"github.com/lzzzzl/basketball-go/modules/request"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/tidwall/gjson"
)

// MAXROUND ...
const MAXROUND = 3

// Playoff ...
type Playoff struct {
	Year int
}

// Bracket Playoff bracket
type Bracket struct {
	RoundNumber    int
	SeriesNumber   int
	Conference     string
	SeriesText     string
	HighSeedCode   string
	HighSeedRank   int
	LowSeedCode    string
	LowSeedRank    int
	DisplayOrder   int
	HighSeedID     string
	DisplayTopTeam string
	SeriesWinner   string
}

// SeriesBrackets series brackets
type SeriesBrackets struct {
	WestConf map[int]map[int]*Bracket
	EastConf map[int]map[int]*Bracket
	Final    Bracket
}

// PlayoffBracketPrinter bracket printer
func (p *Playoff) PlayoffBracketPrinter() error {
	url := fmt.Sprintf(constant.PLAYOFF_BRACKET_URL, p.Year-1)
	log.Println("Playoff Bracket URL: ", url)
	request := &request.HttpRequest{
		Url:     url,
		Headers: constant.LIVE_HEADER,
		TimeOut: 30,
	}
	result, err := request.HttpGet()
	if err != nil {
		return err
	}
	seriesBrackets, err := p.ResultParser(result)
	if err != nil {
		return err
	}
	err = p.PrintTable(seriesBrackets)
	if err != nil {
		return err
	}
	return nil
}

// ResultParser http result parser
func (p *Playoff) ResultParser(result string) (*SeriesBrackets, error) {
	finals := &Bracket{}
	westMap := make(map[int]map[int]*Bracket)
	eastMap := make(map[int]map[int]*Bracket)
	for i := 1; i <= MAXROUND; i++ {
		westMap[i] = make(map[int]*Bracket)
		eastMap[i] = make(map[int]*Bracket)
	}

	series := gjson.Get(result, "bracket.playoffBracketSeries")
	series.ForEach(func(key, value gjson.Result) bool {
		c := value.Get("seriesConference").String()
		r := int(value.Get("roundNumber").Int())
		d := int(value.Get("displayOrderNumber").Int())
		bt := &Bracket{
			RoundNumber:    int(value.Get("roundNumber").Int()),
			SeriesNumber:   int(value.Get("seriesNumber").Int()),
			Conference:     value.Get("seriesConference").String(),
			SeriesText:     value.Get("seriesText").String(),
			HighSeedCode:   value.Get("highSeedTricode").String(),
			HighSeedRank:   int(value.Get("highSeedRank").Int()),
			LowSeedCode:    value.Get("lowSeedTricode").String(),
			LowSeedRank:    int(value.Get("lowSeedRank").Int()),
			DisplayOrder:   int(value.Get("displayOrderNumber").Int()),
			HighSeedID:     value.Get("highSeedID").String(),
			DisplayTopTeam: value.Get("displayTopTeam").String(),
			SeriesWinner:   value.Get("SeriesWinner").String(),
		}
		switch c {
		case "West":
			hc, lc := addWinOrGoHomeEmojiFromBracket(bt, "🔥", "🏠")
			bt.HighSeedCode = hc
			bt.LowSeedCode = lc
			westMap[r][d] = bt
		case "East":
			hc, lc := addWinOrGoHomeEmojiFromBracket(bt, "🔥", "🏠")
			bt.HighSeedCode = hc
			bt.LowSeedCode = lc
			eastMap[r][d] = bt
		default:
			hc, lc := addWinOrGoHomeEmojiFromBracket(bt, "🏆", "🏠")
			bt.HighSeedCode = hc
			bt.LowSeedCode = lc
			finals = bt
		}
		return true
	})
	seriesBrackets := &SeriesBrackets{
		WestConf: westMap,
		EastConf: eastMap,
		Final:    *finals,
	}
	return seriesBrackets, nil
}

// ConfBracket conference bracker
type ConfBracket struct {
	RoundNumber int
	bracketMap  map[int]*Bracket
}

// PrintTable print playoff bracket table
func (p *Playoff) PrintTable(brackets *SeriesBrackets) error {
	if len(brackets.WestConf) == 0 || len(brackets.EastConf) == 0 {
		return errors.New("playoff not esist")
	}
	w := fmt.Sprintf("%d NBA Finals", p.Year)
	print.FigurePrint(w, "", "green")
	fmt.Println()

	var wR1O1H, wR1O1T, wR1O1L, wR1O2H, wR1O2T, wR1O2L, wR1O3H, wR1O3T, wR1O3L, wR1O4H, wR1O4T, wR1O4L string
	if _, ok := brackets.WestConf[1][0]; ok {
		wR1O1H = brackets.WestConf[1][0].HighSeedCode
		wR1O1T = brackets.WestConf[1][0].SeriesText
		wR1O1L = brackets.WestConf[1][0].LowSeedCode
	}
	if _, ok := brackets.WestConf[1][1]; ok {
		wR1O2H = brackets.WestConf[1][1].HighSeedCode
		wR1O2T = brackets.WestConf[1][1].SeriesText
		wR1O2L = brackets.WestConf[1][1].LowSeedCode
	}
	if _, ok := brackets.WestConf[1][2]; ok {
		wR1O3H = brackets.WestConf[1][2].HighSeedCode
		wR1O3T = brackets.WestConf[1][2].SeriesText
		wR1O3L = brackets.WestConf[1][2].LowSeedCode
	}
	if _, ok := brackets.WestConf[1][3]; ok {
		wR1O4H = brackets.WestConf[1][3].HighSeedCode
		wR1O4T = brackets.WestConf[1][3].SeriesText
		wR1O4L = brackets.WestConf[1][3].LowSeedCode
	}

	var eR1O1H, eR1O1T, eR1O1L, eR1O2H, eR1O2L, eR1O2T, eR1O3H, eR1O3T, eR1O3L, eR1O4H, eR1O4T, eR1O4L string
	if _, ok := brackets.EastConf[1][4]; ok {
		eR1O1H = brackets.EastConf[1][4].HighSeedCode
		eR1O1T = brackets.EastConf[1][4].SeriesText
		eR1O1L = brackets.EastConf[1][4].LowSeedCode
	}
	if _, ok := brackets.EastConf[1][5]; ok {
		eR1O2H = brackets.EastConf[1][5].HighSeedCode
		eR1O2T = brackets.EastConf[1][5].SeriesText
		eR1O2L = brackets.EastConf[1][5].LowSeedCode
	}
	if _, ok := brackets.EastConf[1][6]; ok {
		eR1O3H = brackets.EastConf[1][6].HighSeedCode
		eR1O3T = brackets.EastConf[1][6].SeriesText
		eR1O3L = brackets.EastConf[1][6].LowSeedCode
	}
	if _, ok := brackets.EastConf[1][7]; ok {
		eR1O4H = brackets.EastConf[1][7].HighSeedCode
		eR1O4T = brackets.EastConf[1][7].SeriesText
		eR1O4L = brackets.EastConf[1][7].LowSeedCode
	}

	var wR2O1H, wR2O1T, wR2O1L, wR2O2H, wR2O2T, wR2O2L string
	if _, ok := brackets.WestConf[2][0]; ok {
		wR2O1H = brackets.WestConf[2][0].HighSeedCode
		wR2O1T = brackets.WestConf[2][0].SeriesText
		wR2O1L = brackets.WestConf[2][0].LowSeedCode
	}
	if _, ok := brackets.WestConf[2][1]; ok {
		wR2O2H = brackets.WestConf[2][1].HighSeedCode
		wR2O2T = brackets.WestConf[2][1].SeriesText
		wR2O2L = brackets.WestConf[2][1].LowSeedCode
	}

	var eR2O1H, eR2O1T, eR2O1L, eR2O2H, eR2O2T, eR2O2L string
	if _, ok := brackets.EastConf[2][2]; ok {
		eR2O1H = brackets.EastConf[2][2].HighSeedCode
		eR2O1T = brackets.EastConf[2][2].SeriesText
		eR2O1L = brackets.EastConf[2][2].LowSeedCode
	}
	if _, ok := brackets.EastConf[2][3]; ok {
		eR2O2H = brackets.EastConf[2][3].HighSeedCode
		eR2O2T = brackets.EastConf[2][3].SeriesText
		eR2O2L = brackets.EastConf[2][3].LowSeedCode
	}

	var wR3O1H, wR3O1T, wR3O1L string
	if _, ok := brackets.WestConf[3][0]; ok {
		wR3O1H = brackets.WestConf[3][0].HighSeedCode
		wR3O1T = brackets.WestConf[3][0].SeriesText
		wR3O1L = brackets.WestConf[3][0].LowSeedCode
	}

	var eR3O1H, eR3O1T, eR3O1L string
	if _, ok := brackets.EastConf[3][1]; ok {
		eR3O1H = brackets.EastConf[3][1].HighSeedCode
		eR3O1T = brackets.EastConf[3][1].SeriesText
		eR3O1L = brackets.EastConf[3][1].LowSeedCode
	}

	var fH, fT, fL string
	fT = brackets.Final.SeriesText
	if fT == "" {
		fT = "NBA Finals"
	}
	fH = brackets.Final.HighSeedCode
	fL = brackets.Final.LowSeedCode

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"WEST", "WEST", "WEST", " 🏆 ", "EAST", "EAST", "EAST"}, table.RowConfig{AutoMerge: true})
	t.AppendRow(table.Row{"Round1", "Conf.Series", "Conf.Finals", "Finals", "Conf.Finals", "Conf.Series", "Round1"})
	t.AppendRow(table.Row{wR1O1H, " ", " ", " ", " ", " ", eR1O1H})
	t.AppendRow(table.Row{wR1O1T, " ", " ", " ", " ", " ", eR1O1T})
	t.AppendRow(table.Row{wR1O1L, " ", " ", " ", " ", " ", eR1O1L})
	t.AppendRow(table.Row{" ", wR2O1H, " ", " ", " ", eR2O1H, " "})
	t.AppendRow(table.Row{" ", wR2O1T, " ", " ", " ", eR2O1T, " "})
	t.AppendRow(table.Row{" ", wR2O1L, " ", " ", " ", eR2O1L, " "})
	t.AppendRow(table.Row{wR1O2H, " ", " ", " ", " ", " ", eR1O2H})
	t.AppendRow(table.Row{wR1O2T, " ", " ", " ", " ", " ", eR1O2T})
	t.AppendRow(table.Row{wR1O2L, " ", " ", " ", " ", " ", eR1O2L})

	t.AppendRow(table.Row{" ", " ", wR3O1H, fH, eR3O1H, " ", " "})
	t.AppendRow(table.Row{" ", " ", wR3O1T, fT, eR3O1T, " ", " "})
	t.AppendRow(table.Row{" ", " ", wR3O1L, fL, eR3O1L, " ", " "})

	t.AppendRow(table.Row{wR1O3H, " ", " ", " ", " ", " ", eR1O3H})
	t.AppendRow(table.Row{wR1O3T, " ", " ", " ", " ", " ", eR1O3T})
	t.AppendRow(table.Row{wR1O3L, " ", " ", " ", " ", " ", eR1O3L})
	t.AppendRow(table.Row{" ", wR2O2H, " ", " ", " ", eR2O2H, " "})
	t.AppendRow(table.Row{" ", wR2O2T, " ", " ", " ", eR2O2T, " "})
	t.AppendRow(table.Row{" ", wR2O2L, " ", " ", " ", eR2O2L, " "})
	t.AppendRow(table.Row{wR1O4H, " ", " ", " ", " ", " ", eR1O4H})
	t.AppendRow(table.Row{wR1O4T, " ", " ", " ", " ", " ", eR1O4T})
	t.AppendRow(table.Row{wR1O4L, " ", " ", " ", " ", " ", eR1O4L})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true, Align: text.AlignCenter, VAlign: text.VAlignMiddle, AlignHeader: text.AlignCenter},
		{Number: 2, AutoMerge: true, Align: text.AlignCenter, VAlign: text.VAlignMiddle, AlignHeader: text.AlignCenter},
		{Number: 3, AutoMerge: true, Align: text.AlignCenter, VAlign: text.VAlignMiddle, AlignHeader: text.AlignCenter},
		{Number: 4, AutoMerge: true, Align: text.AlignCenter, VAlign: text.VAlignMiddle, AlignHeader: text.AlignCenter},
		{Number: 5, AutoMerge: true, Align: text.AlignCenter, VAlign: text.VAlignMiddle, AlignHeader: text.AlignCenter},
		{Number: 6, AutoMerge: true, Align: text.AlignCenter, VAlign: text.VAlignMiddle, AlignHeader: text.AlignCenter},
		{Number: 7, AutoMerge: true, Align: text.AlignCenter, VAlign: text.VAlignMiddle, AlignHeader: text.AlignCenter},
	})
	t.SetStyle(table.StyleBold)
	t.Style().Options.SeparateRows = true
	t.Render()

	return nil
}

func addWinOrGoHomeEmojiFromBracket(b *Bracket, winE string, loseE string) (string, string) {
	var s1, s2 string
	s1 = "/"
	s2 = "/"

	var hc, lc, cnt string
	var hr, lr int
	hc = b.HighSeedCode
	lc = b.LowSeedCode
	hr = b.HighSeedRank
	lr = b.LowSeedRank
	cnt = b.SeriesText

	switch {
	case strings.Contains(strings.ToUpper(cnt), "WINS") && strings.Contains(cnt, hc):
		s1 = fmt.Sprintf("%d %s %s", hr, hc, winE)
		s2 = fmt.Sprintf("%d %s %s", lr, lc, loseE)
	case strings.Contains(strings.ToUpper(cnt), "WINS") && strings.Contains(cnt, lc):
		s1 = fmt.Sprintf("%d %s %s", hr, hc, loseE)
		s2 = fmt.Sprintf("%d %s %s", lr, lc, winE)
	case cnt != "":
		s1 = fmt.Sprintf("%d %s", hr, hc)
		s2 = fmt.Sprintf("%d %s", lr, lc)
	default:
	}
	if !strings.Contains(s1, winE) {
		return s2, s1
	}

	return s1, s2
}
