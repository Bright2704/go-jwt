package main

import (
	"fmt"
	AuthController "gocode/jwt-api/controller/auth"
	UserController "gocode/jwt-api/controller/user"
	"gocode/jwt-api/middleware"
	"gocode/jwt-api/orm"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// Binding from JSON
type Register struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" bimding:"required"`
	Fullname string `json:"fullname" bimding:"required"`
	Avatar 	 string `json:"avatar" bimding:"required"`
}

type User struct {
	gorm.Model  //ตัวโมเดลสำเร็จรูป มี id, create_at, update_at, deleted_at
	Username	string
	Password	string
	Fullname	string
	Avatar		string
}

func main() {
	err := godotenv.Load(".env")

	if err != nil	{
		fmt.Println("Error loading .env file")

	}

	orm.InitDB()

	r := gin.Default()
  	r.Use(cors.Default())		//ทำให้สามารถภายนอกมาใช้ api ได้ 
	r.POST("/register", AuthController.Register)
	r.POST("/login", AuthController.Login)
	authorized := r.Group("/users", middleware.JWTAuthen())
	authorized.GET("/readall", UserController.ReadAll)
	authorized.GET("/profile", UserController.Profile)
	r.Run("localhost:8080")
}
