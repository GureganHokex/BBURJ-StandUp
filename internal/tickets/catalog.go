package tickets

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Catalog struct {
	client *http.Client
}

func NewCatalog() *Catalog {
	return &Catalog{
		client: &http.Client{Timeout: 20 * time.Second},
	}
}

func (c *Catalog) Providers(cfg Settings) []Provider {
	return []Provider{
		{
			ID:          SourceTimepad,
			Name:        "Timepad",
			Configured:  strings.TrimSpace(cfg.TimepadOrgID) != "",
			Description: "Публичные события организации на timepad.ru",
		},
		{
			ID:          SourceTicketscloud,
			Name:        "TicketsCloud",
			Configured:  strings.TrimSpace(cfg.TicketscloudOrgID) != "" && strings.TrimSpace(cfg.TicketscloudAPIKey) != "",
			Description: "Нужны API-ключ и ID организатора от TicketsCloud",
		},
	}
}

func (c *Catalog) ListEvents(ctx context.Context, source string, cfg Settings) ([]Event, error) {
	switch source {
	case SourceTimepad:
		return fetchTimepad(ctx, c.client, cfg)
	case SourceTicketscloud:
		return fetchTicketscloud(ctx, c.client, cfg)
	default:
		return nil, fmt.Errorf("unknown ticket source: %s", source)
	}
}

func SettingsFromModel(orgTimepad, keyTimepad, orgTC, keyTC, keywords string) Settings {
	return Settings{
		TimepadOrgID:        orgTimepad,
		TimepadAPIKey:       keyTimepad,
		TicketscloudOrgID:   orgTC,
		TicketscloudAPIKey:  keyTC,
		EventImportKeywords: keywords,
	}
}
