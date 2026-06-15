package services

import (
	"html"
	"regexp"
	"strings"
	"time"
)

var (
	tcDescRe = regexp.MustCompile(`(?is)showroom-event-slide__content_desc[^>]*>([\s\S]*?)</article>`)
	tcCityRe = regexp.MustCompile(`(?is)showroom-event-slide_venue[\s\S]*?<address>[\s\S]*?<span>\s*(?:<[^>]+>\s*)?([^,<]+)\s*,`)
	tcTimeRe = regexp.MustCompile(`(?is)showroom-event-slide_venue[\s\S]*?<time>\s*([^<]+?)\s*</time>`)
	timeTagRe = regexp.MustCompile(`(?is)<time[^>]*datetime=["']([^"']+)["'][^>]*>`)
	jsonLDRe  = regexp.MustCompile(`(?is)<script[^>]+type=["']application/ld\+json["'][^>]*>([\s\S]*?)</script>`)
	ruDateTimeRe = regexp.MustCompile(`(?i)(\d{1,2})\s+([а-яё]+)(?:\s+(\d{4}))?(?:,\s*[а-яёё]+)?,?\s*(\d{1,2}):(\d{2})`)
	tagRe     = regexp.MustCompile(`(?is)<[^>]+>`)
	wsRe      = regexp.MustCompile(`\s+`)
)

var ruGenitiveMonths = map[string]time.Month{
	"января":   time.January,
	"февраля":  time.February,
	"марта":    time.March,
	"апреля":   time.April,
	"мая":      time.May,
	"июня":     time.June,
	"июля":     time.July,
	"августа":  time.August,
	"сентября": time.September,
	"октября":  time.October,
	"ноября":   time.November,
	"декабря":  time.December,
}

func enrichPreviewFromHTML(preview *PagePreview, pageHTML, rawTitle string) {
	if preview.City == "" {
		preview.City = extractPageCity(pageHTML)
	}
	if preview.Description == "" {
		preview.Description = extractPageDescription(pageHTML)
	}
	if preview.Date == "" {
		if t, ok := extractPageDateTime(pageHTML); ok {
			preview.Date = t.UTC().Format(time.RFC3339)
		}
	}
	if preview.Date == "" {
		if t, ok := parseDateFromTitleParts(rawTitle); ok {
			preview.Date = t.UTC().Format(time.RFC3339)
		}
	}
}

func extractPageDescription(pageHTML string) string {
	if m := tcDescRe.FindStringSubmatch(pageHTML); len(m) > 1 {
		return cleanHTMLText(m[1])
	}
	return ""
}

func extractPageCity(pageHTML string) string {
	if m := tcCityRe.FindStringSubmatch(pageHTML); len(m) > 1 {
		return strings.TrimSpace(m[1])
	}
	return ""
}

func extractPageDateTime(pageHTML string) (time.Time, bool) {
	if m := tcTimeRe.FindStringSubmatch(pageHTML); len(m) > 1 {
		if t, ok := parseRussianDateTime(strings.TrimSpace(m[1])); ok {
			return t, true
		}
	}
	if m := timeTagRe.FindStringSubmatch(pageHTML); len(m) > 1 {
		if t, err := time.Parse(time.RFC3339, m[1]); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

func parseDateFromTitleParts(rawTitle string) (time.Time, bool) {
	parts := splitTitleParts(rawTitle)
	var dayMonth string
	var clock string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if datePartRe.MatchString(part) {
			dayMonth = part
		}
		if timePartRe.MatchString(part) {
			clock = part
		}
	}
	if dayMonth == "" {
		return time.Time{}, false
	}
	raw := dayMonth
	if clock != "" {
		raw += ", " + clock
	}
	if t, ok := parseRussianDateTime(raw); ok {
		return t, true
	}
	if clock != "" {
		return parseRussianDateTime(dayMonth + " " + clock)
	}
	return time.Time{}, false
}

func parseRussianDateTime(raw string) (time.Time, bool) {
	raw = strings.TrimSpace(raw)
	m := ruDateTimeRe.FindStringSubmatch(raw)
	if len(m) < 6 {
		return time.Time{}, false
	}

	day := atoi(m[1])
	month, ok := ruGenitiveMonths[strings.ToLower(m[2])]
	if !ok {
		return time.Time{}, false
	}
	year := time.Now().Year()
	if m[3] != "" {
		year = atoi(m[3])
	}
	hour := atoi(m[4])
	minute := atoi(m[5])

	loc := time.Local
	t := time.Date(year, month, day, hour, minute, 0, 0, loc)
	if m[3] == "" && t.Before(time.Now().Add(-24*time.Hour)) {
		t = time.Date(year+1, month, day, hour, minute, 0, 0, loc)
	}
	return t, true
}

func cleanHTMLText(raw string) string {
	raw = tagRe.ReplaceAllString(raw, " ")
	raw = html.UnescapeString(raw)
	raw = wsRe.ReplaceAllString(raw, " ")
	return strings.TrimSpace(raw)
}

func atoi(s string) int {
	n := 0
	for _, r := range s {
		if r < '0' || r > '9' {
			continue
		}
		n = n*10 + int(r-'0')
	}
	return n
}
