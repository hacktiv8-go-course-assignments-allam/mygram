package account

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountRole string

const (
	ROLE_ADMIN  AccountRole = "admin"
	ROLE_NORMAL AccountRole = "normal"
)

type Account struct {
	ID        uuid.UUID      `json:"id" gorm:"column:id"`
	Username  string         `json:"username" gorm:"column:username"`
	Password  string         `json:"password" gorm:"password"`
	Role      AccountRole    `json:"role" gorm:"column:role"`

	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`
}

// USER section
type User struct {
	// Id        		uint64         `json:"id" gorm:"column:id;type:integer;primaryKey;autoIncrement"`
	// ID        uuid.UUID      `json:"id" gorm:"column:id"`
	ID        uint64         `json:"id" gorm:"column:id;type:integer;primaryKey;autoIncrement"`
	Username  string         `json:"username" gorm:"column:username"`
	Email  string         	 `json:"email" gorm:"column:email"`
	Password  string         `json:"password" gorm:"password"`
	Age  			uint64         `json:"age" gorm:"column:age"`

	
	CreatedAt 		time.Time      `json:"created_at"`
	UpdatedAt 		time.Time      `json:"updated_at"`
	DeletedAt 		gorm.DeletedAt `json:"-" gorm:"index"`
	// CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	// UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
	// DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`
}

// PHOTO section
type Photo struct {
	// Id        		uint64         `json:"id" gorm:"column:id;type:integer;primaryKey;autoIncrement"`
	
	// ID        uuid.UUID      `json:"id" gorm:"column:id"`
	ID        uint64         `json:"id" gorm:"column:id;type:integer;primaryKey;autoIncrement"`
	UserID    uint64      `json:"user_id" gorm:"column:user_id"`
	Title  string         	 `json:"title" gorm:"column:title"`
	Caption  string          `json:"caption" gorm:"column:caption"`
	PhotoUrl  string         `json:"photo_url" gorm:"column:photo_url"`
	
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`
}

// COMMENT section
type Comment struct {
	// Id        		uint64         `json:"id" gorm:"column:id;type:integer;primaryKey;autoIncrement"`
	// ID        uuid.UUID      `json:"id" gorm:"column:id"`
	ID        uint64         `json:"id" gorm:"column:id;type:integer;primaryKey;autoIncrement"`
	UserID    uint64      `json:"user_id" gorm:"column:user_id"`
	PhotoID    uint64      `json:"photo_id" gorm:"column:photo_id"`
	Message  string         		 `json:"message" gorm:"column:message"`
	
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`
}

// Socmed section
type SocialMedia struct {
	// Id        		uint64         `json:"id" gorm:"column:id;type:integer;primaryKey;autoIncrement"`
	// ID        uuid.UUID      `json:"id" gorm:"column:id"`
	ID        uint64         `json:"id" gorm:"column:id;type:integer;primaryKey;autoIncrement"`
	UserID    uint64      `json:"user_id" gorm:"column:user_id"`
	Name  string         		 `json:"name" gorm:"column:name"`
	SocialMediaUrl  string         		 `json:"social_media_url" gorm:"column:social_media_url"`
	
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`
}

/*

// photo
{
	"user_id" : 1,
	"title" : "Some Cool Title",
	"caption" : "Some Cool Caption",
	"photo_url" : "Some Cool URL"
}
{
	"user_id" : 1,
	"title" : "EHE TE NANDAYO",
	"caption" : "Paimon",
	"photo_url" : "URL Moved"
}

// comment
{
	"user_id": 1,
	"photo_id": 1,
	"message": "I have a really cool motivational comment here"
}
{
	"user_id": 1,
	"photo_id": 1,
	"message": "I am gonna demotivate you haha"
}

// socmed
{
	"user_id": 1,
	"name": "Allam The Savior",
	"social_media_url": "https://twitter.com/masmaserius"
}
{
	"user_id": 1,
	"name": "Allam The Slashslingingslicer",
	"social_media_url": "https://twitter.com/masmaserius"
}
*/