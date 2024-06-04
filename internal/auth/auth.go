// auth/auth.go

package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/nats-io/nats.go"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type AuthService struct {
	db     *sql.DB
	nc     *nats.Conn
	jwtKey []byte
}

func NewAuthService(db *sql.DB, nc *nats.Conn) *AuthService {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	return &AuthService{
		db:     db,
		nc:     nc,
		jwtKey: jwtKey,
	}
}

func (s *AuthService) Register(username, password, email string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := "INSERT INTO Type_Users (Name_Username, Name_Password, Name_Email, Name_Is_Admin) VALUES ($1, $2, $3, $4)"
	_, err = s.db.Exec(query, username, hashedPassword, email, false)
	if err != nil {
		return err
	}

	message := fmt.Sprintf("Зарегистрирован новый пользователь: %s %s", username, email)
	s.nc.Publish("log", []byte(message))

	return nil
}

func (s *AuthService) GenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) Authenticate(username, password string) (string, error) {
	var hashedPassword string
	query := "SELECT Name_Password FROM Type_Users WHERE Name_Username=$1"
	row := s.db.QueryRow(query, username)
	err := row.Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("user not found")
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	tokenString, err := s.GenerateJWT(username)
	if err != nil {
		return "", err
	}

	message := fmt.Sprintf("Авторизован пользователь: %s", username)
	s.nc.Publish("log", []byte(message))

	return tokenString, nil
}
