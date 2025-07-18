package model

import (
	protoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/artist"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/artist/model/repository"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/artist/model/usecase"
)

func ArtistFromRepositoryToUsecase(artist *repoModel.Artist) *usecaseModel.Artist {
	return &usecaseModel.Artist{
		ID:          artist.ID,
		Title:       artist.Title,
		Description: artist.Description,
		Thumbnail:   artist.Thumbnail,
		IsFavorite:  artist.IsFavorite,
	}
}

func ArtistFromUsecaseToProto(artist *usecaseModel.Artist) *protoModel.Artist {
	return &protoModel.Artist{
		Id:          artist.ID,
		Title:       artist.Title,
		Description: artist.Description,
		Thumbnail:   artist.Thumbnail,
		IsFavorite:  artist.IsFavorite,
	}
}

func ArtistDetailedFromRepositoryToUsecase(artist *repoModel.Artist, stats *repoModel.ArtistStats) *usecaseModel.ArtistDetailed {
	return &usecaseModel.ArtistDetailed{
		Artist:         ArtistFromRepositoryToUsecase(artist),
		ListenersCount: stats.ListenersCount,
		FavoritesCount: stats.FavoritesCount,
	}
}

func ArtistDetailedFromUsecaseToProto(artist *usecaseModel.ArtistDetailed) *protoModel.ArtistDetailed {
	return &protoModel.ArtistDetailed{
		Artist:         ArtistFromUsecaseToProto(artist.Artist),
		ListenersCount: artist.ListenersCount,
		FavoritesCount: artist.FavoritesCount,
	}
}

func ArtistListFromRepositoryToUsecase(artists []*repoModel.Artist) *usecaseModel.ArtistList {
	usecaseArtists := make([]*usecaseModel.Artist, len(artists))
	for i, artist := range artists {
		usecaseArtists[i] = ArtistFromRepositoryToUsecase(artist)
	}
	return &usecaseModel.ArtistList{
		Artists: usecaseArtists,
	}
}

func ArtistListFromUsecaseToProto(artists *usecaseModel.ArtistList) *protoModel.ArtistList {
	protoArtists := make([]*protoModel.Artist, len(artists.Artists))
	for i, artist := range artists.Artists {
		protoArtists[i] = ArtistFromUsecaseToProto(artist)
	}
	return &protoModel.ArtistList{
		Artists: protoArtists,
	}
}

func ArtistWithTitleFromRepositoryToUsecase(artist *repoModel.ArtistWithTitle) *usecaseModel.ArtistWithTitle {
	return &usecaseModel.ArtistWithTitle{
		ID:    artist.ID,
		Title: artist.Title,
	}
}

func ArtistWithTitleFromUsecaseToProto(artist *usecaseModel.ArtistWithTitle) *protoModel.ArtistWithTitle {
	return &protoModel.ArtistWithTitle{
		Id:    artist.ID,
		Title: artist.Title,
	}
}

func ArtistWithTitleListFromUsecaseToProto(artists *usecaseModel.ArtistWithTitleList) *protoModel.ArtistWithTitleList {
	protoArtists := make([]*protoModel.ArtistWithTitle, len(artists.Artists))
	for i, artist := range artists.Artists {
		protoArtists[i] = ArtistWithTitleFromUsecaseToProto(artist)
	}
	return &protoModel.ArtistWithTitleList{
		Artists: protoArtists,
	}
}

func ArtistWithTitleListFromRepositoryToUsecase(artists []*repoModel.ArtistWithTitle) *usecaseModel.ArtistWithTitleList {
	usecaseArtists := make([]*usecaseModel.ArtistWithTitle, len(artists))
	for i, artist := range artists {
		usecaseArtists[i] = ArtistWithTitleFromRepositoryToUsecase(artist)
	}
	return &usecaseModel.ArtistWithTitleList{
		Artists: usecaseArtists,
	}
}

