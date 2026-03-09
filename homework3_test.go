package homework01

import (
	"context"
	"homework01/dbinstance"
	"testing"
	"time"

	"gorm.io/gorm"
)

// users 表：存储用户信息，包括 id 、 username 、 password 、 email 等字段。
// posts 表：存储博客文章信息，包括 id 、 title 、 content 、 user_id （关联 users 表的 id ）、 created_at 、 updated_at 等字段。
// comments 表：存储文章评论信息，包括 id 、 content 、 user_id （关联 users 表的 id ）、 post_id （关联 posts 表的 id ）、 created_at 等字段。
type User struct {
	ID        uint           `gorm:"primaryKey"`
	UserName  string         `gorm:"size:64;uniqueIndex;not null"`
	Password  string         `gorm:"size:128;not null"`
	Email     string         `gorm:"size:128;uniqueIndex;not null"`
	Posts     []Post         // Has Many
	Count     int            `gorm:"default:0"` // 文章数量统计字段
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
type AuditFields struct {
	CreatedBy string
	UpdatedBy string
	DeletedBy string
}
type Post struct {
	ID    uint   `gorm:"primaryKey"`
	Title string `gorm:"size:255;not null"`

	Content   string         `gorm:"type:text;not null"`
	UserID    uint           `gorm:"not null"`
	Comments  []Comment      // Has Many
	Audit     AuditFields    `gorm:"embedded"`
	Status    string         `gorm:"size:255"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Comment struct {
	ID        uint           `gorm:"primaryKey"`
	Content   string         `gorm:"type:text;not null"`
	UserID    uint           `gorm:"not null"`
	PostID    uint           `gorm:"not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// 数据初始化
func TestDB(t *testing.T) {
	db := dbinstance.NewTestDB(t, "associations.db")

	// AutoMigrate creates all tables and their relationships
	if err := db.AutoMigrate(&User{}, &Post{}, &Comment{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	// Clean up existing data
	db.Exec("PRAGMA foreign_keys = OFF")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM posts")
	db.Exec("DELETE FROM comments")
	db.Exec("PRAGMA foreign_keys = ON")

	// Create users
	alices := []User{
		User{UserName: "Alice",
			Email:    "alice@example.com",
			Password: "password",
			Posts: []Post{
				Post{Title: "Hello world", Content: "Welcome to my blog!"},
				Post{Title: "About me", Content: "I am a programmer."},
			},
		},
		User{UserName: "Bob", Email: "bob@example.com", Password: "password",
			Posts: []Post{
				Post{Title: "Hello world", Content: "Welcome to my blog!,Bob!"},
			},
		},
		User{UserName: "Charlie", Email: "charlie@example.com", Password: "password",
			Posts: []Post{
				Post{Title: "Hello world", Content: "Welcome to my blog!,Charlie!"},
			},
		},
	}
	if err := db.Session(&gorm.Session{FullSaveAssociations: true}).Create(&alices).Error; err != nil {
		t.Fatalf("create alices: %v", err)
	}

}

func TestFind(t *testing.T) {
	db := dbinstance.NewTestDB(t, "associations.db")

	var userlist []User
	// Associate roles with Alice
	if err := db.Model(&User{}).Find(&userlist).Error; err != nil {
		t.Fatalf("find users: %v", err)
	}

	//编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
	var user1 User
	// Associate roles with Alice
	if err := db.Model(&User{}).Where("user_name = ?", "Alice").Preload("Posts.Comments").Find(&user1).Error; err != nil {
		t.Fatalf("find users: %v", err)
	}
	type result struct {
		ID            uint
		comment_count uint
	}
	var res []result
	//	1. 编写Go代码，使用Gorm查询评论数量最多的文章信息。
	DB := db.Model(&Post{}).Select("posts.id, COUNT(comments.id) as comment_count").
		Joins("left join comments on comments.post_id = posts.id").
		Group("posts.id").
		Order("comment_count desc").
		Limit(1)
	if err := DB.Find(&res).Error; err != nil {
		t.Fatalf("find post: %v", err)
	}
	t.Log(res)

}

type ctxKey string

const ctxKeyOperator ctxKey = "operator"

func currentOperator(tx *gorm.DB) string {
	if tx != nil && tx.Statement != nil && tx.Statement.Context != nil {
		if v, ok := tx.Statement.Context.Value(ctxKeyOperator).(string); ok && v != "" {
			return v
		}
	}
	return "system"
}
func (a *Post) BeforeCreate(tx *gorm.DB) error {
	user := currentOperator(tx)
	a.Audit.CreatedBy = user
	a.Audit.UpdatedBy = user
	// 对于嵌入字段，使用扁平的列名（snake_case）
	tx.Statement.SetColumn("created_by", user)
	tx.Statement.SetColumn("updated_by", user)
	return nil
}

func (a *Post) AfterCreate(tx *gorm.DB) error {
	user := currentOperator(tx)

	// 验证文章创建时自动更新用户的文章数量统计字段
	tx.Model(&User{}).Where("User_Name = ?", user).UpdateColumn("count", gorm.Expr("count + ?", 1))

	return nil
}

func withOperator(name string) context.Context {
	return context.WithValue(context.Background(), ctxKeyOperator, name)
}

// 为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
func TestHook(t *testing.T) {
	db := dbinstance.NewTestDB(t, "associations.db")

	ctx := withOperator("Alice")
	art := Post{
		Title:   "GORM 钩子与软删除",
		Content: "演示如何在 GORM 中使用钩子与乐观锁。",
	}

	if err := db.WithContext(ctx).Create(&art).Error; err != nil {
		t.Fatalf("create article: %v", err)
	}

	// 验证钩子是否正确设置了审计字段
	if art.Audit.CreatedBy != "Alice" {
		t.Errorf("expected created_by to be 'alice', got %s", art.Audit.CreatedBy)
	}
	if art.Audit.UpdatedBy != "Alice" {
		t.Errorf("expected updated_by to be 'alice', got %s", art.Audit.UpdatedBy)
	}

}

func (a *Comment) AfterDelete(tx *gorm.DB) error {
	//	user := currentOperator(tx)
	var count int64
	tx.Model(&Comment{}).Where("post_id = ?", a.PostID).Count(&count)
	if count == 0 {
		tx.Model(&Post{}).Where("ID = ?", a.PostID).UpdateColumn("status", "无评论")
	}

	return nil

	// 为 Comment 模型添加一个钩子函数，在评论删除时自动更新文章的状态字段。
}

func TestCommentHook(t *testing.T) {

	db := dbinstance.NewTestDB(t, "associations.db")

	ctx := withOperator("Alice")
	var user1 User
	// 验证文章创建时自动更新用户的文章数量统计字段
	err := db.Model(&User{}).Where("User_Name = ?", "Alice").First(&user1).Error
	if err != nil {
		db.Logger.Error(context.Background(), "Failed to find user: %v", err)
	}
	art := Comment{
		PostID:  1,
		UserID:  user1.ID,
		Content: "演示如何在 GORM 中使用钩子与乐观锁。",
	}

	if err := db.WithContext(ctx).Create(&art).Error; err != nil {
		t.Fatalf("create article: %v", err)
	}

	if err := db.Delete(&art).Error; err != nil {
		t.Fatalf("delete comment: %v", err)
	}

}
