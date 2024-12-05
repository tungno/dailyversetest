/**
 *  Main entry point for the DailyVerse application. This file sets up the HTTP server,
 *  initializes services, repositories, and handlers, and defines routes for various endpoints.
 *
 *  @file      main.go
 *  @project   DailyVerse
 *  @framework Go HTTP Server
 *  @authors
 *      - Aayush
 *      - Tung
 *      - Boss
 *      - Majd
 */

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"proh2052-group6/internal/repositories"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"proh2052-group6/internal/handlers"
	"proh2052-group6/internal/middleware"
	"proh2052-group6/internal/services"
)

func main() {
	// Load environment variables from a .env file
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	// Create a context for service initialization
	ctx := context.Background()

	// Initialize Firestore client for database access
	dbClient, err := services.NewFirestoreClient(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize Firestore: %v", err)
	}
	defer dbClient.Close() // Ensure Firestore client is closed when the application exits

	// Initialize repositories for data access
	userRepository := repositories.NewFirestoreUserRepository(dbClient)
	friendRepository := repositories.NewFirestoreFriendRepository(dbClient)
	eventRepository := repositories.NewFirestoreEventRepository(dbClient)
	journalRepository := repositories.NewFirestoreJournalRepository(dbClient)

	// Initialize services for business logic
	emailService := services.NewSMTPEmailService()
	userService := services.NewUserService(userRepository, emailService)
	eventService := services.NewEventService(eventRepository)
	friendService := services.NewFriendService(userRepository, friendRepository)
	journalService := services.NewJournalService(journalRepository)
	newsService := services.NewNewsService(userRepository)
	profileService := services.NewProfileService(userRepository)
	cityService := services.NewCityService()
	timetableService := services.NewTimetableService(eventRepository)

	// Initialize HTTP handlers
	userHandler := handlers.NewUserHandler(userService)
	eventHandler := handlers.NewEventHandler(eventService)
	friendHandler := handlers.NewFriendHandler(friendService)
	journalHandler := handlers.NewJournalHandler(journalService)
	newsHandler := handlers.NewNewsHandler(newsService)
	profileHandler := handlers.NewProfileHandler(profileService)
	countryHandler := handlers.NewCountryHandler()
	cityHandler := handlers.NewCityHandler(cityService, userService)
	timetableHandler := handlers.NewTimetableHandler(timetableService)

	// Set up the HTTP router
	router := mux.NewRouter()

	// Define API routes
	// User routes
	router.Handle("/api/signup", middleware.RateLimitMiddleware(http.HandlerFunc(userHandler.Signup))).Methods("POST")
	router.Handle("/api/login", middleware.RateLimitMiddleware(http.HandlerFunc(userHandler.Login))).Methods("POST")
	router.Handle("/api/resend-otp", middleware.RateLimitMiddleware(http.HandlerFunc(userHandler.ResendOTP))).Methods("POST")
	router.HandleFunc("/api/verify-email", userHandler.VerifyEmail).Methods("POST")
	router.HandleFunc("/api/forgot-password", userHandler.ForgotPassword).Methods("POST")
	router.HandleFunc("/api/reset-password", userHandler.ResetPassword).Methods("POST")
	router.Handle("/api/me", middleware.JwtAuthMiddleware(userHandler.GetUserInfo)).Methods("GET")

	// Event routes
	router.Handle("/api/events/create", middleware.JwtAuthMiddleware(eventHandler.CreateEvent)).Methods("POST")
	router.Handle("/api/events/get", middleware.JwtAuthMiddleware(eventHandler.GetEvent)).Methods("GET")
	router.Handle("/api/events/update", middleware.JwtAuthMiddleware(eventHandler.UpdateEvent)).Methods("PUT")
	router.Handle("/api/events/delete", middleware.JwtAuthMiddleware(eventHandler.DeleteEvent)).Methods("DELETE")
	router.Handle("/api/events/all", middleware.JwtAuthMiddleware(eventHandler.GetAllEvents)).Methods("GET")

	// Friend routes
	router.Handle("/api/friends/add", middleware.JwtAuthMiddleware(friendHandler.SendFriendRequest)).Methods("POST")
	router.Handle("/api/friends/accept", middleware.JwtAuthMiddleware(friendHandler.AcceptFriendRequest)).Methods("POST")
	router.Handle("/api/friends/list", middleware.JwtAuthMiddleware(friendHandler.GetFriendsList)).Methods("GET")
	router.Handle("/api/friends/delete", middleware.JwtAuthMiddleware(friendHandler.RemoveFriend)).Methods("DELETE")
	router.Handle("/api/friends/requests", middleware.JwtAuthMiddleware(friendHandler.GetPendingFriendRequests)).Methods("GET")
	router.Handle("/api/friends/decline", middleware.JwtAuthMiddleware(friendHandler.DeclineFriendRequest)).Methods("POST")
	router.Handle("/api/friends/cancel", middleware.JwtAuthMiddleware(friendHandler.CancelFriendRequest)).Methods("POST")

	// User search
	router.Handle("/api/users/search", middleware.JwtAuthMiddleware(userHandler.SearchUsersByUsername)).Methods("GET")

	// Profile routes
	router.Handle("/api/profile", middleware.JwtAuthMiddleware(profileHandler.ProfileHandler)).Methods("GET", "PUT")

	// Country and city routes
	router.HandleFunc("/api/countries", countryHandler.GetCountries).Methods("GET")
	router.HandleFunc("/api/cities", cityHandler.GetCities).Methods("GET")

	// News route
	router.Handle("/api/news", middleware.JwtAuthMiddleware(newsHandler.FetchNews)).Methods("GET")

	// Journal routes
	router.Handle("/api/journal/save", middleware.JwtAuthMiddleware(journalHandler.CreateJournal)).Methods("POST")
	router.Handle("/api/journal", middleware.JwtAuthMiddleware(journalHandler.GetJournal)).Methods("GET")
	router.Handle("/api/journal/update", middleware.JwtAuthMiddleware(journalHandler.UpdateJournal)).Methods("PUT")
	router.Handle("/api/journal/delete", middleware.JwtAuthMiddleware(journalHandler.DeleteJournal)).Methods("DELETE")
	router.Handle("/api/journals", middleware.JwtAuthMiddleware(journalHandler.GetAllJournals)).Methods("GET")

	// Timetable route
	router.Handle("/api/import-ntnu-timetable", middleware.JwtAuthMiddleware(timetableHandler.ImportTimetable)).Methods("POST")

	// Apply CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins for development (adjust in production)
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	// Configure and start the HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	handler := c.Handler(router)
	srv := &http.Server{
		Handler:      handler,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Server running on port %s", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
