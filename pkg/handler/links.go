package handler

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

type inputFrom struct {
	OriginLink string `form:"origin_link"`
}

func (h *Handler) createLink(c *gin.Context) {
	var input inputFrom
	if err := c.Bind(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	shortLink, err := h.services.CreateLink(input.OriginLink)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"short_link": shortLink,
	})
}

func (h *Handler) getLink(c *gin.Context) {
	url := c.Param("url")
	originLink, err := h.services.GetLink(url)
	if err != nil {
		NewErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	c.Redirect(http.StatusPermanentRedirect, originLink)
}

type templateStruct struct {
	ImageBase64 string
	OriginLink  string
	ShortLink   string
}

func (h *Handler) GetQR(c *gin.Context) {
	var buf bytes.Buffer

	url := c.Param("url")

	img, link, err := h.services.GetQR(url)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	templatePath := path.Join(h.staticFiles, "qr.htmpl")
	ts, err := template.ParseFiles(templatePath)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = ts.Execute(&buf, templateStruct{
		ImageBase64: base64.StdEncoding.EncodeToString(img),
		OriginLink:  link,
		ShortLink:   fmt.Sprintf("http://%s/%s", h.domain, url),
	})
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Data(http.StatusOK, "text/html", buf.Bytes())
}
