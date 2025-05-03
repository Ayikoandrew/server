package api

import (
	"os"
	"testing"

	"github.com/Ayikoandrew/server/database"
	"github.com/Ayikoandrew/server/types"
)

func Test_CreateAccount(t *testing.T) {
	os.Setenv("DATABASE_URL", "postgresql://postgres:postgres@localhost:5432/liora?sslmode=disable")
	// url := os.Getenv("TEST_DATABASE_URL")
	store := database.NewStorage()

	if err := store.Init(); err != nil {
		t.Fatalf("Failed to initialize database %s", err)
	}

	defer store.Close()

	testAccount := &types.Account{
		FirstName:   "Test",
		LastName:    "User",
		PhoneNumber: "1234567890",
		Email:       "test1@example.com",
		Password:    []byte("password123"),
	}

	if err := store.CreateAccount(testAccount); err != nil {
		t.Fatalf("Failed to create account: %v", err)
	}
}
