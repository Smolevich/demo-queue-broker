package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"net/http"
	"time"
)

var ctx = context.Background()
var client *http.Client

type RatesResponse struct {
	Date         string  `json:"Вate"`
	PreviousDate string  `json:"PreviousDate"`
	PreviousURL  string  `json:"PreviousURL"`
	Timestamp    string  `json:"Timestamp"`
	Valutes      Valutes `json:"Valute"`
}

type Valutes struct {
	USD Valute `json:"USD"`
}

type Valute struct {
	ID       string  `json:"ID"`
	CharCode string  `json:"CharCode"`
	NumCode  string  `json:"NumCode"`
	Name     string  `json:"Name"`
	Value    float64 `json:"Value"`
}

func main() {
	redisCLient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}
	client = &http.Client{Transport: tr}
	response := RatesResponse{}
	err := getJson("https://www.cbr-xml-daily.ru/daily_json.js", &response)
	fmt.Println(response)
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		err := redisCLient.Set(ctx, "valute_date_updated", response.Date, 0).Err()
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

}

func getJson(url string, target interface{}) error {
	r, err := client.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	fmt.Println(r.Body)

	return json.NewDecoder(r.Body).Decode(target)
}
