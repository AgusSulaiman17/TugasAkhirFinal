package controllers

import (
    "net/http"
    "backend/models"
    "backend/config"
    "github.com/labstack/echo/v4"
)
func CreatePenulis(c echo.Context) error {
    var penulis models.Penulis
    if err := c.Bind(&penulis); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
    }

    if err := config.DB.Create(&penulis).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create penulis"})
    }
    
    return c.JSON(http.StatusCreated, penulis)
}

func GetPenulis(c echo.Context) error {
    id := c.Param("id")
    var penulis models.Penulis

    if err := config.DB.First(&penulis, id).Error; err != nil {
        return c.JSON(http.StatusNotFound, map[string]string{"message": "Penulis not found"})
    }

    return c.JSON(http.StatusOK, penulis)
}
func GetAllPenulis(c echo.Context) error {
    var penulis []models.Penulis

    if err := config.DB.Find(&penulis).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve penulis"})
    }

    return c.JSON(http.StatusOK, penulis)
}

func UpdatePenulis(c echo.Context) error {
    id := c.Param("id")
    var penulis models.Penulis

    if err := config.DB.First(&penulis, id).Error; err != nil {
        return c.JSON(http.StatusNotFound, map[string]string{"message": "Penulis not found"})
    }

    if err := c.Bind(&penulis); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
    }

    if err := config.DB.Save(&penulis).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update penulis"})
    }
    
    return c.JSON(http.StatusOK, penulis)
}

func DeletePenulis(c echo.Context) error {
    id := c.Param("id")
    var penulis models.Penulis

    if err := config.DB.First(&penulis, id).Error; err != nil {
        return c.JSON(http.StatusNotFound, map[string]string{"message": "Penulis not found"})
    }

    if err := config.DB.Delete(&penulis).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete penulis"})
    }

    return c.NoContent(http.StatusNoContent)
}
