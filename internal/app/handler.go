package currencies

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	s "strings"

	readEnv "github.com/spicyt/currencies/pkg"
)

type RapiResponse struct {
	Code    string             `json:"code"`
	Message string             `json:"msg"`
	Base    string             `json:"base"`
	Rates   map[string]float64 `json:"rates"`
}

type FloatRatesRateItem struct {
	Code        string  `json:"code"`
	AlphaCode   string  `json:"alphaCode"`
	NumericCode string  `json:"numericCode"`
	Name        string  `json:"name"`
	Rate        float64 `json:"rate"`
	Date        string  `json:"date"`
	InverseRate float64 `json:"inverseRate"`
}

type FloatRatesResponse map[string]FloatRatesRateItem

type CurrencyExchangeRateInput struct {
	CurrencyPair string `json:"currency-pair"`
}

type CurrencyExchangeRateResponse map[string]float64

type DecodeError struct {
	Err     error
	Context string
}

func (d *DecodeError) Error() string {
	return "Error decoding data: " + d.Err.Error() + " with context " + d.Context
}

type FlotApiResponseError struct {
	Err     error
	Message string
}

func (d *FlotApiResponseError) Error() string {
	return "Error rquesting data: " + d.Err.Error() + d.Message
}

func DecodeResponse(response io.Reader, data interface{}, name string) error {

	val, err := io.ReadAll(response)

	if err != nil {
		return &DecodeError{Err: err, Context: name}
	}

	errJson := json.Unmarshal(val, &data)

	if errJson != nil {
		return &DecodeError{Err: err, Context: name}
	}

	return nil
}

func RequestCurrencyExchangeRateFromRapid(CurrencyOrigin string, CurrencyTarget string, ch chan float64) {
	RapidApiKey := readEnv.ReadVariable("RapidAPIKey")
	RapidUrl := readEnv.ReadVariable("RapidUrl")
	RapidHost := readEnv.ReadVariable("RapidAPIHost")

	requestParams := map[string]string{"from": CurrencyOrigin, "to": CurrencyTarget}
	requestHeaders := map[string]string{"X-RapidAPI-Key": RapidApiKey, "X-RapidAPI-Host": RapidHost}

	requestConfig := RequestConfig{
		Url:     RapidUrl,
		Params:  requestParams,
		Headers: requestHeaders}

	response, _ := Get(requestConfig)

	var responseJson RapiResponse

	DecodeResponse(response.Body, &responseJson, "RapiResponse")
	defer response.Body.Close()

	result := responseJson.Rates[CurrencyTarget]

	if len(ch) > 0 {
		close(ch)
		return
	}

	ch <- result
}

func RequestCurrencyExchangeRateFromFloatRates(CurrencyOrigin string, CurrencyTarget string, ch chan float64) {

	fileExtension := "json"

	floatRateUrl := readEnv.ReadVariable("FloatRatesUrl")

	url := fmt.Sprintf("%s%s.%s",floatRateUrl, s.ToLower(CurrencyOrigin), fileExtension)
	requestConfig := RequestConfig{
		Url: url,
	}

	response, _ := Get(requestConfig)

	var apiRates FloatRatesResponse

	DecodeResponse(response.Body, &apiRates, "FloatRatesResponse")
	defer response.Body.Close()
	result := apiRates[s.ToLower(CurrencyTarget)].Rate

	// TODO: comment later

	if len(ch) > 0 {
		close(ch)
		return
	}
	
	ch <- result
}

func GetCurrencyExchangeHandler(w http.ResponseWriter, r *http.Request) {

	var requestInput CurrencyExchangeRateInput

	err := DecodeResponse(r.Body, &requestInput, "CurrencyExchangeRateInput")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request body"))
		return
	}
	ch := make(chan float64, 2)

	splitedCurrency := s.Split(requestInput.CurrencyPair, "-")

	go RequestCurrencyExchangeRateFromRapid(splitedCurrency[0], splitedCurrency[1], ch)
	go RequestCurrencyExchangeRateFromFloatRates(splitedCurrency[0], splitedCurrency[1], ch)

	resultCounter := 0
	var resultRate float64
	// we create a loop to listen when the first response arrives
	for resultCounter < 1 {
		result, ok := <-ch
		if ok {
			resultCounter += 1
			resultRate = result
		}
	}

	resultKey := s.Join([]string{splitedCurrency[0], splitedCurrency[1]}, "-")

	result := CurrencyExchangeRateResponse{resultKey: resultRate}

	jsonResponse, _ := json.Marshal(result)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonResponse))

}
