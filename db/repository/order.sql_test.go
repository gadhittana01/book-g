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

func TestCheckOrderExists(t *testing.T) {
	mockDB, _ := pgxmock.NewPool()
	defer mockDB.Close()
	q := NewRepository(mockDB)

	orderID := uuid.New()

	t.Run("success query check order", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(checkOrderExists)).
			WithArgs(orderID).
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

		exists, err := q.CheckOrderExists(context.Background(), orderID)
		assert.NoError(t, err)
		assert.Equal(t, true, exists)
	})

	t.Run("failed query check order", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(checkOrderExists)).
			WithArgs(orderID).
			WillReturnError(errQuery)

		exists, err := q.CheckOrderExists(context.Background(), orderID)
		assert.Error(t, err)
		assert.Empty(t, exists)
	})
}

func TestCreateOrder(t *testing.T) {
	mockDB, _ := pgxmock.NewPool()
	defer mockDB.Close()
	q := NewRepository(mockDB)
	userID := uuid.New()
	totalPrice := float64(20)
	now := time.Now()
	status := "pending"

	req := CreateOrderParams{
		UserID:     userID,
		Date:       now,
		TotalPrice: totalPrice,
	}

	expected := Order{
		ID:         uuid.New(),
		UserID:     userID,
		Date:       now,
		TotalPrice: totalPrice,
		Status:     status,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	t.Run("success query create order", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(createOrder)).
			WithArgs(req.UserID, req.Date, req.TotalPrice).
			WillReturnRows(pgxmock.NewRows([]string{
				"id",
				"user_id",
				"date",
				"total_price",
				"status",
				"created_at",
				"updated_at"}).AddRow(
				expected.ID,
				expected.UserID,
				expected.Date,
				expected.TotalPrice,
				expected.Status,
				expected.CreatedAt,
				expected.UpdatedAt,
			))

		res, err := q.CreateOrder(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("failed query create order", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(createOrder)).
			WithArgs(req.UserID, req.Date, req.TotalPrice).
			WillReturnError(errQuery)

		res, err := q.CreateOrder(context.Background(), req)
		assert.Error(t, err)
		assert.Empty(t, res)
	})
}

func TestCreateOrderDetail(t *testing.T) {
	mockDB, _ := pgxmock.NewPool()
	defer mockDB.Close()
	q := NewRepository(mockDB)
	orderID := uuid.New()
	bookID := uuid.New()
	quantity := int32(20)
	now := time.Now()

	req := CreateOrderDetailParams{
		OrderID:  orderID,
		BookID:   bookID,
		Quantity: quantity,
	}

	expected := OrderDetail{
		ID:        uuid.New(),
		OrderID:   orderID,
		BookID:    bookID,
		Quantity:  quantity,
		CreatedAt: now,
		UpdatedAt: now,
	}

	t.Run("success query create order detail", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(createOrderDetail)).
			WithArgs(req.OrderID, req.BookID, req.Quantity).
			WillReturnRows(pgxmock.NewRows([]string{
				"id",
				"order_id",
				"book_id",
				"quantity",
				"created_at",
				"updated_at"}).AddRow(
				expected.ID,
				expected.OrderID,
				expected.BookID,
				expected.Quantity,
				expected.CreatedAt,
				expected.UpdatedAt,
			))

		res, err := q.CreateOrderDetail(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("failed query create order", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(createOrder)).
			WithArgs(req.OrderID, req.BookID, req.Quantity).
			WillReturnError(errQuery)

		res, err := q.CreateOrderDetail(context.Background(), req)
		assert.Error(t, err)
		assert.Empty(t, res)
	})
}

