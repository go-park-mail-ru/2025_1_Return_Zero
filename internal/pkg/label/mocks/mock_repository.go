// Code generated by MockGen. DO NOT EDIT.
// Source: internal/pkg/label/domain/repository.go
//
// Generated by this command:
//
//	mockgen -source=internal/pkg/label/domain/repository.go -destination=internal/pkg/label/mocks/mock_repository.go -package mock_label
//

// Package mock_label is a generated GoMock package.
package mock_label

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

// CheckIsLabelUnique mocks base method.
func (m *MockRepository) CheckIsLabelUnique(ctx context.Context, labelName string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckIsLabelUnique", ctx, labelName)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckIsLabelUnique indicates an expected call of CheckIsLabelUnique.
func (mr *MockRepositoryMockRecorder) CheckIsLabelUnique(ctx, labelName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckIsLabelUnique", reflect.TypeOf((*MockRepository)(nil).CheckIsLabelUnique), ctx, labelName)
}

// CreateLabel mocks base method.
func (m *MockRepository) CreateLabel(ctx context.Context, name string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateLabel", ctx, name)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateLabel indicates an expected call of CreateLabel.
func (mr *MockRepositoryMockRecorder) CreateLabel(ctx, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateLabel", reflect.TypeOf((*MockRepository)(nil).CreateLabel), ctx, name)
}

// GetLabel mocks base method.
func (m *MockRepository) GetLabel(ctx context.Context, labelID int64) (*repository.Label, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLabel", ctx, labelID)
	ret0, _ := ret[0].(*repository.Label)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLabel indicates an expected call of GetLabel.
func (mr *MockRepositoryMockRecorder) GetLabel(ctx, labelID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLabel", reflect.TypeOf((*MockRepository)(nil).GetLabel), ctx, labelID)
}

// MockS3Repository is a mock of S3Repository interface.
type MockS3Repository struct {
	ctrl     *gomock.Controller
	recorder *MockS3RepositoryMockRecorder
	isgomock struct{}
}

// MockS3RepositoryMockRecorder is the mock recorder for MockS3Repository.
type MockS3RepositoryMockRecorder struct {
	mock *MockS3Repository
}

// NewMockS3Repository creates a new mock instance.
func NewMockS3Repository(ctrl *gomock.Controller) *MockS3Repository {
	mock := &MockS3Repository{ctrl: ctrl}
	mock.recorder = &MockS3RepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockS3Repository) EXPECT() *MockS3RepositoryMockRecorder {
	return m.recorder
}
