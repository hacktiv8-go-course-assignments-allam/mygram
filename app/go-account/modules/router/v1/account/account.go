package account

import (
	"github.com/gin-gonic/gin"
	accounthandler "github.com/mygram/go-account/modules/handler/account"
	"github.com/mygram/go-account/pkg/middleware"
)

func NewAccountRouter(v1 *gin.RouterGroup, accountHdl accounthandler.IAccountHandler) {
	gAccount := v1.Group("/account")

	// register all router
	gAccount.POST("",
		accountHdl.CreateAccount)
	gAccount.POST("/login",
		accountHdl.LoginAccount)
	gAccount.GET("",
		middleware.BearerOAuth(),
		accountHdl.GetAccount)

	
	gUser := v1.Group("/user")

	gUser.POST("/register", accountHdl.RegisterUserHdl)
	gUser.POST("/login", accountHdl.LoginUserHdl)
	gUser.GET("",
		middleware.BearerOAuth(),
		accountHdl.GetUser)

	gPhoto := v1.Group("/photo")

	gPhoto.GET("/all", accountHdl.GetAllPhotos)
	gPhoto.GET("", accountHdl.GetPhotoById)
	gPhoto.POST("", 
		middleware.BearerOAuth(), accountHdl.CreatePhoto)
	gPhoto.PUT("/:id", 
		middleware.BearerOAuth(), accountHdl.UpdatePhoto)
	gPhoto.DELETE("/:id", 
		middleware.BearerOAuth(), accountHdl.DeletePhoto)

		
	gComment := v1.Group("/comment")

	gComment.GET("/all", accountHdl.GetAllComments)
	gComment.GET("", accountHdl.GetCommentById)
	gComment.POST("", 
		middleware.BearerOAuth(), accountHdl.CreateComment)
	gComment.PUT("/:id", 
		middleware.BearerOAuth(), accountHdl.UpdateComment)
	gComment.DELETE("/:id", 
		middleware.BearerOAuth(), accountHdl.DeleteComment)
		
	gSocialMedia := v1.Group("/socmed")

	gSocialMedia.GET("/all", accountHdl.GetAllSocialMedias)
	gSocialMedia.GET("", accountHdl.GetSocialMediaById)
	gSocialMedia.POST("", 
		middleware.BearerOAuth(), accountHdl.CreateSocialMedia)
	gSocialMedia.PUT("/:id", 
		middleware.BearerOAuth(), accountHdl.UpdateSocialMedia)
	gSocialMedia.DELETE("/:id", 
		middleware.BearerOAuth(), accountHdl.DeleteSocialMedia)
	// register all router
	// gUser.GET("/all", accountHdl.FindAllUsersHdl)
	// gUser.GET("", accountHdl.FindUserByIdHdl)
	// gUser.POST("", accountHdl.InsertUserHdl)
	// gUser.PUT("/:id", accountHdl.UpdateUserHdl)
	// gUser.DELETE("/:id", accountHdl.DeleteUserByIdHdl)
}
