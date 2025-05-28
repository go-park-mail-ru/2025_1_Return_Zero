package delivery

import (
	"context"

	authProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/auth"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/auth/internal/domain"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/auth/model"
)

type AuthService struct {
	authProto.UnimplementedAuthServiceServer
	authUsecase domain.Usecase
}

func NewAuthService(authUsecase domain.Usecase) *AuthService {
	return &AuthService{
		authUsecase: authUsecase,
	}
}

func (s *AuthService) CreateSession(ctx context.Context, req *authProto.UserID) (*authProto.SessionID, error) {
	sessionID, err := s.authUsecase.CreateSession(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return model.SessionIDFromUsecaseToProto(sessionID), nil
}

func (s *AuthService) DeleteSession(ctx context.Context, req *authProto.SessionID) (*authProto.Nothing, error) {
	err := s.authUsecase.DeleteSession(ctx, req.SessionId)
	if err != nil {
		return nil, err
	}
	return model.NothingFromUsecaseToProto(), nil
}

func (s *AuthService) GetSession(ctx context.Context, req *authProto.SessionID) (*authProto.UserID, error) {
	userID, err := s.authUsecase.GetSession(ctx, req.SessionId)
	if err != nil {
		return nil, err
	}
	return model.UserIDFromUsecaseToProto(userID), nil
}
