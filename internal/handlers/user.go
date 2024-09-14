package handlers

import (
	"log"
	"net/http"
	"yummy_mobile_app_backend/configs"
	"yummy_mobile_app_backend/internal/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword шифрует пароль с использованием bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err

	}
	return string(hashedPassword), nil
}

// CheckPassword проверяет совпадение хэшированного пароля и введённого
func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// SetupRoutes регистрирует маршруты
func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/yummy_app/test") // Префикс для всех маршрутов
	{
		api.POST("/register", func(c *gin.Context) { RegisterUser(c, db) })
		api.POST("/login", func(c *gin.Context) { LoginUser(c, db) })
		api.GET("/user/:id", func(c *gin.Context) { GetUserByID(c, db) })
		api.GET("/users", func(c *gin.Context) { GetAllUsers(c, db) })

	}
	return string(hashedPassword), nil
}


// CheckPassword проверяет совпадение хэшированного пароля и введённого
func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}


// RegisterUser обрабатывает регистрацию нового пользователя
func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Хэшируем пароль перед сохранением в базу данных
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}
	user.Password = hashedPassword

	db := configs.ConnectDB()

	// Сохранение пользователя в базу данных
	if result := db.Create(&user); result.Error != nil {
		log.Printf("Failed to create user: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully!"})
}

// LoginUser обрабатывает аутентификацию пользователя
func LoginUser(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := configs.ConnectDB()
	var user models.User

	// Поиск пользователя по email
	if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Проверка введённого пароля с хэшированным паролем в базе данных
	if err := CheckPassword(user.Password, input.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})

// GetUserByID получает данные пользователя по ID
func GetUserByID(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var user models.User

	// Поиск пользователя по ID
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Возвращаем информацию о пользователе, кроме пароля
	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

// GetAllUsers получает список всех пользователей
func GetAllUsers(c *gin.Context, db *gorm.DB) {
	var users []models.User

	// Получение всех пользователей
	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	// Возвращаем список пользователей без паролей
	userResponses := make([]map[string]interface{}, len(users))
	for i, user := range users {
		userResponses[i] = map[string]interface{}{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		}
	}

	c.JSON(http.StatusOK, gin.H{"users": userResponses})
}
