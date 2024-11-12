package controllers

import (
    "net/http"
    "strconv"

	"time"
	"gorm.io/gorm"
    "github.com/labstack/echo/v4"
    "github.com/dgrijalva/jwt-go"
    "backend/config"
    "backend/models"
)


// GetUser retrieves a user by ID
func GetUser(c echo.Context) error {
    id := c.Param("id")
    var user models.User

    if err := config.DB.Where("id_user = ?", id).First(&user).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
        }
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error fetching user"})
    }

    return c.JSON(http.StatusOK, user)
}

// GetUsers retrieves all users
func GetUsers(c echo.Context) error {
    var users []models.User

    if err := config.DB.Find(&users).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error fetching users"})
    }

    return c.JSON(http.StatusOK, users)
}

// UpdateUser mengupdate informasi pengguna
func UpdateUser(c echo.Context) error {
    id := c.Param("id") // Mengambil ID pengguna dari parameter URL
    var user models.User
    if err := c.Bind(&user); err != nil { // Mengikat data dari request ke struct user
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "Input tidak valid"})
    }

    // Memeriksa apakah pengguna adalah admin atau mencoba memperbarui data mereka sendiri
    claims := c.Get("claims").(jwt.MapClaims)
    role := claims["role"].(string)
    userId := claims["id_user"].(float64) // Melakukan type assertion ke float64
    if role != "admin" && strconv.Itoa(int(userId)) != id {
        return c.JSON(http.StatusForbidden, map[string]string{"message": "Anda tidak berwenang untuk memperbarui pengguna ini"})
    }

    // Perbarui timestamp DiperbaruiPada
    user.DiperbaruiPada = time.Now() // Mengatur waktu sekarang

    // Mengupdate data pengguna di database
    query := "UPDATE users SET nama = ?, email = ?, role = ?, diperbarui_pada = ? WHERE id_user = ?"
    err := config.DB.Exec(query, user.Nama, user.Email, user.Role, user.DiperbaruiPada, id).Error
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Kesalahan saat memperbarui pengguna"})
    }

    return c.JSON(http.StatusOK, map[string]string{"message": "Pengguna berhasil diperbarui"})
}
// DeleteUser deletes a user by ID
func DeleteUser(c echo.Context) error {
    id := c.Param("id")

    // Check if the user is admin
    claims := c.Get("claims").(jwt.MapClaims)
    role := claims["role"].(string)
    if role != "admin" {
        return c.JSON(http.StatusForbidden, map[string]string{"message": "You are not authorized to delete this user"})
    }

    query := "DELETE FROM users WHERE id_user = ?"
    err := config.DB.Exec(query, id).Error
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error deleting user"})
    }

    return c.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
}
