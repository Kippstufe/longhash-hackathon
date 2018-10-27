package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

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

func main() {

	url := "https://apihk.cybex.io"
	reqParams := requestParams{currencyOrigin: "JADE.ETH",
		currencyTarget:        "JADE.BTC",
		dateStart:             time.Now().Format("2006-01-02T15:04:05"),
		dateStop:              time.Now().Add(-100 * time.Hour).Format("2006-01-02T15:04:05"),
		MaxNumberTransactions: 100,
	}

	requestBody := fmt.Sprintf("{\"params\":[\"%s\",\"%s\",\"%s\",\"%s\",%d],\"id\":1,\"jsonrpc\":\"2.0\",\"method\":\"get_trade_history\"}", reqParams.currencyOrigin, reqParams.currencyTarget, reqParams.dateStart, reqParams.dateStop, reqParams.MaxNumberTransactions)

	// fmt.Println(requestBody)
	payload := strings.NewReader(requestBody)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		fmt.Println("Error creating request", err)

	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error making http request", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error closing http request", err)
	}

	// fmt.Println(res)
	// fmt.Println(string(body))

	// Index a transaction (using JSON serialization)
	orders := responseOrders{}

	err = json.Unmarshal(body, &orders)

	if err != nil {
		fmt.Println("Not possible to unmarshal json response", err)
	}

	for _, currentTransaction := range orders.Result {
		stringTransaction, err := json.Marshal(currentTransaction)
		if err != nil {
			fmt.Println("error marshaling struct", err)
		}
		// fmt.Printf("%+v \n", string(stringTransaction))

		client := &http.Client{}
		newID := currentTransaction.Date + currentTransaction.Side1AccountID
		elasticSearchURL := fmt.Sprintf("http://localhost:9200/transactions/_doc/%s", newID)
		request, err := http.NewRequest("POST", elasticSearchURL, strings.NewReader(string(stringTransaction)))
		if err != nil {
			fmt.Println("Error creating request", err)
		}

		request.Header.Set("Content-Type", "application/json")

		res, err := client.Do(request)
		if err != nil {
			fmt.Println("Error making request", err)
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Error closing http request", err)
		}
		fmt.Println(string(body))
	}

}
