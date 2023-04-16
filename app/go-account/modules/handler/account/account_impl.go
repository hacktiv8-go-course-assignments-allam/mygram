package account

import (
	"errors"
	"net/http"
	"net/mail"
	"strconv"

	"github.com/gin-gonic/gin"
	accountmodel "github.com/mygram/go-account/modules/models/account"
	"github.com/mygram/go-account/modules/models/token"
	accountservice "github.com/mygram/go-account/modules/service/account"
	"github.com/mygram/go-account/pkg/middleware"
	"github.com/mygram/go-common/pkg/json"
	"github.com/mygram/go-common/pkg/logger"
	"github.com/mygram/go-common/pkg/response"
)

type AccountHandlerImpl struct {
	accService accountservice.IAccountService
}

func NewAccountHandlerImpl(accService accountservice.IAccountService) IAccountHandler {
	return &AccountHandlerImpl{
		accService: accService,
	}
}

// util SECTION
func (a *AccountHandlerImpl) getIdFromParam(ctx *gin.Context) (idUint uint64, err error) {
	id := ctx.Param("id")
	if id == "" {
		err = errors.New("failed id")
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  "failed to update user",
			Error: response.InvalidParam,
		})
		return
	}
	// transform id string to uint64
	idUint, err = strconv.ParseUint(id, 10, 64)
	if err != nil {
		err = errors.New("failed parse id")
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  "failed to update user",
			Error: response.InvalidParam,
		})
		return
	}
	return
}


// ACCOUNT SECTION
func (a *AccountHandlerImpl) LoginAccount(ctx *gin.Context) {
	// binding payload
	var loginAccount accountmodel.LoginAccount
	if err := ctx.BindJSON(&loginAccount); err != nil {
		logger.Error(ctx, "error binding payload",
			"error", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			response.ErrorResponse{
				Message: response.InvalidBody,
				Error:   "error binding payload",
			},
		)
		return
	}

	tokens, err := a.accService.LoginAccountByUserName(ctx, loginAccount)
	if err != nil {
		logger.Error(ctx, "error create account",
			"error", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			response.ErrorResponse{
				Message: response.InternalServer,
				Error:   response.SomethingWentWrong,
			},
		)
		return
	}
	ctx.JSON(http.StatusAccepted, response.SuccessResponse{
		Message: "success created",
		Data:    tokens,
	})
}

func (a *AccountHandlerImpl) CreateAccount(ctx *gin.Context) {
	// binding payload
	var createAccount accountmodel.CreateAccount
	if err := ctx.BindJSON(&createAccount); err != nil {
		logger.Error(ctx, "error binding payload",
			"error", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			response.ErrorResponse{
				Message: response.InvalidBody,
				Error:   "error binding payload",
			},
		)
		return
	}

	created, err := a.accService.CreateAccount(ctx, createAccount)
	if err != nil {
		logger.Error(ctx, "error create account",
			"error", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			response.ErrorResponse{
				Message: response.InternalServer,
				Error:   response.SomethingWentWrong,
			},
		)
		return
	}
	ctx.JSON(http.StatusAccepted, response.SuccessResponse{
		Message: "success created",
		Data:    created,
	})
}

func (a *AccountHandlerImpl) GetAccount(ctx *gin.Context) {
	// get user_id from context first
	accessClaimI, ok := ctx.Get(middleware.AccessClaim.String())
	if !ok {
		err := errors.New("error get claim from context")
		logger.Error(ctx, "error get payload",
			"error", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse{
			Message: response.InvalidPayload,
			Error:   "invalid user id",
		})
		return
	}

	var accessClaim token.AccessClaim
	if err := json.ObjectMapper(accessClaimI, &accessClaim); err != nil {
		logger.Error(ctx, "error mapping object payload",
			"error", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse{
			Message: response.InvalidPayload,
			Error:   "invalid payload",
		})
		return
	}

	account, err := a.accService.GetAccount(ctx, accessClaim.UserID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse{
			Message: response.InternalServer,
			Error:   "something went wrong",
		})
		return
	}
	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Message: "success",
		Data:    account,
	})
}
// ACCOUNT SECTION

