package print

import (
	"strconv"

	"github.com/common-nighthawk/go-figure"
	"github.com/jedib0t/go-pretty/v6/text"
)

// FigurePrint print string with style and color
func FigurePrint(str string, style string, color string) {
	myFigure := figure.NewColorFigure(str, style, color, true)
	myFigure.Print()
}

// MaxScoreColor max score with color
func MaxScoreColor(s1 string, s2 string) (str1 string, str2 string, err error) {
	n1, err := strconv.Atoi(s1)
	if err != nil {
		return str1, str2, err
	}
	n2, err := strconv.Atoi(s2)
	if err != nil {
		return str1, str2, err
	}
	str1 = text.FgWhite.Sprintf("%3s", s1)
	str2 = text.FgWhite.Sprintf("%-3s", s2)
	if n1 > n2 {
		str1 = text.FgGreen.Sprintf("%3s", s1)
	}
	if n2 > n1 {
		str2 = text.FgGreen.Sprintf("%-3s", s2)
	}

	return str1, str2, nil
}

// MaxTextColor max score text with color
func MaxTextColor(s1 string, s2 string) (str1 string, str2 string, err error) {
	n1, err := strconv.Atoi(s1)
	if err != nil {
		return str1, str2, err
	}
	n2, err := strconv.Atoi(s2)
	if err != nil {
		return str1, str2, err
	}
	str1 = text.FgWhite.Sprintf("%s", s1)
	str2 = text.FgWhite.Sprintf("%s", s2)
	if n1 > n2 {
		str1 = text.FgGreen.Sprintf("%s", s1)
	}
	if n2 > n1 {
		str2 = text.FgGreen.Sprintf("%s", s2)
	}

	return str1, str2, nil
}
