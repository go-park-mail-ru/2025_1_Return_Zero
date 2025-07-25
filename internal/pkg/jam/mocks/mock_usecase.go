// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go
//
// Generated by this command:
//
//	mockgen -source=usecase.go -destination=mocks/mock_usecase.go
//

// Package mock_jam is a generated GoMock package.
package mock_jam

import (
	context "context"
	reflect "reflect"

	usecase "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	gomock "go.uber.org/mock/gomock"
)

// MockUsecase is a mock of Usecase interface.
type MockUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockUsecaseMockRecorder
	isgomock struct{}
}

// MockUsecaseMockRecorder is the mock recorder for MockUsecase.
type MockUsecaseMockRecorder struct {
	mock *MockUsecase
}

// NewMockUsecase creates a new mock instance.
func NewMockUsecase(ctrl *gomock.Controller) *MockUsecase {
	mock := &MockUsecase{ctrl: ctrl}
	mock.recorder = &MockUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUsecase) EXPECT() *MockUsecaseMockRecorder {
	return m.recorder
}

// CreateJam mocks base method.
func (m *MockUsecase) CreateJam(ctx context.Context, request *usecase.CreateJamRequest) (*usecase.CreateJamResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateJam", ctx, request)
	ret0, _ := ret[0].(*usecase.CreateJamResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateJam indicates an expected call of CreateJam.
func (mr *MockUsecaseMockRecorder) CreateJam(ctx, request any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateJam", reflect.TypeOf((*MockUsecase)(nil).CreateJam), ctx, request)
}

// HandleClientMessage mocks base method.
func (m_2 *MockUsecase) HandleClientMessage(ctx context.Context, roomID, userID string, m *usecase.JamMessage) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "HandleClientMessage", ctx, roomID, userID, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleClientMessage indicates an expected call of HandleClientMessage.
func (mr *MockUsecaseMockRecorder) HandleClientMessage(ctx, roomID, userID, m any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleClientMessage", reflect.TypeOf((*MockUsecase)(nil).HandleClientMessage), ctx, roomID, userID, m)
}

// JoinJam mocks base method.
func (m *MockUsecase) JoinJam(ctx context.Context, request *usecase.JoinJamRequest) (*usecase.JamMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "JoinJam", ctx, request)
	ret0, _ := ret[0].(*usecase.JamMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// JoinJam indicates an expected call of JoinJam.
func (mr *MockUsecaseMockRecorder) JoinJam(ctx, request any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "JoinJam", reflect.TypeOf((*MockUsecase)(nil).JoinJam), ctx, request)
}

// LeaveJam mocks base method.
func (m *MockUsecase) LeaveJam(ctx context.Context, roomID, userID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LeaveJam", ctx, roomID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// LeaveJam indicates an expected call of LeaveJam.
func (mr *MockUsecaseMockRecorder) LeaveJam(ctx, roomID, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LeaveJam", reflect.TypeOf((*MockUsecase)(nil).LeaveJam), ctx, roomID, userID)
}

// SubscribeToJamMessages mocks base method.
func (m *MockUsecase) SubscribeToJamMessages(ctx context.Context, roomID string) (<-chan *usecase.JamMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubscribeToJamMessages", ctx, roomID)
	ret0, _ := ret[0].(<-chan *usecase.JamMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubscribeToJamMessages indicates an expected call of SubscribeToJamMessages.
func (mr *MockUsecaseMockRecorder) SubscribeToJamMessages(ctx, roomID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeToJamMessages", reflect.TypeOf((*MockUsecase)(nil).SubscribeToJamMessages), ctx, roomID)
}
