package main

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type Api struct {
    Database *gorm.DB
    Router *gin.Engine
}

type Cat struct {
    gorm.Model
    Name string
}

type User struct {
    gorm.Model
    Username string
    Password string
}
