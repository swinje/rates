package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type rates struct {
	Amount float64            `json:"amount"`
	Base   string             `json:"base"`
	Date   time.Time          `json:"date"`
	Rates  map[string]float64 `json:"rates"`
}

func main() {

	url := "https://api.frankfurter.app/latest"

	ratesClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Go tutorial")

	res, getErr := ratesClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	rates1 := rates{}
	// Use custom since has date
	jsonErr := json.Unmarshal(body, &rates1)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	fmt.Printf("%s er 1 EURO:\n", rates1.Date.Format("01-02-2006"))
	fmt.Printf("Norske kroner %.02f\n", rates1.Rates["NOK"])
	fmt.Printf("Sveitiske Frank %.02f\n", rates1.Rates["CHF"])
	fmt.Printf("US Dollar %.02f\n", rates1.Rates["USD"])
	fmt.Printf("Engelske Pund %.02f\n", rates1.Rates["GBP"])
}

type Unmarshaler interface {
	UnmarshalJSON([]byte) error
}

func (l *rates) UnmarshalJSON(j []byte) error {
	var rawStrings map[string]interface{}

	err := json.Unmarshal(j, &rawStrings)
	if err != nil {
		return err
	}

	for k, v := range rawStrings {
		if strings.ToLower(k) == "amount" {
			fV := v.(float64)
			l.Amount = fV
			if err != nil {
				return err
			}
		}
		if strings.ToLower(k) == "base" {
			sV := v.(string)
			l.Base = sV
		}
		if strings.ToLower(k) == "date" {
			sV := v.(string)
			t, err := time.Parse("2006-01-02", sV)
			if err != nil {
				return err
			}
			l.Date = t
		}
		if strings.ToLower(k) == "rates" {
			data := v.(map[string]interface{})
			l.Rates = make(map[string]float64)
			for key, value := range data {
				sV := value.(float64)
				l.Rates[key] = sV
			}
		}

	}

	return nil
}
