package main

import (
    "github.com/gin-gonic/gin"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

func setupRouter() *gin.Engine {
    r := gin.Default()
    r.GET("/ping", func(c *gin.Context) {
        c.String(200, "pong")
    })
    return r
}

type Cat struct {
    gorm.Model
    Name string
}

func main() {
    db, err := gorm.Open(mysql.Open("cat-api:cat-api@tcp(127.0.0.1:3306)/cat-api?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    db.AutoMigrate(&Cat{})

    r := setupRouter()
    r.Run(":8080")
}
