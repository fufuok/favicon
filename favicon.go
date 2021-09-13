package favicon

import (
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Config defines the config for middleware.
// Code from: github.com/gofiber/fiber/v2/middleware/favicon, thx.
type Config struct {
	// File holds the path to an actual favicon that will be cached
	//
	// Optional. Default: ""
	File string `json:"file"`

	// FileSystem is an optional alternate filesystem to search for the favicon in.
	// An example of this could be an embedded or network filesystem
	// Need to be used with the File parameter
	//
	// Optional. Default: nil
	FileSystem http.FileSystem `json:"-"`

	// FileData is an actual favicon file data that will be cached
	//
	// Optional. Default: nil
	FileData []byte `json:"file_data"`

	// CacheControl defines how the Cache-Control header in the response should be set
	//
	// Optional. Default: "public, max-age=31536000"
	CacheControl string `json:"cache_control"`
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	CacheControl: "public, max-age=31536000",
}

func New(config ...Config) gin.HandlerFunc {
	// Set default config
	cfg := ConfigDefault

	// Override config if provided
	if len(config) > 0 {
		cfg = config[0]

		if cfg.File == "" {
			cfg.File = ConfigDefault.File
		}
		if cfg.CacheControl == "" {
			cfg.CacheControl = ConfigDefault.CacheControl
		}
	}

	// Load icon if provided
	var (
		err     error
		iconLen string
		icon    = cfg.FileData
	)
	if cfg.File != "" {
		// read from configured filesystem if present
		if cfg.FileSystem != nil {
			f, err := cfg.FileSystem.Open(cfg.File)
			if err != nil {
				panic(err)
			}
			if icon, err = io.ReadAll(f); err != nil {
				panic(err)
			}
		} else {
			if icon, err = os.ReadFile(cfg.File); err != nil {
				panic(err)
			}
		}
	}

	iconLen = strconv.Itoa(len(icon))

	// Return new handler
	return func(c *gin.Context) {
		// Only respond to favicon requests
		if len(c.Request.RequestURI) != 12 || c.Request.RequestURI != "/favicon.ico" {
			return
		}

		// Only allow GET, HEAD and OPTIONS requests
		if c.Request.Method != "GET" && c.Request.Method != "HEAD" {
			c.Header("Allow", "GET, HEAD, OPTIONS")
			c.Header("Content-Length", "0")
			if c.Request.Method != "OPTIONS" {
				c.AbortWithStatus(http.StatusMethodNotAllowed)
			} else {
				c.AbortWithStatus(http.StatusOK)
			}
			return
		}

		// No content
		if iconLen == "0" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// Serve cached favicon
		c.Header("Content-Length", iconLen)
		c.Header("Cache-Control", cfg.CacheControl)
		c.Data(http.StatusOK, "image/x-icon", icon)
		c.Abort()
		return
	}
}
