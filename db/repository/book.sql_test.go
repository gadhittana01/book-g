package querier

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
)

var errQuery = errors.New("error")

func TestCheckBookExists(t *testing.T) {
	mockDB, _ := pgxmock.NewPool()
	defer mockDB.Close()
	q := NewRepository(mockDB)

	bookID := uuid.New()

	t.Run("success query check book", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(checkBookExists)).
			WithArgs(bookID).
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

		exists, err := q.CheckBookExists(context.Background(), bookID)
		assert.NoError(t, err)
		assert.Equal(t, true, exists)
	})

	t.Run("failed query check book", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(checkBookExists)).
			WithArgs(bookID).
			WillReturnError(errQuery)

		exists, err := q.CheckBookExists(context.Background(), bookID)
		assert.Error(t, err)
		assert.Empty(t, exists)
	})
}

func TestCreateBook(t *testing.T) {
	mockDB, _ := pgxmock.NewPool()
	defer mockDB.Close()
	q := NewRepository(mockDB)
	title := "Hello"
	description := "World"
	author := "Giri Putra Adhittana"
	price := float64(20)
	now := time.Now()

	req := CreateBookParams{
		Title:       title,
		Description: description,
		Author:      author,
		Price:       price,
	}

	expected := Book{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
		Author:      author,
		Price:       price,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	t.Run("success query create book", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(createBook)).
			WithArgs(req.Title, req.Description, req.Author, req.Price).
			WillReturnRows(pgxmock.NewRows([]string{
				"id",
				"title",
				"description",
				"author",
				"price",
				"created_at",
				"updated_at"}).AddRow(
				expected.ID,
				expected.Title,
				expected.Description,
				expected.Author,
				expected.Price,
				expected.CreatedAt,
				expected.UpdatedAt,
			))

		res, err := q.CreateBook(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("failed query create book", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(createBook)).
			WithArgs(req.Title, req.Description, req.Author, req.Price).
			WillReturnError(errQuery)

		res, err := q.CreateBook(context.Background(), req)
		assert.Error(t, err)
		assert.Empty(t, res)
	})
}

func TestFindBook(t *testing.T) {
	mockDB, _ := pgxmock.NewPool()
	defer mockDB.Close()
	q := NewRepository(mockDB)
	title := "Hello"
	description := "World"
	author := "Giri Putra Adhittana"
	price := float64(20)
	now := time.Now()

	req := FindBookParams{
		Limit:  10,
		Offset: 0,
	}

	expected := []Book{
		{
			ID:          uuid.New(),
			Title:       title,
			Description: description,
			Author:      author,
			Price:       price,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	t.Run("success query find book", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(findBook)).
			WithArgs(req.Limit, req.Offset).
			WillReturnRows(pgxmock.NewRows([]string{
				"id",
				"title",
				"description",
				"author",
				"price",
				"created_at",
				"updated_at"}).AddRow(
				expected[0].ID,
				expected[0].Title,
				expected[0].Description,
				expected[0].Author,
				expected[0].Price,
				expected[0].CreatedAt,
				expected[0].UpdatedAt,
			))

		res, err := q.FindBook(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("failed query find book", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(findBook)).
			WithArgs(req.Limit, req.Offset).
			WillReturnError(errQuery)

		res, err := q.FindBook(context.Background(), req)
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("failed scan find book", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(findBook)).
			WithArgs(req.Limit, req.Offset).
			WillReturnRows(pgxmock.NewRows([]string{
				"id",
				"title",
				"description",
				"author",
				"price",
				"created_at",
				"updated_at"}).AddRow(
				1,
				expected[0].Title,
				expected[0].Description,
				expected[0].Author,
				expected[0].Price,
				expected[0].CreatedAt,
				expected[0].UpdatedAt,
			))

		res, err := q.FindBook(context.Background(), req)
		assert.Error(t, err)
		assert.Empty(t, res)
	})

}

func TestFindBookByID(t *testing.T) {
	mockDB, _ := pgxmock.NewPool()
	defer mockDB.Close()
	q := NewRepository(mockDB)
	title := "Hello"
	description := "World"
	author := "Giri Putra Adhittana"
	price := float64(20)
	now := time.Now()

	req := uuid.New()

	expected := Book{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
		Author:      author,
		Price:       price,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	t.Run("success query find book by ID", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(findBookByID)).
			WithArgs(req).
			WillReturnRows(pgxmock.NewRows([]string{
				"id",
				"title",
				"description",
				"author",
				"price",
				"created_at",
				"updated_at"}).AddRow(
				expected.ID,
				expected.Title,
				expected.Description,
				expected.Author,
				expected.Price,
				expected.CreatedAt,
				expected.UpdatedAt,
			))

		res, err := q.FindBookByID(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("failed query find book by ID", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(findBookByID)).
			WithArgs(req).
			WillReturnError(errQuery)

		res, err := q.FindBookByID(context.Background(), req)
		assert.Error(t, err)
		assert.Empty(t, res)
	})

}

func TestGetBookCount(t *testing.T) {
	mockDB, _ := pgxmock.NewPool()
	defer mockDB.Close()
	q := NewRepository(mockDB)

	expected := int64(10)

	t.Run("success get book count", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(getBookCount)).
			WillReturnRows(pgxmock.NewRows([]string{
				"total",
			}).AddRow(
				expected,
			))

		res, err := q.GetBookCount(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("failed query find book by ID", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(getBookCount)).
			WillReturnError(errQuery)

		res, err := q.GetBookCount(context.Background())
		assert.Error(t, err)
		assert.Empty(t, res)
	})

}
