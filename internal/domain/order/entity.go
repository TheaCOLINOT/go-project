package order

import "time"

type Order struct {
	ID               string
	BusinessID       string
	UserID           string
	CreatedByAdminID *string
	Status           string
	TotalTTCCents    int64
	CancellationReason *string
	CreatedAt        time.Time
}

type OrderItem struct {
	ID                   string
	OrderID              string
	ProductID            string
	ProductNameSnapshot  string
	Quantity             int
	UnitPriceHTCents     int64
	VATRateBasisPoints   int
	UnitPriceTTCCents    int64
}

type OrderStatusHistory struct {
	ID              string
	OrderID         string
	ChangedByUserID *string
	OldStatus       *string
	NewStatus       string
	Reason          *string
	ChangedAt       time.Time
}
