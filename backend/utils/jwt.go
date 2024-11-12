package utils

import (
    "time"
    "os"

    "github.com/golang-jwt/jwt"
    "backend/models"
)

// GenerateJWT generates a new JWT for a given user
func GenerateJWT(user models.User) (string, error) {
    // Get the secret key from environment variable
    secretKey := []byte(os.Getenv("JWT_SECRET"))

    claims := jwt.MapClaims{
        "id_user": user.IDUser,
        "nama":    user.Nama,
        "email":   user.Email,
        "role":    user.Role,
        "exp":     time.Now().Add(time.Hour * 72).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(secretKey)
}
