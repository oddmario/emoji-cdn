package routes

import (
	"io/fs"
	"mario/emoji-cdn/constants"
	"mario/emoji-cdn/utils"
	"net/http"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

func Emoji(c *gin.Context) {
	emoji, ok := c.Params.Get("emoji")
	if !ok || len(emoji) <= 0 {
		c.Status(http.StatusBadRequest)

		return
	}

	style, ok := c.GetQuery("style")
	if !ok || len(style) <= 0 {
		c.Status(http.StatusUnprocessableEntity)
		c.Writer.WriteString("missing the `style` query param")

		return
	}

	if !slices.Contains(constants.EmojipediaPlatforms, style) {
		c.Status(http.StatusUnprocessableEntity)
		c.Writer.WriteString("invalid emoji style specified")

		return
	}

	var foundFiles []string = []string{}

	filepath.WalkDir(filepath.Join(".", ".emojis-db", emoji), func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}

		if strings.HasPrefix(d.Name(), style+".") {
			foundFiles = append(foundFiles, d.Name())
		}

		return nil
	})

	if len(foundFiles) <= 0 {
		c.Status(http.StatusNotFound)
		c.Writer.WriteString("emoji not found")

		return
	}

	oneWeek := int64(60 * 60 * 24 * 7)
	oneMonth := int64(60 * 60 * 24 * 30)

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Cache-Control", "public, s-maxage="+utils.I64ToStr(oneMonth)+", max-age="+utils.I64ToStr(oneWeek))

	c.File(filepath.Join(".", ".emojis-db", emoji, foundFiles[0]))
}
