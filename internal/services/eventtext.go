package services

import (
	"regexp"
	"strings"
)

var (
	ruMonthPattern = `(?i)(января|февраля|марта|апреля|мая|июня|июля|августа|сентября|октября|ноября|декабря)`
	datePartRe     = regexp.MustCompile(`(?i)^\d{1,2}\s+` + ruMonthPattern + `(\s+\d{1,2}:\d{2})?$`)
	timePartRe     = regexp.MustCompile(`^\d{1,2}:\d{2}$`)
	shortDateRe    = regexp.MustCompile(`(?i)^\d{1,2}\.\d{1,2}(\.\d{2,4})?$`)
)

type ParsedEventText struct {
	Title       string
	Description string
}

func ParseEventTitle(rawTitle, rawDescription string) ParsedEventText {
	rawTitle = strings.TrimSpace(rawTitle)
	rawDescription = strings.TrimSpace(rawDescription)

	parts := splitTitleParts(rawTitle)
	kept := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" || isDateTimePart(part) {
			continue
		}
		kept = append(kept, part)
	}

	title := rawTitle
	descParts := make([]string, 0, 2)
	if len(kept) == 0 {
		title = strings.TrimSpace(stripDateTimeSuffix(rawTitle))
	} else if len(kept) == 1 {
		title = kept[0]
	} else {
		title = kept[len(kept)-1]
		descParts = append(descParts, strings.Join(kept[:len(kept)-1], " · "))
	}

	if rawDescription != "" {
		descParts = append(descParts, rawDescription)
	}

	return ParsedEventText{
		Title:       title,
		Description: strings.TrimSpace(strings.Join(descParts, "\n\n")),
	}
}

func splitTitleParts(raw string) []string {
	raw = strings.ReplaceAll(raw, " | ", " / ")
	raw = strings.ReplaceAll(raw, "|", " / ")
	return strings.Split(raw, "/")
}

func isDateTimePart(part string) bool {
	part = strings.TrimSpace(part)
	if part == "" {
		return true
	}
	if datePartRe.MatchString(part) {
		return true
	}
	if timePartRe.MatchString(part) {
		return true
	}
	if shortDateRe.MatchString(part) {
		return true
	}
	return false
}

func stripDateTimeSuffix(raw string) string {
	parts := splitTitleParts(raw)
	kept := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" || isDateTimePart(part) {
			continue
		}
		kept = append(kept, part)
	}
	return strings.Join(kept, " / ")
}

func AttrText(s string) string {
	s = strings.ReplaceAll(s, "\r\n", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	return strings.TrimSpace(s)
}
