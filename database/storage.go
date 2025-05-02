package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Ayikoandrew/server/types"
	_ "github.com/jackc/pgx/v4/stdlib"
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
	query := `CREATE TABLE IF NOT EXISTS users(
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					firstName VARCHAR(255),
					lastName VARCHAR(255),
					phoneNumber VARCHAR(20),
					email VARCHAR(255) UNIQUE,
					passwordHash BYTEA,
					createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`

	if _, err := s.db.Exec(query); err != nil {
		slog.Error("Error executing the query")
		return fmt.Errorf("error executing the query: %w", err)
	}
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
