// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go
//
// Generated by this command:
//
//	mockgen -source=repository.go -destination=mocks/mock_repository.go
//

// Package mock_jam is a generated GoMock package.
package mock_jam

import (
	context "context"
	reflect "reflect"

	repository "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
	gomock "go.uber.org/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
	isgomock struct{}
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// AddUser mocks base method.
func (m *MockRepository) AddUser(ctx context.Context, roomID, userID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUser", ctx, roomID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUser indicates an expected call of AddUser.
func (mr *MockRepositoryMockRecorder) AddUser(ctx, roomID, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUser", reflect.TypeOf((*MockRepository)(nil).AddUser), ctx, roomID, userID)
}

// CheckAllReadyAndPlay mocks base method.
func (m *MockRepository) CheckAllReadyAndPlay(ctx context.Context, roomID string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CheckAllReadyAndPlay", ctx, roomID)
}

// CheckAllReadyAndPlay indicates an expected call of CheckAllReadyAndPlay.
func (mr *MockRepositoryMockRecorder) CheckAllReadyAndPlay(ctx, roomID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckAllReadyAndPlay", reflect.TypeOf((*MockRepository)(nil).CheckAllReadyAndPlay), ctx, roomID)
}

// CreateJam mocks base method.
func (m *MockRepository) CreateJam(ctx context.Context, request *repository.CreateJamRequest) (*repository.CreateJamResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateJam", ctx, request)
	ret0, _ := ret[0].(*repository.CreateJamResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateJam indicates an expected call of CreateJam.
func (mr *MockRepositoryMockRecorder) CreateJam(ctx, request any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateJam", reflect.TypeOf((*MockRepository)(nil).CreateJam), ctx, request)
}

// ExistsRoom mocks base method.
func (m *MockRepository) ExistsRoom(ctx context.Context, roomID string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExistsRoom", ctx, roomID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExistsRoom indicates an expected call of ExistsRoom.
func (mr *MockRepositoryMockRecorder) ExistsRoom(ctx, roomID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExistsRoom", reflect.TypeOf((*MockRepository)(nil).ExistsRoom), ctx, roomID)
}

// GetHostID mocks base method.
func (m *MockRepository) GetHostID(ctx context.Context, roomID string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHostID", ctx, roomID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHostID indicates an expected call of GetHostID.
func (mr *MockRepositoryMockRecorder) GetHostID(ctx, roomID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHostID", reflect.TypeOf((*MockRepository)(nil).GetHostID), ctx, roomID)
}

// GetInitialJamData mocks base method.
func (m *MockRepository) GetInitialJamData(ctx context.Context, roomID string) (*repository.JamMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInitialJamData", ctx, roomID)
	ret0, _ := ret[0].(*repository.JamMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInitialJamData indicates an expected call of GetInitialJamData.
func (mr *MockRepositoryMockRecorder) GetInitialJamData(ctx, roomID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInitialJamData", reflect.TypeOf((*MockRepository)(nil).GetInitialJamData), ctx, roomID)
}

// GetUserInfo mocks base method.
func (m *MockRepository) GetUserInfo(ctx context.Context, roomID, userID string) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserInfo", ctx, roomID, userID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetUserInfo indicates an expected call of GetUserInfo.
func (mr *MockRepositoryMockRecorder) GetUserInfo(ctx, roomID, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserInfo", reflect.TypeOf((*MockRepository)(nil).GetUserInfo), ctx, roomID, userID)
}

// LoadTrack mocks base method.
func (m *MockRepository) LoadTrack(ctx context.Context, roomID, userID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadTrack", ctx, roomID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// LoadTrack indicates an expected call of LoadTrack.
func (mr *MockRepositoryMockRecorder) LoadTrack(ctx, roomID, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadTrack", reflect.TypeOf((*MockRepository)(nil).LoadTrack), ctx, roomID, userID)
}

// MarkUserAsReady mocks base method.
func (m *MockRepository) MarkUserAsReady(ctx context.Context, roomID, userID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkUserAsReady", ctx, roomID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// MarkUserAsReady indicates an expected call of MarkUserAsReady.
func (mr *MockRepositoryMockRecorder) MarkUserAsReady(ctx, roomID, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkUserAsReady", reflect.TypeOf((*MockRepository)(nil).MarkUserAsReady), ctx, roomID, userID)
}

// PauseJam mocks base method.
func (m *MockRepository) PauseJam(ctx context.Context, roomID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PauseJam", ctx, roomID)
	ret0, _ := ret[0].(error)
	return ret0
}

// PauseJam indicates an expected call of PauseJam.
func (mr *MockRepositoryMockRecorder) PauseJam(ctx, roomID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PauseJam", reflect.TypeOf((*MockRepository)(nil).PauseJam), ctx, roomID)
}

// RemoveJam mocks base method.
func (m *MockRepository) RemoveJam(ctx context.Context, roomID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveJam", ctx, roomID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveJam indicates an expected call of RemoveJam.
func (mr *MockRepositoryMockRecorder) RemoveJam(ctx, roomID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveJam", reflect.TypeOf((*MockRepository)(nil).RemoveJam), ctx, roomID)
}

// RemoveUser mocks base method.
func (m *MockRepository) RemoveUser(ctx context.Context, roomID, userID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveUser", ctx, roomID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveUser indicates an expected call of RemoveUser.
func (mr *MockRepositoryMockRecorder) RemoveUser(ctx, roomID, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveUser", reflect.TypeOf((*MockRepository)(nil).RemoveUser), ctx, roomID, userID)
}

// SeekJam mocks base method.
func (m *MockRepository) SeekJam(ctx context.Context, roomID string, position int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SeekJam", ctx, roomID, position)
	ret0, _ := ret[0].(error)
	return ret0
}

// SeekJam indicates an expected call of SeekJam.
func (mr *MockRepositoryMockRecorder) SeekJam(ctx, roomID, position any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SeekJam", reflect.TypeOf((*MockRepository)(nil).SeekJam), ctx, roomID, position)
}

// StoreUserInfo mocks base method.
func (m *MockRepository) StoreUserInfo(ctx context.Context, roomID, userID, username, avatarURL string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreUserInfo", ctx, roomID, userID, username, avatarURL)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreUserInfo indicates an expected call of StoreUserInfo.
func (mr *MockRepositoryMockRecorder) StoreUserInfo(ctx, roomID, userID, username, avatarURL any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreUserInfo", reflect.TypeOf((*MockRepository)(nil).StoreUserInfo), ctx, roomID, userID, username, avatarURL)
}

// SubscribeToJamMessages mocks base method.
func (m *MockRepository) SubscribeToJamMessages(ctx context.Context, roomID string) (<-chan []byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubscribeToJamMessages", ctx, roomID)
	ret0, _ := ret[0].(<-chan []byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubscribeToJamMessages indicates an expected call of SubscribeToJamMessages.
func (mr *MockRepositoryMockRecorder) SubscribeToJamMessages(ctx, roomID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeToJamMessages", reflect.TypeOf((*MockRepository)(nil).SubscribeToJamMessages), ctx, roomID)
}
