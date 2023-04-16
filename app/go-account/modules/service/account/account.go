package account

import (
	"context"

	accountmodel "github.com/mygram/go-account/modules/models/account"
	token "github.com/mygram/go-account/modules/models/token"
)

type IAccountService interface {
	CreateAccount(ctx context.Context, acc accountmodel.CreateAccount) (created accountmodel.AccountResponse, err error)
	LoginAccountByUserName(ctx context.Context, loginAcc accountmodel.LoginAccount) (tokens token.Tokens, err error)
	GetAccount(ctx context.Context, userId string) (account accountmodel.AccountResponse, err error)

	LoginUser(ctx context.Context, loginAcc accountmodel.LoginUser) (tokens token.Tokens, err error)
	RegisterUser(ctx context.Context, acc accountmodel.RegisterUser) (created accountmodel.UserRegisterResponse, err error)
	GetUser(ctx context.Context, userId string) (user accountmodel.User, err error)

	GetAllPhotos(ctx context.Context) (photos []accountmodel.Photo, err error)
	GetPhotoById(ctx context.Context, photoId uint64) (photo accountmodel.Photo, err error)
	CreatePhoto(ctx context.Context, acc accountmodel.Photo) (photo accountmodel.Photo, err error)
	UpdatePhoto(ctx context.Context, acc accountmodel.Photo) (photo accountmodel.Photo, err error)
	DeletePhoto(ctx context.Context, photoId uint64) (photo accountmodel.Photo, err error)

	GetAllComments(ctx context.Context) (comments []accountmodel.Comment, err error)
	GetCommentById(ctx context.Context, commentId uint64) (comment accountmodel.Comment, err error)
	CreateComment(ctx context.Context, com accountmodel.Comment) (comment accountmodel.Comment, err error)
	UpdateComment(ctx context.Context, com accountmodel.Comment) (comment accountmodel.Comment, err error)
	DeleteComment(ctx context.Context, commentId uint64) (comment accountmodel.Comment, err error)
	
	GetAllSocialMedias(ctx context.Context) (socialMedia []accountmodel.SocialMedia, err error)
	GetSocialMediaById(ctx context.Context, socialMediaId uint64) (socialMedia accountmodel.SocialMedia, err error)
	CreateSocialMedia(ctx context.Context, soc accountmodel.SocialMedia) (socialMedia accountmodel.SocialMedia, err error)
	UpdateSocialMedia(ctx context.Context, soc accountmodel.SocialMedia) (socialMedia accountmodel.SocialMedia, err error)
	DeleteSocialMedia(ctx context.Context, socialMediaId uint64) (socialMedia accountmodel.SocialMedia, err error)
	// FindUserByIdSvc(ctx context.Context, userId uint64) (user accountmodel.User, err error)
	// FindAllUsersSvc(ctx context.Context) (users []accountmodel.User, err error)
	// InsertUserSvc(ctx context.Context, userIn accountmodel.User) (user accountmodel.User, err error)
	// UpdateUserSvc(ctx context.Context, userIn accountmodel.User) (err error)
	// DeleteUserByIdSvc(ctx context.Context, userId uint64) (user accountmodel.User, err error)
}
