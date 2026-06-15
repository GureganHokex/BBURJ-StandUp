package api

type ListResponse[T any] struct {
	Data []T  `json:"data"`
	Meta Meta `json:"meta"`
}

type ItemResponse[T any] struct {
	Data T `json:"data"`
}

type Meta struct {
	Total int64 `json:"total"`
}

type ErrorResponse struct {
	Error  string            `json:"error,omitempty"`
	Errors map[string]string `json:"errors,omitempty"`
}

type EventRequest struct {
	Title        string `json:"title"`
	Date         string `json:"date"`
	City         string `json:"city"`
	Description  string `json:"description"`
	TicketURL    string `json:"ticket_url"`
	TicketSource string `json:"ticket_source"`
	ExternalID   string `json:"external_id"`
}

type VideoRequest struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type MerchRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	ImageURL    string `json:"image_url"`
	BuyURL      string `json:"buy_url"`
}
