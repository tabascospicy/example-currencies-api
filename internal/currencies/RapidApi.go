package currencies

import (
	"github.com/spicyt/currencies/pkg/CustomJsonDecoder"
	"github.com/spicyt/currencies/pkg/httpClient"
	readEnv "github.com/spicyt/currencies/pkg/readEnv"
)

type RapiResponse struct {
	Code    string             `json:"code"`
	Message string             `json:"msg"`
	Base    string             `json:"base"`
	Rates   map[string]float64 `json:"rates"`
}


func RequestCurrencyExchangeRateFromRapid(CurrencyOrigin string, CurrencyTarget string, ch chan ChannelData) {

	// read env variables
	RapidApiKey := readEnv.ReadVariable("RapidAPIKey")
	RapidUrl := readEnv.ReadVariable("RapidUrl")
	RapidHost := readEnv.ReadVariable("RapidAPIHost")


	// assign the request parameters and headers
	requestParams := map[string]string{"base": CurrencyOrigin}
	requestHeaders := map[string]string{"X-RapidAPI-Key": RapidApiKey, "X-RapidAPI-Host": RapidHost}

	requestConfig := httpClient.RequestConfig{
		Url:     RapidUrl,
		Params:  requestParams,
		Headers: requestHeaders}


	// make the request
	response, errorResponse := httpClient.Get(requestConfig)

	 // if for some reason there is a problem
	 if  errorResponse != nil{
		result := ChannelData{Err: errorResponse, ExchangeRate: -1}
		ch <- result
		return
	}

	var responseJson RapiResponse

	err := CustomJsonDecoder.DecodeJson(response.Body, &responseJson, "RapiResponse")

 // if for some reason there is a problem
	if err != nil {
		result := ChannelData{Err: err, ExchangeRate: -1}
		ch <- result
		return
	}

	defer response.Body.Close()

	result := responseJson.Rates[CurrencyTarget]


	// if the channel already has a value close it and return to cancel the operation since we just need one value
	if len(ch) > 0 {
		close(ch)
		return
	}

	ch <- ChannelData{Err: nil, ExchangeRate: result}
}
