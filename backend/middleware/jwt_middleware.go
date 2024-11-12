package middleware

import (
    "net/http"
    "os"
    "fmt"
    "strings"
    "github.com/golang-jwt/jwt"
    "github.com/labstack/echo/v4"
)

// JWTMiddleware memeriksa apakah ada JWT yang valid di header permintaan dan mengambil informasi pengguna
func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        tokenString := c.Request().Header.Get("Authorization")
        if tokenString == "" {
            return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Token tidak ditemukan"})
        }

        // Hapus prefix "Bearer "
        tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

        // Parse token dan claims
        claims := jwt.MapClaims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return []byte(os.Getenv("JWT_SECRET")), nil
        })

        if err != nil || !token.Valid {
            // Debug: Log error
            fmt.Println("Error parsing token:", err)
            return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Token tidak valid atau telah kadaluarsa"})
        }

        // Ambil user ID dari claims
        userId, ok := claims["id_user"].(float64)
        if !ok {
            return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Token claims tidak valid"})
        }

        // Simpan userId di context
        c.Set("userId", int(userId))

        return next(c)
    }
}
