package dto

type SignUpRes struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Token    string `json:"token"`
	ExpToken int64  `json:"expToken"`
}

type SignInRes struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Token    string `json:"token"`
	ExpToken int64  `json:"expToken"`
}

type OrderDetail struct {
	OrderDetailID string `json:"orderDetailId"`
	BookID        string `json:"bookId"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Author        string `json:"author"`
	Quantity      int    `json:"quantity"`
}

type CreateOrderRes struct {
	OrderId    string  `json:"orderId"`
	Date       string  `json:"date"`
	TotalPrice float64 `json:"totalPrice"`
	Status     string  `json:"status"`
}

type GetOrderRes struct {
	OrderId    string  `json:"orderId"`
	Date       string  `json:"date"`
	TotalPrice float64 `json:"totalPrice"`
	Status     string  `json:"status"`
}

type GetOrderDetailRes struct {
	OrderId     string        `json:"orderId"`
	Date        string        `json:"date"`
	TotalPrice  float64       `json:"totalPrice"`
	Status      string        `json:"status"`
	OrderDetail []OrderDetail `json:"orderDetail"`
}

type CreateBookRes struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Author      string  `json:"author"`
	Price       float64 `json:"price"`
}

type GetBookRes struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Author      string  `json:"author"`
	Price       float64 `json:"price"`
}
