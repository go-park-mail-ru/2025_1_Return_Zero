// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package track

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// TrackServiceClient is the client API for TrackService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TrackServiceClient interface {
	GetAllTracks(ctx context.Context, in *UserIDWithFilters, opts ...grpc.CallOption) (*TrackList, error)
	GetTrackByID(ctx context.Context, in *TrackIDWithUserID, opts ...grpc.CallOption) (*TrackDetailed, error)
	CreateStream(ctx context.Context, in *TrackStreamCreateData, opts ...grpc.CallOption) (*StreamID, error)
	UpdateStreamDuration(ctx context.Context, in *TrackStreamUpdateData, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetLastListenedTracks(ctx context.Context, in *UserIDWithFilters, opts ...grpc.CallOption) (*TrackList, error)
	GetTracksByIDs(ctx context.Context, in *TrackIDList, opts ...grpc.CallOption) (*TrackList, error)
	GetTracksByIDsFiltered(ctx context.Context, in *TrackIDListWithFilters, opts ...grpc.CallOption) (*TrackList, error)
	GetAlbumIDByTrackID(ctx context.Context, in *TrackID, opts ...grpc.CallOption) (*AlbumID, error)
	GetTracksByAlbumID(ctx context.Context, in *AlbumIDWithUserID, opts ...grpc.CallOption) (*TrackList, error)
	GetMinutesListenedByUserID(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*MinutesListened, error)
	GetTracksListenedByUserID(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*TracksListened, error)
	LikeTrack(ctx context.Context, in *LikeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	SearchTracks(ctx context.Context, in *Query, opts ...grpc.CallOption) (*TrackList, error)
	GetFavoriteTracks(ctx context.Context, in *FavoriteRequest, opts ...grpc.CallOption) (*TrackList, error)
	AddTracksToAlbum(ctx context.Context, in *TracksListWithAlbumID, opts ...grpc.CallOption) (*TrackIdsList, error)
	DeleteTracksByAlbumID(ctx context.Context, in *AlbumID, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetMostLikedTracks(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*TrackList, error)
	GetMostLikedLastWeekTracks(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*TrackList, error)
	GetMostListenedLastMonthTracks(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*TrackList, error)
	GetMostRecentTracks(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*TrackList, error)
}

type trackServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTrackServiceClient(cc grpc.ClientConnInterface) TrackServiceClient {
	return &trackServiceClient{cc}
}

func (c *trackServiceClient) GetAllTracks(ctx context.Context, in *UserIDWithFilters, opts ...grpc.CallOption) (*TrackList, error) {
	out := new(TrackList)
	err := c.cc.Invoke(ctx, "/track.TrackService/GetAllTracks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trackServiceClient) GetTrackByID(ctx context.Context, in *TrackIDWithUserID, opts ...grpc.CallOption) (*TrackDetailed, error) {
	out := new(TrackDetailed)
	err := c.cc.Invoke(ctx, "/track.TrackService/GetTrackByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trackServiceClient) CreateStream(ctx context.Context, in *TrackStreamCreateData, opts ...grpc.CallOption) (*StreamID, error) {
	out := new(StreamID)
	err := c.cc.Invoke(ctx, "/track.TrackService/CreateStream", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trackServiceClient) UpdateStreamDuration(ctx context.Context, in *TrackStreamUpdateData, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/track.TrackService/UpdateStreamDuration", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trackServiceClient) GetLastListenedTracks(ctx context.Context, in *UserIDWithFilters, opts ...grpc.CallOption) (*TrackList, error) {
	out := new(TrackList)
	err := c.cc.Invoke(ctx, "/track.TrackService/GetLastListenedTracks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trackServiceClient) GetTracksByIDs(ctx context.Context, in *TrackIDList, opts ...grpc.CallOption) (*TrackList, error) {
	out := new(TrackList)
	err := c.cc.Invoke(ctx, "/track.TrackService/GetTracksByIDs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trackServiceClient) GetTracksByIDsFiltered(ctx context.Context, in *TrackIDListWithFilters, opts ...grpc.CallOption) (*TrackList, error) {
	out := new(TrackList)
	err := c.cc.Invoke(ctx, "/track.TrackService/GetTracksByIDsFiltered", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trackServiceClient) GetAlbumIDByTrackID(ctx context.Context, in *TrackID, opts ...grpc.CallOption) (*AlbumID, error) {
	out := new(AlbumID)
	err := c.cc.Invoke(ctx, "/track.TrackService/GetAlbumIDByTrackID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trackServiceClient) GetTracksByAlbumID(ctx context.Context, in *AlbumIDWithUserID, opts ...grpc.CallOption) (*TrackList, error) {
	out := new(TrackList)
	err := c.cc.Invoke(ctx, "/track.TrackService/GetTracksByAlbumID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trackServiceClient) GetMinutesListenedByUserID(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*MinutesListened, error) {
	out := new(MinutesListened)
	err := c.cc.Invoke(ctx, "/track.TrackService/GetMinutesListenedByUserID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trackServiceClient) GetTracksListenedByUserID(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*TracksListened, error) {
	out := new(TracksListened)
	err := c.cc.Invoke(ctx, "/track.TrackService/GetTracksListenedByUserID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trackServiceClient) LikeTrack(ctx context.Context, in *LikeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/track.TrackService/LikeTrack", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trackServiceClient) SearchTracks(ctx context.Context, in *Query, opts ...grpc.CallOption) (*TrackList, error) {
	out := new(TrackList)
	err := c.cc.Invoke(ctx, "/track.TrackService/SearchTracks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trackServiceClient) GetFavoriteTracks(ctx context.Context, in *FavoriteRequest, opts ...grpc.CallOption) (*TrackList, error) {
	out := new(TrackList)
	err := c.cc.Invoke(ctx, "/track.TrackService/GetFavoriteTracks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trackServiceClient) AddTracksToAlbum(ctx context.Context, in *TracksListWithAlbumID, opts ...grpc.CallOption) (*TrackIdsList, error) {
	out := new(TrackIdsList)
	err := c.cc.Invoke(ctx, "/track.TrackService/AddTracksToAlbum", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trackServiceClient) DeleteTracksByAlbumID(ctx context.Context, in *AlbumID, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/track.TrackService/DeleteTracksByAlbumID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trackServiceClient) GetMostLikedTracks(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*TrackList, error) {
	out := new(TrackList)
	err := c.cc.Invoke(ctx, "/track.TrackService/GetMostLikedTracks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trackServiceClient) GetMostLikedLastWeekTracks(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*TrackList, error) {
	out := new(TrackList)
	err := c.cc.Invoke(ctx, "/track.TrackService/GetMostLikedLastWeekTracks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trackServiceClient) GetMostListenedLastMonthTracks(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*TrackList, error) {
	out := new(TrackList)
	err := c.cc.Invoke(ctx, "/track.TrackService/GetMostListenedLastMonthTracks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trackServiceClient) GetMostRecentTracks(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*TrackList, error) {
	out := new(TrackList)
	err := c.cc.Invoke(ctx, "/track.TrackService/GetMostRecentTracks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TrackServiceServer is the server API for TrackService service.
// All implementations must embed UnimplementedTrackServiceServer
// for forward compatibility
type TrackServiceServer interface {
	GetAllTracks(context.Context, *UserIDWithFilters) (*TrackList, error)
	GetTrackByID(context.Context, *TrackIDWithUserID) (*TrackDetailed, error)
	CreateStream(context.Context, *TrackStreamCreateData) (*StreamID, error)
	UpdateStreamDuration(context.Context, *TrackStreamUpdateData) (*emptypb.Empty, error)
	GetLastListenedTracks(context.Context, *UserIDWithFilters) (*TrackList, error)
	GetTracksByIDs(context.Context, *TrackIDList) (*TrackList, error)
	GetTracksByIDsFiltered(context.Context, *TrackIDListWithFilters) (*TrackList, error)
	GetAlbumIDByTrackID(context.Context, *TrackID) (*AlbumID, error)
	GetTracksByAlbumID(context.Context, *AlbumIDWithUserID) (*TrackList, error)
	GetMinutesListenedByUserID(context.Context, *UserID) (*MinutesListened, error)
	GetTracksListenedByUserID(context.Context, *UserID) (*TracksListened, error)
	LikeTrack(context.Context, *LikeRequest) (*emptypb.Empty, error)
	SearchTracks(context.Context, *Query) (*TrackList, error)
	GetFavoriteTracks(context.Context, *FavoriteRequest) (*TrackList, error)
	AddTracksToAlbum(context.Context, *TracksListWithAlbumID) (*TrackIdsList, error)
	DeleteTracksByAlbumID(context.Context, *AlbumID) (*emptypb.Empty, error)
	GetMostLikedTracks(context.Context, *UserID) (*TrackList, error)
	GetMostLikedLastWeekTracks(context.Context, *UserID) (*TrackList, error)
	GetMostListenedLastMonthTracks(context.Context, *UserID) (*TrackList, error)
	GetMostRecentTracks(context.Context, *UserID) (*TrackList, error)
	mustEmbedUnimplementedTrackServiceServer()
}

// UnimplementedTrackServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTrackServiceServer struct {
}

func (UnimplementedTrackServiceServer) GetAllTracks(context.Context, *UserIDWithFilters) (*TrackList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllTracks not implemented")
}
func (UnimplementedTrackServiceServer) GetTrackByID(context.Context, *TrackIDWithUserID) (*TrackDetailed, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTrackByID not implemented")
}
func (UnimplementedTrackServiceServer) CreateStream(context.Context, *TrackStreamCreateData) (*StreamID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateStream not implemented")
}
func (UnimplementedTrackServiceServer) UpdateStreamDuration(context.Context, *TrackStreamUpdateData) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateStreamDuration not implemented")
}
func (UnimplementedTrackServiceServer) GetLastListenedTracks(context.Context, *UserIDWithFilters) (*TrackList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLastListenedTracks not implemented")
}
func (UnimplementedTrackServiceServer) GetTracksByIDs(context.Context, *TrackIDList) (*TrackList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTracksByIDs not implemented")
}
func (UnimplementedTrackServiceServer) GetTracksByIDsFiltered(context.Context, *TrackIDListWithFilters) (*TrackList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTracksByIDsFiltered not implemented")
}
func (UnimplementedTrackServiceServer) GetAlbumIDByTrackID(context.Context, *TrackID) (*AlbumID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAlbumIDByTrackID not implemented")
}
func (UnimplementedTrackServiceServer) GetTracksByAlbumID(context.Context, *AlbumIDWithUserID) (*TrackList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTracksByAlbumID not implemented")
}
func (UnimplementedTrackServiceServer) GetMinutesListenedByUserID(context.Context, *UserID) (*MinutesListened, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMinutesListenedByUserID not implemented")
}
func (UnimplementedTrackServiceServer) GetTracksListenedByUserID(context.Context, *UserID) (*TracksListened, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTracksListenedByUserID not implemented")
}
func (UnimplementedTrackServiceServer) LikeTrack(context.Context, *LikeRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LikeTrack not implemented")
}
func (UnimplementedTrackServiceServer) SearchTracks(context.Context, *Query) (*TrackList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchTracks not implemented")
}
func (UnimplementedTrackServiceServer) GetFavoriteTracks(context.Context, *FavoriteRequest) (*TrackList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFavoriteTracks not implemented")
}
func (UnimplementedTrackServiceServer) AddTracksToAlbum(context.Context, *TracksListWithAlbumID) (*TrackIdsList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddTracksToAlbum not implemented")
}
func (UnimplementedTrackServiceServer) DeleteTracksByAlbumID(context.Context, *AlbumID) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTracksByAlbumID not implemented")
}
func (UnimplementedTrackServiceServer) GetMostLikedTracks(context.Context, *UserID) (*TrackList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMostLikedTracks not implemented")
}
func (UnimplementedTrackServiceServer) GetMostLikedLastWeekTracks(context.Context, *UserID) (*TrackList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMostLikedLastWeekTracks not implemented")
}
func (UnimplementedTrackServiceServer) GetMostListenedLastMonthTracks(context.Context, *UserID) (*TrackList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMostListenedLastMonthTracks not implemented")
}
func (UnimplementedTrackServiceServer) GetMostRecentTracks(context.Context, *UserID) (*TrackList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMostRecentTracks not implemented")
}
func (UnimplementedTrackServiceServer) mustEmbedUnimplementedTrackServiceServer() {}

// UnsafeTrackServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TrackServiceServer will
// result in compilation errors.
type UnsafeTrackServiceServer interface {
	mustEmbedUnimplementedTrackServiceServer()
}

func RegisterTrackServiceServer(s grpc.ServiceRegistrar, srv TrackServiceServer) {
	s.RegisterService(&TrackService_ServiceDesc, srv)
}

func _TrackService_GetAllTracks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserIDWithFilters)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrackServiceServer).GetAllTracks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/track.TrackService/GetAllTracks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrackServiceServer).GetAllTracks(ctx, req.(*UserIDWithFilters))
	}
	return interceptor(ctx, in, info, handler)
}

func _TrackService_GetTrackByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TrackIDWithUserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrackServiceServer).GetTrackByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/track.TrackService/GetTrackByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrackServiceServer).GetTrackByID(ctx, req.(*TrackIDWithUserID))
	}
	return interceptor(ctx, in, info, handler)
}

func _TrackService_CreateStream_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TrackStreamCreateData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrackServiceServer).CreateStream(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/track.TrackService/CreateStream",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrackServiceServer).CreateStream(ctx, req.(*TrackStreamCreateData))
	}
	return interceptor(ctx, in, info, handler)
}

func _TrackService_UpdateStreamDuration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TrackStreamUpdateData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrackServiceServer).UpdateStreamDuration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/track.TrackService/UpdateStreamDuration",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrackServiceServer).UpdateStreamDuration(ctx, req.(*TrackStreamUpdateData))
	}
	return interceptor(ctx, in, info, handler)
}

func _TrackService_GetLastListenedTracks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserIDWithFilters)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrackServiceServer).GetLastListenedTracks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/track.TrackService/GetLastListenedTracks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrackServiceServer).GetLastListenedTracks(ctx, req.(*UserIDWithFilters))
	}
	return interceptor(ctx, in, info, handler)
}

