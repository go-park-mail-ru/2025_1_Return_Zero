package model

import (
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/delivery"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	"github.com/stretchr/testify/assert"
)

func TestPaginationConverters(t *testing.T) {
	t.Run("PaginationFromDeliveryToUsecase", func(t *testing.T) {
		deliveryPagination := &delivery.Pagination{
			Offset: 10,
			Limit:  20,
		}

		result := PaginationFromDeliveryToUsecase(deliveryPagination)

		assert.Equal(t, deliveryPagination.Offset, result.Offset)
		assert.Equal(t, deliveryPagination.Limit, result.Limit)
	})

	t.Run("PaginationFromUsecaseToRepository", func(t *testing.T) {
		usecasePagination := &usecase.Pagination{
			Offset: 10,
			Limit:  20,
		}

		result := PaginationFromUsecaseToRepository(usecasePagination)

		assert.Equal(t, usecasePagination.Offset, result.Offset)
		assert.Equal(t, usecasePagination.Limit, result.Limit)
	})
}

func TestAlbumConverters(t *testing.T) {
	t.Run("AlbumFromUsecaseToDelivery", func(t *testing.T) {
		releaseDate := time.Now()
		usecaseAlbum := &usecase.Album{
			ID:          1,
			Title:       "Test Album",
			Type:        usecase.AlbumTypeSingle,
			Thumbnail:   "thumbnail.jpg",
			ReleaseDate: releaseDate,
		}

		usecaseAlbumArtists := []*usecase.AlbumArtist{
			{
				ID:    1,
				Title: "Test Artist",
			},
		}

		result := AlbumFromUsecaseToDelivery(usecaseAlbum, usecaseAlbumArtists)

		assert.Equal(t, usecaseAlbum.ID, result.ID)
		assert.Equal(t, usecaseAlbum.Title, result.Title)
		assert.Equal(t, delivery.AlbumType(usecaseAlbum.Type), result.Type)
		assert.Equal(t, usecaseAlbum.Thumbnail, result.Thumbnail)
		assert.Equal(t, usecaseAlbum.ReleaseDate, result.ReleaseDate)
		assert.Len(t, result.Artists, 1)
		assert.Equal(t, usecaseAlbumArtists[0].ID, result.Artists[0].ID)
		assert.Equal(t, usecaseAlbumArtists[0].Title, result.Artists[0].Title)
	})

	t.Run("AlbumsFromUsecaseToDelivery", func(t *testing.T) {
		releaseDate := time.Now()
		usecaseAlbum := &usecase.Album{
			ID:          1,
			Title:       "Test Album",
			Type:        usecase.AlbumTypeSingle,
			Thumbnail:   "thumbnail.jpg",
			ReleaseDate: releaseDate,
			Artists: []*usecase.AlbumArtist{
				{
					ID:    1,
					Title: "Test Artist",
				},
			},
		}

		usecaseAlbums := []*usecase.Album{usecaseAlbum}

		result := AlbumsFromUsecaseToDelivery(usecaseAlbums)

		assert.Len(t, result, 1)
		assert.Equal(t, usecaseAlbum.ID, result[0].ID)
		assert.Equal(t, usecaseAlbum.Title, result[0].Title)
	})

	t.Run("AlbumFromRepositoryToUsecase", func(t *testing.T) {
		releaseDate := time.Now()
		repositoryAlbum := &repository.Album{
			ID:          1,
			Title:       "Test Album",
			Type:        repository.AlbumTypeSingle,
			Thumbnail:   "thumbnail.jpg",
			ReleaseDate: releaseDate,
		}

		repositoryAlbumArtists := []*repository.ArtistWithTitle{
			{
				ID:    1,
				Title: "Test Artist",
			},
		}

		result := AlbumFromRepositoryToUsecase(repositoryAlbum, repositoryAlbumArtists)

		assert.Equal(t, repositoryAlbum.ID, result.ID)
		assert.Equal(t, repositoryAlbum.Title, result.Title)
		assert.Equal(t, usecase.AlbumType(repositoryAlbum.Type), result.Type)
		assert.Equal(t, repositoryAlbum.Thumbnail, result.Thumbnail)
		assert.Equal(t, repositoryAlbum.ReleaseDate, result.ReleaseDate)
		assert.Len(t, result.Artists, 1)
		assert.Equal(t, repositoryAlbumArtists[0].ID, result.Artists[0].ID)
		assert.Equal(t, repositoryAlbumArtists[0].Title, result.Artists[0].Title)
	})
}

