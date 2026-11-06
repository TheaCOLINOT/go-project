package cart

import "time"

type Cart struct {
	ID        string
	BusinessID string
	UserID    string
	UpdatedAt time.Time
}

type CartItem struct {
	ID        string
	BusinessID string
	CartID    string
	ProductID string
	Quantity  int
}
