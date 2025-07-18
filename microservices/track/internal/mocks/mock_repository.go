// Code generated by MockGen. DO NOT EDIT.
// Source: domain/repository.go
//
// Generated by this command:
//
//	mockgen -source=domain/repository.go -destination=mocks/mock_repository.go
//

// Package mock_domain is a generated GoMock package.
package mock_domain

import (
	context "context"
	reflect "reflect"

	repository "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/track/model/repository"
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

// AddTracksToAlbum mocks base method.
func (m *MockRepository) AddTracksToAlbum(ctx context.Context, tracksList []*repository.Track) ([]int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddTracksToAlbum", ctx, tracksList)
	ret0, _ := ret[0].([]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddTracksToAlbum indicates an expected call of AddTracksToAlbum.
func (mr *MockRepositoryMockRecorder) AddTracksToAlbum(ctx, tracksList any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTracksToAlbum", reflect.TypeOf((*MockRepository)(nil).AddTracksToAlbum), ctx, tracksList)
}

// CheckTrackExists mocks base method.
func (m *MockRepository) CheckTrackExists(ctx context.Context, trackID int64) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckTrackExists", ctx, trackID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckTrackExists indicates an expected call of CheckTrackExists.
func (mr *MockRepositoryMockRecorder) CheckTrackExists(ctx, trackID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckTrackExists", reflect.TypeOf((*MockRepository)(nil).CheckTrackExists), ctx, trackID)
}

// CreateStream mocks base method.
func (m *MockRepository) CreateStream(ctx context.Context, stream *repository.TrackStreamCreateData) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateStream", ctx, stream)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateStream indicates an expected call of CreateStream.
func (mr *MockRepositoryMockRecorder) CreateStream(ctx, stream any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateStream", reflect.TypeOf((*MockRepository)(nil).CreateStream), ctx, stream)
}

// DeleteTracksByAlbumID mocks base method.
func (m *MockRepository) DeleteTracksByAlbumID(ctx context.Context, albumID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTracksByAlbumID", ctx, albumID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTracksByAlbumID indicates an expected call of DeleteTracksByAlbumID.
func (mr *MockRepositoryMockRecorder) DeleteTracksByAlbumID(ctx, albumID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTracksByAlbumID", reflect.TypeOf((*MockRepository)(nil).DeleteTracksByAlbumID), ctx, albumID)
}

// GetAlbumIDByTrackID mocks base method.
func (m *MockRepository) GetAlbumIDByTrackID(ctx context.Context, id int64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAlbumIDByTrackID", ctx, id)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAlbumIDByTrackID indicates an expected call of GetAlbumIDByTrackID.
func (mr *MockRepositoryMockRecorder) GetAlbumIDByTrackID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAlbumIDByTrackID", reflect.TypeOf((*MockRepository)(nil).GetAlbumIDByTrackID), ctx, id)
}

// GetAllTracks mocks base method.
func (m *MockRepository) GetAllTracks(ctx context.Context, filters *repository.TrackFilters, userID int64) ([]*repository.Track, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllTracks", ctx, filters, userID)
	ret0, _ := ret[0].([]*repository.Track)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllTracks indicates an expected call of GetAllTracks.
func (mr *MockRepositoryMockRecorder) GetAllTracks(ctx, filters, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllTracks", reflect.TypeOf((*MockRepository)(nil).GetAllTracks), ctx, filters, userID)
}

// GetFavoriteTracks mocks base method.
func (m *MockRepository) GetFavoriteTracks(ctx context.Context, favoriteRequest *repository.FavoriteRequest) ([]*repository.Track, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFavoriteTracks", ctx, favoriteRequest)
	ret0, _ := ret[0].([]*repository.Track)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFavoriteTracks indicates an expected call of GetFavoriteTracks.
func (mr *MockRepositoryMockRecorder) GetFavoriteTracks(ctx, favoriteRequest any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFavoriteTracks", reflect.TypeOf((*MockRepository)(nil).GetFavoriteTracks), ctx, favoriteRequest)
}

// GetMinutesListenedByUserID mocks base method.
func (m *MockRepository) GetMinutesListenedByUserID(ctx context.Context, userID int64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMinutesListenedByUserID", ctx, userID)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMinutesListenedByUserID indicates an expected call of GetMinutesListenedByUserID.
func (mr *MockRepositoryMockRecorder) GetMinutesListenedByUserID(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMinutesListenedByUserID", reflect.TypeOf((*MockRepository)(nil).GetMinutesListenedByUserID), ctx, userID)
}

// GetMostLikedLastWeekTracks mocks base method.
func (m *MockRepository) GetMostLikedLastWeekTracks(ctx context.Context, userID int64) ([]*repository.Track, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMostLikedLastWeekTracks", ctx, userID)
	ret0, _ := ret[0].([]*repository.Track)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMostLikedLastWeekTracks indicates an expected call of GetMostLikedLastWeekTracks.
func (mr *MockRepositoryMockRecorder) GetMostLikedLastWeekTracks(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMostLikedLastWeekTracks", reflect.TypeOf((*MockRepository)(nil).GetMostLikedLastWeekTracks), ctx, userID)
}

// GetMostLikedTracks mocks base method.
func (m *MockRepository) GetMostLikedTracks(ctx context.Context, userID int64) ([]*repository.Track, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMostLikedTracks", ctx, userID)
	ret0, _ := ret[0].([]*repository.Track)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMostLikedTracks indicates an expected call of GetMostLikedTracks.
func (mr *MockRepositoryMockRecorder) GetMostLikedTracks(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMostLikedTracks", reflect.TypeOf((*MockRepository)(nil).GetMostLikedTracks), ctx, userID)
}

// GetMostListenedLastMonthTracks mocks base method.
func (m *MockRepository) GetMostListenedLastMonthTracks(ctx context.Context, userID int64) ([]*repository.Track, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMostListenedLastMonthTracks", ctx, userID)
	ret0, _ := ret[0].([]*repository.Track)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMostListenedLastMonthTracks indicates an expected call of GetMostListenedLastMonthTracks.
func (mr *MockRepositoryMockRecorder) GetMostListenedLastMonthTracks(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMostListenedLastMonthTracks", reflect.TypeOf((*MockRepository)(nil).GetMostListenedLastMonthTracks), ctx, userID)
}

// GetMostRecentTracks mocks base method.
func (m *MockRepository) GetMostRecentTracks(ctx context.Context, userID int64) ([]*repository.Track, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMostRecentTracks", ctx, userID)
	ret0, _ := ret[0].([]*repository.Track)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMostRecentTracks indicates an expected call of GetMostRecentTracks.
func (mr *MockRepositoryMockRecorder) GetMostRecentTracks(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMostRecentTracks", reflect.TypeOf((*MockRepository)(nil).GetMostRecentTracks), ctx, userID)
}

// GetStreamByID mocks base method.
func (m *MockRepository) GetStreamByID(ctx context.Context, streamID int64) (*repository.TrackStream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStreamByID", ctx, streamID)
	ret0, _ := ret[0].(*repository.TrackStream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStreamByID indicates an expected call of GetStreamByID.
func (mr *MockRepositoryMockRecorder) GetStreamByID(ctx, streamID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStreamByID", reflect.TypeOf((*MockRepository)(nil).GetStreamByID), ctx, streamID)
}

// GetStreamsByUserID mocks base method.
func (m *MockRepository) GetStreamsByUserID(ctx context.Context, userID int64, filters *repository.TrackFilters) ([]*repository.TrackStream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStreamsByUserID", ctx, userID, filters)
	ret0, _ := ret[0].([]*repository.TrackStream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStreamsByUserID indicates an expected call of GetStreamsByUserID.
func (mr *MockRepositoryMockRecorder) GetStreamsByUserID(ctx, userID, filters any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStreamsByUserID", reflect.TypeOf((*MockRepository)(nil).GetStreamsByUserID), ctx, userID, filters)
}

// GetTrackByID mocks base method.
func (m *MockRepository) GetTrackByID(ctx context.Context, id, userID int64) (*repository.TrackWithFileKey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTrackByID", ctx, id, userID)
	ret0, _ := ret[0].(*repository.TrackWithFileKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTrackByID indicates an expected call of GetTrackByID.
func (mr *MockRepositoryMockRecorder) GetTrackByID(ctx, id, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTrackByID", reflect.TypeOf((*MockRepository)(nil).GetTrackByID), ctx, id, userID)
}

// GetTracksByAlbumID mocks base method.
func (m *MockRepository) GetTracksByAlbumID(ctx context.Context, id, userID int64) ([]*repository.Track, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTracksByAlbumID", ctx, id, userID)
	ret0, _ := ret[0].([]*repository.Track)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTracksByAlbumID indicates an expected call of GetTracksByAlbumID.
func (mr *MockRepositoryMockRecorder) GetTracksByAlbumID(ctx, id, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTracksByAlbumID", reflect.TypeOf((*MockRepository)(nil).GetTracksByAlbumID), ctx, id, userID)
}

// GetTracksByIDs mocks base method.
func (m *MockRepository) GetTracksByIDs(ctx context.Context, ids []int64, userID int64) (map[int64]*repository.Track, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTracksByIDs", ctx, ids, userID)
	ret0, _ := ret[0].(map[int64]*repository.Track)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTracksByIDs indicates an expected call of GetTracksByIDs.
func (mr *MockRepositoryMockRecorder) GetTracksByIDs(ctx, ids, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTracksByIDs", reflect.TypeOf((*MockRepository)(nil).GetTracksByIDs), ctx, ids, userID)
}

// GetTracksByIDsFiltered mocks base method.
func (m *MockRepository) GetTracksByIDsFiltered(ctx context.Context, ids []int64, filters *repository.TrackFilters, userID int64) ([]*repository.Track, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTracksByIDsFiltered", ctx, ids, filters, userID)
	ret0, _ := ret[0].([]*repository.Track)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTracksByIDsFiltered indicates an expected call of GetTracksByIDsFiltered.
func (mr *MockRepositoryMockRecorder) GetTracksByIDsFiltered(ctx, ids, filters, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTracksByIDsFiltered", reflect.TypeOf((*MockRepository)(nil).GetTracksByIDsFiltered), ctx, ids, filters, userID)
}

// GetTracksListenedByUserID mocks base method.
func (m *MockRepository) GetTracksListenedByUserID(ctx context.Context, userID int64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTracksListenedByUserID", ctx, userID)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTracksListenedByUserID indicates an expected call of GetTracksListenedByUserID.
func (mr *MockRepositoryMockRecorder) GetTracksListenedByUserID(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTracksListenedByUserID", reflect.TypeOf((*MockRepository)(nil).GetTracksListenedByUserID), ctx, userID)
}

// LikeTrack mocks base method.
func (m *MockRepository) LikeTrack(ctx context.Context, likeRequest *repository.LikeRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LikeTrack", ctx, likeRequest)
	ret0, _ := ret[0].(error)
	return ret0
}

// LikeTrack indicates an expected call of LikeTrack.
func (mr *MockRepositoryMockRecorder) LikeTrack(ctx, likeRequest any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LikeTrack", reflect.TypeOf((*MockRepository)(nil).LikeTrack), ctx, likeRequest)
}

// SearchTracks mocks base method.
func (m *MockRepository) SearchTracks(ctx context.Context, query string, userID int64) ([]*repository.Track, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchTracks", ctx, query, userID)
	ret0, _ := ret[0].([]*repository.Track)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchTracks indicates an expected call of SearchTracks.
func (mr *MockRepositoryMockRecorder) SearchTracks(ctx, query, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchTracks", reflect.TypeOf((*MockRepository)(nil).SearchTracks), ctx, query, userID)
}

// UnlikeTrack mocks base method.
func (m *MockRepository) UnlikeTrack(ctx context.Context, likeRequest *repository.LikeRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnlikeTrack", ctx, likeRequest)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnlikeTrack indicates an expected call of UnlikeTrack.
func (mr *MockRepositoryMockRecorder) UnlikeTrack(ctx, likeRequest any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnlikeTrack", reflect.TypeOf((*MockRepository)(nil).UnlikeTrack), ctx, likeRequest)
}

// UpdateStreamDuration mocks base method.
func (m *MockRepository) UpdateStreamDuration(ctx context.Context, endedStream *repository.TrackStreamUpdateData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStreamDuration", ctx, endedStream)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateStreamDuration indicates an expected call of UpdateStreamDuration.
func (mr *MockRepositoryMockRecorder) UpdateStreamDuration(ctx, endedStream any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStreamDuration", reflect.TypeOf((*MockRepository)(nil).UpdateStreamDuration), ctx, endedStream)
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

// GetPresignedURL mocks base method.
func (m *MockS3Repository) GetPresignedURL(trackKey string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPresignedURL", trackKey)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPresignedURL indicates an expected call of GetPresignedURL.
func (mr *MockS3RepositoryMockRecorder) GetPresignedURL(trackKey any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPresignedURL", reflect.TypeOf((*MockS3Repository)(nil).GetPresignedURL), trackKey)
}

// UploadTrack mocks base method.
func (m *MockS3Repository) UploadTrack(ctx context.Context, fileKey string, file []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadTrack", ctx, fileKey, file)
	ret0, _ := ret[0].(error)
	return ret0
}

// UploadTrack indicates an expected call of UploadTrack.
func (mr *MockS3RepositoryMockRecorder) UploadTrack(ctx, fileKey, file any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadTrack", reflect.TypeOf((*MockS3Repository)(nil).UploadTrack), ctx, fileKey, file)
}

// UploadTrackAvatar mocks base method.
func (m *MockS3Repository) UploadTrackAvatar(ctx context.Context, trackTitle string, file []byte) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadTrackAvatar", ctx, trackTitle, file)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadTrackAvatar indicates an expected call of UploadTrackAvatar.
func (mr *MockS3RepositoryMockRecorder) UploadTrackAvatar(ctx, trackTitle, file any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadTrackAvatar", reflect.TypeOf((*MockS3Repository)(nil).UploadTrackAvatar), ctx, trackTitle, file)
}
