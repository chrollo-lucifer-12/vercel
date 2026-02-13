package tests

import (
	"context"
	"testing"
	"time"

	"github.com/chrollo-lucifer-12/api-server/auth"
	"github.com/chrollo-lucifer-12/shared/db"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"gotest.tools/v3/assert"
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
	assert.NilError(t, err)

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	assert.NilError(t, err)
	t.Log(connStr)

	dbClient, err := db.NewTestDB(connStr, ctx)
	assert.NilError(t, err)
	t.Log("connected to db")

	assert.NilError(t, dbClient.MigrateDB())

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
		user := &db.User{
			Name:     "test",
			Email:    "test@example.com",
			Password: "password123",
		}
		err := service.CreateUser(ctx, user)
		assert.NilError(t, err, "creating user should not fail")
		assert.Assert(t, user.ID != uuid.Nil, "user ID should be generated")
		createdUser = user
	})

	t.Run("CreateUser duplicate email", func(t *testing.T) {
		user := &db.User{
			Name:     "test",
			Email:    "test@example.com",
			Password: "password123",
		}
		err := service.CreateUser(ctx, user)
		assert.ErrorContains(t, err, "duplicate", "should fail due to duplicate email")
	})

	t.Run("GetUser existing", func(t *testing.T) {
		u, err := service.GetUser(ctx, createdUser.Email)
		assert.NilError(t, err)
		assert.Equal(t, u.Email, createdUser.Email)
	})

	t.Run("GetUser non-existing", func(t *testing.T) {
		_, err := service.GetUser(ctx, "doesnotexist@example.com")
		assert.ErrorContains(t, err, "not found")
	})

	t.Run("UpdateUser", func(t *testing.T) {
		createdUser.Name = "updatedName"
		err := service.UpdateUser(ctx, *createdUser)
		assert.NilError(t, err)
		assert.Equal(t, createdUser.Name, "updatedName")
	})

	t.Run("DeleteUser", func(t *testing.T) {
		err := service.DeleteUser(ctx, createdUser.ID)
		assert.NilError(t, err)
	})
}

func TestAuthStore(t *testing.T) {
	service, ctx, cleanup := setupAuthService(t)
	defer cleanup()

	user := &db.User{
		Name:     "authuser",
		Email:    "authuser@example.com",
		Password: "password123",
	}
	assert.NilError(t, service.CreateUser(ctx, user))

	t.Run("CreateSession", func(t *testing.T) {
		session := &db.Session{
			UserID:       user.ID,
			RefreshToken: "token123",
			Revoked:      false,
			ExpiresAt:    time.Now().Add(time.Hour),
		}
		err := service.CreateSession(ctx, session)
		assert.NilError(t, err)
		assert.Equal(t, session.UserID, user.ID)
	})

	t.Run("GetSession", func(t *testing.T) {
		sess, err := service.GetSession(ctx, user.ID)
		assert.NilError(t, err)
		assert.Equal(t, sess.UserID, user.ID)
	})

	t.Run("RevokeSession", func(t *testing.T) {
		sess, err := service.GetSession(ctx, user.ID)
		assert.NilError(t, err)

		sess.Revoked = true
		err = service.RevokeSession(ctx, *sess)
		assert.NilError(t, err)

		updated, err := service.GetSession(ctx, user.ID)
		assert.NilError(t, err)
		assert.Assert(t, updated.Revoked, "session should be revoked")
	})

	t.Run("DeleteSession", func(t *testing.T) {
		err := service.DeleteSession(ctx, user.ID)
		assert.NilError(t, err)

		_, err = service.GetSession(ctx, user.ID)
		assert.ErrorContains(t, err, "not found")
	})
}
func TestJWTMaker(t *testing.T) {
	secretKey := "supersecret"
	maker := auth.NewJWTMaker(secretKey)

	userID := uuid.New()
	email := "test@example.com"
	duration := time.Minute * 1

	t.Run("CreateToken", func(t *testing.T) {
		tokenStr, claims, err := maker.CreateToken(userID, email, duration)
		assert.NilError(t, err)
		assert.Assert(t, tokenStr != "")
		assert.Equal(t, claims.ID, userID)
		assert.Equal(t, claims.Email, email)
	})

	t.Run("VerifyToken", func(t *testing.T) {
		tokenStr, _, err := maker.CreateToken(userID, email, duration)
		assert.NilError(t, err)

		claims, err := maker.VerifyToken(tokenStr)
		assert.NilError(t, err)
		assert.Equal(t, claims.ID, userID)
		assert.Equal(t, claims.Email, email)
	})

	t.Run("ExpiredToken", func(t *testing.T) {
		shortDuration := time.Millisecond * 50
		tokenStr, _, err := maker.CreateToken(userID, email, shortDuration)
		assert.NilError(t, err)

		time.Sleep(time.Millisecond * 60)

		_, err = maker.VerifyToken(tokenStr)
		assert.ErrorContains(t, err, "token is expired")
	})

	t.Run("InvalidSignature", func(t *testing.T) {
		tokenStr, _, err := maker.CreateToken(userID, email, duration)
		assert.NilError(t, err)

		otherMaker := auth.NewJWTMaker("wrongsecret")
		_, err = otherMaker.VerifyToken(tokenStr)
		assert.ErrorContains(t, err, "signature is invalid")
	})

	t.Run("InvalidTokenString", func(t *testing.T) {
		_, err := maker.VerifyToken("not_a_real_token")
		assert.ErrorContains(t, err, "token contains an invalid number of segments")
	})
}
