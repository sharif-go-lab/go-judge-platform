package main

import (
	"html/template"
	"log"
	"strings"
	"time"

	"github.com/sharif-go-lab/go-judge-platform/internal/model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/sharif-go-lab/go-judge-platform/internal/config"
	"github.com/sharif-go-lab/go-judge-platform/internal/db"
	"github.com/sharif-go-lab/go-judge-platform/internal/handler"
)

func main() {
	config.Init()
	db.Init()

	r := gin.Default()
	store := cookie.NewStore([]byte(viper.GetString("session.secret")))
	r.Use(sessions.Sessions("gojudge_session", store))

	r.Use(func(c *gin.Context) {
		if uid := sessions.Default(c).Get("user_id"); uid != nil {
			var u model.User
			if err := db.DB.First(&u, uid).Error; err == nil {
				c.Set("currentUser", &u)
			}
		}
		c.Next()
	})

	r.SetFuncMap(template.FuncMap{
		"current_year": func() int { return time.Now().Year() },
		"lower":        func(s string) string { return strings.ToLower(s) },
	})

	r.Static("/static", "./static")

	templates := template.Must(template.New("").Funcs(template.FuncMap{
		"current_year": func() int { return time.Now().Year() },
		"lower":        func(s string) string { return strings.ToLower(s) },
	}).ParseGlob("templates/**/*.html"))
	r.SetHTMLTemplate(templates)

	handler.RegisterPublicRoutes(r)
	handler.RegisterProtectedRoutes(r)
	handler.RegisterAdminRoutes(r)

	addr := viper.GetString("server.listen")
	log.Printf("Starting server on %s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
