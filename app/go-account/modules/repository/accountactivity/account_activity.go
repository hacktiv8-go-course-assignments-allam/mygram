package account

import (
	"context"

	activitymodel "github.com/mygram/go-account/modules/models/accountactivity"
)

type IAccountActivityRepo interface {
	CreateActivity(ctx context.Context, acc activitymodel.AccountActivity) (created activitymodel.AccountActivity, err error)
	CreateUserActivity(ctx context.Context, acc activitymodel.UserActivity) (created activitymodel.UserActivity, err error)
}