func TestFindOrderByID(t *testing.T) {
	mockDB, _ := pgxmock.NewPool()
	defer mockDB.Close()
	q := NewRepository(mockDB)
	userID := uuid.New()
	orderID := uuid.New()
	totalPrice := float64(20)
	now := time.Now()
	status := "pending"

	req := FindOrderByIDParams{
		ID:     orderID,
		UserID: userID,
	}

	expected := Order{
		ID:         uuid.New(),
		UserID:     userID,
		Date:       now,
		TotalPrice: totalPrice,
		Status:     status,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	t.Run("success query find order by ID", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(findOrderByID)).
			WithArgs(req.UserID, req.ID).
			WillReturnRows(pgxmock.NewRows([]string{
				"id",
				"user_id",
				"date",
				"total_price",
				"status",
				"created_at",
				"updated_at"}).AddRow(
				expected.ID,
				expected.UserID,
				expected.Date,
				expected.TotalPrice,
				expected.Status,
				expected.CreatedAt,
				expected.UpdatedAt,
			))

		res, err := q.FindOrderByID(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("failed query find order by ID", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(findOrderByID)).
			WithArgs(req).
			WillReturnError(errQuery)

		res, err := q.FindOrderByID(context.Background(), req)
		assert.Error(t, err)
		assert.Empty(t, res)
	})
}

func TestFindOrderByUserID(t *testing.T) {
	mockDB, _ := pgxmock.NewPool()
	defer mockDB.Close()
	q := NewRepository(mockDB)
	userID := uuid.New()
	totalPrice := float64(20)
	now := time.Now()
	status := "pending"

	req := FindOrderByUserIDParams{
		Limit:  10,
		Offset: 0,
		UserID: userID,
	}

	expected := []Order{
		{
			ID:         uuid.New(),
			UserID:     userID,
			Date:       now,
			TotalPrice: totalPrice,
			Status:     status,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
	}

	t.Run("success query find order by user ID", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(findOrderByUserID)).
			WithArgs(req.UserID, req.Limit, req.Offset).
			WillReturnRows(pgxmock.NewRows([]string{
				"id",
				"user_id",
				"date",
				"total_price",
				"status",
				"created_at",
				"updated_at"}).AddRow(
				expected[0].ID,
				expected[0].UserID,
				expected[0].Date,
				expected[0].TotalPrice,
				expected[0].Status,
				expected[0].CreatedAt,
				expected[0].UpdatedAt,
			))

		res, err := q.FindOrderByUserID(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("failed query find order by user ID", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(findOrderByUserID)).
			WithArgs(req.UserID, req.Limit, req.Offset).
			WillReturnError(errQuery)

		res, err := q.FindOrderByUserID(context.Background(), req)
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("failed scan find order by user ID", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(findOrderByUserID)).
			WithArgs(req.UserID, req.Limit, req.Offset).
			WillReturnRows(pgxmock.NewRows([]string{
				"id",
				"user_id",
				"date",
				"total_price",
				"status",
				"created_at",
				"updated_at"}).AddRow(
				1,
				expected[0].UserID,
				expected[0].Date,
				expected[0].TotalPrice,
				expected[0].Status,
				expected[0].CreatedAt,
				expected[0].UpdatedAt,
			))

		res, err := q.FindOrderByUserID(context.Background(), req)
		assert.Error(t, err)
		assert.Empty(t, res)
	})
}