func TestArtistConverters(t *testing.T) {
	t.Run("ArtistFromUsecaseToDelivery", func(t *testing.T) {
		usecaseArtist := &usecase.Artist{
			ID:          1,
			Title:       "Test Artist",
			Thumbnail:   "thumbnail.jpg",
			Description: "Test Description",
		}

		result := ArtistFromUsecaseToDelivery(usecaseArtist)

		assert.Equal(t, usecaseArtist.ID, result.ID)
		assert.Equal(t, usecaseArtist.Title, result.Title)
		assert.Equal(t, usecaseArtist.Thumbnail, result.Thumbnail)
		assert.Equal(t, usecaseArtist.Description, result.Description)
	})

	t.Run("ArtistFromRepositoryToUsecase", func(t *testing.T) {
		repositoryArtist := &repository.Artist{
			ID:          1,
			Title:       "Test Artist",
			Thumbnail:   "thumbnail.jpg",
			Description: "Test Description",
		}

		result := ArtistFromRepositoryToUsecase(repositoryArtist)

		assert.Equal(t, repositoryArtist.ID, result.ID)
		assert.Equal(t, repositoryArtist.Title, result.Title)
		assert.Equal(t, repositoryArtist.Thumbnail, result.Thumbnail)
		assert.Equal(t, repositoryArtist.Description, result.Description)
	})

	t.Run("ArtistDetailedFromRepositoryToUsecase", func(t *testing.T) {
		repositoryArtist := &repository.Artist{
			ID:          1,
			Title:       "Test Artist",
			Thumbnail:   "thumbnail.jpg",
			Description: "Test Description",
		}

		repositoryArtistStats := &repository.ArtistStats{
			FavoritesCount: 100,
			ListenersCount: 200,
		}

		result := ArtistDetailedFromRepositoryToUsecase(repositoryArtist, repositoryArtistStats)

		assert.Equal(t, repositoryArtist.ID, result.Artist.ID)
		assert.Equal(t, repositoryArtist.Title, result.Artist.Title)
		assert.Equal(t, repositoryArtist.Thumbnail, result.Artist.Thumbnail)
		assert.Equal(t, repositoryArtist.Description, result.Artist.Description)
		assert.Equal(t, repositoryArtistStats.FavoritesCount, result.Favorites)
		assert.Equal(t, repositoryArtistStats.ListenersCount, result.Listeners)
	})

	t.Run("ArtistDetailedFromUsecaseToDelivery", func(t *testing.T) {
		usecaseArtist := &usecase.Artist{
			ID:          1,
			Title:       "Test Artist",
			Thumbnail:   "thumbnail.jpg",
			Description: "Test Description",
		}

		usecaseArtistDetailed := &usecase.ArtistDetailed{
			Artist:    *usecaseArtist,
			Favorites: 100,
			Listeners: 200,
		}

		result := ArtistDetailedFromUsecaseToDelivery(usecaseArtistDetailed)

		assert.Equal(t, usecaseArtist.ID, result.Artist.ID)
		assert.Equal(t, usecaseArtist.Title, result.Artist.Title)
		assert.Equal(t, usecaseArtist.Thumbnail, result.Artist.Thumbnail)
		assert.Equal(t, usecaseArtist.Description, result.Artist.Description)
		assert.Equal(t, usecaseArtistDetailed.Favorites, result.Favorites)
		assert.Equal(t, usecaseArtistDetailed.Listeners, result.Listeners)
	})
}

