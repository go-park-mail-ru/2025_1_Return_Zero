package repository

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/crypto/argon2"

	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/user"
)

var (
	ErrUsernameExist = errors.New("user with this username already exists")
	ErrEmailExist    = errors.New("user with this email already exists")
	ErrUserNotFound  = errors.New("user not found")
	ErrCreateSalt    = errors.New("failed to create salt")
	ErrWrongPassword = errors.New("wrong password")
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
	uploadAvatarQuery = `
			UPDATE "user"
			SET thumbnail_url = $1
			WHERE username = $2
			`
	getAvatarQuery = `
			SELECT thumbnail_url
			FROM "user"
			WHERE username = $1
			`
	changeUsernameQuery = `
			UPDATE "user"
			SET username = $1
			WHERE username = $2
			`
	changeEmailQuery = `
			UPDATE "user"
			SET email = $1
			WHERE username = $2
			`
	changePasswordQuery = `
			UPDATE "user"
			SET password_hash = $1
			WHERE username = $2
			`
	getPasswordQuery = `
			SELECT password_hash
			FROM "user"
			WHERE username = $1
	`
	deleteUserQuery = `
			DELETE FROM "user"
			WHERE username = $1 AND email = $2
	`
)

func hashPassword(salt []byte, password string) string {
	hashedPass := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	combined := append(salt, hashedPass...)
	return base64.StdEncoding.EncodeToString(combined)
}

func checkPasswordHash(encodedHash string, password string) bool {
	decodedHash, err := base64.StdEncoding.DecodeString(encodedHash)
	if err != nil {
		return false
	}
	salt := decodedHash[:8]
	userPassHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return bytes.Equal(userPassHash, decodedHash[8:])
}

func (r *userPostgresRepository) getPassword(username string) (string, error) {
	row := r.db.QueryRow(getPasswordQuery, username)
	var storedHash string
	err := row.Scan(&storedHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrUserNotFound
		}
		return "", err
	}
	return storedHash, nil
}

func NewUserPostgresRepository(db *sql.DB) user.Repository {
	repo := &userPostgresRepository{
		db: db,
	}

	return repo
}

func createSalt() []byte {
	salt := make([]byte, 8)
	_, err := rand.Read(salt)
	if err != nil {
		return nil
	}
	return salt
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

	salt := createSalt()
	if salt == nil {
		return nil, errors.New("failed to create salt")
	}
	hashedPassword := hashPassword(salt, regData.Password)

	var userID int64
	err = r.db.QueryRow(createUserQuery, regData.Username,
		hashedPassword, regData.Email).Scan(&userID)
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

	if !checkPasswordHash(storedHash, logData.Password) {
		return nil, ErrUserNotFound
	}

	return &user, nil
}

func (r *userPostgresRepository) GetAvatar(username string) (string, error) {
	row := r.db.QueryRow(getAvatarQuery, username)
	var avatarUrl string
	err := row.Scan(&avatarUrl)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrUserNotFound
		}
		return "", err
	}
	return avatarUrl, nil
}

func (r *userPostgresRepository) UploadAvatar(avatarUrl string, username string) error {
	_, err := r.db.Exec(uploadAvatarQuery, avatarUrl, username)
	if err != nil {
		return err
	}
	return nil
}

func (r *userPostgresRepository) changeUsername(username string, newUsername string) error {
	_, err := r.db.Exec(changeUsernameQuery, newUsername, username)
	if err != nil {
		return err
	}
	return nil
}

func (r *userPostgresRepository) changeEmail(username string, newEmail string) error {
	_, err := r.db.Exec(changeEmailQuery, newEmail, username)
	if err != nil {
		return err
	}
	return nil
}

func (r *userPostgresRepository) changePassword(password string, username string, newPassword string) error {
	storedHash, err := r.getPassword(username)
	if err != nil {
		return err
	}
	if !checkPasswordHash(storedHash, password) {
		return ErrWrongPassword
	}
	salt := createSalt()
	newHashedPassword := hashPassword(salt, newPassword)
	_, err = r.db.Exec(changePasswordQuery, newHashedPassword, username)
	if err != nil {
		return err
	}
	return nil
}

func (r *userPostgresRepository) ChangeUserData(changeData *repoModel.ChangeUserData) (*repoModel.User, error) {
	newUser := &repoModel.User{
		Username: changeData.Username,
		Email:    changeData.Email,
	}
	if changeData.NewUsername != "" {
		err := r.changeUsername(changeData.Username, changeData.NewUsername)
		if err != nil {
			return nil, err
		}
		newUser.Username = changeData.NewUsername
	}
	if changeData.NewEmail != "" {
		err := r.changeEmail(changeData.Username, changeData.NewEmail)
		if err != nil {
			return nil, err
		}
		newUser.Email = changeData.NewEmail
	}
	if changeData.NewPassword != "" {
		err := r.changePassword(changeData.Password, changeData.Username, changeData.NewPassword)
		if err != nil {
			return nil, err
		}
	}
	fmt.Println(newUser)
	return newUser, nil
}

func (r *userPostgresRepository) DeleteUser(user *repoModel.User) error {
	storedHash, err := r.getPassword(user.Username)
	if err != nil {
		return err
	}
	if !checkPasswordHash(storedHash, user.Password) {
		return ErrWrongPassword
	}
	_, err = r.db.Exec(deleteUserQuery, user.Username, user.Email)
	if err != nil {
		return err
	}
	return nil
}
