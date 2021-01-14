package main

import (

    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    j "github.com/gin-gonic/contrib/jwt"
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

func main() {
    var api Api
    api.Database = setupDatabase()
    api.Router = setupRouter(api)

    api.Router.Run(":8080")
}

func setupDatabase() *gorm.DB {
    db, err := gorm.Open(mysql.Open("cat-api:cat-api@tcp(127.0.0.1:3306)/cat-api?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    db.AutoMigrate(&Cat{})
    db.AutoMigrate(&User{})

    password := []byte("@#$RF@!718")

    hashed, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
    if err != nil {
        panic(err)
    }

    var user User
    db.FirstOrCreate(&user, User{
        Username: "admin",
    })

    user.Password = string(hashed)
    db.Save(&user)

    return db
}

func setupRouter(api Api) *gin.Engine {
    r := gin.Default()
    r.GET("/ping", ping())
    r.GET("/breeds", j.Auth("secret"), breeds(api))
    r.POST("/login", login(api))

    return r
}
