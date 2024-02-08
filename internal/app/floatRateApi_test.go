package app

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)


func TestRequestCurrencyExchangeRateFromFloatRates(t *testing.T) {
	wantFloatRateResponse := 0.9345

	testFloatRateApi := testServer{
		name: "happy-api-floatrate-server-response",
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"eur":{
						"code": 			"USD",
						"alphaCode": 	"USD",
						"numericCode": 	"840",
						"name": 			"United States Dollar",
						"rate": 			0.9345,
						"date": 			"2021-09-01",
						"inverseRate": 	1.0
					 }
					}`))
		})),
		expectedResponse: &wantFloatRateResponse,
		expectedErr: nil,
}

 t.Run("Test RequestCurrencyExchangeRateFromFloatRates", func(t *testing.T) {
	os.Setenv("FloatRatesUrl", testFloatRateApi.server.URL+"?file=")

	ch := make(chan ChannelData, 1)

	currencyTarget := "EUR"
	currencyOrigin := "USD"

	RequestCurrencyExchangeRateFromFloatRates(currencyOrigin, currencyTarget, ch)

	result := <-ch

	
	if result.ExchangeRate != wantFloatRateResponse {
		t.Errorf("Error requesting data from FloatRates")
	}


 })

 t.Run("Test Error", func(t *testing.T) {
	os.Setenv("FloatRatesUrl", testFloatRateApi.server.URL+"?file=")

	ch := make(chan ChannelData, 1)

	currencyTarget := "EUR"
	currencyOrigin := "SomeBullshit"

	RequestCurrencyExchangeRateFromFloatRates(currencyOrigin, currencyTarget, ch)

	result := <-ch

	if result.Err != nil {
		t.Errorf("Error handling errors from FloatRates")
	}

 })



}