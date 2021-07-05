package handler

import (
	"path"

	"github.com/gin-gonic/gin"
	"github.com/wellWINeo/ShortLink/pkg/service"
)

type Handler struct {
	services    *service.Service
	staticFiles string
	domain      string
}

func NewHandler(services *service.Service, staticFiles string,
	domain string) *Handler {
	return &Handler{
		services:    services,
		staticFiles: staticFiles,
		domain:      domain,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.POST("/create", h.createLink)
	router.GET("/:url", h.getLink)
	router.GET("/qr/:url", h.GetQR)
	router.GET("/", h.getIndex)
	// router.GET("/about", h.getAbout)
	router.Static("/website/css", path.Join(h.staticFiles, "css"))
	router.StaticFile("/about", path.Join(h.staticFiles, "about.html"))

	return router
}
