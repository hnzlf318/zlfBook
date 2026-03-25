package middlewares

import (
	"net/http"
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/core"
)

// TransactionPicturesCors enables CORS for transaction pictures so that browser-side `fetch()` can read image bytes
// and cache them in IndexedDB. `<img>` tags do not require CORS, but `fetch()` does.
func TransactionPicturesCors() core.MiddlewareHandlerFunc {
	return func(c *core.WebContext) {
		if c.Request == nil || c.Request.URL == nil {
			c.Next()
			return
		}

		// Only handle /pictures/<fileName> requests.
		if !strings.HasPrefix(c.Request.URL.Path, "/pictures/") {
			c.Next()
			return
		}

		origin := c.GetHeader("Origin")

		// When origin is provided, echo it back to avoid issues with credentials (we don't use credentials here).
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
		}
		c.Header("Vary", "Origin")
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