func TestTrackConverters(t *testing.T) {
	t.Run("TrackFromUsecaseToDelivery", func(t *testing.T) {
		usecaseTrack := &usecase.Track{
			ID:        1,
			Title:     "Test Track",
			Thumbnail: "thumbnail.jpg",
			Duration:  200,
			Album:     "Test Album",
			AlbumID:   1,
			Artists: []*usecase.TrackArtist{
				{
					ID:    1,
					Title: "Test Artist",
					Role:  "main",
				},
			},
		}

		result := TrackFromUsecaseToDelivery(usecaseTrack)

		assert.Equal(t, usecaseTrack.ID, result.ID)
		assert.Equal(t, usecaseTrack.Title, result.Title)
		assert.Equal(t, usecaseTrack.Thumbnail, result.Thumbnail)
		assert.Equal(t, usecaseTrack.Duration, result.Duration)
		assert.Equal(t, usecaseTrack.Album, result.Album)
		assert.Equal(t, usecaseTrack.AlbumID, result.AlbumID)
		assert.Len(t, result.Artists, 1)
		assert.Equal(t, usecaseTrack.Artists[0].ID, result.Artists[0].ID)
		assert.Equal(t, usecaseTrack.Artists[0].Title, result.Artists[0].Title)
		assert.Equal(t, usecaseTrack.Artists[0].Role, result.Artists[0].Role)
	})

	t.Run("TrackFromRepositoryToUsecase", func(t *testing.T) {
		repositoryTrack := &repository.Track{
			ID:        1,
			Title:     "Test Track",
			Thumbnail: "thumbnail.jpg",
			Duration:  200,
			AlbumID:   1,
		}

		repositoryTrackArtists := []*repository.ArtistWithRole{
			{
				ID:    1,
				Title: "Test Artist",
				Role:  "main",
			},
		}

		albumTitle := "Test Album"

		result := TrackFromRepositoryToUsecase(repositoryTrack, repositoryTrackArtists, albumTitle)

		assert.Equal(t, repositoryTrack.ID, result.ID)
		assert.Equal(t, repositoryTrack.Title, result.Title)
		assert.Equal(t, repositoryTrack.Thumbnail, result.Thumbnail)
		assert.Equal(t, repositoryTrack.Duration, result.Duration)
		assert.Equal(t, albumTitle, result.Album)
		assert.Equal(t, repositoryTrack.AlbumID, result.AlbumID)
		assert.Len(t, result.Artists, 1)
		assert.Equal(t, repositoryTrackArtists[0].ID, result.Artists[0].ID)
		assert.Equal(t, repositoryTrackArtists[0].Title, result.Artists[0].Title)
		assert.Equal(t, repositoryTrackArtists[0].Role, result.Artists[0].Role)
	})

	t.Run("TrackWithFileKeyFromRepositoryToUsecase", func(t *testing.T) {
		repositoryTrack := &repository.TrackWithFileKey{
			Track: repository.Track{
				ID:        1,
				Title:     "Test Track",
				Thumbnail: "thumbnail.jpg",
				Duration:  200,
				AlbumID:   1,
			},
			FileKey: "file_key",
		}

		repositoryTrackArtists := []*repository.ArtistWithRole{
			{
				ID:    1,
				Title: "Test Artist",
				Role:  "main",
			},
		}

		albumTitle := "Test Album"

		result := TrackWithFileKeyFromRepositoryToUsecase(repositoryTrack, repositoryTrackArtists, albumTitle)

		assert.Equal(t, repositoryTrack.ID, result.ID)
		assert.Equal(t, repositoryTrack.Title, result.Title)
		assert.Equal(t, repositoryTrack.Thumbnail, result.Thumbnail)
		assert.Equal(t, repositoryTrack.Duration, result.Duration)
		assert.Equal(t, albumTitle, result.Album)
		assert.Equal(t, repositoryTrack.AlbumID, result.AlbumID)
		assert.Len(t, result.Artists, 1)
		assert.Equal(t, repositoryTrackArtists[0].ID, result.Artists[0].ID)
		assert.Equal(t, repositoryTrackArtists[0].Title, result.Artists[0].Title)
		assert.Equal(t, repositoryTrackArtists[0].Role, result.Artists[0].Role)
	})

	t.Run("TrackDetailedFromRepositoryToUsecase", func(t *testing.T) {
		repositoryTrack := &repository.TrackWithFileKey{
			Track: repository.Track{
				ID:        1,
				Title:     "Test Track",
				Thumbnail: "thumbnail.jpg",
				Duration:  200,
				AlbumID:   1,
			},
			FileKey: "file_key",
		}

		repositoryTrackArtists := []*repository.ArtistWithRole{
			{
				ID:    1,
				Title: "Test Artist",
				Role:  "main",
			},
		}

		albumTitle := "Test Album"
		fileUrl := "http://example.com/file.mp3"

		result := TrackDetailedFromRepositoryToUsecase(repositoryTrack, repositoryTrackArtists, albumTitle, fileUrl)

		assert.Equal(t, repositoryTrack.ID, result.Track.ID)
		assert.Equal(t, repositoryTrack.Title, result.Track.Title)
		assert.Equal(t, repositoryTrack.Thumbnail, result.Track.Thumbnail)
		assert.Equal(t, repositoryTrack.Duration, result.Track.Duration)
		assert.Equal(t, albumTitle, result.Track.Album)
		assert.Equal(t, repositoryTrack.AlbumID, result.Track.AlbumID)
		assert.Equal(t, fileUrl, result.FileUrl)
		assert.Len(t, result.Track.Artists, 1)
		assert.Equal(t, repositoryTrackArtists[0].ID, result.Track.Artists[0].ID)
		assert.Equal(t, repositoryTrackArtists[0].Title, result.Track.Artists[0].Title)
		assert.Equal(t, repositoryTrackArtists[0].Role, result.Track.Artists[0].Role)
	})
}

