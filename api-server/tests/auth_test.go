package tests

import (
	"context"
	"testing"

	"github.com/chrollo-lucifer-12/api-server/auth"
	"github.com/chrollo-lucifer-12/shared/db"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func SetupTestPostgres(t *testing.T) (*db.DB, func()) {
	ctx := context.Background()

	dbName := "test"
	dbUser := "user"
	dbPassword := "password"

	container, err := postgres.Run(ctx,
		"postgres:16",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
	)

	if err != nil {
		t.Fatal(err)
	}

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(connStr)

	dbClient, err := db.NewTestDB(connStr, ctx)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("connected to db")

	if err := dbClient.MigrateDB(); err != nil {
		t.Fatal(err)
	}

	t.Log("connected to db")

	cleanup := func() {
		container.Terminate(ctx)
	}

	return dbClient, cleanup
}

func setupAuthService(t *testing.T) (*auth.AuthService, context.Context, func()) {
	ctx := context.Background()

	dbClient, cleanup := SetupTestPostgres(t)

	service := auth.NewAuthService(
		auth.UserStoreFuncs{
			CreateUserFn: dbClient.CreateUser,
			GetUserFn:    dbClient.GetUser,
			UpdateUserFn: dbClient.UpdateUser,
			DeleteUserFn: dbClient.DeleteUser,
		},
		auth.TokenStoreFuncs{
			CreateSessionFn: dbClient.CreateSession,
			GetSessionFn:    dbClient.GetSession,
			RevokeSessionFn: dbClient.RevokeSession,
			DeleteSessionFn: dbClient.DeleteSession,
		},
	)

	return service, ctx, cleanup
}

func TestUserStore(t *testing.T) {
	service, ctx, cleanup := setupAuthService(t)
	defer cleanup()

	var createdUser *db.User

	t.Run("CreateUser", func(t *testing.T) {
		tests := []struct {
			name        string
			user        *db.User
			expectError bool
		}{
			{
				name: "valid user",
				user: &db.User{
					Name:     "test",
					Email:    "test@example.com",
					Password: "password123",
				},
				expectError: false,
			},
			{
				name: "repeating email",
				user: &db.User{
					Name:     "test",
					Email:    "test@example.com",
					Password: "password123",
				},
				expectError: true,
			},
			{
				name: "missing email",
				user: &db.User{
					Name:     "test1",
					Password: "password123",
				},
				expectError: true,
			},
			{
				name: "missing password",
				user: &db.User{
					Name:  "test2",
					Email: "test2@example.com",
				},
				expectError: true,
			},
			{
				name: "missing name",
				user: &db.User{
					Email:    "test3@example.com",
					Password: "password123",
				},
				expectError: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := service.CreateUser(ctx, tt.user)

				if tt.expectError {
					if err == nil {
						t.Fatalf("expected error but got nil")
					}
					return
				}

				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if tt.user.ID.String() == "" {
					t.Fatalf("expected user ID to be generated")
				}

				if createdUser == nil {
					createdUser = tt.user
				}
			})
		}
	})

	t.Run("GetUser", func(t *testing.T) {
		if createdUser == nil {
			t.Fatal("no user created in CreateUser subtest")
		}

		tests := []struct {
			name        string
			email       string
			expectError bool
		}{
			{
				name:        "existing user",
				email:       createdUser.Email,
				expectError: false,
			},
			{
				name:        "nonexistent user",
				email:       "doesnotexist@example.com",
				expectError: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				u, err := service.GetUser(ctx, tt.email)

				if tt.expectError {
					if err == nil {
						t.Fatalf("expected error but got nil")
					}
					return
				}

				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if u.Email != tt.email {
					t.Fatalf("expected email %s, got %s", tt.email, u.Email)
				}
			})
		}
	})

	t.Run("UpdateUser", func(t *testing.T) {
		if createdUser == nil {
			t.Fatal("no user created in CreateUser subtest")
		}

		tests := []struct {
			name        string
			updatedName string
			expectError bool
		}{
			{
				name:        "valid name",
				updatedName: "random",
				expectError: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				createdUser.Name = tt.updatedName
				err := service.UpdateUser(ctx, *createdUser)

				if tt.expectError {
					if err == nil {
						t.Fatalf("expected error but got nil")
					}
					return
				}

				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if createdUser.Name != tt.updatedName {
					t.Fatalf("name doest not match")
				}
			})
		}
	})

	t.Run("DeleteUser", func(t *testing.T) {
		if createdUser == nil {
			t.Fatal("no user created in CreateUser subtest")
		}

		tests := []struct {
			name        string
			id          uuid.UUID
			expectError bool
		}{

			{
				name:        "valid id",
				id:          createdUser.ID,
				expectError: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := service.DeleteUser(ctx, tt.id)

				if tt.expectError {
					if err == nil {
						t.Fatalf("expected error but got nil")
					}
					return
				}

				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			})
		}
	})
}
