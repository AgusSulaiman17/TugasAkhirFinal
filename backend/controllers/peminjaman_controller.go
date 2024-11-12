package controllers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"fmt"
	"strconv"
	"backend/utils"
	"backend/config"
	"backend/models"
)

// Create Peminjaman
func CreatePeminjaman(c echo.Context) error {
	var peminjaman models.Peminjaman
	if err := c.Bind(&peminjaman); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	// Ambil ID user dari token JWT
	userId := c.Get("userId").(int)
	peminjaman.IDUser = userId

	// Set tanggal pinjam menjadi sekarang
	peminjaman.TanggalPinjam = time.Now()

	// Simpan peminjaman ke database
	if err := config.DB.Create(&peminjaman).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	// Ambil data pengguna untuk email
	var user models.User
	if err := config.DB.First(&user, peminjaman.IDUser).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal mengambil data user"})
	}

	// Kirim email konfirmasi
	userEmail := user.Email
	subject := "Konfirmasi Peminjaman Buku"
	body := "Terima kasih telah meminjam buku. Harap kembalikan sesuai tanggal yang ditentukan."
	utils.SendEmail(userEmail, subject, body)

	return c.JSON(http.StatusCreated, peminjaman)
}


// Get All Peminjaman by User
func GetAllPeminjaman(c echo.Context) error {
	userId := c.Get("userId").(int)

	var peminjaman []models.Peminjaman
	if err := config.DB.Where("id_user = ?", userId).Find(&peminjaman).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, peminjaman)
}

// Update Peminjaman
func UpdatePeminjaman(c echo.Context) error {
	id := c.Param("id")
	var peminjaman models.Peminjaman

	if err := config.DB.First(&peminjaman, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "Peminjaman not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	if err := c.Bind(&peminjaman); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := config.DB.Save(&peminjaman).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, peminjaman)
}

// Delete Peminjaman
func DeletePeminjaman(c echo.Context) error {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Peminjaman{}, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Peminjaman deleted"})
}

// ReturnBook mengelola pengembalian buku oleh user
func ReturnBook(c echo.Context) error {
	idPeminjaman, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "ID peminjaman tidak valid"})
	}

	// Ambil ID user dari JWT
	userId := c.Get("userId").(int)

	var peminjaman models.Peminjaman
	if err := config.DB.First(&peminjaman, idPeminjaman).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "Data peminjaman tidak ditemukan"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal mencari data peminjaman"})
	}

	// Validasi bahwa hanya user yang melakukan peminjaman dapat mengembalikan
	if peminjaman.IDUser != userId {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Anda tidak berhak mengembalikan peminjaman ini"})
	}

	// Cek apakah buku sudah dikembalikan
	if peminjaman.StatusKembali {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Buku sudah dikembalikan"})
	}

	// Set status kembali dan tanggal kembali
	peminjaman.StatusKembali = true
	peminjaman.TanggalKembali = time.Now()

	// Hitung denda jika pengembalian terlambat
	tanggalBatasPengembalian := peminjaman.TanggalPinjam.Add(time.Duration(peminjaman.DurasiHari) * 24 * time.Hour)
	if peminjaman.TanggalKembali.After(tanggalBatasPengembalian) {
		terlambatHari := int(peminjaman.TanggalKembali.Sub(tanggalBatasPengembalian).Hours() / 24)
		peminjaman.Denda = terlambatHari * 5000 // Misalnya denda 5000 per hari
	}

	// Simpan perubahan ke database
	if err := config.DB.Save(&peminjaman).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal mengupdate data peminjaman"})
	}

	// Kirim email setelah pengembalian buku
	var user models.User
	if err := config.DB.First(&user, peminjaman.IDUser).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal mengambil data user"})
	}

	// Kirim email pemberitahuan
	userEmail := user.Email
	if peminjaman.Denda > 0 {
		subject := "Pemberitahuan Denda Pengembalian Buku"
		body := fmt.Sprintf("Buku yang Anda pinjam terlambat dikembalikan. Anda dikenakan denda sebesar %d", peminjaman.Denda)
		utils.SendEmail(userEmail, subject, body)
	} else {
		subject := "Pemberitahuan Pengembalian Buku"
		body := "Buku yang Anda pinjam telah berhasil dikembalikan."
		utils.SendEmail(userEmail, subject, body)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":    "Buku berhasil dikembalikan",
		"peminjaman": peminjaman,
	})
}
