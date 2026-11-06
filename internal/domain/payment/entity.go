package payment

import "time"

type Payment struct {
	ID                 string
	OrderID            string
	Status             string
	Method             string
	AmountCents        int64
	SimulatedReference *string
	CreatedAt          time.Time
}
