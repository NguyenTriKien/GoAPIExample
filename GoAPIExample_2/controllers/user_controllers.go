package controllers

import (
	"net/http"
	"os"
	"time"

	"Module/API/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// Create user input model
type UserInput struct {
	ID       uint   `json:"userid"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

/* GET ALL */
// GET '/user'
// Get all user
func FindAllUser(c *gin.Context) {
	var users []models.User
	var userResponses []UserInput

	models.DB.Find(&users)

	for _, user := range users {
		userResponses = append(userResponses, UserInput{
			ID:       user.ID,
			Username: user.Username,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": userResponses})
}

/* GET */
// GET '/user/todolist'
// Get all user
func FindAllUserWithTodoList(c *gin.Context) {
	var user []models.User
	models.DB.Preload("TodoLists").Find(&user)

	c.JSON(http.StatusOK, gin.H{"data": user})
}

/*CREATE*/
//POST "/user"
// Create new userfunc
func CreateUser(c *gin.Context) {
	// Validate input
	var input UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to hash password"})
	}

	user := models.User{Username: input.Username, Password: string(hash)}

	result := models.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": user})

}

/* LOGIN */
// /user/login
func Login(c *gin.Context) {
	var input UserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := models.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}
	//Compare sent in pass with saved user pass hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})

		return
	}

	// Create JWT
	token, err := createToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func createToken(userID uint) (string, error) { // Đây là khai báo của hàm createToken.
	//Hàm này nhận vào một tham số userID kiểu uint (số nguyên không dấu)
	//và trả về một chuỗi và một lỗi.

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{ // Dòng này tạo ra một token JWT mới với một số yêu cầu (claims).
		// Phương thức ký được sử dụng ở đây là HS256 (HMAC-SHA256).
		//Các yêu cầu bao gồm userID và exp.

		"userID": userID,                                // Đây là ID của người dùng, được trả về khi token được giải mã.
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // Token expiration time (24 hours)
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET"))) // Dòng này ký token với một “bí mật” được lấy từ biến môi trường SECRET.
	//Kết quả là một chuỗi token đã được ký và một lỗi (nếu có).

	if err != nil {
		return "", err
	} //Đây là một khối kiểm tra lỗi. Nếu có lỗi xảy ra khi ký token (tức là err khác nil), hàm sẽ trả về một chuỗi rỗng và lỗi đó.

	return tokenString, nil //Nếu không có lỗi xảy ra, hàm sẽ trả về chuỗi token đã được ký và nil cho lỗi.
}

func getUserIDFromToken(c *gin.Context) (uint, bool) {
	tokenString := c.GetHeader("Authorization")

	if tokenString == "" {
		return 0, false
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil || !token.Valid {
		return 0, false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, false
	} //Tóm lại, đoạn mã này giải mã các claims từ một token JWT
	//và kiểm tra xem việc giải mã có thành công hay không. Nếu không, nó trả về 0 và false

	userID, ok := claims["userID"].(float64)
	if !ok {
		return 0, false
	}

	return uint(userID), true
}
