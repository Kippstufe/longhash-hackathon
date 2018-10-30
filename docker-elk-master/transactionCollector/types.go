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
	CurrencyOrigin string `json:"currency_origin"`
	CurrencyTarget string `json:"currency_target"`
}

// param1: string name(or id) of the first asset

//  param2: string name(or id) of the second asset

//  param3: start time as a UNIX timestamp

//  param4: stop time as a UNIX timestamp

//  param5: number of trasactions to retrieve, capped at 100
type requestParams struct {
	currencyOrigin        string
	currencyTarget        string
	dateStart             string
	dateStop              string
	MaxNumberTransactions int
}

type channelStruct struct {
	orders responseOrders
	params requestParams
}
