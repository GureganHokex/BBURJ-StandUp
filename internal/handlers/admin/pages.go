package admin

import (
	"net/http"
	"strconv"

	"github.com/burj/comic/internal/admin"
	"github.com/burj/comic/internal/render"
	"github.com/burj/comic/internal/middleware"
	"github.com/gin-gonic/gin"
)

type PagesHandler struct {
	csrf   *middleware.CSRF
	render *render.Renderer
}

func NewPagesHandler(csrf *middleware.CSRF, render *render.Renderer) *PagesHandler {
	return &PagesHandler{csrf: csrf, render: render}
}

var modelFields = map[string][]fieldMeta{
	"events": {
		{Name: "title", Label: "Название", Type: "text", Required: true},
		{Name: "date", Label: "Дата", Type: "datetime-local", Required: true},
		{Name: "city", Label: "Город", Type: "text", Required: true},
		{Name: "description", Label: "Описание", Type: "textarea"},
		{Name: "ticket_url", Label: "Ссылка на билеты", Type: "url"},
	},
	"videos": {
		{Name: "title", Label: "Название", Type: "text", Required: true},
		{Name: "url", Label: "URL видео", Type: "url", Required: true},
	},
	"merch": {
		{Name: "title", Label: "Название", Type: "text", Required: true},
		{Name: "description", Label: "Описание", Type: "textarea"},
		{Name: "price", Label: "Цена (копейки)", Type: "number", Required: true},
		{Name: "image_url", Label: "Изображение", Type: "image_upload"},
		{Name: "buy_url", Label: "Ссылка на покупку", Type: "url"},
	},
	"photos": {
		{Name: "title", Label: "Подпись", Type: "text"},
		{Name: "image_url", Label: "Фото", Type: "image_upload", Required: true},
		{Name: "sort_order", Label: "Порядок", Type: "number"},
	},
}

var modelColumns = map[string][]columnMeta{
	"events": {
		{Key: "title", Label: "Название"},
		{Key: "date", Label: "Дата"},
		{Key: "city", Label: "Город"},
	},
	"videos": {
		{Key: "title", Label: "Название"},
		{Key: "platform", Label: "Платформа"},
		{Key: "url", Label: "URL"},
	},
	"merch": {
		{Key: "title", Label: "Название"},
		{Key: "price", Label: "Цена"},
	},
	"photos": {
		{Key: "title", Label: "Подпись"},
		{Key: "sort_order", Label: "Порядок"},
		{Key: "image_url", Label: "URL"},
	},
}

type fieldMeta struct {
	Name     string
	Label    string
	Type     string
	Required bool
}

type columnMeta struct {
	Key   string
	Label string
}

func (h *PagesHandler) List(c *gin.Context) {
	slug := c.Param("model")
	if slug == "settings" {
		c.Redirect(http.StatusFound, "/admin/settings")
		return
	}
	model, ok := admin.FindBySlug(slug)
	if !ok {
		c.Status(http.StatusNotFound)
		return
	}
	h.render.Page(c, 200, "admin/layout", "admin/list_content", gin.H{
		"Title":   model.Name,
		"Model":   model,
		"Columns": modelColumns[slug],
		"CSRF":    h.csrf.TokenFromContext(c),
		"Active":  slug,
	})
}

func (h *PagesHandler) New(c *gin.Context) {
	h.formPage(c, 0)
}

func (h *PagesHandler) Edit(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	h.formPage(c, uint(id))
}

func (h *PagesHandler) formPage(c *gin.Context, id uint) {
	slug := c.Param("model")
	if slug == "settings" {
		c.Redirect(http.StatusFound, "/admin/settings")
		return
	}
	model, ok := admin.FindBySlug(slug)
	if !ok {
		c.Status(http.StatusNotFound)
		return
	}
	title := "Добавить " + model.Name
	if id > 0 {
		title = "Изменить " + model.Name
	}
	h.render.Page(c, 200, "admin/layout", "admin/form_content", gin.H{
		"Title":  title,
		"Model":  model,
		"Fields": modelFields[slug],
		"ID":     id,
		"CSRF":   h.csrf.TokenFromContext(c),
		"Active": slug,
	})
}