func _TrackService_GetTracksByIDs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TrackIDList)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrackServiceServer).GetTracksByIDs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/track.TrackService/GetTracksByIDs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrackServiceServer).GetTracksByIDs(ctx, req.(*TrackIDList))
	}
	return interceptor(ctx, in, info, handler)
}

func _TrackService_GetTracksByIDsFiltered_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TrackIDListWithFilters)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrackServiceServer).GetTracksByIDsFiltered(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/track.TrackService/GetTracksByIDsFiltered",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrackServiceServer).GetTracksByIDsFiltered(ctx, req.(*TrackIDListWithFilters))
	}
	return interceptor(ctx, in, info, handler)
}

func _TrackService_GetAlbumIDByTrackID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TrackID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrackServiceServer).GetAlbumIDByTrackID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/track.TrackService/GetAlbumIDByTrackID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrackServiceServer).GetAlbumIDByTrackID(ctx, req.(*TrackID))
	}
	return interceptor(ctx, in, info, handler)
}

func _TrackService_GetTracksByAlbumID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AlbumIDWithUserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrackServiceServer).GetTracksByAlbumID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/track.TrackService/GetTracksByAlbumID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrackServiceServer).GetTracksByAlbumID(ctx, req.(*AlbumIDWithUserID))
	}
	return interceptor(ctx, in, info, handler)
}

