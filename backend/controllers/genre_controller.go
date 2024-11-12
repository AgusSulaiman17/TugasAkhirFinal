package controllers

import (
    "net/http"
    "backend/config"
    "backend/models"

    "github.com/labstack/echo/v4"
)

// CreateGenre
func CreateGenre(c echo.Context) error {
    var genre models.Genre
    if err := c.Bind(&genre); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
    }

    if err := config.DB.Create(&genre).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error creating genre"})
    }

    return c.JSON(http.StatusCreated, genre)
}

// GetGenres
func GetGenres(c echo.Context) error {
    var genres []models.Genre
    if err := config.DB.Find(&genres).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error fetching genres"})
    }

    return c.JSON(http.StatusOK, genres)
}

// GetGenre
func GetGenre(c echo.Context) error {
    id := c.Param("id")
    var genre models.Genre
    if err := config.DB.First(&genre, id).Error; err != nil {
        return c.JSON(http.StatusNotFound, map[string]string{"message": "Genre not found"})
    }

    return c.JSON(http.StatusOK, genre)
}

// UpdateGenre
func UpdateGenre(c echo.Context) error {
    id := c.Param("id")
    var genre models.Genre
    if err := config.DB.First(&genre, id).Error; err != nil {
        return c.JSON(http.StatusNotFound, map[string]string{"message": "Genre not found"})
    }

    if err := c.Bind(&genre); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
    }

    if err := config.DB.Save(&genre).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error updating genre"})
    }

    return c.JSON(http.StatusOK, genre)
}

// DeleteGenre
func DeleteGenre(c echo.Context) error {
    id := c.Param("id")
    if err := config.DB.Delete(&models.Genre{}, id).Error; err != nil {
        return c.JSON(http.StatusNotFound, map[string]string{"message": "Genre not found"})
    }

    return c.NoContent(http.StatusNoContent)
}
