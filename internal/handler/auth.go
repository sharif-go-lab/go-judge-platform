package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginPageHandler handles the login page display
func LoginPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login",
	})
}

// LoginHandler handles the login form submission
func LoginHandler(c *gin.Context) {
	// For Phase 1, just redirect to the questions page
	// We'll implement actual authentication in Phase 2
	c.Redirect(http.StatusSeeOther, "/questions")
}

// RegisterPageHandler handles the registration page display
func RegisterPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{
		"title": "Register",
	})
}

// RegisterHandler handles the registration form submission
func RegisterHandler(c *gin.Context) {
	// For Phase 1, just redirect to the login page
	// We'll implement actual registration in Phase 2
	c.Redirect(http.StatusSeeOther, "/auth/login")
}

// LogoutHandler handles user logout
func LogoutHandler(c *gin.Context) {
	// For Phase 1, just redirect to the home page
	// We'll implement actual logout in Phase 2
	c.Redirect(http.StatusSeeOther, "/")
}
