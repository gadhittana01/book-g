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
	OrderId string `json:"orderId"`
	Date    string `json:"date"`
}

type GetOrderRes struct {
	OrderId string `json:"orderId"`
	Date    string `json:"date"`
}

type GetOrderDetailRes struct {
	OrderId     string        `json:"orderId"`
	Date        string        `json:"date"`
	OrderDetail []OrderDetail `json:"orderDetail"`
}