// AUTH 
func (a *AccountHandlerImpl) AuthIncomingRequest(ctx *gin.Context) (user accountmodel.User, err error, message string) {
		// get user_id from context first
	accessClaimI, ok := ctx.Get(middleware.AccessClaim.String())
	message = ""
	if !ok {
		message = "error get claim from context"
		err = errors.New(message)
		logger.Error(ctx, message,
			"error", err)
		return
	}

	var accessClaim token.AccessClaim
	if err = json.ObjectMapper(accessClaimI, &accessClaim); err != nil {
		message = "error mapping object payload"
		logger.Error(ctx, message,
			"error", err)
		return
	}

	account, err := a.accService.GetUser(ctx, accessClaim.UserID)
	if err != nil {
		message = "error while getting user"
		logger.Error(ctx, message,
			"error", err)
		return
	}
	return account, err, message
}

func CreateEntityAuth(userId uint64, payloadUserId uint64) (allowed bool, message string) {
	if userId != payloadUserId {
		userIdStr := strconv.FormatUint(uint64(userId), 10)
		payloadUserIdStr := strconv.FormatUint(uint64(payloadUserId), 10)
		return false, "cannot set new userId, stay with your Id, unauthorized " + userIdStr + ":" + payloadUserIdStr
	}
	return true, "aman"
}

func UpdateEntityAuth(userId uint64, payloadUserId uint64, oldEntityUserId uint64) (allowed bool, message string) {
	if userId != payloadUserId {
		userIdStr := strconv.FormatUint(uint64(userId), 10)
		payloadUserIdStr := strconv.FormatUint(uint64(payloadUserId), 10)
		return false, "cannot set new userId, stay with your Id, unauthorized " + userIdStr + ":" + payloadUserIdStr
	}
	if userId != oldEntityUserId {
		userIdStr := strconv.FormatUint(uint64(userId), 10)
		oldEntityUserIdStr := strconv.FormatUint(uint64(oldEntityUserId), 10)
		return false, "cannot set update this entity, not yours, unauthorized" + userIdStr + ":" + oldEntityUserIdStr
	}
	return true, "aman"
}

func DeleteEntityAuth(userId uint64, payloadUserId uint64) (allowed bool, message string) {
	if userId != payloadUserId {
		userIdStr := strconv.FormatUint(uint64(userId), 10)
		payloadUserIdStr := strconv.FormatUint(uint64(payloadUserId), 10)
		return false, "cannot set new userId, stay with your Id, unauthorized " + userIdStr + ":" + payloadUserIdStr
	}
	return true, "aman"
}

func EmailValidation(email string) (con bool, emailAddres string, message string) {
	addr, err := mail.ParseAddress(email)
	if err != nil {
			return false, addr.Address, "invalid email"
	}
	return true, addr.Address, "valid"
}

func PhotoPayloadValidation(pho accountmodel.Photo) (con bool, message string) {
	con = true
	if pho.Title == "" {
		message += "title cannot be empty,"
		con = false
	}
	if pho.PhotoUrl == "" {
		message += "photo url cannot be empty"
		con = false
	}
	return con, message
}

func CommentPayloadValidation(comment accountmodel.Comment) (con bool, message string) {
	con = true
	if comment.Message == "" {
		message += "message cannot be empty,"
		con = false
	}
	return con, message
}

func SocialMediaPayloadValidation(socialMedia accountmodel.SocialMedia) (con bool, message string) {
	con = true
	if socialMedia.Name == "" {
		message += "name cannot be empty,"
		con = false
	}
	if socialMedia.SocialMediaUrl == "" {
		message += "social media url cannot be empty,"
		con = false
	}
	return con, message
}


