package render

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Renderer struct {
	templates *template.Template
}

func New(templates *template.Template) *Renderer {
	return &Renderer{templates: templates}
}

func (r *Renderer) HTML(c *gin.Context, code int, name string, data any) {
	c.Status(code)
	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := r.templates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.String(http.StatusInternalServerError, "template error")
	}
}

// Page renders content block, then wraps it in layout (Go templates require static template names).
func (r *Renderer) Page(c *gin.Context, code int, layoutTpl, contentTpl string, data gin.H) {
	var content bytes.Buffer
	if err := r.templates.ExecuteTemplate(&content, contentTpl, data); err != nil {
		c.String(http.StatusInternalServerError, "template error")
		return
	}
	data["ContentHTML"] = template.HTML(content.String())
	c.Status(code)
	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := r.templates.ExecuteTemplate(c.Writer, layoutTpl, data); err != nil {
		c.String(http.StatusInternalServerError, "template error")
	}
}
