starting:
	echo "Process started"

build-local:
	go build -o bin/main cmd/currencies/main.go

run/build/production:
	@echo "Building the application" 

run/test:
		go test -v ./...

run/local: build-local
	@echo "Running the application" &  PORT=8080 \
	RapidAPIKey= \
	RapidAPIHost=exchange-rate-api1.p.rapidapi.com \
	RapidUrl=https://currency-conversion-and-exchange-rates.p.rapidapi.com/latest \
	FloatRatesUrl=https://www.floatrates.com/daily/ go run cmd/currencies/main.go