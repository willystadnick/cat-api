package main

import (
    "os"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    j "github.com/gin-gonic/contrib/jwt"
)

func main() {
    api := setupApi()

    api.Router.Run(":8080")
}

func setupApi() Api {
    err := godotenv.Load()
    if err != nil {
        panic("failed to load env vars")
    }

    var api Api
    api.Database = setupDatabase()
    api.Router = setupRouter(api)

    return api
}

func setupDatabase() *gorm.DB {
    db_user := os.Getenv("DB_USER")
    db_pass := os.Getenv("DB_PASS")
    db_url := os.Getenv("DB_URL")
    db_port := os.Getenv("DB_PORT")
    db_name := os.Getenv("DB_NAME")
    db_charset := os.Getenv("DB_CHARSET")
    db_parsetime := os.Getenv("DB_PARSETIME")
    db_loc := os.Getenv("DB_LOC")
    dsn := db_user + ":" + db_pass + "@tcp(" + db_url + ":" + db_port + ")/" + db_name + "?charset=" + db_charset + "&parseTime=" + db_parsetime + "&loc=" + db_loc

    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    db.AutoMigrate(&Cat{})
    db.AutoMigrate(&User{})

    password := []byte(os.Getenv("ADMIN_PASS"))
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
    r.GET("/breeds", j.Auth(os.Getenv("JWT_SECRET")), breeds(api))
    r.POST("/login", login(api))

    return r
}
