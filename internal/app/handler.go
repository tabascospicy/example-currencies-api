package currencies

import (
	"encoding/json"
	"io"
	"net/http"
	s "strings"

	readEnv "github.com/spicyt/currencies/pkg"
)


type RapiResponse struct {
	Code string  `json:"code"`
	Message      string `json:"msg"`
	Base   bool   `json:"base"`
	Rates 		map[string]float64 `json:"rates"`
}



type CurrencyExchangeRateInput struct {
	CurrencyPair string `json:"currency-pair"`
}

type CurrencyExchangeRateResponse map[string]float64



type DecodeError struct {
	Err error
	Context string
}

func (d *DecodeError) Error() string {
	return "Error decoding data: " + d.Err.Error() + " with context " + d.Context
}


func DecodeResponse(response io.Reader, data interface{}, name string) error {
	err := json.NewDecoder(response).Decode(&data)
	if err != nil {
		return &DecodeError{Err: err, Context: name}
	}

//	fmt.Printf("Data: %v\n", data)
	return nil
}


func RequestCurrencyExchangeRateFromRapid(CurrencyOrigin string , CurrencyTarget string) (RapiResponse, error) {
	RapidApiKey := readEnv.ReadVariable("RapidAPIKey")
	RapidUrl := readEnv.ReadVariable("RapidUrl")
	RapidHost := readEnv.ReadVariable("RapidAPIHost")

	requestParams := map[string]string{"from": CurrencyOrigin, "to": CurrencyTarget}
	requestHeaders := map[string]string{"X-RapidAPI-Key": RapidApiKey, "X-RapidAPI-Host": RapidHost}

	requestConfig := RequestConfig{ 
		Url: RapidUrl, 
		Params: requestParams, 
		Headers: requestHeaders}

	response, err := Get(requestConfig)
	
	if err != nil {
		return RapiResponse{} , err
	}

	defer response.Body.Close()

	var responseJson RapiResponse

 	errReq := DecodeResponse(response.Body, &responseJson, "RapiResponse") 
	
  
	return responseJson, errReq
}



func GetCurrencyExchangeHandler(w http.ResponseWriter, r *http.Request){
	readEnv.Init()
	var requestInput CurrencyExchangeRateInput

	err := DecodeResponse(r.Body, &requestInput, "CurrencyExchangeRateInput")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request body"))
		return
	}

	splitedCurrency := s.Split(requestInput.CurrencyPair, "-") 

	
	response, err := RequestCurrencyExchangeRateFromRapid(splitedCurrency[0], splitedCurrency[1])

	if(err != nil){
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error getting exchange rate: " + err.Error()))
		return
	}

	resultKey := s.Join([]string{splitedCurrency[0], splitedCurrency[1]}, "-")

	result := CurrencyExchangeRateResponse{resultKey: response.Rates[splitedCurrency[1]]}
  
	jsonResponse, _ := json.Marshal(result)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonResponse))
}
