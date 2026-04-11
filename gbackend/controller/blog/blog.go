package blog

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/luoye-g/webgate/controller/response"
	"github.com/luoye-g/webgate/model"
	"github.com/luoye-g/webgate/services/blog"
)

type BlogParam struct {
	Title    string `json:"title" binding:"required,min=1,max=255"`
	Content  string `json:"content" binding:"required"`
	Category string `json:"category"`
	Tags     string `json:"tags"`
	Status   int    `json:"status"`
	Scope    int    `json:"scope"`
}
type PublishParam struct {
	BlogID uint64 `json:"blog_id" binding:"required"`
}

// CreateBlog 创建博客
func CreateBlog(ctx *gin.Context) {
	var param BlogParam
	if err := ctx.ShouldBindBodyWithJSON(&param); err != nil {
		response.InvalidParam(ctx, "invalid_param")
		return
	}

	// 获取当前用户ID（从session中获取）
	userID := ctx.GetUint64("user_id")
	if userID == 0 {
		response.RequireLogin(ctx)
		return
	}

	// 设置默认scope为公开
	if param.Scope == 0 {
		param.Scope = model.BlogScopePublic
	}

	blogService := blog.GetBlogService()
	blogItem, err := blogService.CreateBlogWithScope(userID, param.Title, param.Content, param.Category, param.Tags, param.Scope)
	if err != nil {
		response.Failed(ctx, err)
		return
	}

	response.Success(ctx, blogItem)
}

// UpdateBlog 编辑博客
func UpdateBlog(ctx *gin.Context) {
	blogID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response.InvalidParam(ctx, "invalid_blog_id")
		return
	}

	var param BlogParam
	if err := ctx.ShouldBindBodyWithJSON(&param); err != nil {
		response.InvalidParam(ctx, "invalid_param")
		return
	}

	// 获取当前用户ID
	userID := ctx.GetUint64("user_id")
	if userID == 0 {
		response.RequireLogin(ctx)
		return
	}

	blogService := blog.GetBlogService()
	blogItem, appErr := blogService.UpdateBlogWithScope(blogID, userID, param.Title, param.Content, param.Category, param.Tags, param.Status, param.Scope)
	if appErr != nil {
		response.Failed(ctx, appErr)
		return
	}

	response.Success(ctx, blogItem)
}

// DeleteBlog 删除博客
func DeleteBlog(ctx *gin.Context) {
	blogID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response.InvalidParam(ctx, "invalid_blog_id")
		return
	}

	// 获取当前用户ID
	userID := ctx.GetUint64("user_id")
	if userID == 0 {
		response.RequireLogin(ctx)
		return
	}

	blogService := blog.GetBlogService()
	appErr := blogService.DeleteBlog(blogID, userID)
	if appErr != nil {
		response.Failed(ctx, appErr)
		return
	}

	response.Success(ctx, nil)
}

// PublishBlog 发布博客
func PublishBlog(ctx *gin.Context) {
	var param PublishParam
	if err := ctx.ShouldBindBodyWithJSON(&param); err != nil {
		response.InvalidParam(ctx, "invalid_param")
		return
	}

	// 获取当前用户ID
	userID := ctx.GetUint64("user_id")
	if userID == 0 {
		response.RequireLogin(ctx)
		return
	}

	blogService := blog.GetBlogService()
	err := blogService.PublishBlog(param.BlogID, userID)
	if err != nil {
		response.Failed(ctx, err)
		return
	}

	response.Success(ctx, nil)
}

// GetBlog 获取博客详情
func GetBlog(ctx *gin.Context) {
	blogID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response.InvalidParam(ctx, "invalid_blog_id")
		return
	}

	blogService := blog.GetBlogService()
	blogItem, appErr := blogService.GetBlog(blogID)
	if appErr != nil {
		response.Failed(ctx, appErr)
		return
	}

	response.Success(ctx, blogItem)
}

// ListBlogs 获取博客列表
func ListBlogs(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	blogService := blog.GetBlogService()
	blogs, total, err := blogService.ListBlogs(offset, pageSize)
	if err != nil {
		response.Failed(ctx, err)
		return
	}

	response.Success(ctx, gin.H{
		"blogs": blogs,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// ListMyBlogs 获取我的博客列表
func ListMyBlogs(ctx *gin.Context) {
	userID := ctx.GetUint64("user_id")
	if userID == 0 {
		response.RequireLogin(ctx)
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	blogService := blog.GetBlogService()
	blogs, total, err := blogService.ListMyBlogs(userID, offset, pageSize)
	if err != nil {
		response.Failed(ctx, err)
		return
	}

	response.Success(ctx, gin.H{
		"blogs": blogs,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// ListBlogsByScope 获取指定范围的博客列表
func ListBlogsByScope(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	status, _ := strconv.Atoi(ctx.DefaultQuery("status", "1")) // 默认只显示已发布的
	scope, _ := strconv.Atoi(ctx.DefaultQuery("scope", "1"))   // 默认只显示公开的

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	blogService := blog.GetBlogService()
	blogs, total, err := blogService.ListBlogsByScope(offset, pageSize, status, scope)
	if err != nil {
		response.Failed(ctx, err)
		return
	}

	response.Success(ctx, gin.H{
		"blogs": blogs,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}
