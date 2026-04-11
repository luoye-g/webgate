package blog

import (
	"context"
	"time"

	"github.com/luoye-g/webgate/model"
	"gorm.io/gorm"
)

type Blog interface {
	Create(authorID uint64, title, content, category, tags string, status int) (*model.Blog, error)
	CreateWithScope(authorID uint64, title, content, category, tags string, status int, scope int) (*model.Blog, error)
	Update(id uint64, title, content, category, tags string, status int) (*model.Blog, error)
	UpdateWithScope(id uint64, title, content, category, tags string, status int, scope int) (*model.Blog, error)
	GetByID(id uint64) (*model.Blog, error)
	List(offset, pageSize int, status int, scope *int) ([]*model.Blog, int64, error)
	ListByAuthor(authorID uint64, offset, pageSize int) ([]*model.Blog, int64, error)
	UpdateStatus(id uint64, status int) error
	Publish(id uint64, publishedAt time.Time) error
	UpdatePublishedTime(id uint64, publishedAt time.Time) error
	IncrementViewCount(id uint64) error
	Delete(id uint64) error
}

type blog struct {
	BlogDao model.BlogDAO
}

var blogRepo Blog

// GetBlogRepo returns the singleton instance of the blog repository.
// If the repository hasn't been initialized, it creates a new instance with a BlogDAO.
func GetBlogRepo() Blog {
	if blogRepo == nil {
		blogRepo = &blog{
			BlogDao: model.NewBlogDAO(),
		}
	}
	return blogRepo
}

func SetBlogRepo(repo Blog) {
	blogRepo = repo
}

func NewBlogRepo(db *gorm.DB) Blog {
	return &blog{
		BlogDao: model.NewBlogDAO(),
	}
}

func (b *blog) Create(authorID uint64, title, content, category, tags string, status int) (*model.Blog, error) {
	return b.CreateWithScope(authorID, title, content, category, tags, status, model.BlogScopePublic)
}

func (b *blog) CreateWithScope(authorID uint64, title, content, category, tags string, status int, scope int) (*model.Blog, error) {
	blog := &model.Blog{
		UserID:    authorID,
		Title:     title,
		Content:   content,
		Category:  category,
		Tags:      tags,
		Status:    status,
		Scope:     scope,
		ViewCount: 0,
		LikeCount: 0,
	}

	err := b.BlogDao.Insert(context.Background(), blog)
	if err != nil {
		return nil, err
	}

	// 获取作者信息
	userDAO := model.NewUserDAO()
	author, err := userDAO.GetByUserID(authorID)
	if err == nil && author != nil {
		blog.AuthorName = author.UserName
	}

	return blog, nil
}

func (b *blog) Update(id uint64, title, content, category, tags string, status int) (*model.Blog, error) {
	return b.UpdateWithScope(id, title, content, category, tags, status, -1) // -1 表示不更新scope
}

func (b *blog) UpdateWithScope(id uint64, title, content, category, tags string, status int, scope int) (*model.Blog, error) {
	blog, err := b.BlogDao.FindByID(context.Background(), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	if blog == nil {
		return nil, nil
	}

	blog.Title = title
	blog.Content = content
	blog.Category = category
	blog.Tags = tags
	blog.Status = status
	
	// 只有当scope不等于-1时才更新scope字段
	if scope != -1 {
		blog.Scope = scope
	}

	err = b.BlogDao.Update(context.Background(), blog)
	if err != nil {
		return nil, err
	}

	return b.GetByID(id)
}

func (b *blog) GetByID(id uint64) (*model.Blog, error) {
	blog, err := b.BlogDao.FindByID(context.Background(), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	if blog == nil || blog.Status == model.BlogStatusDeleted {
		return nil, nil
	}

	// 获取作者信息
	userDAO := model.NewUserDAO()
	author, err := userDAO.GetByUserID(blog.UserID)
	if err == nil && author != nil {
		blog.AuthorName = author.UserName
	}

	return blog, nil
}

func (b *blog) List(offset, pageSize int, status int, scope *int) ([]*model.Blog, int64, error) {
	ctx := context.Background()
	var blogs []*model.Blog
	var total int64
	var err error

	// 处理scope参数
	scopeValue := model.BlogScopePublic // 默认为公开
	if scope != nil {
		scopeValue = *scope
	}

	// 根据status和scope查询博客
	if status > 0 {
		// 查询特定状态和范围的博客
		total, err = b.BlogDao.CountByStatusAndScope(ctx, status, scopeValue)
		if err != nil {
			return nil, 0, err
		}
		blogs, err = b.BlogDao.FindByStatusAndScope(ctx, status, scopeValue, offset, pageSize)
		if err != nil {
			return nil, 0, err
		}
	} else {
		// 查询所有非删除状态且指定范围的博客
		allBlogs, err := b.BlogDao.FindByScope(ctx, scopeValue, offset, pageSize)
		if err != nil {
			return nil, 0, err
		}
		
		// 过滤已删除的博客
		for _, blog := range allBlogs {
			if blog.Status != model.BlogStatusDeleted {
				blogs = append(blogs, blog)
			}
		}
		
		// 计算非删除博客的总数
		scopeCount, err := b.BlogDao.CountByScope(ctx, scopeValue)
		if err != nil {
			return nil, 0, err
		}
		deletedCount, err := b.BlogDao.CountByStatusAndScope(ctx, model.BlogStatusDeleted, scopeValue)
		if err != nil {
			return nil, 0, err
		}
		total = scopeCount - deletedCount
	}

	// 批量获取作者信息
	if len(blogs) > 0 {
		userDAO := model.NewUserDAO()
		for _, blog := range blogs {
			author, err := userDAO.GetByUserID(blog.UserID)
			if err == nil && author != nil {
				blog.AuthorName = author.UserName
			}
		}
	}

	return blogs, total, nil
}

func (b *blog) ListByAuthor(authorID uint64, offset, pageSize int) ([]*model.Blog, int64, error) {
	ctx := context.Background()
	total, err := b.BlogDao.CountByUserID(ctx, authorID)
	if err != nil {
		return nil, 0, err
	}

	blogs, err := b.BlogDao.FindByUserID(ctx, authorID, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return blogs, total, nil
}

func (b *blog) UpdateStatus(id uint64, status int) error {
	return b.BlogDao.UpdateStatus(context.Background(), id, status)
}

func (b *blog) Publish(id uint64, publishedAt time.Time) error {
	return b.BlogDao.Publish(context.Background(), id, publishedAt)
}

func (b *blog) UpdatePublishedTime(id uint64, publishedAt time.Time) error {
	// DAO没有直接更新发布时间的方法，使用UpdateStatus
	blog, err := b.BlogDao.FindByID(context.Background(), id)
	if err != nil {
		return err
	}
	if blog == nil {
		return gorm.ErrRecordNotFound
	}
	blog.PublishedAt = &publishedAt
	return b.BlogDao.Update(context.Background(), blog)
}

func (b *blog) IncrementViewCount(id uint64) error {
	return b.BlogDao.IncrementViewCount(context.Background(), id)
}

func (b *blog) Delete(id uint64) error {
	return b.BlogDao.Delete(context.Background(), id)
}
