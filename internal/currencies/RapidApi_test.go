package currencies

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)


func TestRequestCurrencyExchangeRateFromRapid(t *testing.T) {
	wantFloatRateResponse := 0.9345
	wantRapidApiResponse := 0.92
	testRapidApi:= testServer{
		name: "happy-api-rapid-server-response",
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"code": 			"USD",
				"msg": "Success",
				"rates":{ "EUR": 0.92},
				"base": "USD"
			}`))
	})),
	expectedResponse: &wantRapidApiResponse,
	expectedErr: nil,
	}

 t.Run("Test RequestCurrencyExchangeRateFromRapid", func(t *testing.T) {
	os.Setenv("RapidUrl", testRapidApi.server.URL)

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
	os.Setenv("RapidUrl", testRapidApi.server.URL)

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