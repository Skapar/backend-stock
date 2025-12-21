package handler

// =========================
// Common responses
// =========================

type ErrorResponse struct {
	Error string `json:"error" example:"invalid input"`
}

type MessageResponse struct {
	Message string `json:"message" example:"ok"`
}

type IDResponse struct {
	Message string `json:"message" example:"created"`
	ID      int64  `json:"id" example:"1"`
}

// =========================
// Auth
// =========================

type RegisterRequest struct {
	Email    string `json:"email" example:"test@mail.com"`
	Password string `json:"password" example:"123456"`
	Role     string `json:"role" example:"TRADER"`
}

type RegisterResponse struct {
	Message string `json:"message" example:"user created successfully"`
	UserID  int64  `json:"user_id" example:"1"`
}

type LoginRequest struct {
	Email    string `json:"email" example:"test@mail.com"`
	Password string `json:"password" example:"123456"`
}

type LoginResponse struct {
	Token     string `json:"token" example:"eyJhbGciOi..."`
	ExpiresIn int64  `json:"expiresIn" example:"1730000000"`
}

// =========================
// Users
// =========================

type GetMeResponse struct {
	Email   string  `json:"email" example:"test@mail.com"`
	Balance float64 `json:"balance" example:"1000"`
}

type UpdateUserRequest struct {
	Email    string  `json:"email" example:"new@mail.com"`
	Password string  `json:"password" example:"newpass123"`
	Role     string  `json:"role" example:"ADMIN"`
	Balance  float64 `json:"balance" example:"5000"`
}

// =========================
// Orders
// =========================

type UpdateOrderStatusRequest struct {
	Status string `json:"status" example:"COMPLETED"`
}

type OrderCreatedResponse struct {
	Message string `json:"message" example:"order executed successfully"`
	OrderID int64  `json:"order_id" example:"123"`
}

// =========================
// Portfolio
// =========================

type CreateOrUpdatePortfolioRequest struct {
	UserID   int64   `json:"user_id" example:"1"`
	StockID  int64   `json:"stock_id" example:"10"`
	Quantity float64 `json:"quantity" example:"2"`
}

// =========================
// History
// =========================

type HistoryCreatedResponse struct {
	HistoryID int64 `json:"history_id" example:"55"`
}
type CreateOrderRequest struct {
	StockID  int64  `json:"stock_id" example:"4"`
	Quantity int    `json:"quantity" example:"2"`
	Type     string `json:"type" example:"BUY"` // BUY или SELL
}
