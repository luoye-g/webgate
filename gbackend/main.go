package main

import (
	"github.com/luoye-g/webgate/controller/blog"
	"github.com/luoye-g/webgate/controller/user"
	"github.com/luoye-g/webgate/middlewares"
	"github.com/luoye-g/webgate/pkg/log"
	"github.com/luoye-g/webgate/system"

	"github.com/gin-gonic/gin"
)

func main() {
	system.SystemInit()
	defer func() { _ = log.Sync() }()

	r := gin.New()

	// 全局中间件：Recovery -> Logger -> NoCache
	r.Use(middlewares.Recovery())
	r.Use(middlewares.Logger())
	r.Use(middlewares.NoCache())

	// 公开API
	// login
	r.POST("/api/user/login", user.Login)
	// register
	r.POST("/api/user/register", user.Register)
	// 获取博客列表（公开）
	r.GET("/api/blogs", blog.ListBlogs)
	r.GET("/api/blogs/:id", blog.GetBlog)

	// 需要认证的API
	authGroup := r.Group("/api")
	authGroup.Use(middlewares.Auth())
	{
		// user info
		authGroup.GET("/user/info", user.GetUserInfo)
		authGroup.PUT("/user/info", user.UpdateUserInfo)
		authGroup.POST("/user/logout", user.Logout)

		// 博客相关API
		authGroup.POST("/blogs", blog.CreateBlog)
		authGroup.PUT("/blogs/:id", blog.UpdateBlog)
		authGroup.DELETE("/blogs/:id", blog.DeleteBlog)
		authGroup.POST("/blogs/publish", blog.PublishBlog)
		authGroup.GET("/my-blogs", blog.ListMyBlogs)
	}

	log.Info("server_starting", "addr", ":9001")
	if err := r.Run(":9001"); err != nil {
		log.Error("server_exited", "err", err.Error())
	}
}
