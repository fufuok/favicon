# Favicon middleware for Gin

Ignore the favicon.ico request or cache the provided icon in memory to improve performance by skipping disk access.

```go
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
```

## Examples

Ref: [example](example)

### 1. Use go:embed in Go 1.16+

```go
package main

import (
	"embed"
	"net/http"

	"github.com/fufuok/favicon"
	"github.com/gin-gonic/gin"
)

//go:embed favicon.ico
var fav embed.FS

func main() {
	r := gin.Default()
	r.Use(favicon.New(favicon.Config{
		File:       "favicon.ico",
		FileSystem: http.FS(fav),
	}))

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hi, favicon.ico")
	})

	_ = r.Run()
}

// go run use_embed.go

// curl -I http://127.0.0.1:8080/favicon.ico
// HTTP/1.1 200 OK
// Cache-Control: public, max-age=31536000
// Content-Length: 15086
// Content-Type: image/x-icon
// Date: Mon, 13 Sep 2021 07:08:02 GMT

// curl -I -XOPTIONS http://127.0.0.1:8080/favicon.ico
// HTTP/1.1 200 OK
// Allow: GET, HEAD, OPTIONS
// Content-Length: 0
// Date: Mon, 13 Sep 2021 07:11:00 GMT

// curl http://127.0.0.1:8080/favicon.ico
// Warning: Binary output can mess up your terminal. Use "--output -" to tell
// Warning: curl to output it to your terminal anyway, or consider "--output
// Warning: <FILE>" to save to a file.
```

### 2. Use favicon.ico file

```go
package main

import (
	"net/http"

	"github.com/fufuok/favicon"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(favicon.New(favicon.Config{
		File: "favicon.ico",
	}))

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hi, favicon.ico")
	})

	_ = r.Run()
}

// go run use_file.go
```

### 3. Use favicon.ico file data

```go
package main

import (
	_ "embed"
	"net/http"

	"github.com/fufuok/favicon"
	"github.com/gin-gonic/gin"
)

//go:embed favicon.ico
var favData []byte

func main() {
	r := gin.Default()
	r.Use(favicon.New(favicon.Config{
		FileData: favData,
	}))

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hi, favicon.ico")
	})

	_ = r.Run()
}

// go run use_data.go
```

