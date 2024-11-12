package main

import (
    "backend/config"
    "backend/routes"
    "backend/middleware"
    "github.com/labstack/echo/v4"
    "github.com/go-co-op/gocron" 
    "backend/utils" 
    "backend/models"
    "time"
    "fmt"
)

func main() {
    config.ConnectDB()
    e := echo.New()
    if config.DB == nil {
        panic("Database connection failed. Please check your configuration.")
    }
    e.Use(middleware.CORS())
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    routes.AuthRoutes(e)
    routes.UserRoutes(e)
    routes.GenreRoutes(e)
    routes.PenulisRoutes(e)
    routes.BookRoutes(e)
    routes.PeminjamanRoutes(e)

    // Setup scheduler untuk pengiriman email notifikasi
    go scheduleEmailNotifications()

    e.Logger.Fatal(e.Start(":8080"))
}

// Fungsi untuk mengatur scheduler
func scheduleEmailNotifications() {
    s := gocron.NewScheduler(time.UTC)

    // Atur agar fungsi notifikasi email dijalankan setiap 24 jam
    s.Every(24).Hours().Do(notifyBeforeReturn)

    // Menjalankan scheduler
    s.StartAsync()
}

// Fungsi untuk notifikasi pengembalian buku
func notifyBeforeReturn() {
    // Ambil peminjaman buku yang harus segera dikembalikan (misalnya dalam 1 hari)
    var peminjaman []models.Peminjaman

    if err := config.DB.Where("status_kembali = ? AND tanggal_kembali < ?", false, time.Now().Add(24*time.Hour)).Find(&peminjaman).Error; err != nil {
        fmt.Println("Error fetching peminjaman:", err)
        return
    }

    // Kirim email notifikasi untuk setiap peminjaman
    for _, pinjam := range peminjaman {
        userEmail := "user-email@example.com" // Ganti sesuai email user terkait
        subject := "Peringatan Pengembalian Buku"
        body := fmt.Sprintf("Buku '%d' yang Anda pinjam akan jatuh tempo dalam 1 hari. Harap segera mengembalikan buku tersebut.", pinjam.IDBuku)

        if err := utils.SendEmail(userEmail, subject, body); err != nil {
            fmt.Println("Gagal mengirim email:", err)
        }
    }
}
