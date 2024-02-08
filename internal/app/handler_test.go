package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)



type testServer struct {
	name             string
	server           *httptest.Server
	expectedResponse *float64
	expectedErr      error
}

func TestGetCurrencyExchangeHandler(t *testing.T) {

	 wantFloatRateResponse := 0.9345
	 wantRapidApiResponse := 0.92

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


  t.Run("Test all integration", func(t *testing.T) {
		os.Setenv("RapidAPIKey", "TEST-KEY")
		os.Setenv("RapidAPIHost", "TEST-HOST")
		os.Setenv("RapidUrl", testRapidApi.server.URL)
		os.Setenv("FloatRatesUrl", testFloatRateApi.server.URL+"?file=")
		req, err := http.NewRequest("GET", "/", bytes.NewReader([]byte(`{"currency-pair": "USD-EUR"}`)))

		if err != nil {
			t.Fatal(err)
		}
	
		req.Header.Set("Content-Type", "application/json")
	
		rr := httptest.NewRecorder()
	
		GetCurrencyExchangeHandler(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	
		var gotResponse CurrencyExchangeRateResponse
	
		if err := json.Unmarshal(rr.Body.Bytes(), &gotResponse); err != nil {
			t.Errorf("Handler returned unexpected  invalid response %v", err)
		}
		fmt.Print(gotResponse["USD-EUR"])
	
		if gotResponse["USD-EUR"] !=  *testRapidApi.expectedResponse && gotResponse["USD-EUR"] !=  *testFloatRateApi.expectedResponse{
			t.Errorf("Handler returned unexpected currencyPair: got %v want number above 0", gotResponse["currency-pair"])
		}
	
		expectedContentType := "application/json"
	
		if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
			t.Errorf("Handler returned wrong content type: got %v want %v", contentType, expectedContentType)
		}
	})

	t.Run("Test with the two APIs failing", func(t *testing.T) {
		os.Setenv("RapidAPIKey", "TEST-KEY")
		os.Setenv("RapidAPIHost", "TEST-HOST")
		//  just add some invalid url to make the request fail in purpose
		os.Setenv("RapidUrl", "http://testUrl.com")
		os.Setenv("FloatRatesUrl", "http://testUrl2.com")
		req, err := http.NewRequest("GET", "/", bytes.NewReader([]byte(`{"currency-pair": "USD-EUR"}`)))

		if err != nil {
			t.Fatal(err)
		}
	
		req.Header.Set("Content-Type", "application/json")
	
		rr := httptest.NewRecorder()
	
		GetCurrencyExchangeHandler(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	
		errorMessage := string(rr.Body.Bytes())

		if errorMessage != "Error requesting data from the APIs" {
			t.Errorf("Handler returned unexpected  invalid response %v", err)
		}
	
	})










}
