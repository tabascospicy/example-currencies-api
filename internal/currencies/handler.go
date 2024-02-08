package currencies

import (
	"encoding/json"
	"net/http"
	s "strings"

	"github.com/spicyt/currencies/pkg/CustomJsonDecoder"
)

type CurrencyExchangeRateInput struct {
	CurrencyPair string `json:"currency-pair"`
}

type CurrencyExchangeRateResponse map[string]float64


type ChannelData struct {
	ExchangeRate float64
	Err error
}






func GetCurrencyExchangeHandler(w http.ResponseWriter, r *http.Request) {

	var requestInput CurrencyExchangeRateInput

	err := CustomJsonDecoder.DecodeJson(r.Body, &requestInput, "CurrencyExchangeRateInput")

	// if there is no body then send a bad request response
	if err != nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Invalid request body"))
		return
	}
	ch := make(chan ChannelData, 2)

	splitedCurrency := s.Split(requestInput.CurrencyPair, "-")

	// run the two requests in parallel
	go RequestCurrencyExchangeRateFromRapid(splitedCurrency[0], splitedCurrency[1], ch)
	go RequestCurrencyExchangeRateFromFloatRates(splitedCurrency[0], splitedCurrency[1], ch)

	resultCounter := 0
	var resultRate float64
	erroredCounter := 0

	// we create a loop to listen when the first response arrives
	for (resultCounter < 1 && erroredCounter < 2) {
		result, ok := <-ch
		if ok && result.ExchangeRate != -1 {
			// once the first response arrives we end the loop and save the result
			resultCounter += 1
			resultRate = result.ExchangeRate
		}
		// if the two responses gave an error then we return a bad request response
		if (result.Err != nil){
			erroredCounter += 1
			continue 
		}
	}

	if(erroredCounter >= 2){
		  close(ch)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error requesting data from the APIs"))
			return
	}

	// here we form the string key for the response map like "USD-EUR"
	resultKey := s.Join([]string{splitedCurrency[0], splitedCurrency[1]}, "-")

	// then create the map with the result
	result := CurrencyExchangeRateResponse{resultKey: resultRate}

	jsonResponse, _ := json.Marshal(result)


	// and finally we write the response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonResponse))

}
