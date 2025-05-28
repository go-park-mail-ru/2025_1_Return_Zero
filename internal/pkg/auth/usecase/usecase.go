package usecase

import (
	"context"

	authProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/auth"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/auth"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model"
)

func NewAuthUsecase(authClient *authProto.AuthServiceClient) auth.Usecase {
	return &AuthUsecase{authClient: authClient}
}

type AuthUsecase struct {
	authClient *authProto.AuthServiceClient
}

func (u *AuthUsecase) CreateSession(ctx context.Context, ID int64) (string, error) {
	sessionID, err := (*u.authClient).CreateSession(ctx, &authProto.UserID{Id: ID})
	if err != nil {
		return "", err
	}
	return model.SessionIDFromProtoToUsecase(sessionID), nil
}

func (u *AuthUsecase) DeleteSession(ctx context.Context, SID string) error {
	_, err := (*u.authClient).DeleteSession(ctx, &authProto.SessionID{SessionId: SID})
	if err != nil {
		return err
	}
	return nil
}

func (u *AuthUsecase) GetSession(ctx context.Context, SID string) (int64, error) {
	id, err := (*u.authClient).GetSession(ctx, &authProto.SessionID{SessionId: SID})
	if err != nil {
		return -1, err
	}
	return model.UserIDFromProtoToUsecase(id), nil
}
