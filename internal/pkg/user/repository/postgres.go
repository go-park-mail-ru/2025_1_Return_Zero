package repository

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"

	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/user"
)

var (
	ErrUsernameExist = errors.New("user with this username already exists")
	ErrEmailExist    = errors.New("user with this email already exists")
	ErrUserNotFound  = errors.New("user not found")
)

type userPostgresRepository struct {
	db *sql.DB
}

const (
	checkUserExist = `
			SELECT 1 
			FROM "user"
			WHERE username = $1 OR email = $2
			`
	createUserQuery = `
			INSERT INTO "user" (username, password_hash, email) 
			VALUES ($1, $2, $3) 
			RETURNING id
			`
	getUserByIDQuery = `
			SELECT id, username, email
			FROM "user"
			WHERE id = $1
			`
	loginUserQuery = `
			SELECT id, username, email, password_hash
			FROM "user"
			WHERE username = $1 OR email = $2
			`
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func NewUserPostgresRepository(db *sql.DB) user.Repository {
	repo := &userPostgresRepository{
		db: db,
	}

	return repo
}

func (r *userPostgresRepository) CreateUser(regData *repoModel.User) (*repoModel.User, error) {
	var exists bool
	err := r.db.QueryRow(checkUserExist, regData.Username, regData.Email).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if exists {
		return nil, ErrUsernameExist
	}

	hashedPassword, err := HashPassword(regData.Password)
	if err != nil {
		return nil, err
	}

	var userID int64
	err = r.db.QueryRow(createUserQuery, regData.Username, hashedPassword, regData.Email).Scan(&userID)
	if err != nil {
		return nil, err
	}

	return &repoModel.User{
		ID:       userID,
		Username: regData.Username,
		Email:    regData.Email,
	}, nil
}

func (r *userPostgresRepository) GetUserByID(ID int64) (*repoModel.User, error) {
	row := r.db.QueryRow(getUserByIDQuery, ID)
	var user repoModel.User
	err := row.Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userPostgresRepository) LoginUser(logData *repoModel.User) (*repoModel.User, error) {
	var storedHash string
	row := r.db.QueryRow(loginUserQuery, logData.Username, logData.Email)
	var user repoModel.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &storedHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if !CheckPasswordHash(logData.Password, storedHash) {
		return nil, ErrUserNotFound
	}

	return &user, nil
}

