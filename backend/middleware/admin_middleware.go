package middleware

import (
    "net/http"
    "strings"

    "github.com/dgrijalva/jwt-go"
    "github.com/labstack/echo/v4"
)

func AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        token := c.Request().Header.Get("Authorization")
        if token == "" {
            return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Missing token"})
        }

        token = strings.Replace(token, "Bearer ", "", 1)

        claims := &jwt.MapClaims{}
        t, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
            return []byte("JmySuperSecretKey12345"), nil
        })
        if err != nil || !t.Valid {
            return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token"})
        }

        role := (*claims)["role"].(string)
        if role != "admin" {
            return c.JSON(http.StatusForbidden, map[string]string{"message": "Access denied"})
        }

        return next(c)
    }
}
