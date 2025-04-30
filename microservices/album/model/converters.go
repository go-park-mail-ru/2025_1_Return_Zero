package model

import (
	albumProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/album"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/model/repository"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/album/model/usecase"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func PaginationFromUsecaseToRepository(pagination *usecaseModel.Pagination) *repoModel.Pagination {
	return &repoModel.Pagination{
		Offset: pagination.Offset,
		Limit:  pagination.Limit,
	}
}

func FiltersFromUsecaseToRepository(filters *usecaseModel.AlbumFilters) *repoModel.AlbumFilters {
	return &repoModel.AlbumFilters{
		Pagination: PaginationFromUsecaseToRepository(filters.Pagination),
	}
}

func PaginationFromProtoToUsecase(pagination *albumProto.Pagination) *usecaseModel.Pagination {
	return &usecaseModel.Pagination{
		Offset: pagination.Offset,
		Limit:  pagination.Limit,
	}
}

func AlbumFiltersFromProtoToUsecase(filters *albumProto.Filters) *usecaseModel.AlbumFilters {
	return &usecaseModel.AlbumFilters{
		Pagination: PaginationFromProtoToUsecase(filters.Pagination),
	}
}

func AlbumTypeFromRepositoryToUsecase(albumType repoModel.AlbumType) usecaseModel.AlbumType {
	return usecaseModel.AlbumType(albumType)
}

func AlbumTypeFromUsecaseToProto(albumType usecaseModel.AlbumType) albumProto.AlbumType {
	switch albumType {
	case usecaseModel.AlbumTypeAlbum:
		return albumProto.AlbumType_AlbumTypeAlbum
	case usecaseModel.AlbumTypeEP:
		return albumProto.AlbumType_AlbumTypeEP
	case usecaseModel.AlbumTypeSingle:
		return albumProto.AlbumType_AlbumTypeSingle
	case usecaseModel.AlbumTypeCompilation:
		return albumProto.AlbumType_AlbumTypeCompilation
	}
	return albumProto.AlbumType_AlbumTypeAlbum
}

func AlbumFromRepositoryToUsecase(album *repoModel.Album) *usecaseModel.Album {
	return &usecaseModel.Album{
		ID:          album.ID,
		Title:       album.Title,
		Type:        AlbumTypeFromRepositoryToUsecase(album.Type),
		Thumbnail:   album.Thumbnail,
		ReleaseDate: album.ReleaseDate,
	}
}

func AlbumFromUsecaseToProto(album *usecaseModel.Album) *albumProto.Album {
	return &albumProto.Album{
		Id:          album.ID,
		Title:       album.Title,
		Type:        AlbumTypeFromUsecaseToProto(album.Type),
		Thumbnail:   album.Thumbnail,
		ReleaseDate: timestamppb.New(album.ReleaseDate),
	}
}

func AlbumListFromRepositoryToUsecase(albums []*repoModel.Album) []*usecaseModel.Album {
	albumList := make([]*usecaseModel.Album, len(albums))
	for i, album := range albums {
		albumList[i] = AlbumFromRepositoryToUsecase(album)
	}
	return albumList
}

func AlbumListFromUsecaseToProto(albums []*usecaseModel.Album) *albumProto.AlbumList {
	albumList := make([]*albumProto.Album, len(albums))
	for i, album := range albums {
		albumList[i] = AlbumFromUsecaseToProto(album)
	}
	return &albumProto.AlbumList{
		Albums: albumList,
	}
}

func AlbumTitleFromUsecaseToProto(albumTitle *usecaseModel.AlbumTitle) *albumProto.AlbumTitle {
	return &albumProto.AlbumTitle{
		Title: albumTitle.Title,
	}
}

func AlbumTitleMapFromRepositoryToUsecase(albumTitles map[int64]string) *usecaseModel.AlbumTitleMap {
	albumTitleMap := make(map[int64]*usecaseModel.AlbumTitle)
	for key, albumTitle := range albumTitles {
		albumTitleMap[key] = &usecaseModel.AlbumTitle{
			Title: albumTitle,
		}
	}

	return &usecaseModel.AlbumTitleMap{
		Titles: albumTitleMap,
	}
}

func AlbumTitleMapFromUsecaseToProto(albumTitles *usecaseModel.AlbumTitleMap) *albumProto.AlbumTitleMap {
	albumTitleMap := make(map[int64]*albumProto.AlbumTitle)
	for key, albumTitle := range albumTitles.Titles {
		albumTitleMap[key] = AlbumTitleFromUsecaseToProto(albumTitle)
	}
	return &albumProto.AlbumTitleMap{
		Titles: albumTitleMap,
	}
}

func AlbumStreamCreateDataFromUsecaseToProto(albumStreamCreateData *usecaseModel.AlbumStreamCreateData) *albumProto.AlbumStreamCreateData {
	return &albumProto.AlbumStreamCreateData{
		AlbumId: &albumProto.AlbumID{Id: albumStreamCreateData.AlbumID},
		UserId:  &albumProto.UserID{Id: albumStreamCreateData.UserID},
	}
}

func AlbumStreamCreateDataFromProtoToUsecase(albumStreamCreateData *albumProto.AlbumStreamCreateData) *usecaseModel.AlbumStreamCreateData {
	return &usecaseModel.AlbumStreamCreateData{
		AlbumID: albumStreamCreateData.AlbumId.Id,
		UserID:  albumStreamCreateData.UserId.Id,
	}
}