func _TrackService_GetMinutesListenedByUserID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrackServiceServer).GetMinutesListenedByUserID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/track.TrackService/GetMinutesListenedByUserID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrackServiceServer).GetMinutesListenedByUserID(ctx, req.(*UserID))
	}
	return interceptor(ctx, in, info, handler)
}

func _TrackService_GetTracksListenedByUserID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrackServiceServer).GetTracksListenedByUserID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/track.TrackService/GetTracksListenedByUserID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrackServiceServer).GetTracksListenedByUserID(ctx, req.(*UserID))
	}
	return interceptor(ctx, in, info, handler)
}

func _TrackService_LikeTrack_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LikeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrackServiceServer).LikeTrack(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/track.TrackService/LikeTrack",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrackServiceServer).LikeTrack(ctx, req.(*LikeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TrackService_SearchTracks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Query)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrackServiceServer).SearchTracks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/track.TrackService/SearchTracks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrackServiceServer).SearchTracks(ctx, req.(*Query))
	}
	return interceptor(ctx, in, info, handler)
}

func _TrackService_GetFavoriteTracks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FavoriteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrackServiceServer).GetFavoriteTracks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/track.TrackService/GetFavoriteTracks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrackServiceServer).GetFavoriteTracks(ctx, req.(*FavoriteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TrackService_AddTracksToAlbum_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TracksListWithAlbumID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrackServiceServer).AddTracksToAlbum(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/track.TrackService/AddTracksToAlbum",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrackServiceServer).AddTracksToAlbum(ctx, req.(*TracksListWithAlbumID))
	}
	return interceptor(ctx, in, info, handler)
}

