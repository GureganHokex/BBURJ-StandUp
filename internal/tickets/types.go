package tickets

import "time"

const (
	SourceManual       = "manual"
	SourceTimepad      = "timepad"
	SourceTicketscloud = "ticketscloud"
)

type Provider struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Configured  bool   `json:"configured"`
	Description string `json:"description,omitempty"`
}

type Event struct {
	Source      string    `json:"source"`
	ExternalID  string    `json:"external_id"`
	Title       string    `json:"title"`
	Date        time.Time `json:"date"`
	City        string    `json:"city"`
	Description    string    `json:"description"`
	TicketURL      string    `json:"ticket_url"`
	PosterImageURL string    `json:"poster_image_url"`
	AlreadyAdded   bool      `json:"already_added"`
}

type Settings struct {
	TimepadOrgID        string
	TimepadAPIKey       string
	TicketscloudOrgID   string
	TicketscloudAPIKey  string
	EventImportKeywords string
}
