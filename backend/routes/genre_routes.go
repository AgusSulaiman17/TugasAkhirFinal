package routes

import (
    "github.com/labstack/echo/v4"
    "backend/controllers"
    "backend/middleware"
)

func GenreRoutes(e *echo.Echo) {
    adminGroup := e.Group("/genres")
    adminGroup.Use(middleware.JWTMiddleware) // Menggunakan middleware JWT untuk otorisasi

    adminGroup.POST("", controllers.CreateGenre)
    adminGroup.GET("", controllers.GetGenres)
    adminGroup.GET("/:id", controllers.GetGenre)
    adminGroup.PUT("/:id", controllers.UpdateGenre)
    adminGroup.DELETE("/:id", controllers.DeleteGenre)
}
