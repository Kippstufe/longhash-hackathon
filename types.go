package main

type responseOrders struct {
	ID      int           `json:"id"`
	Jsonrpc string        `json:"jsonrpc"`
	Result  []transaction `json:"result"`
}

type transaction struct {
	Sequence       int    `json:"sequence"`
	Date           string `json:"date"`
	Price          string `json:"price"`
	Amount         string `json:"amount"`
	Value          string `json:"value"`
	Side1AccountID string `json:"side1_account_id"`
	Side2AccountID string `json:"side2_account_id"`
}