// USER SECTION
func (a *AccountHandlerImpl) LoginUserHdl(ctx *gin.Context) {
	// binding payload
	var loginAccount accountmodel.LoginUser
	if err := ctx.BindJSON(&loginAccount); err != nil {
		logger.Error(ctx, "error binding payload",
			"error", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			response.ErrorResponse{
				Message: response.InvalidBody,
				Error:   "error binding payload",
			},
		)
		return
	}

	tokens, err := a.accService.LoginUser(ctx, loginAccount)
	if err != nil {
		logger.Error(ctx, "error create account",
			"error", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			response.ErrorResponse{
				Message: response.InternalServer,
				Error:   response.SomethingWentWrong,
			},
		)
		return
	}
	ctx.JSON(http.StatusAccepted, response.SuccessResponse{
		Message: "success created",
		Data:    tokens,
	})
}

func (a *AccountHandlerImpl) RegisterUserHdl(ctx *gin.Context) {
	// binding payload
	var createAccount accountmodel.RegisterUser
	if err := ctx.BindJSON(&createAccount); err != nil {
		logger.Error(ctx, "error binding payload",
			"error", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			response.ErrorResponse{
				Message: response.InvalidBody,
				Error:   "error binding payload",
			},
		)
		return
	}

	// not null validation
	isEmailEmpty := createAccount.Email == ""
	isUsernameEmpty := createAccount.Username == ""
	isPasswordEmpty := createAccount.Password == ""
	isAgeEmpty := createAccount.Age == ""

	if isEmailEmpty || isUsernameEmpty || isPasswordEmpty || isAgeEmpty {
		message := "empty fields "
		if isEmailEmpty {
			message += "email "
		}
		if isUsernameEmpty {
			message += "username "
		}
		if isPasswordEmpty {
			message += "password "
		}
		if isAgeEmpty {
			message += "age "
		}
		logger.Error(ctx, "error create account",
			"error", message)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			response.ErrorResponse{
				Message: message,
				Error:   response.InvalidPayload,
			},
		)
		return
	}

	// password >= 6
	isPasswordLenEnough := len(createAccount.Password) >= 6

	if !isPasswordLenEnough {
		logger.Error(ctx, "error create account",
			"error", "password length insufficient")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			response.ErrorResponse{
				Message: "password length insufficient",
				Error:   response.InvalidPayload,
			},
		)
		return
	}

	// age > 8
	age, err := strconv.ParseInt(createAccount.Age, 10, 64)
	if err != nil {
		message := "age conversion failed"
		logger.Error(ctx, "error create account",
			"error", message)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			response.ErrorResponse{
				Message: message,
				Error:   response.InvalidPayload,
			},
		)
		return
	}
	isAgeEnough := age > 8

	if !isAgeEnough {
		logger.Error(ctx, "error create account",
			"error", "age insufficient")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			response.ErrorResponse{
				Message: "age insufficient",
				Error:   response.InvalidPayload,
			},
		)
		return
	}

	con, emailAddress, message := EmailValidation(createAccount.Email)
	if !con {
		logger.Error(ctx, "error create account",
			"error", message)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			response.ErrorResponse{
				Message: message + " " + emailAddress,
				Error:   response.InvalidPayload,
			},
		)
		return
	}

	created, err := a.accService.RegisterUser(ctx, createAccount)
	if err != nil {
		logger.Error(ctx, "error create account",
			"error", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			response.ErrorResponse{
				Message: response.InternalServer,
				Error:   response.SomethingWentWrong,
			},
		)
		return
	}
	ctx.JSON(http.StatusAccepted, response.SuccessResponse{
		Message: "success created",
		Data:    created,
	})
}

func (a *AccountHandlerImpl) GetUser(ctx *gin.Context) {
	// get user_id from context first
	accessClaimI, ok := ctx.Get(middleware.AccessClaim.String())
	if !ok {
		err := errors.New("error get claim from context")
		logger.Error(ctx, "error get payload",
			"error", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse{
			Message: response.InvalidPayload,
			Error:   "invalid user id",
		})
		return
	}

	var accessClaim token.AccessClaim
	if err := json.ObjectMapper(accessClaimI, &accessClaim); err != nil {
		logger.Error(ctx, "error mapping object payload",
			"error", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse{
			Message: response.InvalidPayload,
			Error:   "invalid payload",
		})
		return
	}

	account, err := a.accService.GetUser(ctx, accessClaim.UserID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse{
			Message: response.InternalServer,
			Error:   "something went wrong",
		})
		return
	}
	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Message: "success",
		Data:    account,
	})
}

