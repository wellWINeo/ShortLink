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

func (h *Handler) createLink(c *gin.Context) {
	fmt.Print("Handler execution started\n")

	link := c.PostForm("origin_link")

	if link == "" {
		NewErrorResponse(c, http.StatusBadRequest, "Can't parse value from Form")
		return
	}

	fmt.Printf("Value: %s\n", c.PostForm("origin_link"))

	shortLink, err := h.services.CreateLink(link)
	if err != nil {
		fmt.Printf("Error after CreateLink service: %s\n", err.Error())
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.Request.Method = "GET"
	c.Redirect(http.StatusMovedPermanently,
		fmt.Sprintf("http://%s/qr/%s", h.domain, shortLink))
}

func (h *Handler) getLink(c *gin.Context) {
	fmt.Printf("Path: %s\n", c.Request.URL)
	url := c.Param("url")
	fmt.Printf("Url: %s\n", url)

	if url == "" {
		NewErrorResponse(c, http.StatusBadRequest, "Empty URL")
	}

	originLink, err := h.services.GetLink(url)
	if err != nil {
		NewErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	c.Redirect(http.StatusMovedPermanently, originLink)
}

type templateStruct struct {
	ImageBase64 string
	OriginLink  string
	ShortLink   string
}

func (h *Handler) getQR(c *gin.Context) {
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
