package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Ayikoandrew/server/database"
	api "github.com/Ayikoandrew/server/functions"
	"github.com/Ayikoandrew/server/middleware"
	"github.com/Ayikoandrew/server/security"
	"github.com/Ayikoandrew/server/types"
	"github.com/Ayikoandrew/server/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type Server struct {
	listenAddr string
	store      database.DBHandler
}

func NewServer(listenAddr string, store database.DBHandler) *Server {
	return &Server{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *Server) Run() {
	router := mux.NewRouter()

	// registry := prometheus.NewRegistry()
	// promMiddleware := security.New(registry, nil)
	router.Use(middleware.LoggingMiddleware)
	// router.Use(func(h http.Handler) http.Handler {
	// 	return promMiddleware.WrapHandler("api", h)
	// })

	// router.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{})).Methods(http.MethodGet)

	router.Handle("/signup",
		middleware.RateLimitMiddlewareTokenBucket(makeHTTPHandlerFunc(s.createAccount))).Methods(http.MethodPost)
	router.Handle("/login",
		middleware.RateLimitMiddlewareTokenBucket(makeHTTPHandlerFunc(s.loginAccount))).Methods(http.MethodPost)
	router.Handle("/health",
		makeHTTPHandlerFunc(s.handleHealth)).Methods(http.MethodGet)
	router.Handle("/logout", middleware.RateLimitMiddlewareTokenBucket(makeHTTPHandlerFunc(s.logoutHandler))).Methods(http.MethodPost)

	router.Handle("/auth/refresh", makeHTTPHandlerFunc(s.refreshTokenHandler)).Methods(http.MethodPost)
	router.Handle("/expense", security.ValidateAccessTokenMiddleware(makeHTTPHandlerFunc(s.uploadExpenses))).Methods(http.MethodGet)
	router.Handle("/", security.ValidateAccessTokenMiddleware(makeHTTPHandlerFunc(s.retriveExpenses))).Methods(http.MethodGet)

	serve := &http.Server{
		Addr:         s.listenAddr,
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  20 * time.Second,
	}

	slog.Info("server is running on", "listendAddr", serve.Addr)

	errChan := make(chan error)
	go func() {
		if err := serve.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	select {
	case err := <-errChan:
		fmt.Printf("Error starting server %+v\n", err)
	case sig := <-sigChan:
		fmt.Printf("Received shutdown signal %+v\n", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := serve.Shutdown(ctx); err != nil {
		slog.Error("Error shutting down server", "error", err)
		return
	}
}

func (s *Server) loginAccount(w http.ResponseWriter, r *http.Request) error {
	account := new(types.LoginRequest)
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return nil
	}

	if account.Password == "" && account.Username == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return nil
	}

	user, err := s.store.Authenticate(account.Password, account.Username)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) || errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return nil
		}
		return err
	}

	database.Set(
		user.User.ID,
		user.AccessToken,
		30*time.Minute,
		context.Background(),
	)

	token := &types.RefreshToken{
		UserID:       user.User.ID,
		RefreshToken: user.RefreshToken,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	}

	s.store.StoreRefreshToken(token)

	security.SetTokenCookies(w, user.AccessToken, user.RefreshToken)
	return writeJSON(w, http.StatusOK, user)
}

func (s *Server) createAccount(w http.ResponseWriter, r *http.Request) error {
	account := new(types.Account)

	if err := json.NewDecoder(r.Body).Decode(account); err != nil {
		return err
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), 12)
	if err != nil {
		return err
	}

	account.Password = string(hashPassword)
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return writeJSON(w, http.StatusCreated, account)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) error {
	err := s.store.Ping()
	if err != nil {
		slog.Error("Health check failed", "error", err)
		return writeJSON(w, http.StatusServiceUnavailable, map[string]string{
			"status": "unhealthy",
			"error":  err.Error(),
		})
	}

	return writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) refreshTokenHandler(w http.ResponseWriter, r *http.Request) error {
	refreshSecret := os.Getenv("REFRESH_TOKEN")
	if refreshSecret == "" {
		log.Println("Warning: REFRESH_TOKEN environment variable is not set")
	}

	var refreshTokenValue string
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		refreshTokenValue = strings.TrimPrefix(authHeader, "Bearer ")
	} else {
		refreshToken, err := r.Cookie("refresh_token")
		if err != nil {
			if err == http.ErrNoCookie {
				return writeJSON(w, http.StatusBadRequest, "Token required")
			}
			return writeJSON(w, http.StatusBadRequest, "Invalid request")
		}
		refreshTokenValue = refreshToken.Value
	}

	if refreshTokenValue == "" {
		return writeJSON(w, http.StatusBadRequest, "Invalid token")
	}

	claims, err := s.validateRefreshToken(refreshTokenValue, refreshSecret)

	if err != nil {
		return writeJSON(w, http.StatusUnauthorized, "Validation failed!")
	}

	hashedToken := utils.HashToken(refreshTokenValue)
	userID, err := s.store.ValidateRefreshToken(hashedToken)

	if err != nil {
		return writeJSON(w, http.StatusUnauthorized, "Invalid refresh token")
	}

	newAccessToken, err := api.CreateAccessToken(claims.Subject)
	if err != nil {
		return writeJSON(w, http.StatusInternalServerError, "Failed to generate token")
	}

	database.Set(userID, newAccessToken, 30*time.Minute, context.Background())

	security.SetTokenCookies(w, newAccessToken, refreshTokenValue)

	return writeJSON(w, http.StatusCreated, &types.TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: refreshTokenValue,
	})

}

func (s *Server) logoutHandler(w http.ResponseWriter, r *http.Request) error {
	accessSecret := os.Getenv("ACCESS_TOKEN")

	accessToken, err := r.Cookie("access_token")
	if err != nil {
		if err == http.ErrNoCookie {
			return writeJSON(w, http.StatusBadRequest, "Token required")
		}
		return writeJSON(w, http.StatusBadRequest, "Invalid request")
	}

	claims, err := s.validateAccessToken(accessToken.Value, accessSecret)
	if err != nil {
		return writeJSON(w, http.StatusUnauthorized, "Validation failed!")
	}

	userID := claims.Subject

	if err := s.store.RevokeAllUserTokens(userID); err != nil {
		return writeJSON(w, http.StatusInternalServerError, "Logout failed")
	}

	database.Delete(userID, context.Background())

	security.ClearTokenCookies(w)

	return writeJSON(w, http.StatusOK, map[string]string{
		"message": "Successfully logged out",
	})

}

func (s *Server) validateRefreshToken(tokenString string, refreshToken string) (*types.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &types.CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(refreshToken), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token expired")
		}
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*types.CustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func (s *Server) validateAccessToken(tokenString string, accessToken string) (*types.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &types.CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(accessToken), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token expired")
		}
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*types.CustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func (s *Server) StartTokenCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		for range ticker.C {
			if err := s.store.CleanupExpiredTokens(); err != nil {
				log.Printf("Token cleanup failed: %v", err)
			}
		}
	}()
}
