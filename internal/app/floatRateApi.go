package app

import (
	"fmt"
	s "strings"

	"github.com/spicyt/currencies/pkg/CustomJsonDecoder"
	"github.com/spicyt/currencies/pkg/httpClient"
	readEnv "github.com/spicyt/currencies/pkg/readEnv"
)

// types

type FlotApiResponseError struct {
	Err     error
	Message string
}

// float rates api response single element
type FloatRatesRateItem struct {
	Code        string  `json:"code"`
	AlphaCode   string  `json:"alphaCode"`
	NumericCode string  `json:"numericCode"`
	Name        string  `json:"name"`
	Rate        float64 `json:"rate"`
	Date        string  `json:"date"`
	InverseRate float64 `json:"inverseRate"`
}


// complete float rates api response it encapsulates all the rates in a map like ratesData[moneyName] example ratesData["USD"]
type FloatRatesResponse map[string]FloatRatesRateItem


func (d *FlotApiResponseError) Error() string {
	return "Error requesting data: " + d.Err.Error() + d.Message
}


func RequestCurrencyExchangeRateFromFloatRates(CurrencyOrigin string, CurrencyTarget string, ch chan ChannelData) {

	// define the file extension
	fileExtension := "json"

	// read the float rates url from the environment
	floatRateUrl := readEnv.ReadVariable("FloatRatesUrl")

	// create the url string by concatenating the floatRateUrl with the currency origin name and the file extension 
	// final result shoul be: "https://www.floatrates.com/daily/eur.json"
	url := fmt.Sprintf("%s%s.%s",floatRateUrl, s.ToLower(CurrencyOrigin), fileExtension)

	// create the request config
	requestConfig := httpClient.RequestConfig{
		Url: url,
	}

	//TODO: add error handling

	response, errorResponse := httpClient.Get(requestConfig)

	if errorResponse != nil{
		result := ChannelData{Err: errorResponse, ExchangeRate: -1}

		ch <- result
		return
	}

	var apiRates FloatRatesResponse

	err := CustomJsonDecoder.DecodeJson(response.Body, &apiRates, "FloatRatesResponse")

	if err != nil {
		result := ChannelData{Err: err, ExchangeRate: -1}
		ch <- result
		return
	}


	defer response.Body.Close()

	// read the rate from the map using the target currency as the key
	result := apiRates[s.ToLower(CurrencyTarget)].Rate


	// if the channel already has a value close it and return to cancel the operation since we just need one value
	if len(ch) > 0 {
		close(ch)
		return
	}


	
	ch <- ChannelData{ExchangeRate: result, Err: nil}
}