// PHOTO SECTION
func (a *AccountHandlerImpl) GetAllPhotos(ctx *gin.Context) {
	photos, err := a.accService.GetAllPhotos(ctx)
	if err != nil {
		// bad code, should be wrapped in other package
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Message: response.SomethingWentWrong,
			Error: "failed to get photos",
		})
		return
	}
	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Message: "success get photos",
		Data:    photos,
	})
}
func (a *AccountHandlerImpl) GetPhotoById(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  "failed to find photo",
			Error: response.InvalidQuery,
		})
		return
	}
	// transform id string to uint64
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  "failed to find photo",
			Error: response.InvalidParam,
		})
		return
	}

	// call service
	photo, err := a.accService.GetPhotoById(ctx, idUint)
	if err != nil {
		if err.Error() == "photo is not found" {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{
				Message:  "failed to find photo",
				Error: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Message:  "failed to find photo",
			Error: response.SomethingWentWrong,
		})
		return
	}
	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Message: "success find photo",
		Data:    photo,
	})
}
func (a *AccountHandlerImpl) CreatePhoto(ctx *gin.Context) {
	// mendapatkan body
	var photoIn accountmodel.Photo

	logger.Info(ctx, "otw bind photo")
	if err := ctx.Bind(&photoIn); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  "failed to insert photo",
			Error: response.InvalidBody,
		})
		return
	}

	user, err, message := a.AuthIncomingRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  message,
			Error: response.InternalServer,
		})
		return
	}
	con, message := CreateEntityAuth(user.ID, photoIn.UserID)
	if con != true {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  message,
			Error: response.Unauthorized,
		})
		return
	}

	logger.Info(ctx, "otw validate photo")
	con, message = PhotoPayloadValidation(photoIn)
	if !con {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  message,
			Error: response.InvalidParam,
		})
		return
	}

	logger.Info(ctx, "otw insert photo")
	insertedPhoto, err := a.accService.CreatePhoto(ctx, photoIn)
	if err != nil {
		// bad code, should be wrapped in other package
		if err.Error() == "error duplication email" {
			ctx.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
				Message:  "failed to insert photo",
				Error: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Message:  "failed to insert photo",
			Error: response.SomethingWentWrong,
		})
		return
	}

	ctx.JSON(http.StatusAccepted, response.SuccessResponse{
		Message: "success create photo",
		Data:    insertedPhoto,
	})
}
func (a *AccountHandlerImpl) UpdatePhoto(ctx *gin.Context) {
	idUint, err := a.getIdFromParam(ctx)
	if err != nil {
		return
	}
	// binding payload
	var photoIn accountmodel.Photo
	if err := ctx.Bind(&photoIn); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  "failed to update user",
			Error: response.InvalidBody,
		})
		return
	}
	photoIn.ID = idUint

	user, err, message := a.AuthIncomingRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  message,
			Error: response.InternalServer,
		})
		return
	}

	toBeUpdatedPhoto, err := a.accService.GetPhotoById(ctx, photoIn.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  "error before update, while getting photo",
			Error: response.InternalServer,
		})
		return
	}

	con, message := UpdateEntityAuth(user.ID, photoIn.UserID, toBeUpdatedPhoto.UserID)
	if !con {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  message,
			Error: response.Unauthorized,
		})
		return
	}

	// validate name
	// if photoIn.Name == "" {
	// 	ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
	// 		Message:  "failed to update user",
	// 		Error: response.InvalidParam,
	// 	})
	// 	return
	// }
	logger.Info(ctx, "otw validate photo")
	con, message = PhotoPayloadValidation(photoIn)
	if !con {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  message,
			Error: response.InvalidParam,
		})
		return
	}

	updatedPhoto, err := a.accService.UpdatePhoto(ctx, photoIn);
	if err != nil {
		if err.Error() == "user is not found" {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{
				Message:  "failed to update user",
				Error: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Message:  "failed to update user",
			Error: response.SomethingWentWrong,
		})
		return
	}
	ctx.JSON(http.StatusAccepted, response.SuccessResponse{
		Message: "success update user",
		Data: updatedPhoto,
	})
}
func (a *AccountHandlerImpl) DeletePhoto(ctx *gin.Context) {
	idUint, err := a.getIdFromParam(ctx)
	if err != nil {
		return
	}
	
	photo, err := a.accService.GetPhotoById(ctx, idUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  "error get photo",
			Error: response.InternalServer,
		})
		return
	}

	user, err, message := a.AuthIncomingRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  message,
			Error: response.InternalServer,
		})
		return
	}
	
	con, message := DeleteEntityAuth(photo.UserID, user.ID)
	if !con {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  message,
			Error: response.Unauthorized,
		})
		return
	}
	
	deletedPhoto, err := a.accService.DeletePhoto(ctx, idUint)
	if err != nil {
		if err.Error() == "photo is not found" {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{
				Message:  "failed to delete photo",
				Error: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Message:  "failed to delete photo",
			Error: response.SomethingWentWrong,
		})
		return
	}
	ctx.JSON(http.StatusAccepted, response.SuccessResponse{
		Message: "success delete photo",
		Data:    deletedPhoto,
	})
}

