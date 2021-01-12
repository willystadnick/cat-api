package main

import (
    "encoding/json"
    "io/ioutil"
    "net/http"

    "github.com/gin-gonic/gin"
    "gorm.io/driver/mysql"
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

    return db
}

func setupRouter(api Api) *gin.Engine {
    r := gin.Default()

    r.GET("/ping", func(c *gin.Context) {
        c.String(200, "pong")
    })

    r.GET("/breeds", func(c *gin.Context) {
        name := c.Query("name")
        if name == "" {
            c.String(400, "invalid breed name")
            return
        }

        var cats []Cat
        api.Database.Where("name LIKE ?", "%" + name + "%").Find(&cats)
        if len(cats) == 0 {
            cli := &http.Client{}
            req, err := http.NewRequest("GET", "https://api.thecatapi.com/v1/breeds/search?q=" + name, nil)
            if err != nil {
                panic(err)
            }
            req.Header.Add("x-api-key", "DEMO-API-KEY")
            res, err := cli.Do(req)
            if err != nil {
                panic(err)
            }
            defer res.Body.Close()
            body, err := ioutil.ReadAll(res.Body)
            if err != nil {
                panic(err)
            }
            json.Unmarshal(body, &cats)
            if len(cats) == 0 {
                c.String(200, "no records")
                return
            }
            api.Database.Create(&cats)
        }

        c.JSON(200, cats)
    })

    return r
}
