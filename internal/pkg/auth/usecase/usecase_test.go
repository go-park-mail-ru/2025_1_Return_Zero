package usecase

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	"go.uber.org/mock/gomock"

// 	mocks "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/auth/mocks"
// )

// func TestAuthUsecase_CreateSession(t *testing.T) {
// 	type mockBehavior func(r *mocks.MockRepository, ctx context.Context, userID int64)

// 	tests := []struct {
// 		name          string
// 		userID        int64
// 		mockBehavior  mockBehavior
// 		expectedID    string
// 		expectedError error
// 	}{
// 		{
// 			name:   "Success",
// 			userID: 1,
// 			mockBehavior: func(r *mocks.MockRepository, ctx context.Context, userID int64) {
// 				r.EXPECT().CreateSession(ctx, userID).Return("session123", nil)
// 			},
// 			expectedID:    "session123",
// 			expectedError: nil,
// 		},
// 		{
// 			name:   "Repository Error",
// 			userID: 1,
// 			mockBehavior: func(r *mocks.MockRepository, ctx context.Context, userID int64) {
// 				r.EXPECT().CreateSession(ctx, userID).Return("", errors.New("repo error"))
// 			},
// 			expectedID:    "",
// 			expectedError: errors.New("repo error"),
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			repo := mocks.NewMockRepository(ctrl)
// 			usecase := NewAuthUsecase(repo)
// 			ctx := context.Background()

// 			tt.mockBehavior(repo, ctx, tt.userID)

// 			sessionID, err := usecase.CreateSession(ctx, tt.userID)

// 			if tt.expectedError != nil && err == nil {
// 				t.Errorf("expected error: %v, got: nil", tt.expectedError)
// 			}
// 			if tt.expectedError == nil && err != nil {
// 				t.Errorf("expected no error, got: %v", err)
// 			}
// 			if tt.expectedError != nil && err != nil && tt.expectedError.Error() != err.Error() {
// 				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
// 			}
// 			if sessionID != tt.expectedID {
// 				t.Errorf("expected session ID: %s, got: %s", tt.expectedID, sessionID)
// 			}
// 		})
// 	}
// }

// func TestAuthUsecase_DeleteSession(t *testing.T) {
// 	type mockBehavior func(r *mocks.MockRepository, ctx context.Context, sessionID string)

// 	tests := []struct {
// 		name          string
// 		sessionID     string
// 		mockBehavior  mockBehavior
// 		expectedError error
// 	}{
// 		{
// 			name:      "Success",
// 			sessionID: "session123",
// 			mockBehavior: func(r *mocks.MockRepository, ctx context.Context, sessionID string) {
// 				r.EXPECT().DeleteSession(ctx, sessionID).Return(nil)
// 			},
// 			expectedError: nil,
// 		},
// 		{
// 			name:      "Repository Error",
// 			sessionID: "session123",
// 			mockBehavior: func(r *mocks.MockRepository, ctx context.Context, sessionID string) {
// 				r.EXPECT().DeleteSession(ctx, sessionID).Return(errors.New("repo error"))
// 			},
// 			expectedError: errors.New("repo error"),
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			repo := mocks.NewMockRepository(ctrl)
// 			usecase := NewAuthUsecase(repo)
// 			ctx := context.Background()

// 			tt.mockBehavior(repo, ctx, tt.sessionID)

// 			err := usecase.DeleteSession(ctx, tt.sessionID)

// 			if tt.expectedError != nil && err == nil {
// 				t.Errorf("expected error: %v, got: nil", tt.expectedError)
// 			}
// 			if tt.expectedError == nil && err != nil {
// 				t.Errorf("expected no error, got: %v", err)
// 			}
// 			if tt.expectedError != nil && err != nil && tt.expectedError.Error() != err.Error() {
// 				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
// 			}
// 		})
// 	}
// }

// func TestAuthUsecase_GetSession(t *testing.T) {
// 	type mockBehavior func(r *mocks.MockRepository, ctx context.Context, sessionID string)

// 	tests := []struct {
// 		name          string
// 		sessionID     string
// 		mockBehavior  mockBehavior
// 		expectedID    int64
// 		expectedError error
// 	}{
// 		{
// 			name:      "Success",
// 			sessionID: "session123",
// 			mockBehavior: func(r *mocks.MockRepository, ctx context.Context, sessionID string) {
// 				r.EXPECT().GetSession(ctx, sessionID).Return(int64(1), nil)
// 			},
// 			expectedID:    1,
// 			expectedError: nil,
// 		},
// 		{
// 			name:      "Repository Error",
// 			sessionID: "session123",
// 			mockBehavior: func(r *mocks.MockRepository, ctx context.Context, sessionID string) {
// 				r.EXPECT().GetSession(ctx, sessionID).Return(int64(0), errors.New("repo error"))
// 			},
// 			expectedID:    -1,
// 			expectedError: errors.New("repo error"),
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			repo := mocks.NewMockRepository(ctrl)
// 			usecase := NewAuthUsecase(repo)
// 			ctx := context.Background()

// 			tt.mockBehavior(repo, ctx, tt.sessionID)

// 			id, err := usecase.GetSession(ctx, tt.sessionID)

// 			if tt.expectedError != nil && err == nil {
// 				t.Errorf("expected error: %v, got: nil", tt.expectedError)
// 			}
// 			if tt.expectedError == nil && err != nil {
// 				t.Errorf("expected no error, got: %v", err)
// 			}
// 			if tt.expectedError != nil && err != nil && tt.expectedError.Error() != err.Error() {
// 				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
// 			}
// 			if id != tt.expectedID {
// 				t.Errorf("expected ID: %d, got: %d", tt.expectedID, id)
// 			}
// 		})
// 	}
// }
