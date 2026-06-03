package web

import (
	"fmt"
	"strings"
	"time"
)

var ruMonths = []string{
	"", "ЯНВАРЯ", "ФЕВРАЛЯ", "МАРТА", "АПРЕЛЯ", "МАЯ", "ИЮНЯ",
	"ИЮЛЯ", "АВГУСТА", "СЕНТЯБРЯ", "ОКТЯБРЯ", "НОЯБРЯ", "ДЕКАБРЯ",
}

func formatEventDateCard(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return fmt.Sprintf("%d %s %02d:%02d", t.Day(), ruMonths[t.Month()], t.Hour(), t.Minute())
}

func formatEventDay(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return fmt.Sprintf("%d", t.Day())
}

func formatEventMeta(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return fmt.Sprintf("%s %02d:%02d", ruMonths[t.Month()], t.Hour(), t.Minute())
}

func upperASCII(s string) string {
	return strings.ToUpper(s)
}
