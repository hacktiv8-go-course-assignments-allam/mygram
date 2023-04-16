package account

import (
	"time"

	"github.com/google/uuid"
)

type AccountResponse struct {
	ID        uuid.UUID   `json:"id"`
	Username  string      `json:"username"`
	Role      AccountRole `json:"role"`
	CreatedAt time.Time   `json:"created_at"`
}

type AccountResponseWithPassword struct {
	AccountResponse
	Password string
}

type UserResponse struct {
	// - di json, menandakan golang akan mengabaikan properti
	// - di gorm, menandakan gorm akan mengabaikan properti
	ID        uint64      `json:"id" gorm:"column:id"`
	Username  string         `json:"username" gorm:"column:username"`
	Email  string         	 `json:"email" gorm:"column:email"`
	Password  string         `json:"password" gorm:"password"`
	Age  			uint64         `json:"age" gorm:"column:age"`

	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
	// gorm.Model
}

type UserRegisterResponse struct {
	ID        uint64      `json:"id" gorm:"column:id"`
	Username  string         `json:"username" gorm:"column:username"`
	Email  string         	 `json:"email" gorm:"column:email"`
	Password  string         `json:"password" gorm:"password"`
	Age  			uint64         `json:"age" gorm:"column:age"`

	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
	// gorm.Model
}