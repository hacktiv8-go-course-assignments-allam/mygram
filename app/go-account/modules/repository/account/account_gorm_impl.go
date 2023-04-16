package account

import (
	"context"
	"errors"
	"fmt"

	accountmodel "github.com/mygram/go-account/modules/models/account"
	"github.com/mygram/go-common/pkg/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AccountRepoGormImpl struct {
	master *gorm.DB
}

func NewAccountRepoGormImpl(master *gorm.DB) IAccountRepo {
	return &AccountRepoGormImpl{
		master: master,
	}
}

// ACCOUNT SECTION
func (a *AccountRepoGormImpl) CreateAccount(ctx context.Context, acc accountmodel.Account) (created accountmodel.Account, err error) {
	logCtx := fmt.Sprintf("%T - CreateAccount", a)
	logger.Info(ctx, "%v invoked", "logCtx", logCtx)

	err = a.master.
		Table("accounts").
		Create(&acc).Error
	if err != nil {
		return
	}

	return acc, err
}

func (a *AccountRepoGormImpl) GetAccountByUserName(ctx context.Context, username string) (account accountmodel.Account, err error) {
	logCtx := fmt.Sprintf("%T - GetAccountByUserName", a)
	logger.Info(ctx, "%v invoked", "logCtx", logCtx)

	err = a.master.
		Table("accounts").
		Where("username = ?", username).
		Find(&account).Error
	if err != nil {
		return
	}

	return account, err
}

func (a *AccountRepoGormImpl) GetAccountByUserID(ctx context.Context, userId string) (account accountmodel.Account, err error) {
	logCtx := fmt.Sprintf("%T - GetAccountByUserID", a)
	logger.Info(ctx, "%v invoked", "logCtx", logCtx)

	err = a.master.
		Table("accounts").
		Where("id = ?", userId).
		Find(&account).Error

	if err != nil {
		return
	}
	return account, err
}

// USER SECTION
func (a *AccountRepoGormImpl) CreateUser(ctx context.Context, acc accountmodel.User) (created accountmodel.User, err error) {
	logCtx := fmt.Sprintf("%T - CreateAccount", a)
	logger.Info(ctx, "%v invoked", "logCtx", logCtx)

	err = a.master.
		Table("user").
		Create(&acc).Error
	if err != nil {
		return
	}

	return acc, err
}

func (a *AccountRepoGormImpl) GetUserByUserName(ctx context.Context, username string) (account accountmodel.User, err error) {
	logCtx := fmt.Sprintf("%T - GetAccountByUserName", a)
	logger.Info(ctx, "%v invoked", "logCtx", logCtx)

	err = a.master.
		Table("user").
		Where("username = ?", username).
		Find(&account).Error
	if err != nil {
		return
	}

	return account, err
}

func (a *AccountRepoGormImpl) GetUserById(ctx context.Context, userId string) (account accountmodel.User, err error) {
	logCtx := fmt.Sprintf("%T - GetUserById", a)
	logger.Info(ctx, "%v invoked", "logCtx", logCtx)

	err = a.master.
		Table("user").
		Where("id = ?", userId).
		Find(&account).Error
	if err != nil {
		return
	}

	return account, err
}

// func (u *UserGormRepoImpl) FindAllUsers(ctx context.Context) (users []model.User, err error) {
// 	tx := u.db.
// 		Model(&model.User{}).
// 		Find(&users).
// 		Order("created_at ASC")

// 	if err = tx.Error; err != nil {
// 		return
// 	}

// 	return
// }

func (a *AccountRepoGormImpl) GetAllPhotos(ctx context.Context) (photo []accountmodel.Photo, err error) {
	logCtx := fmt.Sprintf("%T - GetAllPhotos", a)
	logger.Info(ctx, "%v invoked", "logCtx", logCtx)

	err = a.master.
		Table("photo").
		Limit(20).
		Find(&photo).
		Order("created_at DESC").Error
	if err != nil {
		return
	}

	return photo, err
}

func (a *AccountRepoGormImpl) GetPhotoById(ctx context.Context, photoId uint64) (photo accountmodel.Photo, err error) {
	logCtx := fmt.Sprintf("%T - GetPhotoById", a)
	logger.Info(ctx, "%v invoked", "logCtx", logCtx)

	err = a.master.
		Table("photo").
		Where("id = ?", photoId).
		Find(&photo).Error

	if err != nil {
		return
	}
	return photo, err
}
func (a *AccountRepoGormImpl) CreatePhoto(ctx context.Context, pho accountmodel.Photo) (photo accountmodel.Photo, err error){
	logCtx := fmt.Sprintf("%T - CreatePhoto", a)
	logger.Info(ctx, "%v invoked", "logCtx", logCtx)

	err = a.master.
		Table("photo").
		Create(&pho).Error
	if err != nil {
		return
	}

	return pho, err
}
func (a *AccountRepoGormImpl) UpdatePhoto(ctx context.Context, pho accountmodel.Photo) (photo accountmodel.Photo, err error) {
	logCtx := fmt.Sprintf("%T - UpdatePhoto", a)
	logger.Info(ctx, "%v invoked", "logCtx", logCtx)
	
	tx := a.master.
		Model(&photo).
		Table("photo").
		Where("id = ?", pho.ID).
		Updates(&pho)

	if err = tx.Error; err != nil {
		return
	}

	if tx.RowsAffected <= 0 {
		err = errors.New("book is not found")
		return
	}

	return
}
func (a *AccountRepoGormImpl) DeletePhoto(ctx context.Context, photoId uint64) (photo accountmodel.Photo, err error) {
	tx := a.master.
		Model(&photo).
		Table("photo").
		// clause to return data after delete
		Clauses(clause.Returning{}).
		Where("id = ?", photoId).
		Delete(&photo)
	if err = tx.Error; err != nil {
		return
	}

	if tx.RowsAffected <= 0 {
		err = errors.New("book is not found")
		return
	}
	return
}

