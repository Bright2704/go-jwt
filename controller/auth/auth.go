package auth

import (
	"fmt"
	"gocode/jwt-api/orm"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var hmacSampleSecret []byte
// Binding from JSON
type RegisterBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" bimding:"required"`
	Fullname string `json:"fullname" bimding:"required"`
	Avatar 	 string `json:"avatar" bimding:"required"`
}

func Register(c *gin.Context) {
	var json RegisterBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//Check User Exists
	var userExist orm.User 
	orm.Db.Where("username = ?", json.Username).First(&userExist)
	if userExist.ID > 0 {
		c.JSON(http.StatusOK, gin.H{"status": "error", "massage": "User Exists"})
		return
	}


	// Create User
	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(json.Password), 10)  //สร้างUserPassword เอาไปเก็บไว้
	user := orm.User{Username: json.Username, Password: string(encryptedPassword),		//สร้างตัว User
	 Fullname: json.Fullname, Avatar: json.Avatar,}
	orm.Db.Create(&user)
	if (user.ID > 0) {		//ตรวจสอบว่ามีไอดีซ้ำหรือไม่
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Create Success", "userId": user.ID})
	}else {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "User Create Failed",})
	}
}

// Binding from JSON
type LoginBody struct {			//ล็อคอิน/ตัวโดยใช้ username-password
	Username string `json:"username" binding:"required"`
	Password string `json:"password" bimding:"required"`
	
}

func Login(c *gin.Context) {			//เช็คว่าที่ล็อคอินมามีข้อมูลครบไหม
	var json LoginBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Check User Exists
	var userExist orm.User 			//กดหนดถ้า User เป็น 0 หรือไม่มีก็จะแสดงข้อความแล้วออกมาเลย
 	orm.Db.Where("username = ?", json.Username).First(&userExist)
	if userExist.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"status": "error", "massage": "User Does Not Exists"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(userExist.Password), []byte(json.Password))
	if  err == nil {
		hmacSampleSecret = []byte(os.Getenv("JWT_SECRET_KET"))
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userId" : userExist.ID,
			"exp": 	time.Now().Add(time.Minute * 1).Unix(), 
		})
		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString(hmacSampleSecret)
		fmt.Println(tokenString, err)


		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Login Success", "token": tokenString})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "Login Failed"})
	}

}