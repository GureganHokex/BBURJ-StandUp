package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var (
	ErrPreviewURLInvalid = errors.New("invalid preview url")
	ErrPreviewURLBlocked = errors.New("preview url is not allowed")
	ErrPreviewNoImage    = errors.New("no poster image found on page")
)

type PagePreview struct {
	PosterImageURL string `json:"poster_image_url"`
	Title          string `json:"title,omitempty"`
	Description    string `json:"description,omitempty"`
	City           string `json:"city,omitempty"`
	Date           string `json:"date,omitempty"`
}

var (
	ogImagePatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)<meta[^>]+property=["']og:image(?::url)?["'][^>]+content=["']([^"']+)["']`),
		regexp.MustCompile(`(?i)<meta[^>]+content=["']([^"']+)["'][^>]+property=["']og:image(?::url)?["']`),
		regexp.MustCompile(`(?i)<meta[^>]+name=["']twitter:image["'][^>]+content=["']([^"']+)["']`),
		regexp.MustCompile(`(?i)<meta[^>]+content=["']([^"']+)["'][^>]+name=["']twitter:image["']`),
	}
	ogTitlePatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)<meta[^>]+property=["']og:title["'][^>]+content=["']([^"']+)["']`),
		regexp.MustCompile(`(?i)<meta[^>]+content=["']([^"']+)["'][^>]+property=["']og:title["']`),
	}
	ogDescPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)<meta[^>]+property=["']og:description["'][^>]+content=["']([^"']+)["']`),
		regexp.MustCompile(`(?i)<meta[^>]+content=["']([^"']+)["'][^>]+property=["']og:description["']`),
	}
)

type URLPreviewService struct {
	client *http.Client
}

func NewURLPreviewService() *URLPreviewService {
	return &URLPreviewService{
		client: &http.Client{
			Timeout: 12 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 4 {
					return errors.New("too many redirects")
				}
				if err := validatePreviewURL(req.URL.String()); err != nil {
					return err
				}
				return nil
			},
		},
	}
}

func (s *URLPreviewService) FetchPagePreview(ctx context.Context, rawURL string) (PagePreview, error) {
	if err := validatePreviewURL(rawURL); err != nil {
		return PagePreview{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return PagePreview{}, err
	}
	req.Header.Set("User-Agent", "BBURJ-StandUp/1.0 (+poster-preview)")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")

	resp, err := s.client.Do(req)
	if err != nil {
		return PagePreview{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return PagePreview{}, fmt.Errorf("page returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return PagePreview{}, err
	}

	html := string(body)
	ogTitle := extractMeta(html, ogTitlePatterns)
	parsed := ParseEventTitle(ogTitle, "")

	preview := PagePreview{
		Title: parsed.Title,
	}

	enrichPreviewFromHTML(&preview, html, ogTitle)

	if preview.Description == "" {
		if ogDesc := extractMeta(html, ogDescPatterns); ogDesc != "" {
			preview.Description = ogDesc
		} else {
			preview.Description = parsed.Description
		}
	}

	imageURL := extractMeta(html, ogImagePatterns)
	if imageURL == "" {
		return preview, ErrPreviewNoImage
	}

	imageURL = resolveRelativeURL(resp.Request.URL, imageURL)
	if err := validatePreviewURL(imageURL); err != nil {
		return preview, err
	}
	if errs := validateOptionalURL("poster_image_url", imageURL); errs.HasErrors() {
		return preview, ErrPreviewURLInvalid
	}
	preview.PosterImageURL = imageURL
	return preview, nil
}

func (s *URLPreviewService) FetchPosterImage(ctx context.Context, rawURL string) (string, error) {
	preview, err := s.FetchPagePreview(ctx, rawURL)
	if err != nil {
		if errors.Is(err, ErrPreviewNoImage) {
			return "", err
		}
		return "", err
	}
	if preview.PosterImageURL == "" {
		return "", ErrPreviewNoImage
	}
	return preview.PosterImageURL, nil
}

func validatePreviewURL(rawURL string) error {
	rawURL = strings.TrimSpace(rawURL)
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return ErrPreviewURLInvalid
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return ErrPreviewURLInvalid
	}
	if strings.Contains(parsed.Host, "@") {
		return ErrPreviewURLInvalid
	}

	host := strings.ToLower(parsed.Hostname())
	if host == "localhost" || strings.HasSuffix(host, ".local") {
		return ErrPreviewURLBlocked
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return nil
	}
	for _, ip := range ips {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() {
			return ErrPreviewURLBlocked
		}
	}
	return nil
}

func extractMeta(html string, patterns []*regexp.Regexp) string {
	for _, re := range patterns {
		if m := re.FindStringSubmatch(html); len(m) > 1 {
			return strings.TrimSpace(m[1])
		}
	}
	return ""
}

func resolveRelativeURL(base *url.URL, raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	if parsed.IsAbs() {
		return parsed.String()
	}
	return base.ResolveReference(parsed).String()
}
