package main

import (
	"flag"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/sharif-go-lab/go-judge-platform/internal/config"
	"github.com/sharif-go-lab/go-judge-platform/internal/db"
	"github.com/sharif-go-lab/go-judge-platform/internal/handler"
	"github.com/spf13/viper"
	"html/template"
	"log"
	"strings"
	"time"
)

func main() {

	// Parse command line flags
	listenAddr := flag.String("listen", ":8080", "HTTP listen address")
	flag.Parse()

	// Initialize configuration
	config.Init()
	db.Init()

	// session middleware
	store := cookie.NewStore([]byte(viper.GetString("session.secret")))

	// load templates and routes as before...

	viper.SetDefault("server.listen", *listenAddr)

	// Set up gin router
	r := gin.Default()
	r.Use(sessions.Sessions("go-judge_session", store))
	r.SetFuncMap(template.FuncMap{
		"current_year": func() int {
			return time.Now().Year()
		},
		"lower": func(s string) string {
			return strings.ToLower(s)
		},
		// add more helpers here as needed…
	})
	// Static files
	r.Static("/static", "./static")

	// Load HTML templates
	r.LoadHTMLGlob("templates/**/*")

	// Register routes
	handler.RegisterRoutes(r)

	// Start HTTP server
	serverAddr := viper.GetString("server.listen")
	fmt.Printf("Starting server on %s\n", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
