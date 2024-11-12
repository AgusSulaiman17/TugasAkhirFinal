package controllers

import (
    "net/http"
    "time"
    "gorm.io/gorm"
    "golang.org/x/crypto/bcrypt"
    "github.com/labstack/echo/v4"
    "backend/models"
    "backend/config"
    "backend/utils"
)

// Fungsi Register menangani registrasi pengguna
func Register(c echo.Context) error {
    var user models.User

    // Bind data input ke variabel user
    if err := c.Bind(&user); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "Input tidak valid"})
    }

    // Jika role tidak diisi, set default ke 'user'
    if user.Role == "" {
        user.Role = "user"
    }

    // Cek apakah email sudah terdaftar
    var existingUser models.User
    if err := config.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
        return c.JSON(http.StatusConflict, map[string]string{"message": "Email sudah terdaftar"})
    }

    // Hash password menggunakan bcrypt untuk keamanan
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.KataSandi), bcrypt.DefaultCost)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Tidak dapat hash password"})
    }
    user.KataSandi = string(hashedPassword)

    // Tentukan timestamp dibuat dan diperbarui
    user.DibuatPada = time.Now()
    user.DiperbaruiPada = time.Now()

    // Simpan data pengguna ke dalam database
    if err := config.DB.Create(&user).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Tidak dapat membuat pengguna"})
    }

    // Kirim email konfirmasi ke pengguna
    subject := "Registrasi Berhasil - Selamat Datang!"
    body := "Halo " + user.Nama + ",\n\nTerima kasih telah mendaftar! Selamat datang di aplikasi kami."
    err = utils.SendEmail(user.Email, subject, body)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Pendaftaran berhasil, namun email gagal dikirim"})
    }

    return c.JSON(http.StatusCreated, map[string]string{"message": "Registrasi berhasil, silakan cek email Anda untuk konfirmasi"})
}

// Fungsi Login menangani login pengguna
func Login(c echo.Context) error {
    var user models.User
    var dbUser models.User

    // Bind data input ke variabel user
    if err := c.Bind(&user); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "Input tidak valid"})
    }

    // Cari pengguna berdasarkan email
    if err := config.DB.Where("email = ?", user.Email).First(&dbUser).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Kredensial tidak valid"})
        }
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Kesalahan saat mencari pengguna"})
    }

    // Bandingkan password hash
    err := bcrypt.CompareHashAndPassword([]byte(dbUser.KataSandi), []byte(user.KataSandi))
    if err != nil {
        return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Kredensial tidak valid"})
    }

    // Buat token JWT untuk autentikasi
    token, err := utils.GenerateJWT(dbUser)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Kesalahan saat membuat token"})
    }

    // Berikan respons yang menyertakan token dan data user
    return c.JSON(http.StatusOK, map[string]interface{}{
        "token": token,
        "user": map[string]interface{}{
            "id_user": dbUser.IDUser,
            "nama":    dbUser.Nama,
            "email":   dbUser.Email,
            "role":    dbUser.Role,
        },
    })
}