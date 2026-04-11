package blog

import (
	"fmt"
	"time"

	"github.com/luoye-g/webgate/errors"
	"github.com/luoye-g/webgate/model"
	"github.com/luoye-g/webgate/repository/blog"
)

type BlogService struct {
	blogRepo blog.Blog
}

var blogService *BlogService

func GetBlogService() *BlogService {
	if blogService == nil {
		blogService = &BlogService{blogRepo: blog.GetBlogRepo()}
	}
	return blogService
}

// BlogResponse 博客响应DTO

type BlogResponse struct {
	ID          uint64     `json:"id"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	AuthorID    uint64     `json:"author_id"`
	AuthorName  string     `json:"author_name"`
	Status      int        `json:"status"`
	Scope       int        `json:"scope"`
	Category    string     `json:"category"`
	Tags        string     `json:"tags"`
	ViewCount   uint32     `json:"view_count"`
	LikeCount   uint32     `json:"like_count"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
}

// CreateBlog 创建博客
func (service *BlogService) CreateBlog(authorID uint64, title, content, category, tags string) (*BlogResponse, *errors.AppError) {
	return service.CreateBlogWithScope(authorID, title, content, category, tags, model.BlogScopePublic)
}

// CreateBlogWithScope 创建博客（指定范围）
func (service *BlogService) CreateBlogWithScope(authorID uint64, title, content, category, tags string, scope int) (*BlogResponse, *errors.AppError) {
	// 验证参数
	if title == "" {
		return nil, errors.NewError(errors.ValidationError, fmt.Errorf("标题不能为空"))
	}
	if content == "" {
		return nil, errors.NewError(errors.ValidationError, fmt.Errorf("内容不能为空"))
	}
	if len(title) > 255 {
		return nil, errors.NewError(errors.ValidationError, fmt.Errorf("标题长度不能超过255个字符"))
	}

	// 创建博客
	blogModel, err := service.blogRepo.CreateWithScope(authorID, title, content, category, tags, model.BlogStatusDraft, scope)
	if err != nil {
		return nil, errors.NewError(errors.InternalError, err)
	}

	return service.convertToResponse(blogModel), nil
}

// UpdateBlog 更新博客
func (service *BlogService) UpdateBlog(blogID, authorID uint64, title, content, category, tags string, status int) (*BlogResponse, *errors.AppError) {
	return service.UpdateBlogWithScope(blogID, authorID, title, content, category, tags, status, -1) // -1 表示不更新scope
}

// UpdateBlogWithScope 更新博客（包含范围）
func (service *BlogService) UpdateBlogWithScope(blogID, authorID uint64, title, content, category, tags string, status int, scope int) (*BlogResponse, *errors.AppError) {
	// 验证参数
	if title == "" {
		return nil, errors.NewError(errors.ValidationError, fmt.Errorf("标题不能为空"))
	}
	if content == "" {
		return nil, errors.NewError(errors.ValidationError, fmt.Errorf("内容不能为空"))
	}
	if len(title) > 255 {
		return nil, errors.NewError(errors.ValidationError, fmt.Errorf("标题长度不能超过255个字符"))
	}

	// 检查博客是否存在且属于该作者
	blogModel, err := service.blogRepo.GetByID(blogID)
	if err != nil {
		return nil, errors.NewError(errors.InternalError, err)
	}
	if blogModel == nil {
		return nil, errors.NewError(errors.NotFoundError, fmt.Errorf("博客不存在"))
	}
	if blogModel.UserID != authorID {
		return nil, errors.NewError(errors.PermissionError, fmt.Errorf("无权限修改此博客"))
	}

	// 更新博客
	blogModel, err = service.blogRepo.UpdateWithScope(blogID, title, content, category, tags, status, scope)
	if err != nil {
		return nil, errors.NewError(errors.InternalError, err)
	}

	// 如果状态变更为发布，设置发布时间
	if status == model.BlogStatusPublished && (blogModel.PublishedAt == nil || blogModel.PublishedAt.IsZero()) {
		now := time.Now()
		blogModel.PublishedAt = &now
		err = service.blogRepo.UpdatePublishedTime(blogID, now)
		if err != nil {
			return nil, errors.NewError(errors.InternalError, err)
		}
	}

	return service.convertToResponse(blogModel), nil
}

// DeleteBlog 删除博客
func (service *BlogService) DeleteBlog(blogID, authorID uint64) *errors.AppError {
	// 检查博客是否存在且属于该作者
	blogModel, err := service.blogRepo.GetByID(blogID)
	if err != nil {
		return errors.NewError(errors.InternalError, err)
	}
	if blogModel == nil {
		return errors.NewError(errors.NotFoundError, fmt.Errorf("博客不存在"))
	}
	if blogModel.UserID != authorID {
		return errors.NewError(errors.PermissionError, fmt.Errorf("无权限删除此博客"))
	}

	// 软删除：将状态设为已删除
	err = service.blogRepo.UpdateStatus(blogID, model.BlogStatusDeleted)
	if err != nil {
		return errors.NewError(errors.InternalError, err)
	}

	return nil
}

