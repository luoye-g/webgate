package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

func main() {
	GinHttps()
}

func GinHttps() error {

	r := gin.Default()
	r.GET("/api", func(c *gin.Context) {
		c.String(200, "test for 【%s】", "https")
	})
	return r.Run(":8080")
	// r.Use(TLSHandler(443))

	// return r.RunTLS(":"+strconv.Itoa(443), "./ssh/luoye-g.top.pem", "./ssh/luoye-g.top.rsa.key")
}

func TLSHandler(port int) gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     ":" + strconv.Itoa(port),
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		// If there was an error, do not continue.
		if err != nil {
			return
		}
		c.Next()
	}
}
