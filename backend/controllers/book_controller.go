package controllers

import (
    "net/http"
    "os"
    "path/filepath"
    "io"
    "strconv"

    "github.com/labstack/echo/v4"
    "gorm.io/gorm"
    "backend/models"
    "backend/config"
)

var db *gorm.DB

func init() {
    db = config.DB
}


// CreateBuku untuk menambahkan buku baru
func CreateBuku(c echo.Context) error {
    var buku models.Buku
    file, err := c.FormFile("gambar")
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "No image provided"})
    }

    // Simpan gambar
    src, err := file.Open()
    if err != nil {
        return err
    }
    defer src.Close()

    // Pastikan folder "uploads" ada
    if _, err := os.Stat("uploads"); os.IsNotExist(err) {
        os.Mkdir("uploads", os.ModePerm)
    }

    dst := filepath.Join("uploads", file.Filename)
    out, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer out.Close()

    if _, err = io.Copy(out, src); err != nil {
        return err
    }

    buku.Judul = c.FormValue("judul")
    buku.IdPenulis, _ = strconv.Atoi(c.FormValue("id_penulis"))
    buku.IdGenre, _ = strconv.Atoi(c.FormValue("id_genre"))
    buku.Deskripsi = c.FormValue("deskripsi")
    buku.Jumlah, _ = strconv.Atoi(c.FormValue("jumlah"))
    buku.Gambar = dst // Simpan path gambar
    buku.Status = true

    // Gunakan config.DB secara langsung untuk menghindari error variabel nil
    if err := config.DB.Create(&buku).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create book"})
    }

    return c.JSON(http.StatusCreated, buku)
}


// GetBuku untuk mendapatkan daftar semua buku
func GetBuku(c echo.Context) error {
    var bukus []models.Buku
    if err := db.Find(&bukus).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve books"})
    }
    return c.JSON(http.StatusOK, bukus)
}

// GetBukuByID untuk mendapatkan buku berdasarkan ID
func GetBukuByID(c echo.Context) error {
    id := c.Param("id")
    var buku models.Buku
    if err := db.First(&buku, id).Error; err != nil {
        return c.JSON(http.StatusNotFound, map[string]string{"message": "Book not found"})
    }
    return c.JSON(http.StatusOK, buku)
}

// UpdateBuku untuk memperbarui informasi buku
func UpdateBuku(c echo.Context) error {
    id := c.Param("id")
    var buku models.Buku
    if err := db.First(&buku, id).Error; err != nil {
        return c.JSON(http.StatusNotFound, map[string]string{"message": "Book not found"})
    }

    if err := c.Bind(&buku); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
    }

    // Hanya update kolom yang diizinkan
    if c.FormValue("judul") != "" {
        buku.Judul = c.FormValue("judul")
    }
    if c.FormValue("id_penulis") != "" {
        buku.IdPenulis, _ = strconv.Atoi(c.FormValue("id_penulis"))
    }
    if c.FormValue("id_genre") != "" {
        buku.IdGenre, _ = strconv.Atoi(c.FormValue("id_genre"))
    }
    if c.FormValue("deskripsi") != "" {
        buku.Deskripsi = c.FormValue("deskripsi")
    }
    if c.FormValue("jumlah") != "" {
        buku.Jumlah, _ = strconv.Atoi(c.FormValue("jumlah"))
    }
    if c.FormValue("gambar") != "" {
        // Simpan gambar baru jika diupload
        file, err := c.FormFile("gambar")
        if err == nil {
            src, err := file.Open()
            if err == nil {
                defer src.Close()
                dst := filepath.Join("uploads", file.Filename)
                out, err := os.Create(dst)
                if err == nil {
                    defer out.Close()
                    io.Copy(out, src)
                    buku.Gambar = dst // Update path gambar
                }
            }
        }
    }

    if err := db.Save(&buku).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update book"})
    }

    return c.JSON(http.StatusOK, buku)
}

// DeleteBuku untuk menghapus buku berdasarkan ID
func DeleteBuku(c echo.Context) error {
    id := c.Param("id")
    var buku models.Buku
    if err := db.First(&buku, id).Error; err != nil {
        return c.JSON(http.StatusNotFound, map[string]string{"message": "Book not found"})
    }

    if err := db.Delete(&buku).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete book"})
    }

    return c.JSON(http.StatusOK, map[string]string{"message": "Book deleted"})
}
