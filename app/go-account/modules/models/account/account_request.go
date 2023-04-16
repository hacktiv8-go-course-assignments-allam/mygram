package account

type CreateAccount struct {
	Username string      `json:"username" binding:"required"`
	Password string      `json:"password" binding:"required"`
	Role     AccountRole `json:"role" binding:"required"`
}

type LoginAccount struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUser struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterUser struct {
	Username string      `json:"username" gorm:"column:username;not null"`
	Email string      	 `json:"email" gorm:"column:email;not null"`
	Password string      `json:"password" gorm:"column:password;not null"`
	Age string      		 `json:"age" gorm:"column:age;not null"`
}
