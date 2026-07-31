package domain

type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// ponytail: float64 money, switch to integer minor units if rounding ever matters
	Price float64 `json:"price"`
}