func (a *AccountHandlerImpl) GetAllComments(ctx *gin.Context) {
	comments, err := a.accService.GetAllComments(ctx)
	if err != nil {
		// bad code, should be wrapped in other package
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Message: response.SomethingWentWrong,
			Error: "failed to get comments",
		})
		return
	}
	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Message: "success get comments",
		Data:    comments,
	})
}
func (a *AccountHandlerImpl) GetCommentById(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  "failed to find comment",
			Error: response.InvalidQuery,
		})
		return
	}
	// transform id string to uint64
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  "failed to find comment",
			Error: response.InvalidParam,
		})
		return
	}

	// call service
	comment, err := a.accService.GetCommentById(ctx, idUint)
	if err != nil {
		if err.Error() == "comment is not found" {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{
				Message:  "failed to find comment",
				Error: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Message:  "failed to find comment",
			Error: response.SomethingWentWrong,
		})
		return
	}
	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Message: "success find comment",
		Data:    comment,
	})
}
func (a *AccountHandlerImpl) CreateComment(ctx *gin.Context) {
	// mendapatkan body
	var commentIn accountmodel.Comment

	logger.Info(ctx, "otw bind comment")
	if err := ctx.Bind(&commentIn); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  "failed to insert comment",
			Error: response.InvalidBody,
		})
		return
	}

	logger.Info(ctx, "otw validate comment")
	// validate name and email
	// if commentIn.Title == "" || commentIn.PhotoUrl == "" {
	// 	ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
	// 		Message:  "failed to insert comment",
	// 		Error: response.InvalidParam,
	// 	})
	// 	return
	// }
	
	user, err, message := a.AuthIncomingRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  message,
			Error: response.InternalServer,
		})
		return
	}
	con, message := CreateEntityAuth(user.ID, commentIn.UserID)
	if con != true {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  message,
			Error: response.Unauthorized,
		})
		return
	}
	con, message = CommentPayloadValidation(commentIn)
	if !con {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  message,
			Error: response.InvalidParam,
		})
		return
	}

	logger.Info(ctx, "otw insert comment")
	insertedComment, err := a.accService.CreateComment(ctx, commentIn)
	if err != nil {
		// bad code, should be wrapped in other package
		if err.Error() == "error duplication email" {
			ctx.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
				Message:  "failed to insert comment",
				Error: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Message:  "failed to insert comment",
			Error: response.SomethingWentWrong,
		})
		return
	}

	ctx.JSON(http.StatusAccepted, response.SuccessResponse{
		Message: "success create comment",
		Data:    insertedComment,
	})
}
func (a *AccountHandlerImpl) UpdateComment(ctx *gin.Context) {
	idUint, err := a.getIdFromParam(ctx)
	if err != nil {
		return
	}
	// binding payload
	var commentIn accountmodel.Comment
	if err := ctx.Bind(&commentIn); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  "failed to update user",
			Error: response.InvalidBody,
		})
		return
	}
	commentIn.ID = idUint

	// validate name
	// if commentIn.Name == "" {
	// 	ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
	// 		Message:  "failed to update user",
	// 		Error: response.InvalidParam,
	// 	})
	// 	return
	// }
	user, err, message := a.AuthIncomingRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  message,
			Error: response.InternalServer,
		})
		return
	}

	toBeUpdatedComment, err := a.accService.GetCommentById(ctx, commentIn.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  "error before update, while getting comment",
			Error: response.InternalServer,
		})
		return
	}

	con, message := UpdateEntityAuth(user.ID, commentIn.UserID, toBeUpdatedComment.UserID)
	if con != true {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  message,
			Error: response.Unauthorized,
		})
		return
	}
	con, message = CommentPayloadValidation(commentIn)
	if !con {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  message,
			Error: response.InvalidParam,
		})
		return
	}

	updatedComment, err := a.accService.UpdateComment(ctx, commentIn);
	if err != nil {
		if err.Error() == "user is not found" {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{
				Message:  "failed to update user",
				Error: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Message:  "failed to update user",
			Error: response.SomethingWentWrong,
		})
		return
	}
	ctx.JSON(http.StatusAccepted, response.SuccessResponse{
		Message: "success update user",
		Data: updatedComment,
	})
}
func (a *AccountHandlerImpl) DeleteComment(ctx *gin.Context) {
	idUint, err := a.getIdFromParam(ctx)
	if err != nil {
		return
	}
	
	comment, err := a.accService.GetCommentById(ctx, idUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  "error get comment",
			Error: response.InternalServer,
		})
		return
	}

	user, err, message := a.AuthIncomingRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  message,
			Error: response.InternalServer,
		})
		return
	}
	
	con, message := DeleteEntityAuth(comment.UserID, user.ID)
	if !con {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  message,
			Error: response.Unauthorized,
		})
		return
	}

	deletedComment, err := a.accService.DeleteComment(ctx, idUint)
	if err != nil {
		if err.Error() == "comment is not found" {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{
				Message:  "failed to delete comment",
				Error: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Message:  "failed to delete comment",
			Error: response.SomethingWentWrong,
		})
		return
	}
	ctx.JSON(http.StatusAccepted, response.SuccessResponse{
		Message: "success delete comment",
		Data:    deletedComment,
	})
}

