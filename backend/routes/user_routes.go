package routes

import (
    "github.com/labstack/echo/v4"
    "backend/controllers"
    "backend/middleware"
)

func UserRoutes(e *echo.Echo) {
    e.GET("/user/:id", controllers.GetUser, middleware.JWTMiddleware)
    e.GET("/users", controllers.GetUsers, middleware.JWTMiddleware)
    e.PUT("/user/:id", controllers.UpdateUser, middleware.JWTMiddleware)
    e.DELETE("/user/:id", controllers.DeleteUser, middleware.JWTMiddleware)
}
