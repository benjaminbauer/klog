package cli

import (
	"errors"
	"fmt"
	. "klog/lib/jotaen/tf"
	"klog/parser/engine"
	"klog/record"
	"strings"
)

func prettifyError(err error) error {
	if pErrs, isParserErrors := err.(engine.Errors); isParserErrors {
		message := ""
		INDENT := "    "
		for _, e := range pErrs.Get() {
			message += fmt.Sprintf(
				Style{Background: "160", Color: "015"}.Format(" Error in line %d: "),
				e.Context().LineNumber,
			) + "\n"
			message += fmt.Sprintf(
				Style{Color: "247"}.Format(INDENT+"%s"),
				string(e.Context().Value),
			) + "\n"
			message += fmt.Sprintf(
				Style{Color: "160"}.Format(INDENT+"%s%s"),
				strings.Repeat(" ", e.Position()), strings.Repeat("^", e.Length()),
			) + "\n"
			message += fmt.Sprintf(
				Style{Color: "227"}.Format(INDENT+"%s"),
				strings.Join(breakLines(e.Message(), 60), "\n"+INDENT),
			) + "\n\n"
		}
		return errors.New(message)
	} else {
		return err
	}
}

func breakLines(text string, maxLength int) []string {
	SPACE := " "
	words := strings.Split(text, SPACE)
	lines := []string{""}
	for i, w := range words {
		lastLine := lines[len(lines)-1]
		isLastWord := i == len(words)-1
		if !isLastWord && len(lastLine)+len(words[i+1]) > maxLength {
			lines = append(lines, "")
		}
		lines[len(lines)-1] += w + SPACE
	}
	return lines
}

type stylerT struct{}

var styler stylerT

func (h stylerT) PrintDate(d record.Date) string {
	return Style{Color: "098", IsUnderlined: true}.Format(d.ToString())
}
func (h stylerT) PrintShouldTotal(d record.Duration, symbol string) string {
	return Style{Color: "213"}.Format(d.ToString()) + Style{Color: "201"}.Format(symbol)
}
func (h stylerT) PrintSummary(s record.Summary) string {
	txt := s.ToString()
	style := Style{Color: "249"}
	hashStyle := style.ChangedBold(true).ChangedColor("251")
	txt = record.HashTagPattern.ReplaceAllStringFunc(txt, func(h string) string {
		return hashStyle.FormatAndRestore(h, style)
	})
	return style.Format(txt)
}
func (h stylerT) PrintRange(r record.Range) string {
	return Style{Color: "117"}.Format(r.ToString())
}
func (h stylerT) PrintOpenRange(or record.OpenRange) string {
	return Style{Color: "027"}.Format(or.ToString())
}
func (h stylerT) PrintDuration(d record.Duration) string {
	f := Style{Color: "120"}
	if d.InMinutes() < 0 {
		f.Color = "167"
	}
	return f.Format(d.ToString())
}
