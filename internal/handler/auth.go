package handler

import (
	"log"
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
	Username string `form:"username" binding:"required"` // can be username or email
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
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"title": "Login",
			"error": "All fields are required",
		})
		return
	}

	var user model.User
	// ─── KEY CHANGE #1: allow lookup by username OR email ────────────────
	if err := db.DB.
		Where("username = ? OR email = ?", form.Username, form.Username).
		First(&user).Error; err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"title": "Login",
			"error": "Invalid credentials",
		})
		return
	}

	// compare bcrypt hashes
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(form.Password),
	); err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"title": "Login",
			"error": "Invalid credentials",
		})
		return
	}

	// set session values
	sess := sessions.Default(c)
	sess.Set("username", user.Username)
	sess.Set("user_id", user.ID)
	sess.Set("is_admin", user.IsAdmin)
	sess.Save()

	c.Redirect(http.StatusSeeOther, "/profile")
}

// RegisterPageHandler renders the registration form
func RegisterPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{"title": "Register"})
}

// RegisterHandler processes registration submissions
func RegisterHandler(c *gin.Context) {
	// ─── KEY CHANGE #2: bind to your RegisterForm type ───────────────────
	var form RegisterForm
	if err := c.ShouldBind(&form); err != nil {
		log.Println(err.Error())
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"title": "Register",
			"error": "All fields are required and must be valid",
		})
		return
	}

	// hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{
			"title": "Register",
			"error": "Internal error, please try again",
		})
		return
	}

	// construct the new user record
	user := model.User{
		Username:     form.Username,
		Email:        form.Email,
		PasswordHash: string(hash),
		IsAdmin:      false,
		CreatedAt:    time.Now(),
	}

	// write to database
	if err := db.DB.Create(&user).Error; err != nil {
		log.Println(err.Error())
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{
			"title": "Register",
			"error": "Could not create user (maybe username or email already exists)",
		})
		return
	}

	log.Printf("[INFO] new user created: id=%d username=%q email=%q\n",
		user.ID, user.Username, user.Email)

	// redirect to login page on success
	// set session values
	sess := sessions.Default(c)
	sess.Set("username", user.Username)
	sess.Set("user_id", user.ID)
	sess.Set("is_admin", user.IsAdmin)
	sess.Save()
	c.Redirect(http.StatusSeeOther, "/profile")
}

// LogoutHandler clears the session and redirects home
func LogoutHandler(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Clear()
	sess.Save()
	c.Redirect(http.StatusSeeOther, "/")
}
