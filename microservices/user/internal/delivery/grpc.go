package delivery

import (
	"context"

	userProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/user"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/user/internal/domain"
	model "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/user/model"
)

type UserService struct {
	userProto.UnimplementedUserServiceServer
	userUsecase domain.Usecase
}

func NewUserService(userUsecase domain.Usecase) *UserService {
	return &UserService{
		userUsecase: userUsecase,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *userProto.RegisterData) (*userProto.UserFront, error) {
	data := model.RegisterDataFromProtoToUsecase(req)
	user, err := s.userUsecase.CreateUser(ctx, data)
	if err != nil {
		return nil, err
	}
	return model.UserFrontFromUsecaseToProto(user), nil
}

func (s *UserService) LoginUser(ctx context.Context, req *userProto.LoginData) (*userProto.UserFront, error) {
	data := model.LoginDataFromProtoToUsecase(req)
	user, err := s.userUsecase.LoginUser(ctx, data)
	if err != nil {
		return nil, err
	}
	return model.UserFrontFromUsecaseToProto(user), nil
}

func (s *UserService) GetUserByID(ctx context.Context, req *userProto.UserID) (*userProto.UserFront, error) {
	user, err := s.userUsecase.GetUserByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return model.UserFrontFromUsecaseToProto(user), nil
}

func (s *UserService) UploadAvatar(ctx context.Context, req *userProto.AvatarData) (*userProto.Nothing, error) {
	err := s.userUsecase.UploadAvatar(ctx, req.AvatarPath, req.Id)
	if err != nil {
		return nil, err
	}
	return &userProto.Nothing{Dummy: true}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *userProto.UserDelete) (*userProto.Nothing, error) {
	err := s.userUsecase.DeleteUser(ctx, model.UserDeleteFromProtoToUsecase(req))
	if err != nil {
		return nil, err
	}
	return &userProto.Nothing{Dummy: true}, nil
}

func (s *UserService) ChangeUserData(ctx context.Context, req *userProto.ChangeUserDataMessage) (*userProto.Nothing, error) {
	err := s.userUsecase.ChangeUserData(ctx, req.Username, model.ChangeUserDataFromProtoToUsecase(req))
	if err != nil {
		return nil, err
	}
	return &userProto.Nothing{Dummy: true}, nil
}

func (s *UserService) ChangeUserPrivacySettings(ctx context.Context, req *userProto.PrivacySettings) (*userProto.Nothing, error) {
	err := s.userUsecase.ChangeUserPrivacySettings(ctx, req.Username, model.PrivacySettingsFromProtoToUsecase(req))
	if err != nil {
		return nil, err
	}
	return &userProto.Nothing{Dummy: true}, nil
}

func (s *UserService) GetUserFullData(ctx context.Context, req *userProto.Username) (*userProto.UserFullData, error) {
	user, err := s.userUsecase.GetFullUserData(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	return model.UserFullDataFromUsecaseToProto(user), nil
}
