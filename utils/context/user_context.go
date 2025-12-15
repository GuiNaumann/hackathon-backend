package context

import (
	"context"
	"hackathon-backend/domain/entities"
)

const CtxUserKey = "auth-ctx-user-data"

func GetUserFromContext(ctx context.Context) (*entities.User, bool) {
	user, ok := ctx.Value(CtxUserKey).(*entities.User)
	return user, ok
}

func SetUserInContext(ctx context.Context, user *entities.User) context.Context {
	return context.WithValue(ctx, CtxUserKey, user)
}
