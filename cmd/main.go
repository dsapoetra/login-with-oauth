package main

import (
	"database/sql"
	"log"
	"login-with-oauth/internal/configs"
	"login-with-oauth/internal/database"
	"login-with-oauth/internal/handlers"
	"login-with-oauth/internal/logger"
	"login-with-oauth/internal/repository"
	"login-with-oauth/internal/services"
	"net/http"

	"github.com/spf13/viper"
)

func main() {
	// Initialize Viper across the application
	configs.InitializeViper()

	// Initialize Logger across the application
	logger.InitializeZapCustomLogger()

	// Initialize Database
	db, err := sql.Open("postgres", viper.GetString("database.dsn"))
	if err != nil {
		logger.Log.Fatal("Failed to connect to database:" + err.Error())
	}
	defer db.Close()

	// Initialize Repository
	userRepo := repository.NewUserRepository(db)

	// Initialize Oauth2 Services
	googleService := services.NewGoogleService(
		viper.GetString("google.clientID"),
		viper.GetString("google.clientSecret"),
		userRepo,
	)

	githubService := services.NewGitHubService(
		viper.GetString("github.clientID"),
		viper.GetString("github.clientSecret"),
		userRepo,
	)
	// Initialize Oauth2 Services
	googleHandler := handlers.NewGoogleHandler(googleService)

	// Initialize Handlers
	authHandler := handlers.NewAuthHandler(githubService)

	if err := database.RunMigrations(); err != nil {
		logger.Log.Fatal("Failed to run migrations:" + err.Error())
	}

	// Routes for the application
	http.HandleFunc("/", services.HandleMain)
	http.HandleFunc("/login-gl", googleHandler.GoogleLogin)
	http.HandleFunc("/callback-gl", googleHandler.GoogleCallback)
	http.HandleFunc("/login-gh", authHandler.GitHubLogin)
	http.HandleFunc("/gh-cb", authHandler.GitHubCallback)

	logger.Log.Info("Started running on http://localhost:" + viper.GetString("port"))
	log.Fatal(http.ListenAndServe(":"+viper.GetString("port"), nil))
}
