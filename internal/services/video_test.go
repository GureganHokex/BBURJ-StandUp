package services

import (
	"testing"

	"github.com/burj/comic/internal/models"
)

func TestYouTubeVideoID(t *testing.T) {
	cases := []struct {
		url  string
		want string
	}{
		{"https://www.youtube.com/watch?v=dQw4w9WgXcQ", "dQw4w9WgXcQ"},
		{"https://youtu.be/dQw4w9WgXcQ", "dQw4w9WgXcQ"},
		{"https://www.youtube.com/shorts/dQw4w9WgXcQ", "dQw4w9WgXcQ"},
		{"https://www.youtube.com/embed/dQw4w9WgXcQ", "dQw4w9WgXcQ"},
		{"https://www.youtube.com/watch?v=dQw4w9WgXcQ&t=42", "dQw4w9WgXcQ"},
		{"https://www.youtube.com/channel/UCxxxx", ""},
		{"https://example.com/watch?v=dQw4w9WgXcQ", ""},
	}

	for _, tc := range cases {
		got := youtubeVideoID(tc.url)
		if got != tc.want {
			t.Errorf("youtubeVideoID(%q) = %q, want %q", tc.url, got, tc.want)
		}
	}
}

func TestEmbedURLYouTube(t *testing.T) {
	got := EmbedURL(models.PlatformYouTube, "https://www.youtube.com/shorts/dQw4w9WgXcQ")
	want := "https://www.youtube-nocookie.com/embed/dQw4w9WgXcQ"
	if got != want {
		t.Errorf("EmbedURL() = %q, want %q", got, want)
	}
}
