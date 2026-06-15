package tickets

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func fetchTimepad(ctx context.Context, client *http.Client, cfg Settings) ([]Event, error) {
	orgID := strings.TrimSpace(cfg.TimepadOrgID)
	if orgID == "" {
		return nil, fmt.Errorf("timepad org id is not configured")
	}

	q := url.Values{}
	q.Add("organization_ids[]", orgID)
	q.Add("limit", "50")
	q.Add("sort", "+starts_at")
	q.Add("starts_at_min", time.Now().UTC().Format(time.RFC3339))
	for _, kw := range splitKeywords(cfg.EventImportKeywords) {
		q.Add("keywords[]", kw)
	}
	q.Add("moderation_statuses[]", "featured")
	q.Add("moderation_statuses[]", "shown")
	q.Add("moderation_statuses[]", "not_moderated")

	reqURL := "https://api.timepad.ru/v1/events.json?" + q.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if key := strings.TrimSpace(cfg.TimepadAPIKey); key != "" {
		req.Header.Set("Authorization", "Bearer "+key)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("timepad api: %s", strings.TrimSpace(string(body)))
	}

	var payload struct {
		Values []struct {
			ID               int    `json:"id"`
			Name             string `json:"name"`
			DescriptionShort string `json:"description_short"`
			StartsAt         string `json:"starts_at"`
			URL              string `json:"url"`
			Location         struct {
				City string `json:"city"`
			} `json:"location"`
		} `json:"values"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}

	out := make([]Event, 0, len(payload.Values))
	for _, item := range payload.Values {
		date, err := time.Parse(time.RFC3339, item.StartsAt)
		if err != nil {
			date, _ = time.Parse("2006-01-02", item.StartsAt)
		}
		out = append(out, Event{
			Source:      SourceTimepad,
			ExternalID:  fmt.Sprintf("%d", item.ID),
			Title:       item.Name,
			Date:        date,
			City:        strings.TrimSpace(item.Location.City),
			Description: item.DescriptionShort,
			TicketURL:   item.URL,
		})
	}
	return out, nil
}
