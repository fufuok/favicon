package main

import (
	_ "embed"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fufuok/favicon"
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