func ArtistWithTitleMapFromRepositoryToUsecase(artists map[int64][]*repoModel.ArtistWithTitle) *usecaseModel.ArtistWithTitleMap {
	usecaseArtists := make(map[int64]*usecaseModel.ArtistWithTitleList)
	for id, artist := range artists {
		usecaseArtists[id] = ArtistWithTitleListFromRepositoryToUsecase(artist)
	}
	return &usecaseModel.ArtistWithTitleMap{
		Artists: usecaseArtists,
	}
}

func ArtistWithTitleMapFromUsecaseToProto(artists *usecaseModel.ArtistWithTitleMap) *protoModel.ArtistWithTitleMap {
	protoArtists := make(map[int64]*protoModel.ArtistWithTitleList)
	for id, artist := range artists.Artists {
		protoArtists[id] = ArtistWithTitleListFromUsecaseToProto(artist)
	}
	return &protoModel.ArtistWithTitleMap{
		Artists: protoArtists,
	}
}

func ArtistWithRoleFromRepositoryToUsecase(artist *repoModel.ArtistWithRole) *usecaseModel.ArtistWithRole {
	return &usecaseModel.ArtistWithRole{
		ID:    artist.ID,
		Title: artist.Title,
		Role:  artist.Role,
	}
}

func ArtistWithRoleFromUsecaseToProto(artist *usecaseModel.ArtistWithRole) *protoModel.ArtistWithRole {
	return &protoModel.ArtistWithRole{
		Id:    artist.ID,
		Title: artist.Title,
		Role:  artist.Role,
	}
}

func ArtistWithRoleListFromUsecaseToProto(artists *usecaseModel.ArtistWithRoleList) *protoModel.ArtistWithRoleList {
	protoArtists := make([]*protoModel.ArtistWithRole, len(artists.Artists))
	for i, artist := range artists.Artists {
		protoArtists[i] = ArtistWithRoleFromUsecaseToProto(artist)
	}
	return &protoModel.ArtistWithRoleList{
		Artists: protoArtists,
	}
}

func ArtistWithRoleMapFromRepositoryToUsecase(artists map[int64][]*repoModel.ArtistWithRole) *usecaseModel.ArtistWithRoleMap {
	usecaseArtists := make(map[int64]*usecaseModel.ArtistWithRoleList)
	for id, artist := range artists {
		usecaseArtists[id] = ArtistWithRoleListFromRepositoryToUsecase(artist)
	}
	return &usecaseModel.ArtistWithRoleMap{
		Artists: usecaseArtists,
	}
}

func ArtistWithRoleMapFromUsecaseToProto(artists *usecaseModel.ArtistWithRoleMap) *protoModel.ArtistWithRoleMap {
	protoArtists := make(map[int64]*protoModel.ArtistWithRoleList)
	for id, artist := range artists.Artists {
		protoArtists[id] = ArtistWithRoleListFromUsecaseToProto(artist)
	}
	return &protoModel.ArtistWithRoleMap{
		Artists: protoArtists,
	}
}
func ArtistWithRoleListFromRepositoryToUsecase(artists []*repoModel.ArtistWithRole) *usecaseModel.ArtistWithRoleList {
	usecaseArtists := make([]*usecaseModel.ArtistWithRole, len(artists))
	for i, artist := range artists {
		usecaseArtists[i] = ArtistWithRoleFromRepositoryToUsecase(artist)
	}
	return &usecaseModel.ArtistWithRoleList{
		Artists: usecaseArtists,
	}
}

func PaginationFromUsecaseToRepository(pagination *usecaseModel.Pagination) *repoModel.Pagination {
	return &repoModel.Pagination{
		Offset: pagination.Offset,
		Limit:  pagination.Limit,
	}
}

func ArtistFiltersFromUsecaseToRepository(filters *usecaseModel.Filters) *repoModel.Filters {
	return &repoModel.Filters{
		Pagination: PaginationFromUsecaseToRepository(filters.Pagination),
	}
}

func TrackIDListFromProtoToUsecase(ids []*protoModel.TrackID) []int64 {
	trackIDs := make([]int64, len(ids))
	for i, id := range ids {
		trackIDs[i] = id.Id
	}
	return trackIDs
}

