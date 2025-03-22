package repository

import (
	"cafeteria/internal/helpers"
	"cafeteria/internal/models"
	"cafeteria/pkg/config"
	"cafeteria/pkg/jtoken"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type UserRepository struct {
	Db *sql.DB
}

const (
	expirationJWT = time.Hour * 5
)

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		Db: db,
	}
}

func (r *UserRepository) Register(ctx context.Context, user *models.User) (string, error) {
	query := `
		INSERT INTO users (username, age, sex, password, allergens) 
		VALUES ($1, $2, $3, $4, $5)`

	vals := []any{user.Username, user.Age, user.Sex, helpers.CreateMd5Hash(user.Password), user.Allergens}

	_, err := r.Db.QueryContext(ctx, query, vals...)
	if err != nil {
		return "", err
	}

	payload := LoadPayload(user)
	token, err := jtoken.GenerateAccessToken(payload, config.GetJWTSecret())
	if err != nil {
		return "", err
	}

	return token, nil
}

func (r *UserRepository) GetToken(ctx context.Context, username, pass string) (string, error) {
	query := `SELECT username, password, is_admin, age, sex, registration_date, allergens FROM users WHERE username = $1`

	user := &models.User{}
	err := r.Db.QueryRowContext(ctx, query, username).Scan(
		&user.Username,
		&user.Password,
		&user.IsAdmin,
		&user.Age,
		&user.Sex,
		&user.FirstOrder,
		&user.Allergens,
	)
	if err != nil {
		return "", fmt.Errorf("error fetching user: %v", err)
	}

	hashedPass := helpers.CreateMd5Hash(pass)
	if strings.Trim(hashedPass, " ") != strings.Trim(user.Password, " ") {
		fmt.Println(hashedPass, user.Password)
		return "", fmt.Errorf("invalid password")
	}

	payload := LoadPayload(user)
	token, err := jtoken.GenerateAccessToken(payload, config.GetJWTSecret())
	if err != nil {
		return "", fmt.Errorf("eror generating jwt token: %v", err)
	}

	return token, nil
}

func LoadPayload(user *models.User) *jtoken.Payload {
	payload := &jtoken.Payload{}
	payload.IsAdmin = user.IsAdmin
	payload.Username = user.Username
	payload.ExpiresAt = time.Now().Add(expirationJWT)

	return payload
}
