package services

import "testing"

func TestParseRussianDateTime(t *testing.T) {
	tm, ok := parseRussianDateTime("30 июня 2026, вторник, 20:00")
	if !ok {
		t.Fatal("expected parse ok")
	}
	if tm.Day() != 30 || tm.Month() != 6 || tm.Year() != 2026 || tm.Hour() != 20 || tm.Minute() != 0 {
		t.Fatalf("unexpected time: %v", tm)
	}
}

func TestParseDateFromTitleParts(t *testing.T) {
	raw := "Денис Антипин / Илья Буржинский / Балуемся / 30 июня / 20:00"
	tm, ok := parseDateFromTitleParts(raw)
	if !ok {
		t.Fatal("expected parse ok")
	}
	if tm.Day() != 30 || tm.Month() != 6 || tm.Hour() != 20 {
		t.Fatalf("unexpected time: %v", tm)
	}
}

func TestExtractPageCity(t *testing.T) {
	html := `<section class="showroom-event-slide showroom-event-slide_venue"><address><span>Санкт-Петербург,<br/>ул. Восстания</span></address></section>`
	city := extractPageCity(html)
	if city != "Санкт-Петербург" {
		t.Fatalf("got %q", city)
	}
}