func _TrackService_DeleteTracksByAlbumID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AlbumID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrackServiceServer).DeleteTracksByAlbumID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/track.TrackService/DeleteTracksByAlbumID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrackServiceServer).DeleteTracksByAlbumID(ctx, req.(*AlbumID))
	}
	return interceptor(ctx, in, info, handler)
}

func _TrackService_GetMostLikedTracks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrackServiceServer).GetMostLikedTracks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/track.TrackService/GetMostLikedTracks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrackServiceServer).GetMostLikedTracks(ctx, req.(*UserID))
	}
	return interceptor(ctx, in, info, handler)
}

func _TrackService_GetMostLikedLastWeekTracks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrackServiceServer).GetMostLikedLastWeekTracks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/track.TrackService/GetMostLikedLastWeekTracks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrackServiceServer).GetMostLikedLastWeekTracks(ctx, req.(*UserID))
	}
	return interceptor(ctx, in, info, handler)
}

func _TrackService_GetMostListenedLastMonthTracks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrackServiceServer).GetMostListenedLastMonthTracks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/track.TrackService/GetMostListenedLastMonthTracks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrackServiceServer).GetMostListenedLastMonthTracks(ctx, req.(*UserID))
	}
	return interceptor(ctx, in, info, handler)
}

func _TrackService_GetMostRecentTracks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrackServiceServer).GetMostRecentTracks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/track.TrackService/GetMostRecentTracks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrackServiceServer).GetMostRecentTracks(ctx, req.(*UserID))
	}
	return interceptor(ctx, in, info, handler)
}

