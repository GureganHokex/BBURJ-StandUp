package tickets

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var tcDTStart = regexp.MustCompile(`DTSTART[^:]*:([0-9]{8}T[0-9]{6}Z?)`)

func fetchTicketscloud(ctx context.Context, client *http.Client, cfg Settings) ([]Event, error) {
	orgID := strings.TrimSpace(cfg.TicketscloudOrgID)
	apiKey := strings.TrimSpace(cfg.TicketscloudAPIKey)
	if orgID == "" || apiKey == "" {
		return nil, fmt.Errorf("ticketscloud org id and api key are required")
	}

	q := url.Values{}
	q.Set("org", orgID)
	q.Set("status", "public")
	q.Set("removed", "false")
	q.Set("fields-schema", "id,title,lifetime,venue,status")

	reqURL := "https://ticketscloud.com/v1/resources/events?" + q.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Key "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ticketscloud api: %s", strings.TrimSpace(string(body)))
	}

	var raw []tcEvent
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	venueCities := map[string]string{}
	out := make([]Event, 0, len(raw))
	keywords := splitKeywords(cfg.EventImportKeywords)

	for _, item := range raw {
		title, desc := parseTCTitle(item.Title)
		if len(keywords) > 0 && !matchesKeywords(title, keywords) {
			continue
		}

		date := parseTCDate(item.Lifetime)
		city := ""
		if item.Venue != "" {
			if c, ok := venueCities[item.Venue]; ok {
				city = c
			} else {
				city, _ = fetchTCVenueCity(ctx, client, apiKey, item.Venue)
				venueCities[item.Venue] = city
			}
		}

		out = append(out, Event{
			Source:      SourceTicketscloud,
			ExternalID:  item.ID,
			Title:       title,
			Date:        date,
			City:        city,
			Description: desc,
			TicketURL:   "https://ticketscloud.com/event/" + item.ID,
		})
	}
	return out, nil
}

type tcEvent struct {
	ID       string          `json:"id"`
	Title    json.RawMessage `json:"title"`
	Lifetime string          `json:"lifetime"`
	Venue    string          `json:"venue"`
}

func parseTCTitle(raw json.RawMessage) (title, desc string) {
	if len(raw) == 0 {
		return "", ""
	}
	var asString string
	if err := json.Unmarshal(raw, &asString); err == nil {
		return asString, ""
	}
	var asObject struct {
		Text string `json:"text"`
		Desc string `json:"desc"`
	}
	if err := json.Unmarshal(raw, &asObject); err == nil {
		return asObject.Text, asObject.Desc
	}
	return strings.Trim(string(raw), `"`), ""
}

func parseTCDate(lifetime string) time.Time {
	match := tcDTStart.FindStringSubmatch(lifetime)
	if len(match) < 2 {
		return time.Time{}
	}
	layouts := []string{"20060102T150405Z", "20060102T150405"}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, match[1]); err == nil {
			return t
		}
	}
	return time.Time{}
}

func fetchTCVenueCity(ctx context.Context, client *http.Client, apiKey, venueID string) (string, error) {
	reqURL := "https://ticketscloud.com/v1/resources/venues/" + url.PathEscape(venueID) + "?fields-schema=name,address"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Key "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("venue lookup failed")
	}

	var venue struct {
		Address struct {
			City string `json:"city"`
		} `json:"address"`
	}
	if err := json.Unmarshal(body, &venue); err != nil {
		return "", err
	}
	return strings.TrimSpace(venue.Address.City), nil
}

func splitKeywords(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func matchesKeywords(title string, keywords []string) bool {
	lower := strings.ToLower(title)
	for _, kw := range keywords {
		if strings.Contains(lower, strings.ToLower(kw)) {
			return true
		}
	}
	return false
}
