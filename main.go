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

	currencies := map[string][]string{
		"JADE.BTC":  []string{"CYB", "JADE.LTC", "JADE.LHT"},
		"JADE.ETH":  []string{"JADE.INK", "CYB", "JADE.LHT", "JADE.MT", "JADE.EOS", "JADE.DPY", "JADE.PPT", "JADE.TCT", "JADE.GNX", "JADE.MVP", "JADE.GNT", "JADE.MKR", "JADE.FUN"},
		"JADE.USDT": []string{"CYB", "JADE.ETH", "JADE.BTC", "JADE.EOS", "JADE.LTC"},
	}
	currentTime := time.Now()
	for origin, possibleTargets := range currencies {
		for _, target := range possibleTargets {
			for hour := 0; hour < 100; hour++ {
				fmt.Println(origin, target, hour)
				substractedTime := time.Duration(hour) * time.Hour
				dateStart := currentTime.Add(-substractedTime)
				dateStop := dateStart.Add(-10 * time.Hour)
				fmt.Printf("dateStart: %s, dateStop: %s, substractedTime: %s, hour: %d \n", dateStart, dateStop, substractedTime, hour)
				queryOrder(origin, target, dateStart, dateStop)
			}
		}
	}
}

func queryOrder(origin, target string, dateStart, dateStop time.Time) {

	reqParams := requestParams{currencyOrigin: origin,
		currencyTarget:        target,
		dateStart:             dateStart.Format("2006-01-02T15:04:05"),
		dateStop:              dateStop.Format("2006-01-02T15:04:05"),
		MaxNumberTransactions: 100,
	}
	// fmt.Printf("dateStart: %s, dateStop: %s \n", dateStart, dateStop)
	body := getDataFromCybex(reqParams)

	orders := responseOrders{}

	err := json.Unmarshal(body, &orders)
	if err != nil {
		fmt.Println("Not possible to unmarshal json response", err)
	}

	postTransactionsToES(orders, reqParams)

}

func postTransactionsToES(orders responseOrders, reqParams requestParams) {
	for _, currentTransaction := range orders.Result {

		currentTransaction.CurrencyOrigin = reqParams.currencyOrigin
		currentTransaction.CurrencyTarget = reqParams.currencyTarget

		fmt.Printf("%+v \n ", currentTransaction)

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

func getDataFromCybex(reqParams requestParams) []byte {
	url := "https://apihk.cybex.io"

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
	return body
}
