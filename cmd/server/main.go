package main

import (
	"flag"
	"fmt"
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
	// 1) Parse flags & load config
	listen := flag.String("listen", ":8080", "HTTP listen address")
	flag.Parse()
	config.Init()
	db.Init()

	// 2) Set up Gin and session store
	r := gin.Default()
	store := cookie.NewStore([]byte(viper.GetString("session.secret")))
	r.Use(sessions.Sessions("gojudge_session", store))

	// 3) Load currentUser into context if logged in
	r.Use(func(c *gin.Context) {
		sess := sessions.Default(c)
		if uid := sess.Get("user_id"); uid != nil {
			var u model.User
			if err := db.DB.First(&u, uid).Error; err == nil {
				c.Set("currentUser", &u)
			}
		}
		c.Next()
	})

	// 4) Template helpers
	r.SetFuncMap(template.FuncMap{
		"current_year": func() int { return time.Now().Year() },
		"lower":        func(s string) string { return strings.ToLower(s) },
	})

	// 5) Static files & templates
	r.Static("/static", "./static")

	// Load templates with proper inheritance
	templates := template.Must(template.New("").Funcs(template.FuncMap{
		"current_year": func() int { return time.Now().Year() },
		"lower":        func(s string) string { return strings.ToLower(s) },
	}).ParseGlob("templates/**/*.html"))
	r.SetHTMLTemplate(templates)

	// 6a) Public routes
	handler.RegisterPublicRoutes(r)

	// 6b) Protected routes (login required)
	handler.RegisterProtectedRoutes(r)

	// 6c) Admin routes (login + admin required)
	handler.RegisterAdminRoutes(r)

	// 7) Start server
	viper.SetDefault("server.listen", *listen)
	addr := viper.GetString("server.listen")
	fmt.Printf("Starting server on %s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