func TestFindOrderDetailByOrderID(t *testing.T) {
	mockDB, _ := pgxmock.NewPool()
	defer mockDB.Close()
	q := NewRepository(mockDB)
	userID := uuid.New()
	orderID := uuid.New()
	bookID := uuid.New()
	totalPrice := float64(20)
	now := time.Now()
	title := "Hello"
	description := "World"
	author := "Giri Putra Adhittana"
	quantity := int32(5)
	price := float64(12)
	status := "pending"

	req := FindOrderDetailByOrderIDParams{
		ID:     orderID,
		UserID: userID,
	}

	expected := []FindOrderDetailByOrderIDRow{
		{
			ID:          uuid.New(),
			Date:        now,
			TotalPrice:  totalPrice,
			Status:      status,
			BookID:      bookID,
			Title:       title,
			Description: description,
			Author:      author,
			Quantity:    quantity,
			Price:       price,
		},
	}

	t.Run("success query find order detail by order ID", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(findOrderDetailByOrderID)).
			WithArgs(req.UserID, req.ID).
			WillReturnRows(pgxmock.NewRows([]string{
				"id",
				"date",
				"book_id",
				"title",
				"total_price",
				"status",
				"description",
				"author",
				"quantity",
				"price"}).AddRow(
				expected[0].ID,
				expected[0].Date,
				expected[0].BookID,
				expected[0].Title,
				expected[0].TotalPrice,
				expected[0].Status,
				expected[0].Description,
				expected[0].Author,
				expected[0].Quantity,
				expected[0].Price,
			))

		res, err := q.FindOrderDetailByOrderID(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("failed query find order detail by order ID", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(findOrderDetailByOrderID)).
			WithArgs(req.UserID, req.ID).
			WillReturnError(errQuery)

		res, err := q.FindOrderDetailByOrderID(context.Background(), req)
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("failed scan find order detail by order ID", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(findOrderDetailByOrderID)).
			WithArgs(req.UserID, req.ID).
			WillReturnRows(pgxmock.NewRows([]string{
				"id",
				"date",
				"book_id",
				"title",
				"total_price",
				"status",
				"description",
				"author",
				"quantity",
				"price"}).AddRow(
				1,
				expected[0].Date,
				expected[0].BookID,
				expected[0].Title,
				expected[0].TotalPrice,
				expected[0].Status,
				expected[0].Description,
				expected[0].Author,
				expected[0].Quantity,
				expected[0].Price,
			))

		res, err := q.FindOrderDetailByOrderID(context.Background(), req)
		assert.Error(t, err)
		assert.Empty(t, res)
	})
}

func TestGetOrderCountByUserId(t *testing.T) {
	mockDB, _ := pgxmock.NewPool()
	defer mockDB.Close()
	q := NewRepository(mockDB)

	req := uuid.New()
	expected := int64(10)

	t.Run("success query find order detail by order ID", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(getOrderCountByUserId)).
			WithArgs(req).
			WillReturnRows(pgxmock.NewRows([]string{
				"total",
			}).AddRow(
				expected,
			))

		res, err := q.GetOrderCountByUserId(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("failed query find order detail by order ID", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(getOrderCountByUserId)).
			WithArgs(req).
			WillReturnError(errQuery)

		res, err := q.GetOrderCountByUserId(context.Background(), req)
		assert.Error(t, err)
		assert.Empty(t, res)
	})
}

func TestUpdateOrderByID(t *testing.T) {
	mockDB, _ := pgxmock.NewPool()
	defer mockDB.Close()
	q := NewRepository(mockDB)
	orderID := uuid.New()
	totalPrice := float64(10)
	userID := uuid.New()
	now := time.Now()
	status := "pending"

	req := UpdateOrderByIDParams{
		ID:         orderID,
		TotalPrice: totalPrice,
	}

	expected := Order{
		ID:         orderID,
		UserID:     userID,
		Date:       now,
		TotalPrice: totalPrice,
		Status:     status,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	t.Run("success query update order by ID", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(updateOrderByID)).
			WithArgs(req.ID, req.TotalPrice).
			WillReturnRows(pgxmock.NewRows([]string{
				"id",
				"user_id",
				"date",
				"total_price",
				"status",
				"created_at",
				"updated_at",
			}).AddRow(
				expected.ID,
				expected.UserID,
				expected.Date,
				expected.TotalPrice,
				expected.Status,
				expected.CreatedAt,
				expected.UpdatedAt,
			))

		res, err := q.UpdateOrderByID(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("failed query update order order ID", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(updateOrderByID)).
			WithArgs(req).
			WillReturnError(errQuery)

		res, err := q.UpdateOrderByID(context.Background(), req)
		assert.Error(t, err)
		assert.Empty(t, res)
	})
}
