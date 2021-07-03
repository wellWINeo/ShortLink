package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getIndex(c *gin.Context) {
	newPath := path.Join(h.staticFiles, "new.html")
	newHTML, err := ioutil.ReadFile(newPath)
	if err != nil {
		NewErrorResponse(c, http.StatusNotFound,
			fmt.Sprintf("HTML file not found: %s", newPath))
		return
	}

	c.Data(http.StatusOK, "text/html", newHTML)
}

func (h *Handler) getAbout(c *gin.Context) {

}
