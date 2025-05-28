package types

type Expense struct {
	Amount      float64 `json:"amount"`
	Date        string  `json:"date"`
	Description string  `json:"description"`
}

type PaymentMethods struct {
	Name string `json:"name"`
}

type Category struct {
	Name   string  `json:"name"`
	Budget float64 `json:"budget"`
}
