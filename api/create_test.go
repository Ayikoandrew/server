package api

import (
	"os"
	"testing"

	"github.com/Ayikoandrew/server/database"
	"github.com/Ayikoandrew/server/types"
)

func Test_CreateAccount(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping database test in CI environment")
	}
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
		Email:       "test10@example.com",
		Password:    "password123",
	}

	if err := store.CreateAccount(testAccount); err != nil {
		t.Fatalf("Failed to create account: %v", err)
	}
}
