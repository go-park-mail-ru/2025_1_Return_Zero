package auth

type Repository interface {
	CreateSession(ID int64) string
	DeleteSession(SID string)
	GetSession(SID string) (int64, error)
}