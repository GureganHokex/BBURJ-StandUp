package web

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"strings"
)

//go:embed templates/*
var templateFS embed.FS

//go:embed static/css/* static/js/* static/img/*
var embeddedStatic embed.FS

func Templates() *template.Template {
	return template.Must(template.New("").Funcs(template.FuncMap{
		"formatPrice":        formatPriceRub,
		"formatEventDate":    formatEventDateCard,
		"formatEventDay":     formatEventDay,
		"formatEventMeta":    formatEventMeta,
		"upper":              upperASCII,
		"attr":               attrText,
		"telegramContactURL": telegramContactURL,
	}).ParseFS(templateFS,
		"templates/layouts/*.html",
		"templates/public/*.html",
		"templates/admin/*.html",
		"templates/admin/partials/*.html",
		"templates/public/partials/*.html",
	))
}

func StaticFS() fs.FS {
	sub, err := fs.Sub(embeddedStatic, "static")
	if err != nil {
		panic(err)
	}
	return sub
}

func formatPriceRub(kopecks int) string {
	rub := kopecks / 100
	kop := kopecks % 100
	if kop == 0 {
		return fmt.Sprintf("%d ₽", rub)
	}
	return fmt.Sprintf("%d,%02d ₽", rub, kop)
}

func telegramContactURL(handle string) string {
	h := strings.TrimSpace(handle)
	if h == "" {
		return ""
	}
	if strings.HasPrefix(h, "http://") || strings.HasPrefix(h, "https://") {
		return h
	}
	return "https://t.me/" + strings.TrimPrefix(h, "@")
}
