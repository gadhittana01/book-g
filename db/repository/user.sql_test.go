package querier

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
)

func TestCheckEmailExists(t *testing.T) {
	mockDB, _ := pgxmock.NewPool()
	defer mockDB.Close()
	q := NewRepository(mockDB)

	req := "test@gmail.com"

	t.Run("success query check email exists", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(checkEmailExists)).
			WithArgs(req).
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

		exists, err := q.CheckEmailExists(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, true, exists)
	})

	t.Run("failed query check email exists", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(checkEmailExists)).
			WithArgs(req).
			WillReturnError(errQuery)

		exists, err := q.CheckEmailExists(context.Background(), req)
		assert.Error(t, err)
		assert.Empty(t, exists)
	})
}

func TestCreateUser(t *testing.T) {
	mockDB, _ := pgxmock.NewPool()
	defer mockDB.Close()
	q := NewRepository(mockDB)
	name := "Giri Putra Adhittana"
	email := "test@gmail.com"
	pwd := "123"
	now := time.Now()

	req := CreateUserParams{
		Name:     name,
		Email:    email,
		Password: pwd,
	}

	expected := User{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		Password:  pwd,
		CreatedAt: now,
		UpdatedAt: now,
	}

	t.Run("success query create user", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(createUser)).
			WithArgs(req.Name, req.Email, req.Password).
			WillReturnRows(pgxmock.NewRows([]string{
				"id",
				"name",
				"email",
				"password",
				"created_at",
				"updated_at",
			}).AddRow(
				expected.ID,
				expected.Name,
				expected.Email,
				expected.Password,
				expected.CreatedAt,
				expected.UpdatedAt,
			))

		res, err := q.CreateUser(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("failed query create user", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(createUser)).
			WithArgs(req.Name, req.Email, req.Password).
			WillReturnError(errQuery)

		res, err := q.CreateUser(context.Background(), req)
		assert.Error(t, err)
		assert.Empty(t, res)
	})
}

func TestFindUserByEmail(t *testing.T) {
	mockDB, _ := pgxmock.NewPool()
	defer mockDB.Close()
	q := NewRepository(mockDB)
	name := "Giri Putra Adhittana"
	email := "test@gmail.com"
	pwd := "123"
	now := time.Now()

	req := email

	expected := User{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		Password:  pwd,
		CreatedAt: now,
		UpdatedAt: now,
	}

	t.Run("success query find user by email", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(findUserByEmail)).
			WithArgs(req).
			WillReturnRows(pgxmock.NewRows([]string{
				"id",
				"name",
				"email",
				"password",
				"created_at",
				"updated_at",
			}).AddRow(
				expected.ID,
				expected.Name,
				expected.Email,
				expected.Password,
				expected.CreatedAt,
				expected.UpdatedAt,
			))

		res, err := q.FindUserByEmail(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("failed query find user by email", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(findUserByEmail)).
			WithArgs(req).
			WillReturnError(errQuery)

		res, err := q.FindUserByEmail(context.Background(), req)
		assert.Error(t, err)
		assert.Empty(t, res)
	})
}
