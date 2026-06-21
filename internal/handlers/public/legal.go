package public

import (
	"github.com/burj/comic/internal/render"
	"github.com/burj/comic/internal/services"
	"github.com/gin-gonic/gin"
)

type LegalHandler struct {
	settings *services.SiteSettingsService
	render   *render.Renderer
}

func NewLegalHandler(settings *services.SiteSettingsService, render *render.Renderer) *LegalHandler {
	return &LegalHandler{settings: settings, render: render}
}

func (h *LegalHandler) Privacy(c *gin.Context) {
	h.render.Page(c, 200, "public/layout", "public/legal_privacy_content", mergeSite(h.settings, gin.H{
		"Title":     "Политика конфиденциальности",
		"PageTitle": "Политика конфиденциальности",
	}))
}

func (h *LegalHandler) Terms(c *gin.Context) {
	h.render.Page(c, 200, "public/layout", "public/legal_terms_content", mergeSite(h.settings, gin.H{
		"Title":     "Пользовательское соглашение",
		"PageTitle": "Пользовательское соглашение",
	}))
}

func (h *LegalHandler) Consent(c *gin.Context) {
	h.render.Page(c, 200, "public/layout", "public/legal_consent_content", mergeSite(h.settings, gin.H{
		"Title":     "Согласие на обработку персональных данных",
		"PageTitle": "Согласие на обработку персональных данных",
	}))
}