func TestUserConverters(t *testing.T) {
	t.Run("PrivacyRepositoryToUsecase", func(t *testing.T) {
		repositoryPrivacy := &repository.UserPrivacySettings{
			IsPublicPlaylists:       true,
			IsPublicMinutesListened: false,
			IsPublicFavoriteArtists: true,
			IsPublicTracksListened:  false,
			IsPublicFavoriteTracks:  true,
			IsPublicArtistsListened: false,
		}

		result := PrivacyRepositoryToUsecase(repositoryPrivacy)

		assert.Equal(t, repositoryPrivacy.IsPublicPlaylists, result.IsPublicPlaylists)
		assert.Equal(t, repositoryPrivacy.IsPublicMinutesListened, result.IsPublicMinutesListened)
		assert.Equal(t, repositoryPrivacy.IsPublicFavoriteArtists, result.IsPublicFavoriteArtists)
		assert.Equal(t, repositoryPrivacy.IsPublicTracksListened, result.IsPublicTracksListened)
		assert.Equal(t, repositoryPrivacy.IsPublicFavoriteTracks, result.IsPublicFavoriteTracks)
		assert.Equal(t, repositoryPrivacy.IsPublicArtistsListened, result.IsPublicArtistsListened)
	})

	t.Run("StatisticsRepositoryToUsecase", func(t *testing.T) {
		repositoryStatistics := &repository.UserStats{
			MinutesListened: 100,
			TracksListened:  50,
			ArtistsListened: 20,
		}

		result := StatisticsRepositoryToUsecase(repositoryStatistics)

		assert.Equal(t, repositoryStatistics.MinutesListened, result.MinutesListened)
		assert.Equal(t, repositoryStatistics.TracksListened, result.TracksListened)
		assert.Equal(t, repositoryStatistics.ArtistsListened, result.ArtistsListened)
	})

	t.Run("UserFullDataRepositoryToUsecase", func(t *testing.T) {
		repositoryPrivacy := &repository.UserPrivacySettings{
			IsPublicPlaylists:       true,
			IsPublicMinutesListened: false,
			IsPublicFavoriteArtists: true,
			IsPublicTracksListened:  false,
			IsPublicFavoriteTracks:  true,
			IsPublicArtistsListened: false,
		}

		repositoryStatistics := &repository.UserStats{
			MinutesListened: 100,
			TracksListened:  50,
			ArtistsListened: 20,
		}

		repositoryUser := &repository.UserFullData{
			Username:   "testuser",
			Email:      "test@example.com",
			Privacy:    repositoryPrivacy,
			Statistics: repositoryStatistics,
		}

		result := UserFullDataRepositoryToUsecase(repositoryUser)

		assert.Equal(t, repositoryUser.Username, result.Username)
		assert.Equal(t, repositoryUser.Email, result.Email)
		assert.Equal(t, repositoryPrivacy.IsPublicPlaylists, result.Privacy.IsPublicPlaylists)
		assert.Equal(t, repositoryStatistics.MinutesListened, result.Statistics.MinutesListened)
	})

	t.Run("UserFullDataUsecaseToDelivery", func(t *testing.T) {
		usecasePrivacy := &usecase.UserPrivacy{
			IsPublicPlaylists:       true,
			IsPublicMinutesListened: false,
			IsPublicFavoriteArtists: true,
			IsPublicTracksListened:  false,
			IsPublicFavoriteTracks:  true,
			IsPublicArtistsListened: false,
		}

		usecaseStatistics := &usecase.UserStatistics{
			MinutesListened: 100,
			TracksListened:  50,
			ArtistsListened: 20,
		}

		usecaseUser := &usecase.UserFullData{
			Username:   "testuser",
			Email:      "test@example.com",
			AvatarUrl:  "avatar.jpg",
			Privacy:    usecasePrivacy,
			Statistics: usecaseStatistics,
		}

		result := UserFullDataUsecaseToDelivery(usecaseUser)

		assert.Equal(t, usecaseUser.Username, result.Username)
		assert.Equal(t, usecaseUser.Email, result.Email)
		assert.Equal(t, usecaseUser.AvatarUrl, result.AvatarUrl)
		assert.Equal(t, usecasePrivacy.IsPublicPlaylists, result.Privacy.IsPublicPlaylists)
		assert.Equal(t, usecaseStatistics.MinutesListened, result.Statistics.MinutesListened)
	})
}