func AlbumIDListFromProtoToUsecase(ids []*protoModel.AlbumID) []int64 {
	albumIDs := make([]int64, len(ids))
	for i, id := range ids {
		albumIDs[i] = id.Id
	}
	return albumIDs
}

func PaginationFromProtoToUsecase(pagination *protoModel.Pagination) *usecaseModel.Pagination {
	return &usecaseModel.Pagination{
		Offset: pagination.Offset,
		Limit:  pagination.Limit,
	}
}

func ArtistFiltersFromProtoToUsecase(filters *protoModel.Filters) *usecaseModel.Filters {
	return &usecaseModel.Filters{
		Pagination: PaginationFromProtoToUsecase(filters.Pagination),
	}
}

func ArtistStreamCreateDataFromProtoToUsecase(data *protoModel.ArtistStreamCreateDataList) *usecaseModel.ArtistStreamCreateDataList {
	artistIDs := make([]int64, len(data.ArtistIds.Ids))
	for i, id := range data.ArtistIds.Ids {
		artistIDs[i] = id.Id
	}
	return &usecaseModel.ArtistStreamCreateDataList{
		ArtistIDs: artistIDs,
		UserID:    data.UserId.Id,
	}
}

func ArtistStreamCreateDataFromUsecaseToRepository(data *usecaseModel.ArtistStreamCreateDataList) *repoModel.ArtistStreamCreateDataList {
	return &repoModel.ArtistStreamCreateDataList{
		ArtistIDs: data.ArtistIDs,
		UserID:    data.UserID,
	}
}

func LikeRequestFromProtoToUsecase(request *protoModel.LikeRequest) *usecaseModel.LikeRequest {
	return &usecaseModel.LikeRequest{
		ArtistID: request.ArtistId.Id,
		UserID:   request.UserId.Id,
		IsLike:   request.IsLike,
	}
}

func LikeRequestFromUsecaseToRepository(request *usecaseModel.LikeRequest) *repoModel.LikeRequest {
	return &repoModel.LikeRequest{
		ArtistID: request.ArtistID,
		UserID:   request.UserID,
	}
}

func ArtistLoadFromUsecaseToRepository(artist *usecaseModel.ArtistLoad) *repoModel.Artist {
	return &repoModel.Artist{
		Title:   artist.Title,
		LabelID: artist.LabelID,
	}
}

func ArtistLoadFromProtoToUsecase(artist *protoModel.ArtistLoad) *usecaseModel.ArtistLoad {
	return &usecaseModel.ArtistLoad{
		Title:   artist.Title,
		Image:   artist.Image,
		LabelID: artist.LabelId,
	}
}

func ArtistEditFromProtoToUsecase(artist *protoModel.ArtistEdit) *usecaseModel.ArtistEdit {
	return &usecaseModel.ArtistEdit{
		ArtistID: artist.ArtistId,
		NewTitle: artist.NewTitle,
		Image:    artist.Image,
		LabelID:  artist.LabelId,
	}
}
func ArtistEditFromUsecaseToRepository(artist *usecaseModel.ArtistEdit) *repoModel.ArtistEdit {
	return &repoModel.ArtistEdit{
		ArtistID: artist.ArtistID,
		NewTitle: artist.NewTitle,
		Image:    artist.Image,
		LabelID:  artist.LabelID,
	}
}

func ArtistDeleteFromProtoToUsecase(artist *protoModel.ArtistDelete) *usecaseModel.ArtistDelete {
	return &usecaseModel.ArtistDelete{
		ArtistID: artist.ArtistId,
		LabelID:  artist.LabelId,
	}
}

func ArtistsIDWithAlbumIDFromProtoToUsecase(artists *protoModel.ArtistsIDWithAlbumID) ([]int64, int64, []int64) {
	var artistIDs []int64
	for _, id := range artists.ArtistIds.Ids {
		artistIDs = append(artistIDs, id.Id)
	}
	var tracksIDs []int64
	for _, id := range artists.TrackIds.Ids {
		tracksIDs = append(tracksIDs, id.Id)
	}
	return artistIDs, artists.AlbumId.Id, tracksIDs
}
