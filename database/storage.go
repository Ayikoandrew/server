package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Ayikoandrew/server/security"
	"github.com/Ayikoandrew/server/types"
	_ "github.com/jackc/pgx/v4/stdlib"
	"golang.org/x/crypto/bcrypt"
)

type Storage struct {
	db *sql.DB
}

func NewStorage() DBHandler {
	dns := os.Getenv("DATABASE_URL")

	if dns == "" {
		slog.Error("DATABASE_URL environment variable is not seen")
	}
	db, err := sql.Open("pgx", dns)
	if err != nil {
		slog.Error("Error opening pgx")

	}

	err = db.Ping()
	if err != nil {
		slog.Error("Error pinging the database")
	}

	db.SetMaxOpenConns(25)
	db.SetConnMaxIdleTime(10)
	db.SetConnMaxLifetime(5 * time.Second)

	return &Storage{
		db: db,
	}
}

func (s *Storage) Init() error {
	return s.CreateTable()
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) Ping() error {
	return s.db.Ping()
}

func (s *Storage) CreateTable() error {
	query := `
	-- Create tables if they don't exist
	CREATE TABLE IF NOT EXISTS users(
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					firstName VARCHAR(255),
					lastName VARCHAR(255),
					phoneNumber VARCHAR(20),
					email VARCHAR(255) UNIQUE,
					passwordHash BYTEA,
					createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	CREATE TABLE IF NOT EXISTS user_sessions (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        user_id UUID NOT NULL,
        refresh_token TEXT NOT NULL,
        device_info TEXT,
        expires_at TIMESTAMPTZ NOT NULL,
        revoked BOOLEAN DEFAULT FALSE,
        created_at TIMESTAMPTZ DEFAULT NOW (),
        FOREIGN KEY (user_id) REFERENCES users (id)
    );

	CREATE INDEX IF NOT EXISTS idx_refresh_token ON user_sessions (refresh_token);
	`

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(query); err != nil {
		slog.Error("Error executing schema creation", "error", err)
		return fmt.Errorf("error creating database schema: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit schema transaction: %w", err)
	}

	slog.Info("Database schema created successfully")
	return nil
}

func (s *Storage) CreateAccount(account *types.Account) error {

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	query := `INSERT INTO users 
	(firstName, lastName, phoneNumber, email, passwordHash) 
	VALUEs ($1, $2, $3, $4, $5) RETURNING id`

	var id string
	err = tx.QueryRow(query,
		account.FirstName,
		account.LastName,
		account.PhoneNumber,
		account.Email,
		account.Password,
	).Scan(&id)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	slog.Info("Account created successfully", "id", id)
	return nil
}

func (s *Storage) Authenticate(password, email string) (types.LoginResponse, error) {
	query := `SELECT id, firstName, lastName, phoneNumber, email, passwordhash FROM users
	WHERE email=$1`
	var user types.User
	err := s.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.Email,
		&user.Password,
	)
	if err != nil {
		return types.LoginResponse{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return types.LoginResponse{}, err
	}

	accessToken, err := security.CreateAccessToken(&user)
	if err != nil {
		return types.LoginResponse{}, err
	}

	refreshToken, err := security.CreateRefreshToken(&user)
	if err != nil {
		return types.LoginResponse{}, err
	}

	return types.LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
