starting:
	echo "Process started"

build-local:
	go build -o bin/main cmd/currencies/main.go

run/build/production:
	@echo "Building the application" 

run/test:
		go test -v ./...

run/api/local: build-local
	@echo "Running the application" &  PORT=8080 \
	RapidAPIKey=87cc7dc600msh21712c972e7b14ap19d15fjsn629c3ff9dee9 \
	RapidAPIHost=exchange-rate-api1.p.rapidapi.com \
	RapidUrl=https://currency-conversion-and-exchange-rates.p.rapidapi.com/latest \
	FloatRatesUrl=https://www.floatrates.com/daily/ go run cmd/currencies/main.go

run/script/local: 
	@echo "Running the application" & CurrencyExchangeApiUrl=http://localhost:8080 go run cmd/persons/main.go 