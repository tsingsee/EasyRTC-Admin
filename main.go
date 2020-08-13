package main

import (
	"fmt"

	"github.com/unrolled/secure"

	"github.com/gin-gonic/gin"
	"tsingsee.com/adminserver/app"
	"tsingsee.com/adminserver/routes"
)

func main() {
	r := gin.Default()
	https := gin.Default()
	app := app.NewApp()

	if app.Config().HttpsPort > 0 {
		httpsPort := fmt.Sprintf(":%d", app.Config().HttpsPort)
		https.Use(TlsHandler(httpsPort))
		routes.Setup(https, app)
		go https.RunTLS(httpsPort, app.Config().CertPath, app.Config().KeyPath)
	}

	routes.Setup(r, app)

	r.Run(fmt.Sprintf(":%d", app.Config().Port))
}

// 初始 TLS
func TlsHandler(httpsPort string) gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     httpsPort,
		})
		err := secureMiddleware.Process(c.Writer, c.Request)
		if err != nil {
			return
		}
		c.Next()
	}
}
