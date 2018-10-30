package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func main() {

	// currencies := []string{"JADE.LHT", "JADE.INK", "CYB", "JADE.LHT", "JADE.MT", "JADE.DPY", "JADE.PPT", "JADE.TCT", "JADE.GNX", "JADE.MVP", "JADE.GNT", "JADE.MKR", "JADE.FUN", "JADE.ETH", "JADE.BTC", "JADE.EOS", "JADE.LTC"}

	currencies := []string{"JADE.ETH", "JADE.BTC"}
	c1 := make(chan requestParams, 2)
	c2 := make(chan channelStruct, 2)
	c3 := make(chan bool, 1)
	go func() {
		for {
			params := <-c1
			c2 <- channelStruct{queryOrder(params), params}
		}
	}()

	go func() {
		for {
			channelresponse := <-c2
			postTransactionsToES(channelresponse.orders, channelresponse.params, c3)
			<-c3
		}
	}()

	currentTime := time.Now()
	for numTimeBlocks := 0; numTimeBlocks < 50000; numTimeBlocks++ {
		for _, origin := range currencies {
			for _, target := range currencies {
				if origin != target {
					fmt.Println(origin, target, numTimeBlocks)

					substractedTime := time.Duration(numTimeBlocks*2) * time.Minute
					dateStart := currentTime.Add(-substractedTime)
					dateStop := dateStart.Add(-2 * time.Minute)

					reqParams := requestParams{currencyOrigin: origin,
						currencyTarget:        target,
						dateStart:             dateStart.Format("2006-01-02T15:04:05"),
						dateStop:              dateStop.Format("2006-01-02T15:04:05"),
						MaxNumberTransactions: 100,
					}

					c1 <- reqParams
					// postTransactionsToES(orders, reqParams)
				}
			}
		}
	}
}

func queryOrder(reqParams requestParams) responseOrders {

	body := getDataFromCybex(reqParams)

	orders := responseOrders{}

	err := json.Unmarshal(body, &orders)
	if err != nil {
		fmt.Println("Not possible to unmarshal json response", err)
	}
	return orders

}

func postTransactionsToES(orders responseOrders, reqParams requestParams, c3 chan bool) {
	fmt.Println(orders.Result)
	for _, currentTransaction := range orders.Result {

		currentTransaction.CurrencyOrigin = reqParams.currencyOrigin
		currentTransaction.CurrencyTarget = reqParams.currencyTarget

		stringTransaction, err := json.Marshal(currentTransaction)
		if err != nil {
			fmt.Println("error marshaling struct", err)
		}

		client := &http.Client{}
		newID := currentTransaction.Date + currentTransaction.Side1AccountID
		elasticSearchURL := fmt.Sprintf("http://elasticsearch:9200/transactions/_doc/%s", newID)
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
		fmt.Println("Shouldprint")
		fmt.Println(res.StatusCode, string(body))
	}
	c3 <- true
}

func getDataFromCybex(reqParams requestParams) []byte {
	url := "https://apihk.cybex.io"

	requestBody := fmt.Sprintf(
		`{
			"params":["%s","%s","%s","%s",%d],
			"id":1,
			"jsonrpc":"2.0",
			"method":"get_trade_history"
		}`,
		reqParams.currencyOrigin,
		reqParams.currencyTarget,
		reqParams.dateStart,
		reqParams.dateStop,
		reqParams.MaxNumberTransactions)

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
	// fmt.Println(string(body))
	return body
}
