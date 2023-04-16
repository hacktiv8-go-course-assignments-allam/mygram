package account

import "github.com/gin-gonic/gin"

type IAccountHandler interface {
	LoginAccount(ctx *gin.Context)
	CreateAccount(ctx *gin.Context)
	GetAccount(ctx *gin.Context)

	RegisterUserHdl(ctx *gin.Context)
	LoginUserHdl(ctx *gin.Context)
	GetUser(ctx *gin.Context)

	GetAllPhotos(ctx *gin.Context)
	GetPhotoById(ctx *gin.Context)
	CreatePhoto(ctx *gin.Context)
	UpdatePhoto(ctx *gin.Context)
	DeletePhoto(ctx *gin.Context)

	GetAllComments(ctx *gin.Context)
	GetCommentById(ctx *gin.Context)
	CreateComment(ctx *gin.Context)
	UpdateComment(ctx *gin.Context)
	DeleteComment(ctx *gin.Context)
	
	GetAllSocialMedias(ctx *gin.Context)
	GetSocialMediaById(ctx *gin.Context)
	CreateSocialMedia(ctx *gin.Context)
	UpdateSocialMedia(ctx *gin.Context)
	DeleteSocialMedia(ctx *gin.Context)
	// FindAllUsersHdl(ctx *gin.Context)
	// FindUserByIdHdl(ctx *gin.Context)
	// InsertUserHdl(ctx *gin.Context)
	// UpdateUserHdl(ctx *gin.Context)
	// DeleteUserByIdHdl(ctx *gin.Context)
}
