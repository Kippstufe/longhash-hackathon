package main

import (
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
		dateStop:              time.Now().Add(-10 * time.Hour).Format("2006-01-02T15:04:05"),
		MaxNumberTransactions: 100,
	}

	requestBody := fmt.Sprintf("{\"params\":[\"%s\",\"%s\",\"%s\",\"%s\",%d],\"id\":1,\"jsonrpc\":\"2.0\",\"method\":\"get_trade_history\"}", reqParams.currencyOrigin, reqParams.currencyTarget, reqParams.dateStart, reqParams.dateStop, reqParams.MaxNumberTransactions)

	fmt.Println(requestBody)
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
	fmt.Println(string(body))

}
