package homework01

import (
	"homework01/dbinstance"
	"testing"
	"time"

	"gorm.io/gorm"
)

// users 表：存储用户信息，包括 id 、 username 、 password 、 email 等字段。
// posts 表：存储博客文章信息，包括 id 、 title 、 content 、 user_id （关联 users 表的 id ）、 created_at 、 updated_at 等字段。
// comments 表：存储文章评论信息，包括 id 、 content 、 user_id （关联 users 表的 id ）、 post_id （关联 posts 表的 id ）、 created_at 等字段。
type User struct {
	ID       uint   `gorm:"primaryKey"`
	UserName string `gorm:"size:64;uniqueIndex;not null"`
	Password string `gorm:"size:128;not null"`
	Email    string `gorm:"size:128;uniqueIndex;not null"`
	Posts    []Post // Has Many

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Post struct {
	ID    uint   `gorm:"primaryKey"`
	Title string `gorm:"size:255;not null"`

	Content   string    `gorm:"type:text;not null"`
	UserID    uint      `gorm:"not null"`
	Comments  []Comment // Has Many
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Comment struct {
	ID        uint      `gorm:"primaryKey"`
	Content   string    `gorm:"type:text;not null"`
	UserID    uint      `gorm:"not null"`
	PostID    uint      `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

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
				Post{Title: "Hello world", Content: "Welcome to my blog!",
					Comments: []Comment{
						Comment{Content: "Nice blog!"},
					},
				},
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
	var userlist []User
	// Associate roles with Alice
	if err := db.Model(&User{}).Find(&userlist).Error; err != nil {
		t.Fatalf("find users: %v", err)
	}

	//编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
	var user1 User
	// Associate roles with Alice
	if err := db.Model(&User{}).Where("UserName = ?", "Alice").Preload("Posts.Comments").Find(&user1).Error; err != nil {
		t.Fatalf("find users: %v", err)
	}
	//	1. 编写Go代码，使用Gorm查询评论数量最多的文章信息。
	DB := db.Model(&Post{}).Select("posts.*, COUNT(comments.id) as comment_count").
		Joins("left join comments on comments.post_id = posts.id").
		Group("posts.id").
		Order("comment_count desc").
		Limit(1)
	var post Post
	if err := DB.Find(&post).Error; err != nil {
		t.Fatalf("find post: %v", err)
	}
	t.Log(post)
	//	2. 编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息，并按照评论数量倒序排列。
	DB = db.Model(&User{}).Where("UserName = ?", "Alice").
		Preload("Posts.Comments").
		Select("users.*, posts.*, COUNT(comments.id) as comment_count").
		Joins("left join posts on posts.user_id = users.id").
		Joins("left join comments on comments.post_id = posts.id").
		Group("posts.id").
		Order("comment_count desc")

	//	为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
	//为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。
}
