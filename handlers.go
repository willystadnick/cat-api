package main

import (
    "encoding/json"
    "io/ioutil"
    "net/http"

    "github.com/dgrijalva/jwt-go"
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
)

func ping() func(c *gin.Context) {
    return func(c *gin.Context) {
        c.String(200, "pong")
    }
}

func breeds(api Api) func(c *gin.Context) {
    return func(c *gin.Context) {
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
    }
}

func login(api Api) func(c *gin.Context) {
    return func(c *gin.Context) {
        var body, user User
        c.BindJSON(&body)

        api.Database.Where("username = ?", body.Username).Find(&user)

        err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
        if err != nil {
            c.String(400, "invalid credentials")
            return
        }

        token, err := jwt.New(jwt.SigningMethodHS256).SignedString([]byte("secret"))
        if err != nil {
            c.String(500, "failed to generate jwt")
            return
        }

        c.String(200, token)
    }
}
