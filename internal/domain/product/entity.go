package product

type Category struct {
	ID          string
	Name        string
	Description *string
}

type Product struct {
	ID                 string
	BusinessID         string
	CategoryID         string
	Name               string
	Description        *string
	PriceHTCents       int64
	VATRateBasisPoints int
	Stock              int
	IsActive           bool
}
