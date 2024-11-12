package routes

import (
    "backend/controllers"
    "backend/middleware"
    "github.com/labstack/echo/v4"
)

func PeminjamanRoutes(e *echo.Echo) {
    // Gunakan JWTMiddleware untuk setiap rute peminjaman yang membutuhkan autentikasi
    e.POST("/peminjaman", middleware.JWTMiddleware(controllers.CreatePeminjaman))
    e.GET("/peminjaman", middleware.JWTMiddleware(controllers.GetAllPeminjaman))
    e.PUT("/peminjaman/:id", middleware.JWTMiddleware(controllers.UpdatePeminjaman))
    e.DELETE("/peminjaman/:id", middleware.JWTMiddleware(controllers.DeletePeminjaman))
    e.POST("/peminjaman/:id/kembalikan", middleware.JWTMiddleware(controllers.ReturnBook))
}