func (a *AccountRepoGormImpl) GetAllComments(ctx context.Context) (comment []accountmodel.Comment, err error){
	logCtx := fmt.Sprintf("%T - GetAllComments", a)
	logger.Info(ctx, "%v invoked", "logCtx", logCtx)

	err = a.master.
		Table("comment").
		Limit(20).
		Find(&comment).
		Order("created_at DESC").Error
	if err != nil {
		return
	}

	return comment, err
}
func (a *AccountRepoGormImpl) GetCommentById(ctx context.Context, commentId uint64) (comment accountmodel.Comment, err error) {
	logCtx := fmt.Sprintf("%T - GetCommentById", a)
	logger.Info(ctx, "%v invoked", "logCtx", logCtx)

	err = a.master.
		Table("comment").
		Where("id = ?", commentId).
		Find(&comment).Error

	if err != nil {
		return
	}
	return comment, err
}
func (a *AccountRepoGormImpl) CreateComment(ctx context.Context, com accountmodel.Comment) (comment accountmodel.Comment, err error) {
	logCtx := fmt.Sprintf("%T - CreateComment", a)
	logger.Info(ctx, "%v invoked", "logCtx", logCtx)

	err = a.master.
		Table("comment").
		Create(&com).Error
	if err != nil {
		return
	}

	return com, err
}
func (a *AccountRepoGormImpl) UpdateComment(ctx context.Context, com accountmodel.Comment) (comment accountmodel.Comment, err error) {
	tx := a.master.
		Model(&comment).
		Table("comment").
		Where("id = ?", com.ID).
		Updates(&com)

	if err = tx.Error; err != nil {
		return
	}

	if tx.RowsAffected <= 0 {
		err = errors.New("book is not found")
		return
	}

	return
}
func (a *AccountRepoGormImpl) DeleteComment(ctx context.Context, commentId uint64) (comment accountmodel.Comment, err error) {
	tx := a.master.
		Model(&comment).
		Table("comment").
		// clause to return data after delete
		Clauses(clause.Returning{}).
		Where("id = ?", commentId).
		Delete(&comment)
	if err = tx.Error; err != nil {
		return
	}

	if tx.RowsAffected <= 0 {
		err = errors.New("book is not found")
		return
	}
	return
}

func (a *AccountRepoGormImpl) GetAllSocialMedias(ctx context.Context) (socialMedia []accountmodel.SocialMedia, err error) {
	logCtx := fmt.Sprintf("%T - GetAllSocialMedias", a)
	logger.Info(ctx, "%v invoked", "logCtx", logCtx)

	err = a.master.
		Table("socialmedia").
		Limit(20).
		Find(&socialMedia).
		Order("created_at DESC").Error
	if err != nil {
		return
	}

	return socialMedia, err
}
func (a *AccountRepoGormImpl) GetSocialMediaById(ctx context.Context, socialMediaId uint64) (socialMedia accountmodel.SocialMedia, err error){
	logCtx := fmt.Sprintf("%T - GetSocialMediaById", a)
	logger.Info(ctx, "%v invoked", "logCtx", logCtx)

	err = a.master.
		Table("socialmedia").
		Where("id = ?", socialMediaId).
		Find(&socialMedia).Error

	if err != nil {
		return
	}
	return socialMedia, err
}
func (a *AccountRepoGormImpl) CreateSocialMedia(ctx context.Context, soc accountmodel.SocialMedia) (socialMedia accountmodel.SocialMedia, err error) {
	logCtx := fmt.Sprintf("%T - CreateSocialMedia", a)
	logger.Info(ctx, "%v invoked", "logCtx", logCtx)

	err = a.master.
		Table("socialmedia").
		Create(&soc).Error
	if err != nil {
		return
	}

	return soc, err
}
func (a *AccountRepoGormImpl) UpdateSocialMedia(ctx context.Context, soc accountmodel.SocialMedia) (socialMedia accountmodel.SocialMedia, err error) {
	tx := a.master.
		Model(&socialMedia).
		Table("socialmedia").
		Where("id = ?", soc.ID).
		Updates(&soc)

	if err = tx.Error; err != nil {
		return
	}

	if tx.RowsAffected <= 0 {
		err = errors.New("book is not found")
		return
	}

	return
}
func (a *AccountRepoGormImpl) DeleteSocialMedia(ctx context.Context, socialMediaId uint64) (socialMedia accountmodel.SocialMedia, err error) {
	tx := a.master.
		Model(&socialMedia).
		Table("socialmedia").
		// clause to return data after delete
		Clauses(clause.Returning{}).
		Where("id = ?", socialMediaId).
		Delete(&socialMedia)
	if err = tx.Error; err != nil {
		return
	}

	if tx.RowsAffected <= 0 {
		err = errors.New("book is not found")
		return
	}
	return
}