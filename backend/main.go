package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/chaitin/SafeLine/docs"
	"github.com/chaitin/SafeLine/internal/handler"
	"github.com/chaitin/SafeLine/internal/service"
)

func main() {
	viper.AutomaticEnv()

	// variables that must be set
	viper.SetDefault("GITHUB_TOKEN", "")

	// optional variables to set
	viper.SetDefault("GITHUB_CACHE_TTL", 10) // cache timeout in minutes
	viper.SetDefault("LISTEN_ADDR", ":8080") // api server addr

	githubToken := viper.GetString("GITHUB_TOKEN")
	if githubToken == "" {
		log.Fatal("GITHUB_TOKEN must be set")
	}
	listenAddr := viper.GetString("LISTEN_ADDR")

	r := gin.Default()

	// Initialize the GitHubService.
	gitHubService := service.NewGitHubService(&service.GithubConfig{
		Token:    githubToken,
		Owner:    "chaitin",
		Repo:     "SafeLine",
		CacheTTL: viper.GetDuration("GITHUB_CACHE_TTL"),
	})

	// Create a new instance of GitHubHandler.
	gitHubHandler := handler.NewGitHubHandler(gitHubService)

	v1 := r.Group("/api")
	v1.GET("/repos/issues", gitHubHandler.GetIssues)
	v1.GET("/repos/discussions", gitHubHandler.GetDiscussions)

	docs.SwaggerInfo.BasePath = v1.BasePath()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Start the Gin server on the specified port.
	if err := r.Run(listenAddr); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
