package trackFile

type Repository interface {
	GetPresignedURL(trackKey string) (string, error)
}
