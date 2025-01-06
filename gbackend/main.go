package main

import (
	"github.com/luoye-g/webgate/controller/login"
	"github.com/luoye-g/webgate/system"

	"github.com/gin-gonic/gin"
)

func main() {
	system.SystemInit()

	r := gin.Default()

	// login
	r.GET("/api/login", login.Login)

	r.Run(":9001")
}
