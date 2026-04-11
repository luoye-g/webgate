package model

import (
	"context"
	"time"

	"github.com/luoye-g/webgate/pkg/mysql"
	"gorm.io/gorm"
)

const (
	BlogStatusDraft     = 0 // 草稿
	BlogStatusPublished = 1 // 已发布
	BlogStatusDeleted   = 2 // 已删除

	BlogScopePrivate = 0 // 私密
	BlogScopePublic  = 1 // 公开
)

type Blog struct {
	ID          uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Title       string     `gorm:"column:title;type:varchar(255);not null;default:''" json:"title"`
	Content     string     `gorm:"column:content;type:text;not null" json:"content"`
	UserID      uint64     `gorm:"column:user_id;not null" json:"user_id"`
	AuthorName  string     `gorm:"-" json:"author_name,omitempty"`
	Status      int        `gorm:"column:status;type:tinyint unsigned;not null;default:0" json:"status"`
	Scope       int        `gorm:"column:scope;type:tinyint unsigned;not null;default:1" json:"scope"`
	Category    string     `gorm:"column:category;type:varchar(50);not null;default:''" json:"category"`
	Tags        string     `gorm:"column:tags;type:varchar(255);not null;default:''" json:"tags"`
	ViewCount   uint32     `gorm:"column:view_count;type:int unsigned;not null;default:0" json:"view_count"`
	LikeCount   uint32     `gorm:"column:like_count;type:int unsigned;not null;default:0" json:"like_count"`
	CreatedAt   time.Time  `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	PublishedAt *time.Time `gorm:"column:published_at;type:datetime" json:"published_at"`
}

// TableName 指定表名
func (Blog) TableName() string {
	return "blogs"
}

// IsPublished 检查是否已发布
func (b *Blog) IsPublished() bool {
	return b.Status == BlogStatusPublished && b.PublishedAt != nil
}

// IsDraft 检查是否为草稿
func (b *Blog) IsDraft() bool {
	return b.Status == BlogStatusDraft
}

// IsDeleted 检查是否已删除
func (b *Blog) IsDeleted() bool {
	return b.Status == BlogStatusDeleted
}

// BlogDAO 博客数据访问对象接口
type BlogDAO interface {
	// 基础CRUD操作
	Insert(ctx context.Context, blog *Blog) error
	Update(ctx context.Context, blog *Blog) error
	Delete(ctx context.Context, id uint64) error
	FindByID(ctx context.Context, id uint64) (*Blog, error)
	FindByIDs(ctx context.Context, ids []uint64) ([]*Blog, error)
	
	// 查询操作
	FindAll(ctx context.Context, offset, limit int) ([]*Blog, error)
	FindByUserID(ctx context.Context, userID uint64, offset, limit int) ([]*Blog, error)
	FindByStatus(ctx context.Context, status int, offset, limit int) ([]*Blog, error)
	FindPublished(ctx context.Context, offset, limit int) ([]*Blog, error)
	FindByScope(ctx context.Context, scope int, offset, limit int) ([]*Blog, error)
	FindByStatusAndScope(ctx context.Context, status int, scope int, offset, limit int) ([]*Blog, error)
	CountAll(ctx context.Context) (int64, error)
	CountByUserID(ctx context.Context, userID uint64) (int64, error)
	CountByStatus(ctx context.Context, status int) (int64, error)
	CountPublished(ctx context.Context) (int64, error)
	CountByScope(ctx context.Context, scope int) (int64, error)
	CountByStatusAndScope(ctx context.Context, status int, scope int) (int64, error)
	
	// 业务操作
	Publish(ctx context.Context, id uint64, publishedAt time.Time) error
	IncrementViewCount(ctx context.Context, id uint64) error
	UpdateStatus(ctx context.Context, id uint64, status int) error
	FindByCategory(ctx context.Context, category string, offset, limit int) ([]*Blog, error)
	Search(ctx context.Context, keyword string, offset, limit int) ([]*Blog, error)
	
	// 批量操作
	BatchUpdateStatus(ctx context.Context, ids []uint64, status int) error
	BatchDelete(ctx context.Context, ids []uint64) error
}

// BlogDao 博客数据访问对象
type BlogDao struct{}

// NewBlogDAO 创建博客DAO实例
func NewBlogDAO() BlogDAO {
	return &BlogDao{}
}

// Insert 插入博客
func (dao *BlogDao) Insert(ctx context.Context, blog *Blog) error {
	return mysql.GetDB().WithContext(ctx).Create(blog).Error
}

// Update 更新博客
func (dao *BlogDao) Update(ctx context.Context, blog *Blog) error {
	return mysql.GetDB().WithContext(ctx).Save(blog).Error
}

// Delete 删除博客（硬删除）
func (dao *BlogDao) Delete(ctx context.Context, id uint64) error {
	return mysql.GetDB().WithContext(ctx).Where("id = ?", id).Delete(&Blog{}).Error
}

// FindByID 根据ID查找博客
func (dao *BlogDao) FindByID(ctx context.Context, id uint64) (*Blog, error) {
	var blog Blog
	err := mysql.GetDB().WithContext(ctx).Where("id = ?", id).First(&blog).Error
	if err != nil {
		return nil, err
	}
	return &blog, nil
}

// FindByIDs 根据多个ID查找博客
func (dao *BlogDao) FindByIDs(ctx context.Context, ids []uint64) ([]*Blog, error) {
	var blogs []*Blog
	err := mysql.GetDB().WithContext(ctx).Where("id IN ?", ids).Find(&blogs).Error
	return blogs, err
}

// FindByUserID 根据用户ID查找博客（排除已删除）
func (dao *BlogDao) FindByUserID(ctx context.Context, userID uint64, offset, limit int) ([]*Blog, error) {
	var blogs []*Blog
	query := mysql.GetDB().WithContext(ctx).Where("user_id = ? AND status != ?", userID, BlogStatusDeleted)
	if offset > 0 {
		query = query.Offset(offset)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Order("created_at DESC").Find(&blogs).Error
	return blogs, err
}

// FindByStatus 根据状态查找博客
func (dao *BlogDao) FindByStatus(ctx context.Context, status int, offset, limit int) ([]*Blog, error) {
	var blogs []*Blog
	query := mysql.GetDB().WithContext(ctx).Where("status = ?", status)
	if offset > 0 {
		query = query.Offset(offset)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Order("created_at DESC").Find(&blogs).Error
	return blogs, err
}

// FindPublished 查找已发布的博客
func (dao *BlogDao) FindPublished(ctx context.Context, offset, limit int) ([]*Blog, error) {
	var blogs []*Blog
	query := mysql.GetDB().WithContext(ctx).Where("status = ?", BlogStatusPublished)
	if offset > 0 {
		query = query.Offset(offset)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Order("published_at DESC, created_at DESC").Find(&blogs).Error
	return blogs, err
}

// FindByScope 根据范围查找博客
func (dao *BlogDao) FindByScope(ctx context.Context, scope int, offset, limit int) ([]*Blog, error) {
	var blogs []*Blog
	query := mysql.GetDB().WithContext(ctx).Where("scope = ?", scope)
	if offset > 0 {
		query = query.Offset(offset)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Order("created_at DESC").Find(&blogs).Error
	return blogs, err
}

// FindByStatusAndScope 根据状态和范围查找博客
func (dao *BlogDao) FindByStatusAndScope(ctx context.Context, status int, scope int, offset, limit int) ([]*Blog, error) {
	var blogs []*Blog
	query := mysql.GetDB().WithContext(ctx).Where("status = ? AND scope = ?", status, scope)
	if offset > 0 {
		query = query.Offset(offset)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Order("published_at DESC, created_at DESC").Find(&blogs).Error
	return blogs, err
}

// CountByUserID 统计用户的博客数量（排除已删除）
func (dao *BlogDao) CountByUserID(ctx context.Context, userID uint64) (int64, error) {
	var count int64
	err := mysql.GetDB().WithContext(ctx).Model(&Blog{}).Where("user_id = ? AND status != ?", userID, BlogStatusDeleted).Count(&count).Error
	return count, err
}

// CountByStatus 统计指定状态的博客数量
func (dao *BlogDao) CountByStatus(ctx context.Context, status int) (int64, error) {
	var count int64
	err := mysql.GetDB().WithContext(ctx).Model(&Blog{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

// CountPublished 统计已发布的博客数量
func (dao *BlogDao) CountPublished(ctx context.Context) (int64, error) {
	return dao.CountByStatus(ctx, BlogStatusPublished)
}

// CountByScope 统计指定范围的博客数量
func (dao *BlogDao) CountByScope(ctx context.Context, scope int) (int64, error) {
	var count int64
	err := mysql.GetDB().WithContext(ctx).Model(&Blog{}).Where("scope = ?", scope).Count(&count).Error
	return count, err
}

// CountByStatusAndScope 统计指定状态和范围的博客数量
func (dao *BlogDao) CountByStatusAndScope(ctx context.Context, status int, scope int) (int64, error) {
	var count int64
	err := mysql.GetDB().WithContext(ctx).Model(&Blog{}).Where("status = ? AND scope = ?", status, scope).Count(&count).Error
	return count, err
}

// Publish 发布博客
func (dao *BlogDao) Publish(ctx context.Context, id uint64, publishedAt time.Time) error {
	updates := map[string]interface{}{
		"status":        BlogStatusPublished,
		"published_at":  publishedAt,
	}
	return mysql.GetDB().WithContext(ctx).Model(&Blog{}).Where("id = ?", id).Updates(updates).Error
}

// IncrementViewCount 增加浏览次数
func (dao *BlogDao) IncrementViewCount(ctx context.Context, id uint64) error {
	return mysql.GetDB().WithContext(ctx).Model(&Blog{}).Where("id = ?", id).
		Update("view_count", gorm.Expr("view_count + 1")).Error
}

// UpdateStatus 更新博客状态
func (dao *BlogDao) UpdateStatus(ctx context.Context, id uint64, status int) error {
	return mysql.GetDB().WithContext(ctx).Model(&Blog{}).Where("id = ?", id).Update("status", status).Error
}

// FindByCategory 根据分类查找博客
func (dao *BlogDao) FindByCategory(ctx context.Context, category string, offset, limit int) ([]*Blog, error) {
	var blogs []*Blog
	query := mysql.GetDB().WithContext(ctx).Where("category = ?", category)
	if offset > 0 {
		query = query.Offset(offset)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Order("published_at DESC, created_at DESC").Find(&blogs).Error
	return blogs, err
}

// Search 搜索博客（标题和内容）
func (dao *BlogDao) Search(ctx context.Context, keyword string, offset, limit int) ([]*Blog, error) {
	var blogs []*Blog
	if keyword == "" {
		return blogs, nil
	}
	
	query := mysql.GetDB().WithContext(ctx).
		Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	
	if offset > 0 {
		query = query.Offset(offset)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	err := query.Order("published_at DESC, created_at DESC").Find(&blogs).Error
	return blogs, err
}

// BatchUpdateStatus 批量更新状态
func (dao *BlogDao) BatchUpdateStatus(ctx context.Context, ids []uint64, status int) error {
	return mysql.GetDB().WithContext(ctx).Model(&Blog{}).Where("id IN ?", ids).Update("status", status).Error
}

// BatchDelete 批量删除
func (dao *BlogDao) BatchDelete(ctx context.Context, ids []uint64) error {
	return mysql.GetDB().WithContext(ctx).Where("id IN ?", ids).Delete(&Blog{}).Error
}

// FindAll 查找所有博客
func (dao *BlogDao) FindAll(ctx context.Context, offset, limit int) ([]*Blog, error) {
	var blogs []*Blog
	query := mysql.GetDB().WithContext(ctx)
	if offset > 0 {
		query = query.Offset(offset)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Order("created_at DESC").Find(&blogs).Error
	return blogs, err
}

// CountAll 统计所有博客数量
func (dao *BlogDao) CountAll(ctx context.Context) (int64, error) {
	var count int64
	err := mysql.GetDB().WithContext(ctx).Model(&Blog{}).Count(&count).Error
	return count, err
}

// GlobalDAO 全局DAO实例（单例模式）
var globalBlogDAO BlogDAO

// GetGlobalBlogDAO 获取全局博客DAO
func GetGlobalBlogDAO() BlogDAO {
	return globalBlogDAO
}

// SetGlobalBlogDAO 设置全局博客DAO
func SetGlobalBlogDAO(dao BlogDAO) {
	globalBlogDAO = dao
}