func TestStreamConverters(t *testing.T) {
	t.Run("TrackStreamCreateDataFromDeliveryToUsecase", func(t *testing.T) {
		deliveryTrackStream := &delivery.TrackStreamCreateData{
			TrackID: 1,
			UserID:  2,
		}

		result := TrackStreamCreateDataFromDeliveryToUsecase(deliveryTrackStream)

		assert.Equal(t, deliveryTrackStream.TrackID, result.TrackID)
		assert.Equal(t, deliveryTrackStream.UserID, result.UserID)
	})

	t.Run("TrackStreamCreateDataFromUsecaseToRepository", func(t *testing.T) {
		usecaseTrackStream := &usecase.TrackStreamCreateData{
			TrackID: 1,
			UserID:  2,
		}

		result := TrackStreamCreateDataFromUsecaseToRepository(usecaseTrackStream)

		assert.Equal(t, usecaseTrackStream.TrackID, result.TrackID)
		assert.Equal(t, usecaseTrackStream.UserID, result.UserID)
	})

	t.Run("TrackStreamUpdateDataFromUsecaseToRepository", func(t *testing.T) {
		usecaseTrackStream := &usecase.TrackStreamUpdateData{
			StreamID: 1,
			Duration: 200,
			UserID:   2,
		}

		result := TrackStreamUpdateDataFromUsecaseToRepository(usecaseTrackStream)

		assert.Equal(t, usecaseTrackStream.StreamID, result.StreamID)
		assert.Equal(t, usecaseTrackStream.Duration, result.Duration)
	})

	t.Run("TrackStreamUpdateDataFromDeliveryToUsecase", func(t *testing.T) {
		deliveryTrackStream := &delivery.TrackStreamUpdateData{
			Duration: 200,
		}

		userID := int64(2)
		streamID := int64(1)

		result := TrackStreamUpdateDataFromDeliveryToUsecase(deliveryTrackStream, userID, streamID)

		assert.Equal(t, streamID, result.StreamID)
		assert.Equal(t, deliveryTrackStream.Duration, result.Duration)
		assert.Equal(t, userID, result.UserID)
	})
}
