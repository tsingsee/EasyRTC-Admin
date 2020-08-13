package routes

import (
	"net/http"
	"time"

	"github.com/gin-contrib/static"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"tsingsee.com/adminserver/app"
	"tsingsee.com/adminserver/server"
)

func Setup(r *gin.Engine, app *app.App) {
	gin.SetMode(gin.DebugMode)

	r.Use(static.Serve("/admin", static.LocalFile("./www", true)))

	admin := r.Group("/admin")
	admin.Use(errorMiddleware, timeoutMiddleware(5*time.Second))
	{

		admin.POST("/captcha-id", handleCaptchaId)
		admin.GET("/captcha-id", handleCaptchaId)
		admin.GET("/captcha/:id", gin.WrapH(captcha.Server(captcha.StdWidth, captcha.StdHeight)))

		passport := admin.Group("/passport")
		{
			server := server.NewPassportServer(app)
			passport.POST("/signup", server.Signup)
			passport.POST("/login", server.Login)
			passport.POST("/logout", server.Logout)
			passport.POST("/info", authMiddleware(app), server.Info)
			passport.POST("/modify", authMiddleware(app), server.Modify)
		}

		roomGroup := admin.Group("/room", authMiddleware(app))
		{
			roomServer := server.NewRoomServer(app)
			roomGroup.POST("/info", roomServer.Info)
			roomGroup.POST("/list", roomServer.List)
			roomGroup.POST("/create", roomServer.Create)
			roomGroup.POST("/modify", roomServer.Modify)
			roomGroup.POST("/delete", roomServer.Delete)
			roomGroup.POST("/token", roomServer.Token)
		}

		conferenceGroup := admin.Group("/conference", authMiddleware(app))
		{
			conferenceServer := server.NewConferenceServer(app)
			conferenceGroup.POST("/info", conferenceServer.Info)
			conferenceGroup.POST("/runing", conferenceServer.Runing)
			conferenceGroup.POST("/dispose", conferenceServer.Dispose)
			conferenceGroup.POST("/lock", conferenceServer.Lock)
			conferenceGroup.POST("/unlock", conferenceServer.Unlock)
			conferenceGroup.POST("/history", conferenceServer.History)
			conferenceGroup.POST("/action", conferenceServer.Action)
		}

		recordGroup := admin.Group("/record", authMiddleware(app))
		{
			recordServer := server.NewRecordServer(app)
			recordGroup.POST("/info", recordServer.Info)
			recordGroup.POST("/list", recordServer.List)
			recordGroup.POST("/delete", recordServer.Delete)
		}
	}
}

func handleCaptchaId(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"id": captcha.NewLen(6),
	})
}