// TrackService_ServiceDesc is the grpc.ServiceDesc for TrackService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TrackService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "track.TrackService",
	HandlerType: (*TrackServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAllTracks",
			Handler:    _TrackService_GetAllTracks_Handler,
		},
		{
			MethodName: "GetTrackByID",
			Handler:    _TrackService_GetTrackByID_Handler,
		},
		{
			MethodName: "CreateStream",
			Handler:    _TrackService_CreateStream_Handler,
		},
		{
			MethodName: "UpdateStreamDuration",
			Handler:    _TrackService_UpdateStreamDuration_Handler,
		},
		{
			MethodName: "GetLastListenedTracks",
			Handler:    _TrackService_GetLastListenedTracks_Handler,
		},
		{
			MethodName: "GetTracksByIDs",
			Handler:    _TrackService_GetTracksByIDs_Handler,
		},
		{
			MethodName: "GetTracksByIDsFiltered",
			Handler:    _TrackService_GetTracksByIDsFiltered_Handler,
		},
		{
			MethodName: "GetAlbumIDByTrackID",
			Handler:    _TrackService_GetAlbumIDByTrackID_Handler,
		},
		{
			MethodName: "GetTracksByAlbumID",
			Handler:    _TrackService_GetTracksByAlbumID_Handler,
		},
		{
			MethodName: "GetMinutesListenedByUserID",
			Handler:    _TrackService_GetMinutesListenedByUserID_Handler,
		},
		{
			MethodName: "GetTracksListenedByUserID",
			Handler:    _TrackService_GetTracksListenedByUserID_Handler,
		},
		{
			MethodName: "LikeTrack",
			Handler:    _TrackService_LikeTrack_Handler,
		},
		{
			MethodName: "SearchTracks",
			Handler:    _TrackService_SearchTracks_Handler,
		},
		{
			MethodName: "GetFavoriteTracks",
			Handler:    _TrackService_GetFavoriteTracks_Handler,
		},
		{
			MethodName: "AddTracksToAlbum",
			Handler:    _TrackService_AddTracksToAlbum_Handler,
		},
		{
			MethodName: "DeleteTracksByAlbumID",
			Handler:    _TrackService_DeleteTracksByAlbumID_Handler,
		},
		{
			MethodName: "GetMostLikedTracks",
			Handler:    _TrackService_GetMostLikedTracks_Handler,
		},
		{
			MethodName: "GetMostLikedLastWeekTracks",
			Handler:    _TrackService_GetMostLikedLastWeekTracks_Handler,
		},
		{
			MethodName: "GetMostListenedLastMonthTracks",
			Handler:    _TrackService_GetMostListenedLastMonthTracks_Handler,
		},
		{
			MethodName: "GetMostRecentTracks",
			Handler:    _TrackService_GetMostRecentTracks_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "track/track.proto",
}