func (a *AccountHandlerImpl) GetAllSocialMedias(ctx *gin.Context) {
	socialMedias, err := a.accService.GetAllSocialMedias(ctx)
	if err != nil {
		// bad code, should be wrapped in other package
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Message: response.SomethingWentWrong,
			Error: "failed to get socialMedias",
		})
		return
	}
	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Message: "success get socialMedias",
		Data:    socialMedias,
	})
}
func (a *AccountHandlerImpl) GetSocialMediaById(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  "failed to find socialMedia",
			Error: response.InvalidQuery,
		})
		return
	}
	// transform id string to uint64
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  "failed to find socialMedia",
			Error: response.InvalidParam,
		})
		return
	}

	// call service
	socialMedia, err := a.accService.GetSocialMediaById(ctx, idUint)
	if err != nil {
		if err.Error() == "socialMedia is not found" {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{
				Message:  "failed to find socialMedia",
				Error: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Message:  "failed to find socialMedia",
			Error: response.SomethingWentWrong,
		})
		return
	}
	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Message: "success find socialMedia",
		Data:    socialMedia,
	})
}
func (a *AccountHandlerImpl) CreateSocialMedia(ctx *gin.Context) {
	// mendapatkan body
	var socialMediaIn accountmodel.SocialMedia

	logger.Info(ctx, "otw bind socialmedia")
	if err := ctx.Bind(&socialMediaIn); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  "failed to insert socialmedia",
			Error: response.InvalidBody,
		})
		return
	}

	logger.Info(ctx, "otw validate socialmedia")
	// validate name and email
	// if socialMediaIn.Title == "" || socialMediaIn.PhotoUrl == "" {
	// 	ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
	// 		Message:  "failed to insert socialmedia",
	// 		Error: response.InvalidParam,
	// 	})
	// 	return
	// }

	user, err, message := a.AuthIncomingRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  message,
			Error: response.InternalServer,
		})
		return
	}
	con, message := CreateEntityAuth(user.ID, socialMediaIn.UserID)
	if con != true {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  message,
			Error: response.Unauthorized,
		})
		return
	}

	con, message = SocialMediaPayloadValidation(socialMediaIn)
	if !con {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  message,
			Error: response.InvalidParam,
		})
		return
	}

	logger.Info(ctx, "otw insert socialmedia")
	insertedSocialMedia, err := a.accService.CreateSocialMedia(ctx, socialMediaIn)
	if err != nil {
		// bad code, should be wrapped in other package
		if err.Error() == "error duplication email" {
			ctx.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
				Message:  "failed to insert socialmedia",
				Error: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Message:  "failed to insert socialmedia",
			Error: response.SomethingWentWrong,
		})
		return
	}

	ctx.JSON(http.StatusAccepted, response.SuccessResponse{
		Message: "success create socialmedia",
		Data:    insertedSocialMedia,
	})
}
func (a *AccountHandlerImpl) UpdateSocialMedia(ctx *gin.Context) {
	idUint, err := a.getIdFromParam(ctx)
	if err != nil {
		return
	}
	// binding payload
	var socialMediaIn accountmodel.SocialMedia
	if err := ctx.Bind(&socialMediaIn); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  "failed to update user",
			Error: response.InvalidBody,
		})
		return
	}
	socialMediaIn.ID = idUint

	// validate name
	// if socialMediaIn.Name == "" {
	// 	ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
	// 		Message:  "failed to update user",
	// 		Error: response.InvalidParam,
	// 	})
	// 	return
	// }
	user, err, message := a.AuthIncomingRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  message,
			Error: response.InternalServer,
		})
		return
	}
	
	toBeUpdatedSocialMedia, err := a.accService.GetSocialMediaById(ctx, socialMediaIn.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  "error before update, while getting socialMedia",
			Error: response.InternalServer,
		})
		return
	}

	con, message := UpdateEntityAuth(user.ID, socialMediaIn.UserID, toBeUpdatedSocialMedia.UserID)
	if con != true {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  message,
			Error: response.Unauthorized,
		})
		return
	}
	con, message = SocialMediaPayloadValidation(socialMediaIn)
	if !con {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  message,
			Error: response.InvalidParam,
		})
		return
	}

	updatedSocialMedia, err := a.accService.UpdateSocialMedia(ctx, socialMediaIn);
	if err != nil {
		if err.Error() == "user is not found" {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{
				Message:  "failed to update user",
				Error: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Message:  "failed to update user",
			Error: response.SomethingWentWrong,
		})
		return
	}
	ctx.JSON(http.StatusAccepted, response.SuccessResponse{
		Message: "success update user",
		Data: updatedSocialMedia,
	})
}
func (a *AccountHandlerImpl) DeleteSocialMedia(ctx *gin.Context) {
	idUint, err := a.getIdFromParam(ctx)
	if err != nil {
		return
	}
	socialMedia, err := a.accService.GetSocialMediaById(ctx, idUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  "error get socialMedia",
			Error: response.InternalServer,
		})
		return
	}

	user, err, message := a.AuthIncomingRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  message,
			Error: response.InternalServer,
		})
		return
	}
	
	con, message := DeleteEntityAuth(socialMedia.UserID, user.ID)
	if !con {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message:  message,
			Error: response.Unauthorized,
		})
		return
	}


	deletedSocialMedia, err := a.accService.DeleteSocialMedia(ctx, idUint)
	if err != nil {
		if err.Error() == "social media is not found" {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{
				Message:  "failed to delete social media",
				Error: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Message:  "failed to delete social media",
			Error: response.SomethingWentWrong,
		})
		return
	}
	ctx.JSON(http.StatusAccepted, response.SuccessResponse{
		Message: "success delete social media",
		Data:    deletedSocialMedia,
	})
}
