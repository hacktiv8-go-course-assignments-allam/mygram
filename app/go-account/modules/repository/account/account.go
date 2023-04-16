package account

import (
	"context"

	accountmodel "github.com/mygram/go-account/modules/models/account"
)

type IAccountRepo interface {
	CreateAccount(ctx context.Context, acc accountmodel.Account) (created accountmodel.Account, err error)
	GetAccountByUserName(ctx context.Context, username string) (account accountmodel.Account, err error)
	GetAccountByUserID(ctx context.Context, userId string) (account accountmodel.Account, err error)

	CreateUser(ctx context.Context, acc accountmodel.User) (created accountmodel.User, err error)
	GetUserByUserName(ctx context.Context, username string) (account accountmodel.User, err error)
	GetUserById(ctx context.Context, userId string) (account accountmodel.User, err error)

	GetAllPhotos(ctx context.Context) (account []accountmodel.Photo, err error)
	GetPhotoById(ctx context.Context, photoId uint64) (account accountmodel.Photo, err error)
	CreatePhoto(ctx context.Context, acc accountmodel.Photo) (account accountmodel.Photo, err error)
	UpdatePhoto(ctx context.Context, acc accountmodel.Photo) (account accountmodel.Photo, err error)
	DeletePhoto(ctx context.Context, photoId uint64) (account accountmodel.Photo, err error)

	GetAllComments(ctx context.Context) (account []accountmodel.Comment, err error)
	GetCommentById(ctx context.Context, commentId uint64) (comment accountmodel.Comment, err error)
	CreateComment(ctx context.Context, com accountmodel.Comment) (comment accountmodel.Comment, err error)
	UpdateComment(ctx context.Context, com accountmodel.Comment) (comment accountmodel.Comment, err error)
	DeleteComment(ctx context.Context, commentId uint64) (account accountmodel.Comment, err error)
	
	GetAllSocialMedias(ctx context.Context) (account []accountmodel.SocialMedia, err error)
	GetSocialMediaById(ctx context.Context, socialMediaId uint64) (socialMedia accountmodel.SocialMedia, err error)
	CreateSocialMedia(ctx context.Context, soc accountmodel.SocialMedia) (socialMedia accountmodel.SocialMedia, err error)
	UpdateSocialMedia(ctx context.Context, soc accountmodel.SocialMedia) (socialMedia accountmodel.SocialMedia, err error)
	DeleteSocialMedia(ctx context.Context, socialMediaId uint64) (account accountmodel.SocialMedia, err error)
	
	// FindBookById(ctx context.Context, bookId uint64) (book model.Book, err error)
	// FindAllBooks(ctx context.Context) (books []model.Book, err error)
	// InsertBook(ctx context.Context, bookIn model.Book) (book model.Book, err error)
	// UpdateBook(ctx context.Context, bookIn model.Book) (err error)
	// DeleteBookById(ctx context.Context, bookId uint64) (book model.Book, err error)
}
