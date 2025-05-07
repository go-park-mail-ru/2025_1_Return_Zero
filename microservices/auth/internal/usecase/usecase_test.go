package usecase

import (
	"context"
	"testing"

	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	mock_domain "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/auth/internal/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func setupTest(t *testing.T) (*mock_domain.MockRepository, context.Context) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_domain.NewMockRepository(ctrl)

	logger := zap.NewNop().Sugar()
	ctx := loggerPkg.LoggerToContext(context.Background(), logger)

	return mockRepo, ctx
}

func TestCreateSession(t *testing.T) {
    mockRepo, ctx := setupTest(t)
    usecase := NewAuthUsecase(mockRepo)
    
    userID := int64(1)
    expectedSessionID := "test-session-id"
    
    mockRepo.EXPECT().CreateSession(ctx, userID).Return(expectedSessionID, nil)
    
    sessionID, err := usecase.CreateSession(ctx, userID)
    
    assert.NoError(t, err)
    assert.Equal(t, expectedSessionID, sessionID)
}

func TestGetSession(t *testing.T) {
    mockRepo, ctx := setupTest(t)
    usecase := NewAuthUsecase(mockRepo)
    
    sessionID := "test-session-id"
    expectedUserID := int64(1)
    
    mockRepo.EXPECT().GetSession(ctx, sessionID).Return(expectedUserID, nil)
    
    userID, err := usecase.GetSession(ctx, sessionID)
    
    assert.NoError(t, err)
    assert.Equal(t, expectedUserID, userID)
}