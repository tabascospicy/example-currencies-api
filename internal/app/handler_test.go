package currencies

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func Init() {
	// load .env file
	os.Setenv("testing", "true")
	err := godotenv.Load("../../api.env")

	if err != nil {
		fmt.Printf(" %v\n", err)
		log.Fatalf("Error loading .env file")
	}

}

func TestGetCurrencyExchangeHandler(t *testing.T) {
	Init()
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

	if gotResponse["USD-EUR"] <= 0 {
		t.Errorf("Handler returned unexpected currencyPair: got %v want number above 0", gotResponse["currency-pair"])
	}

	expectedContentType := "application/json"

	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("Handler returned wrong content type: got %v want %v", contentType, expectedContentType)
	}
}
