package handler

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/sharif-go-lab/go-judge-platform/internal/db"
	"github.com/sharif-go-lab/go-judge-platform/internal/model"
)

// LoginForm binds login inputs
type LoginForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// RegisterForm binds registration inputs
type RegisterForm struct {
	Username string `form:"username" binding:"required,alphanum"`
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required,min=6"`
}

// LoginPageHandler renders the login form
func LoginPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{"title": "Login"})
}

// LoginHandler processes login submissions
func LoginHandler(c *gin.Context) {
	var form LoginForm
	if err := c.ShouldBind(&form); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": "All fields required"})
		return
	}

	var user model.User
	if err := db.DB.Where("username = ?", form.Username).First(&user).Error; err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid credentials"})
		return
	}

	// set session values
	sess := sessions.Default(c)
	sess.Set("user_id", user.ID)
	sess.Set("is_admin", user.IsAdmin)
	sess.Save()

	c.Redirect(http.StatusSeeOther, "/questions")
}

// RegisterPageHandler renders the registration form
func RegisterPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{"title": "Register"})
}

// RegisterHandler processes registration submissions
func RegisterHandler(c *gin.Context) {
	var form RegisterForm
	if err := c.ShouldBind(&form); err != nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{"error": "All fields required"})
		return
	}

	// check existing user/email
	var count int64
	db.DB.Model(&model.User{}).
		Where("username = ? OR email = ?", form.Username, form.Email).
		Count(&count)
	if count > 0 {
		c.HTML(http.StatusConflict, "register.html", gin.H{"error": "Username or email already taken"})
		return
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{"error": "Server error"})
		return
	}

	user := model.User{
		Username:  form.Username,
		Email:     form.Email,
		Password:  string(hash),
		IsAdmin:   false,
		CreatedAt: time.Now(),
	}
	if err := db.DB.Create(&user).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{"error": "Could not create user"})
		return
	}

	// log in new user
	sess := sessions.Default(c)
	sess.Set("user_id", user.ID)
	sess.Save()

	c.Redirect(http.StatusSeeOther, "/questions")
}

// LogoutHandler clears the session and redirects home
func LogoutHandler(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Clear()
	sess.Save()
	c.Redirect(http.StatusSeeOther, "/")
}
