package account

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mygram/go-common/pkg/logger"

	accountmodel "github.com/mygram/go-account/modules/models/account"
	"github.com/mygram/go-account/modules/models/accountactivity"
	token "github.com/mygram/go-account/modules/models/token"
	accountrepo "github.com/mygram/go-account/modules/repository/account"
	activityrepo "github.com/mygram/go-account/modules/repository/accountactivity"
	crypto "github.com/mygram/go-account/pkg/crypto"
)

type AccountServiceImpl struct {
	accountRepo  accountrepo.IAccountRepo
	activityRepo activityrepo.IAccountActivityRepo
}

func NewAccountServiceImpl(
	accountRepo accountrepo.IAccountRepo,
	activityRepo activityrepo.IAccountActivityRepo,
) IAccountService {
	return &AccountServiceImpl{
		accountRepo:  accountRepo,
		activityRepo: activityRepo,
	}
}

func (a *AccountServiceImpl) CreateAccount(ctx context.Context, acc accountmodel.CreateAccount) (created accountmodel.AccountResponse, err error) {
	logCtx := fmt.Sprintf("%T - CreatedAccount", a)
	logger.Info(ctx, "invoked", "logCtx", logCtx)

	// need to hash password
	hashedPassowrd, err := crypto.GenerateHash(acc.Password)
	if err != nil {
		logger.Error(ctx, "error when hashing password",
			"logCtx", logCtx,
			"error", err)
		return
	}
	// update passowrd with hashed password
	acc.Password = hashedPassowrd
	// store to db
	createdAcc, err := a.accountRepo.CreateAccount(ctx, accountmodel.Account{
		ID:       uuid.New(),
		Username: acc.Username,
		Password: acc.Password,
		Role:     acc.Role,
	})
	if err != nil {
		logger.Error(ctx, "error when storing account",
			"logCtx", logCtx,
			"error", err)
		return
	}

	return accountmodel.AccountResponse{
		ID:        createdAcc.ID,
		Username:  createdAcc.Username,
		Role:      createdAcc.Role,
		CreatedAt: createdAcc.CreatedAt,
	}, err
}

func (a *AccountServiceImpl) LoginAccountByUserName(ctx context.Context, loginAcc accountmodel.LoginAccount) (tokens token.Tokens, err error) {
	logCtx := fmt.Sprintf("%T - LoginAccountByUserName", a)
	logger.Info(ctx, "invoked", "logCtx", logCtx)

	// get account by username
	acc, err := a.getAccountWithPassword(ctx, loginAcc.Username)
	if err != nil {
		logger.Error(ctx, "error when fetching account by username",
			"logCtx", logCtx,
			"error", err)
		return
	}

	// compare password
	// password acc -> hashed password
	// password login acc -> plain password
	if err = crypto.CompareHash(acc.Password, loginAcc.Password); err != nil {
		logger.Error(ctx, "error when comparing password",
			"logCtx", logCtx,
			"error", err)
		return
	}

	// record activity
	createdActivity, err := a.activityRepo.CreateActivity(ctx, accountactivity.AccountActivity{
		ID:     uuid.New(),
		UserID: acc.ID,
		Type:   accountactivity.ACTIVITY_LOGIN,
	})
	if err != nil {
		logger.Error(ctx, "error when creating activity",
			"logCtx", logCtx,
			"error", err)
		return
	}

	idToken, accessToken, refreshToken, err := a.generateAllTokensConcurrent(ctx,
		acc.ID.String(),
		acc.Username,
		string(acc.Role),
		createdActivity.ID.String())
	if err != nil {
		return
	}

	return token.Tokens{
		IDToken:      (idToken),
		AccessToken:  (accessToken),
		RefreshToken: (refreshToken),
	}, err
}

func (a *AccountServiceImpl) GetAccount(ctx context.Context, userId string) (account accountmodel.AccountResponse, err error) {
	logCtx := fmt.Sprintf("%T - GetAccount", a)
	logger.Info(ctx, "invoked", "logCtx", logCtx)
	// get account from database
	acc, err := a.accountRepo.GetAccountByUserID(ctx, userId)
	if err != nil {
		return
	}
	return accountmodel.AccountResponse{
		ID:        acc.ID,
		Username:  acc.Username,
		Role:      acc.Role,
		CreatedAt: acc.CreatedAt,
	}, err
}

