# Go API to get currencies exchange rate and script to list and manage persons data

## API currencies 

#### Run API 

Install Dependencies

    RUN go mod download

to run api you can use makefile

    make run/api/local

in order to have the correct environment remenber to add the RapidAPIKey to the makefile
example

    run/api/local:
	@echo "Run API" & PORT=8080 \
	RapidAPIKey=API_KEY_HERE \
	RapidAPIHost=exchange-rate-api1.p.rapidapi.com \
	RapidUrl=https://currency-conversion-and-exchange-rates.p.rapidapi.com/latest \
	FloatRatesUrl=https://www.floatrates.com/daily/ go run cmd/currencies/main.go

expected output

    Server running on port 8080

#### Run Script

with this script you can use multiple arguments 
`-file=data/persons.json`: path to .json with the list of persons
`-order=ASC`: ASC or DESC order the person list 
`-salary=100`: Filter the persons by people with more than 100 USD in salary
`outputDir=./`: if included it will create a csv file with the filter result 

also it needs the environment variable `CurrencyExchangeApiUrl=http://localhost:8080` included in the makefile

Hope You like it and left a star
