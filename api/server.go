package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ayikoandrew/server/database"
	"github.com/Ayikoandrew/server/middleware"
	"github.com/Ayikoandrew/server/security"
	"github.com/Ayikoandrew/server/types"
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

	router.Use(middleware.LoggingMiddleware)

	router.HandleFunc("/signup",
		middleware.RateLimitMiddlewareTokenBucket(makeHTTPHandlerFunc(s.createAccount))).Methods(http.MethodPost)
	router.HandleFunc("/login",
		middleware.RateLimitMiddlewareTokenBucket(makeHTTPHandlerFunc(s.loginAccount))).Methods(http.MethodPost)
	router.HandleFunc("/health",
		makeHTTPHandlerFunc(s.handleHealth)).Methods(http.MethodGet)

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

	client := database.GetRedisClient()
	ctx := context.Background()
	userKey := fmt.Sprintf("user:%s:accessToken", user.User.ID)
	if err := client.Set(ctx, userKey, user.AccessToken, 15*time.Minute).Err(); err != nil {
		slog.Error("Failed to store user token in Redis", "error", err, "userId", user.User.ID)
	}

	token := &types.RefreshToken{
		UserID:       user.User.ID,
		RefreshToken: user.RefreshToken,
		ExpiresAt:    time.Now().Add(7 * 24 * 60 * 60),
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
		return err
	}

	return writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