func (a *AccountServiceImpl) getAccountWithPassword(ctx context.Context, username string) (account accountmodel.AccountResponseWithPassword, err error) {
	logCtx := fmt.Sprintf("%T - getAccountWithPassword", a)
	logger.Info(ctx, "invoked", "logCtx", logCtx)

	// get account from database
	acc, err := a.accountRepo.GetAccountByUserName(ctx, username)
	if err != nil {
		return
	}
	return accountmodel.AccountResponseWithPassword{
		AccountResponse: accountmodel.AccountResponse{
			ID:        acc.ID,
			Username:  acc.Username,
			Role:      acc.Role,
			CreatedAt: acc.CreatedAt,
		},
		Password: acc.Password,
	}, err
}

func (a *AccountServiceImpl) generateAllTokensConcurrent(ctx context.Context, userid, username, role, jti string) (idToken, accessToken, refreshToken string, err error) {
	logCtx := fmt.Sprintf("%T - generateAllTokens", a)
	logger.Info(ctx, "invoked", "logCtx", logCtx)

	// https://github.com/kataras/jwt
	timeNow := time.Now()
	defaultClaim := token.DefaultClaim{
		Expired:   int(timeNow.Add(24 * time.Hour).Unix()),
		NotBefore: int(timeNow.Unix()),
		IssuedAt:  int(timeNow.Unix()),
		Issuer:    "http://go-account",
		Audience:  "http://dts-07",
		JTI:       jti,
		Type:      token.ID_TOKEN,
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func(defaultClaim_ token.DefaultClaim) {
		defer wg.Done()
		// generate id token
		idTokenClaim := struct {
			token.DefaultClaim
			token.IDClaim
		}{
			DefaultClaim: defaultClaim_,
			IDClaim: token.IDClaim{
				Username: username,
				Role:     role,
			},
		}
		idToken, err = crypto.SignJWT(idTokenClaim)
		if err != nil {
			logger.Error(ctx, "error when creating id token",
				"logCtx", logCtx,
				"error", err)
			return
		}
	}(defaultClaim)

	go func(defaultClaim_ token.DefaultClaim) {
		defer wg.Done()
		// generate access token
		defaultClaim_.Expired = int(timeNow.Add(20 * time.Minute).UnixMilli())
		defaultClaim_.Type = token.ACCESS_TOKEN
		accessTokenClaim := struct {
			token.DefaultClaim
			token.AccessClaim
		}{
			DefaultClaim: defaultClaim_,
			AccessClaim: token.AccessClaim{
				Role:   role,
				UserID: userid,
			},
		}
		accessToken, err = crypto.SignJWT(accessTokenClaim)
		if err != nil {
			logger.Error(ctx, "error when creating access token",
				"logCtx", logCtx,
				"error", err)
			return
		}
	}(defaultClaim)

	go func(defaultClaim_ token.DefaultClaim) {
		defer wg.Done()
		// generate refresh token
		defaultClaim_.Expired = int(timeNow.Add(time.Hour).UnixMilli())
		defaultClaim_.Type = token.REFRESH_TOKEN
		refreshTokenClaim := struct {
			token.DefaultClaim
		}{
			DefaultClaim: defaultClaim_,
		}
		refreshToken, err = crypto.SignJWT(refreshTokenClaim)
		if err != nil {
			logger.Error(ctx, "error when creating refresh token",
				"logCtx", logCtx,
				"error", err)
			return
		}
	}(defaultClaim)

	wg.Wait()
	return
}

func (a *AccountServiceImpl) generateAllTokens(ctx context.Context, userid, username, role, jti string) (idToken, accessToken, refreshToken string, err error) {
	logCtx := fmt.Sprintf("%T - generateAllTokens", a)
	logger.Info(ctx, "invoked", "logCtx", logCtx)

	// https://github.com/kataras/jwt
	timeNow := time.Now()
	defaultClaim_ := token.DefaultClaim{
		Expired:   int(timeNow.Add(24 * time.Hour).Unix()),
		NotBefore: int(timeNow.Unix()),
		IssuedAt:  int(timeNow.Unix()),
		Issuer:    "http://go-account",
		Audience:  "http://dts-07",
		JTI:       jti,
		Type:      token.ID_TOKEN,
	}

	var wg sync.WaitGroup
	wg.Add(3)

	// generate id token
	idTokenClaim := struct {
		token.DefaultClaim
		token.IDClaim
	}{
		DefaultClaim: defaultClaim_,
		IDClaim: token.IDClaim{
			Username: username,
			Role:     role,
		},
	}
	idToken, err = crypto.SignJWT(idTokenClaim)
	if err != nil {
		logger.Error(ctx, "error when creating id token",
			"logCtx", logCtx,
			"error", err)
		return
	}

	// generate access token
	defaultClaim_.Expired = int(timeNow.Add(20 * time.Minute).UnixMilli())
	defaultClaim_.Type = token.ACCESS_TOKEN
	accessTokenClaim := struct {
		token.DefaultClaim
		token.AccessClaim
	}{
		DefaultClaim: defaultClaim_,
		AccessClaim: token.AccessClaim{
			Role:   role,
			UserID: userid,
		},
	}
	accessToken, err = crypto.SignJWT(accessTokenClaim)
	if err != nil {
		logger.Error(ctx, "error when creating access token",
			"logCtx", logCtx,
			"error", err)
		return
	}

	// generate refresh token
	defaultClaim_.Expired = int(timeNow.Add(time.Hour).UnixMilli())
	defaultClaim_.Type = token.REFRESH_TOKEN
	refreshTokenClaim := struct {
		token.DefaultClaim
	}{
		DefaultClaim: defaultClaim_,
	}
	refreshToken, err = crypto.SignJWT(refreshTokenClaim)
	if err != nil {
		logger.Error(ctx, "error when creating refresh token",
			"logCtx", logCtx,
			"error", err)
		return
	}

	return
}


// USER SECTION
func (a *AccountServiceImpl) getUserWithPassword(ctx context.Context, username string) (user accountmodel.User, err error) {
	logCtx := fmt.Sprintf("%T - getUserWithPassword", a)
	logger.Info(ctx, "invoked", "logCtx", logCtx)

	// get account from database
	acc, err := a.accountRepo.GetUserByUserName(ctx, username)
	if err != nil {
		return
	}
	return acc, err
}

func (a *AccountServiceImpl) LoginUser(ctx context.Context, loginAcc accountmodel.LoginUser) (tokens token.Tokens, err error) {
	logCtx := fmt.Sprintf("%T - LoginAccountByUserName", a)
	logger.Info(ctx, "invoked", "logCtx", logCtx)

	// get account by username
	acc, err := a.getUserWithPassword(ctx, loginAcc.Username)
	if err != nil {
		logger.Error(ctx, "error when fetching account by username",
			"logCtx", logCtx,
			"error", err)
		return
	}

	// compare password
	// password acc -> hashed password
	// password login acc -> plain password
	if err = crypto.CompareHash(acc.Password, loginAcc.Password); err != nil {
		logger.Error(ctx, "error when comparing password",
			"logCtx", logCtx,
			"error", err)
		logger.Error(ctx, "acc pass: ", acc.Password, "login pass: ", loginAcc.Password)
		return
	}

	// record activity
	createdActivity, err := a.activityRepo.CreateUserActivity(ctx, accountactivity.UserActivity{
		ID:     uuid.New(),
		UserID: acc.ID,
		Type:   accountactivity.ACTIVITY_LOGIN,
	})
	if err != nil {
		logger.Error(ctx, "error when creating activity",
			"logCtx", logCtx,
			"error", err)
		return
	}

	idToken, accessToken, refreshToken, err := a.generateAllTokensConcurrent(ctx,
		strconv.FormatUint(acc.ID, 10),
		acc.Username,
		string("normal"),
		createdActivity.ID.String())
	if err != nil {
		return
	}

	return token.Tokens{
		IDToken:      (idToken),
		AccessToken:  (accessToken),
		RefreshToken: (refreshToken),
	}, err
}

func (a *AccountServiceImpl) RegisterUser(ctx context.Context, acc accountmodel.RegisterUser) (created accountmodel.UserRegisterResponse, err error) {
	logCtx := fmt.Sprintf("%T - CreatedAccount", a)
	logger.Info(ctx, "invoked", "logCtx", logCtx)

	// need to hash password
	hashedPassowrd, err := crypto.GenerateHash(acc.Password)
	if err != nil {
		logger.Error(ctx, "error when hashing password",
			"logCtx", logCtx,
			"error", err)
		return
	}
	// update passowrd with hashed password
	acc.Password = hashedPassowrd
	// store to db
	
	// transform id string to uint64
	accAge, err := strconv.ParseUint(acc.Age, 10, 64)
	if err != nil {
		logger.Error(ctx, "error when transforming age string to uint64",
			"logCtx", logCtx,
			"error", err)
		return
	}

	createdAcc, err := a.accountRepo.CreateUser(ctx, accountmodel.User{
		// ID:       uuid.New(),
		Username: acc.Username,
		Email: acc.Email,
		Password: acc.Password,
		Age:     accAge,
	})
	if err != nil {
		logger.Error(ctx, "error when storing account",
			"logCtx", logCtx,
			"error", err)
		return
	}

	return accountmodel.UserRegisterResponse{
		ID:        createdAcc.ID,
		Username:  createdAcc.Username,
		Email:  createdAcc.Email,
		Age  :  createdAcc.Age  ,
		CreatedAt: createdAcc.CreatedAt,
	}, err
}

func (a *AccountServiceImpl) GetUser(ctx context.Context, userId string) (user accountmodel.User, err error) {
	logCtx := fmt.Sprintf("%T - GetUser", a)
	logger.Info(ctx, "invoked", "logCtx", logCtx)
	// get account from database
	acc, err := a.accountRepo.GetUserById(ctx, userId)
	if err != nil {
		return
	}
	return acc, err
}


func (a *AccountServiceImpl) GetAllPhotos(ctx context.Context) (photos []accountmodel.Photo, err error){
	logCtx := fmt.Sprintf("%T - GetAllPhotos", a)
	logger.Info(ctx, "invoked", "logCtx", logCtx)
	if photos, err = a.accountRepo.GetAllPhotos(ctx); err != nil {
		logger.Error(ctx, "error GetAllPhotos",
			"logCtx", logCtx,
			"error", err)
	}
	return
}
func (a *AccountServiceImpl) GetPhotoById(ctx context.Context, photoId uint64) (photo accountmodel.Photo, err error) {
	logCtx := fmt.Sprintf("%T - GetPhotoById", a)
	logger.Info(ctx, "invoked", "logCtx", logCtx)
	if photo, err = a.accountRepo.GetPhotoById(ctx, photoId); err != nil {
		logger.Error(ctx, "error GetPhotoById",
			"logCtx", logCtx,
			"error", err)
	}
	return
}
func (a *AccountServiceImpl) CreatePhoto(ctx context.Context, acc accountmodel.Photo) (photo accountmodel.Photo, err error) {
	logCtx := fmt.Sprintf("%T - CreatePhoto", a)
	logger.Info(ctx, "invoked", "logCtx", logCtx)
	if photo, err = a.accountRepo.CreatePhoto(ctx, acc); err != nil {
		logger.Error(ctx, "error CreatePhoto",
			"logCtx", logCtx,
			"error", err)
	}
	return
}
func (a *AccountServiceImpl) UpdatePhoto(ctx context.Context, acc accountmodel.Photo) (photo accountmodel.Photo, err error) {
	logCtx := fmt.Sprintf("%T - UpdatePhoto", a)
	logger.Info(ctx, "invoked", "logCtx", logCtx)
	if photo, err = a.accountRepo.UpdatePhoto(ctx, acc); err != nil {
		logger.Error(ctx, "error UpdatePhoto",
			"logCtx", logCtx,
			"error", err)
	}
	return
}
func (a *AccountServiceImpl) DeletePhoto(ctx context.Context, photoId uint64) (photo accountmodel.Photo, err error){
	logCtx := fmt.Sprintf("%T - DeletePhoto", a)
	logger.Info(ctx, "invoked", "logCtx", logCtx)
	if photo, err = a.accountRepo.DeletePhoto(ctx, photoId); err != nil {
		logger.Error(ctx, "error DeletePhoto",
			"logCtx", logCtx,
			"error", err)
	}
	return
}

func (a *AccountServiceImpl) GetAllComments(ctx context.Context) (comments []accountmodel.Comment, err error) {
	logCtx := fmt.Sprintf("%T - GetAllComments", a)
	logger.Info(ctx, "invoked", "logCtx", logCtx)
	if comments, err = a.accountRepo.GetAllComments(ctx); err != nil {
		logger.Error(ctx, "error GetAllComments",
			"logCtx", logCtx,
			"error", err)
	}
	return
}
func (a *AccountServiceImpl) GetCommentById(ctx context.Context, commentId uint64) (comment accountmodel.Comment, err error){
	logCtx := fmt.Sprintf("%T - GetCommentById", a)
	logger.Info(ctx, "invoked", "logCtx", logCtx)
	if comment, err = a.accountRepo.GetCommentById(ctx, commentId); err != nil {
		logger.Error(ctx, "error GetCommentById",
			"logCtx", logCtx,
			"error", err)
	}
	return
}
func (a *AccountServiceImpl) CreateComment(ctx context.Context, com accountmodel.Comment) (comment accountmodel.Comment, err error) {
	logCtx := fmt.Sprintf("%T - CreateComment", a)
	logger.Info(ctx, "invoked", "logCtx", logCtx)
	if comment, err = a.accountRepo.CreateComment(ctx, com); err != nil {
		logger.Error(ctx, "error CreateComment",
			"logCtx", logCtx,
			"error", err)
	}
	return
}
func (a *AccountServiceImpl) UpdateComment(ctx context.Context, com accountmodel.Comment) (comment accountmodel.Comment, err error) {
	logCtx := fmt.Sprintf("%T - UpdateComment", a)
	logger.Info(ctx, "invoked", "logCtx", logCtx)
	if comment, err = a.accountRepo.UpdateComment(ctx, com); err != nil {
		logger.Error(ctx, "error UpdateComment",
			"logCtx", logCtx,
			"error", err)
	}
	return
}
func (a *AccountServiceImpl) DeleteComment(ctx context.Context, commentId uint64) (account accountmodel.Comment, err error){
	logCtx := fmt.Sprintf("%T - DeleteComment", a)
	logger.Info(ctx, "invoked", "logCtx", logCtx)
	if account, err = a.accountRepo.DeleteComment(ctx, commentId); err != nil {
		logger.Error(ctx, "error DeleteComment",
			"logCtx", logCtx,
			"error", err)
	}
	return
}

func (a *AccountServiceImpl) GetAllSocialMedias(ctx context.Context) (socialMedias []accountmodel.SocialMedia, err error) {
	logCtx := fmt.Sprintf("%T - GetAllSocialMedias", a)
	logger.Info(ctx, "invoked", "logCtx", logCtx)
	if socialMedias, err = a.accountRepo.GetAllSocialMedias(ctx); err != nil {
		logger.Error(ctx, "error GetAllSocialMedias",
			"logCtx", logCtx,
			"error", err)
	}
	return
}
func (a *AccountServiceImpl) GetSocialMediaById(ctx context.Context, socialMediaId uint64) (socialMedia accountmodel.SocialMedia, err error){
	logCtx := fmt.Sprintf("%T - GetSocialMediaById", a)
	logger.Info(ctx, "invoked", "logCtx", logCtx)
	if socialMedia, err = a.accountRepo.GetSocialMediaById(ctx, socialMediaId); err != nil {
		logger.Error(ctx, "error GetSocialMediaById",
			"logCtx", logCtx,
			"error", err)
	}
	return
}
func (a *AccountServiceImpl) CreateSocialMedia(ctx context.Context, soc accountmodel.SocialMedia) (socialMedia accountmodel.SocialMedia, err error){
	logCtx := fmt.Sprintf("%T - CreateSocialMedia", a)
	logger.Info(ctx, "invoked", "logCtx", logCtx)
	if socialMedia, err = a.accountRepo.CreateSocialMedia(ctx, soc); err != nil {
		logger.Error(ctx, "error CreateSocialMedia",
			"logCtx", logCtx,
			"error", err)
	}
	return
}
func (a *AccountServiceImpl) UpdateSocialMedia(ctx context.Context, soc accountmodel.SocialMedia) (socialMedia accountmodel.SocialMedia, err error){
	logCtx := fmt.Sprintf("%T - UpdateSocialMedia", a)
	logger.Info(ctx, "invoked", "logCtx", logCtx)
	if socialMedia, err = a.accountRepo.UpdateSocialMedia(ctx, soc); err != nil {
		logger.Error(ctx, "error UpdateSocialMedia",
			"logCtx", logCtx,
			"error", err)
	}
	return
}
func (a *AccountServiceImpl) DeleteSocialMedia(ctx context.Context, socialMediaId uint64) (socialMedia accountmodel.SocialMedia, err error){
	logCtx := fmt.Sprintf("%T - DeleteSocialMedia", a)
	logger.Info(ctx, "invoked", "logCtx", logCtx)
	if socialMedia, err = a.accountRepo.DeleteSocialMedia(ctx, socialMediaId); err != nil {
		logger.Error(ctx, "error DeleteSocialMedia",
			"logCtx", logCtx,
			"error", err)
	}
	return
}

/*
func (u *AccountServiceImpl) FindUserByIdSvc(ctx context.Context, userId uint64) (user accountmodel.User, err error) {
	log.Printf("[INFO] %T FindUserById invoked\n", u)
	if user, err = u.accountRepo.FindUserById(ctx, userId); err != nil {
		log.Printf("[ERROR] error FindUserById :%v\n", err)
	}
	return
}

func (u *AccountServiceImpl) FindAllUsersSvc(ctx context.Context) (users []accountmodel.User, err error) {
	log.Printf("[INFO] %T FindAllUsers invoked\n", u)
	if users, err = u.accountRepo.FindAllUsers(ctx); err != nil {
		log.Printf("[ERROR] error FindAllUsers :%v\n", err)
	}
	return
}

func (u *AccountServiceImpl) InsertUserSvc(ctx context.Context, userIn accountmodel.User) (user accountmodel.User, err error) {
	log.Printf("[INFO] %T InsertUser invoked\n", u)
	if user, err = u.accountRepo.InsertUser(ctx, userIn); err != nil {
		log.Printf("[ERROR] error InsertUser :%v\n", err)
	}
	return
}

func (u *AccountServiceImpl) UpdateUserSvc(ctx context.Context, userIn accountmodel.User) (err error) {
	log.Printf("[INFO] %T UpdateUser invoked\n", u)
	if err = u.accountRepo.UpdateUser(ctx, userIn); err != nil {
		log.Printf("[ERROR] error InsertUser :%v\n", err)
	}
	return
}

func (u *AccountServiceImpl) DeleteUserByIdSvc(ctx context.Context, userId uint64) (deletedUser accountmodel.User, err error) {
	log.Printf("[INFO] %T DeleteUserById invoked\n", u)
	if deletedUser, err = u.accountRepo.DeleteUserById(ctx, userId); err != nil {
		log.Printf("[ERROR] error DeleteUserById :%v\n", err)
	}
	return
}
*/

// mockgen -source=modules/repository/account/account.go \ interface kita
// -destination=modules/repository/account/mock/account_mock.go \ mock kita mau diletakin mana
// -package=mock
