package api

import (
	"embed"
	"io/fs"
	"net/http"

	"workemailchecker/internal/config"
	"workemailchecker/internal/validator"

	"github.com/gin-gonic/gin"
)

//go:embed static/*
var staticFiles embed.FS

func SetupRouter(cfg *config.Config) *gin.Engine {
	go func() {
		if err := validator.LoadFreeProviders(cfg.FreeProvidersURL); err != nil {
		}
	}()

	validator.SetOverrides(cfg.CorporateOverrides, cfg.PersonalOverrides)

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	rateLimiter := NewRateLimiter(float64(cfg.RateLimitRPS), cfg.RateLimitBurst)
	aiLimiter := NewRateLimiter(cfg.AIRateLimitRPS, cfg.AIRateLimitBurst)

	api := router.Group("/api")
	{
		api.POST("/check", toGin(rateLimiter.RateLimit(EmailCheckHandler(cfg, aiLimiter))))
		api.GET("/health", toGin(HealthCheckHandler))
	}

	staticFS, err := fs.Sub(staticFiles, "static")
	if err == nil {
		router.GET("/", func(c *gin.Context) {
			b, rerr := fs.ReadFile(staticFS, "index.html")
			if rerr != nil {
				c.Status(http.StatusNotFound)
				return
			}
			c.Data(http.StatusOK, "text/html; charset=utf-8", b)
		})
		router.HEAD("/", func(c *gin.Context) {
			if _, rerr := fs.ReadFile(staticFS, "index.html"); rerr != nil {
				c.Status(http.StatusNotFound)
				return
			}
			c.Status(http.StatusOK)
		})
		router.GET("/docs", func(c *gin.Context) {
			b, rerr := fs.ReadFile(staticFS, "docs.html")
			if rerr != nil {
				c.Status(http.StatusNotFound)
				return
			}
			c.Data(http.StatusOK, "text/html; charset=utf-8", b)
		})
		router.HEAD("/docs", func(c *gin.Context) {
			if _, rerr := fs.ReadFile(staticFS, "docs.html"); rerr != nil {
				c.Status(http.StatusNotFound)
				return
			}
			c.Status(http.StatusOK)
		})
	} else {
		router.GET("/", func(c *gin.Context) { c.File("./static/index.html") })
		router.HEAD("/", func(c *gin.Context) { c.Status(http.StatusOK) })
		router.GET("/docs", func(c *gin.Context) { c.File("./static/docs.html") })
		router.HEAD("/docs", func(c *gin.Context) { c.Status(http.StatusOK) })
	}

	return router
}

func toGin(h http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(c.Writer, c.Request)
	}
}
