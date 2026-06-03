package public

import (
	"github.com/burj/comic/internal/models"
	"github.com/burj/comic/internal/render"
	"github.com/burj/comic/internal/services"
	"github.com/gin-gonic/gin"
)

type VideoView struct {
	models.Video
	EmbedURL string
}

type VideosHandler struct {
	videos   *services.VideoService
	settings *services.SiteSettingsService
	render   *render.Renderer
}

func NewVideosHandler(videos *services.VideoService, settings *services.SiteSettingsService, render *render.Renderer) *VideosHandler {
	return &VideosHandler{videos: videos, settings: settings, render: render}
}

func (h *VideosHandler) List(c *gin.Context) {
	items, _, _ := h.videos.List(100, 0)
	views := make([]VideoView, 0, len(items))
	for _, v := range items {
		views = append(views, VideoView{
			Video:    v,
			EmbedURL: services.EmbedURL(v.Platform, v.URL),
		})
	}
	h.render.Page(c, 200, "public/layout", "public/videos_content", mergeSite(h.settings, gin.H{
		"Title":     "Видео",
		"PageTitle": "Видео",
		"Videos":    views,
		"Active":    "videos",
	}))
}
