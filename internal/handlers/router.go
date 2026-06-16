package handlers

import (
	"html/template"
	"io/fs"
	"net/http"
	"time"

	"github.com/burj/comic/internal/config"
	"github.com/burj/comic/internal/handlers/admin"
	"github.com/burj/comic/internal/handlers/api"
	"github.com/burj/comic/internal/handlers/public"
	"github.com/burj/comic/internal/middleware"
	"github.com/burj/comic/internal/render"
	"github.com/burj/comic/internal/services"
	"github.com/burj/comic/internal/storage"
	"github.com/burj/comic/internal/tickets"
	"github.com/gin-gonic/gin"
)

type Deps struct {
	Config    config.Config
	Auth      *services.AuthService
	Events    *services.EventService
	Videos    *services.VideoService
	Merch     *services.MerchService
	Photos    *services.PhotoService
	Settings    *services.SiteSettingsService
	URLPreview  *services.URLPreviewService
	Uploader    *storage.Uploader
	UploadDir string
	Templates *template.Template
	StaticFS  fs.FS
}

func NewRouter(deps Deps) *gin.Engine {
	if deps.Config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(middleware.SecurityHeaders(deps.Config.IsProduction()))
	r.Use(func(c *gin.Context) {
		c.Set("cfg", deps.Config)
		c.Next()
	})

	if len(deps.Config.TrustedProxies) > 0 {
		_ = r.SetTrustedProxies(deps.Config.TrustedProxies)
	}

	limiter := middleware.NewRateLimiter()
	renderer := render.New(deps.Templates)
	csrf := middleware.NewCSRF(deps.Config)

	r.StaticFS("/static", http.FS(deps.StaticFS))
	if deps.UploadDir != "" {
		r.Static("/uploads", deps.UploadDir)
	}

	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	r.GET("/yandex_30040fe237d3d610.html", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=UTF-8")
		c.String(http.StatusOK, "<html>\n<head>\n<meta http-equiv=\"Content-Type\" content=\"text/html; charset=UTF-8\">\n</head>\n<body>Verification: 30040fe237d3d610</body>\n</html>")
	})

	pub := r.Group("/")
	{
		home := public.NewHomeHandler(deps.Events, deps.Videos, deps.Photos, deps.Merch, deps.Settings, renderer)
		events := public.NewEventsHandler(deps.Events, deps.Settings, renderer)
		videos := public.NewVideosHandler(deps.Videos, deps.Settings, renderer)
		photos := public.NewPhotosHandler(deps.Photos, deps.Settings, renderer)
		merch := public.NewMerchHandler(deps.Merch, deps.Settings, renderer)

		pub.GET("/", home.Index)
		pub.GET("/events", events.List)
		pub.GET("/videos", videos.List)
		pub.GET("/photos", photos.List)
		pub.GET("/merch", merch.List)
	}

	authHandler := admin.NewAuthHandler(deps.Auth, csrf, deps.Config, renderer)
	r.GET("/admin/login", authHandler.LoginPage)

	adminGroup := r.Group("/admin")
	adminGroup.Use(csrf.EnsureToken())
	adminGroup.POST("/login",
		middleware.RateLimit(limiter, 10, 15*time.Minute),
		csrf.Protect(),
		authHandler.Login,
	)

	protectedAdmin := adminGroup.Group("")
	protectedAdmin.Use(middleware.AdminAuth(deps.Auth))
	protectedAdmin.Use(csrf.EnsureToken())
	{
		dashboard := admin.NewDashboardHandler(
			deps.Events, deps.Videos, deps.Merch, deps.Photos, deps.Settings, csrf, renderer,
		)
		pages := admin.NewPagesHandler(csrf, renderer)
		settingsPage := admin.NewSettingsHandler(csrf, renderer)

		accountPage := admin.NewAccountHandler(csrf, renderer)

		protectedAdmin.POST("/logout", csrf.Protect(), authHandler.Logout)
		protectedAdmin.GET("", dashboard.Index)
		protectedAdmin.GET("/account", accountPage.Page)
		protectedAdmin.GET("/settings", settingsPage.Page)
		protectedAdmin.GET("/:model/new", pages.New)
		protectedAdmin.GET("/:model/:id/edit", pages.Edit)
		protectedAdmin.GET("/:model", pages.List)
	}

	apiGroup := r.Group("/api")
	apiGroup.Use(middleware.AdminAuth(deps.Auth))
	apiGroup.Use(csrf.EnsureToken())
	{
		eventAPI := api.NewEventHandler(deps.Events, deps.URLPreview, deps.Uploader)
		videoAPI := api.NewVideoHandler(deps.Videos)
		merchAPI := api.NewMerchHandler(deps.Merch)
		photoAPI := api.NewPhotoHandler(deps.Photos)
		settingsAPI := api.NewSettingsHandler(deps.Settings)
		accountAPI := api.NewAccountHandler(deps.Auth)
		uploadAPI := api.NewUploadHandler(deps.Uploader)
		ticketCatalogAPI := api.NewTicketCatalogHandler(tickets.NewCatalog(), deps.Settings, deps.Events)

		mutating := apiGroup.Group("")
		mutating.Use(csrf.Protect())
		{
			mutating.POST("/upload",
				middleware.RateLimit(limiter, 60, time.Hour),
				uploadAPI.Upload,
			)
			mutating.PUT("/account/password", accountAPI.ChangePassword)
			mutating.POST("/events", eventAPI.Create)
			mutating.PUT("/events/:id", eventAPI.Update)
			mutating.DELETE("/events/:id", eventAPI.Delete)
			mutating.POST("/videos", videoAPI.Create)
			mutating.PUT("/videos/:id", videoAPI.Update)
			mutating.DELETE("/videos/:id", videoAPI.Delete)
			mutating.POST("/merch", merchAPI.Create)
			mutating.PUT("/merch/:id", merchAPI.Update)
			mutating.DELETE("/merch/:id", merchAPI.Delete)
			mutating.POST("/photos", photoAPI.Create)
			mutating.PUT("/photos/:id", photoAPI.Update)
			mutating.DELETE("/photos/:id", photoAPI.Delete)
			mutating.PUT("/settings", settingsAPI.Update)
		}

		apiGroup.GET("/events", eventAPI.List)
		apiGroup.GET("/events/preview-ticket", eventAPI.PreviewTicket)
		apiGroup.GET("/events/import-poster", eventAPI.ImportPoster)
		apiGroup.GET("/events/:id", eventAPI.Get)
		apiGroup.GET("/videos", videoAPI.List)
		apiGroup.GET("/videos/:id", videoAPI.Get)
		apiGroup.GET("/merch", merchAPI.List)
		apiGroup.GET("/merch/:id", merchAPI.Get)
		apiGroup.GET("/photos", photoAPI.List)
		apiGroup.GET("/photos/:id", photoAPI.Get)
		apiGroup.GET("/settings", settingsAPI.Get)
		apiGroup.GET("/ticket-catalog/providers", ticketCatalogAPI.Providers)
		apiGroup.GET("/ticket-catalog/:source/events", ticketCatalogAPI.Events)
	}

	return r
}