// PublishBlog 发布博客
func (service *BlogService) PublishBlog(blogID, authorID uint64) *errors.AppError {
	// 检查博客是否存在且属于该作者
	blogModel, err := service.blogRepo.GetByID(blogID)
	if err != nil {
		return errors.NewError(errors.InternalError, err)
	}
	if blogModel == nil {
		return errors.NewError(errors.NotFoundError, fmt.Errorf("博客不存在"))
	}
	if blogModel.UserID != authorID {
		return errors.NewError(errors.PermissionError, fmt.Errorf("无权限发布此博客"))
	}

	// 检查是否已经发布
	if blogModel.Status == model.BlogStatusPublished {
		return errors.NewError(errors.ValidationError, fmt.Errorf("博客已经发布"))
	}

	// 发布博客
	now := time.Now()
	err = service.blogRepo.Publish(blogID, now)
	if err != nil {
		return errors.NewError(errors.InternalError, err)
	}

	return nil
}

// GetBlog 获取博客详情
func (service *BlogService) GetBlog(blogID uint64) (*BlogResponse, *errors.AppError) {
	blogModel, err := service.blogRepo.GetByID(blogID)
	if err != nil {
		return nil, errors.NewError(errors.InternalError, err)
	}
	if blogModel == nil {
		return nil, errors.NewError(errors.NotFoundError, fmt.Errorf("博客不存在"))
	}

	// 检查博客是否为公开状态
	if blogModel.Scope != model.BlogScopePublic {
		return nil, errors.NewError(errors.NotFoundError, fmt.Errorf("博客不存在"))
	}

	// 增加浏览次数
	err = service.blogRepo.IncrementViewCount(blogID)
	if err != nil {
		// 浏览次数增加失败不影响返回博客内容
		// 记录日志即可
	}

	return service.convertToResponse(blogModel), nil
}

// ListBlogs 获取博客列表
func (service *BlogService) ListBlogs(offset, pageSize int) ([]*BlogResponse, int64, *errors.AppError) {
	// 默认只查询公开博客
	scope := model.BlogScopePublic
	blogs, total, err := service.blogRepo.List(offset, pageSize, model.BlogStatusPublished, &scope)
	if err != nil {
		return nil, 0, errors.NewError(errors.InternalError, err)
	}

	responses := make([]*BlogResponse, len(blogs))
	for i, blog := range blogs {
		responses[i] = service.convertToResponse(blog)
	}

	return responses, total, nil
}

// ListMyBlogs 获取我的博客列表
func (service *BlogService) ListMyBlogs(authorID uint64, offset, pageSize int) ([]*BlogResponse, int64, *errors.AppError) {
	blogs, total, err := service.blogRepo.ListByAuthor(authorID, offset, pageSize)
	if err != nil {
		return nil, 0, errors.NewError(errors.InternalError, err)
	}

	responses := make([]*BlogResponse, len(blogs))
	for i, blog := range blogs {
		responses[i] = service.convertToResponse(blog)
	}

	return responses, total, nil
}

// ListBlogsByScope 获取指定范围的博客列表
func (service *BlogService) ListBlogsByScope(offset, pageSize int, status int, scope int) ([]*BlogResponse, int64, *errors.AppError) {
	blogs, total, err := service.blogRepo.List(offset, pageSize, status, &scope)
	if err != nil {
		return nil, 0, errors.NewError(errors.InternalError, err)
	}

	responses := make([]*BlogResponse, len(blogs))
	for i, blog := range blogs {
		responses[i] = service.convertToResponse(blog)
	}

	return responses, total, nil
}

// convertToResponse 转换为响应格式
func (service *BlogService) convertToResponse(blogModel *model.Blog) *BlogResponse {
	return &BlogResponse{
		ID:          blogModel.ID,
		Title:       blogModel.Title,
		Content:     blogModel.Content,
		AuthorID:    blogModel.UserID,
		AuthorName:  blogModel.AuthorName,
		Status:      blogModel.Status,
		Scope:       blogModel.Scope,
		Category:    blogModel.Category,
		Tags:        blogModel.Tags,
		ViewCount:   blogModel.ViewCount,
		LikeCount:   blogModel.LikeCount,
		CreatedAt:   blogModel.CreatedAt,
		UpdatedAt:   blogModel.UpdatedAt,
		PublishedAt: blogModel.PublishedAt,
	}
}
