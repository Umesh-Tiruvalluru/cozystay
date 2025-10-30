package helper

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func WriteJSON(w http.ResponseWriter, data any, status int) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateToken(id uuid.UUID, role string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":   id,
			"role": role,
		},
	)

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenString string) (map[string]any, error) {
	token, err := jwt.Parse(
		tokenString,
		func(t *jwt.Token) (any, error) {
			return []byte("simba"), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return nil, errors.New("Invalid Token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return map[string]any{}, errors.New("Invalid token claims")
	}

	userID, ok := claims["id"].(string)
	if !ok {
		return map[string]any{}, errors.New("user_id missing in token")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return map[string]any{}, errors.New("role is missing in token")
	}

	data := map[string]any{
		"userID": userID,
		"role":   role,
	}

	return data, nil
}
