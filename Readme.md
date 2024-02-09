# Go API to get currencies exchange rate and script to list and manage persons data

## API currencies 

#### Run API 

Install Dependencies

    RUN go mod download

to run api you can use makefile

    make run/api/local

in order to have the correct environment remember to add the RapidAPIKey to the makefile <br>
you can find a makefile with the API Key on the email I've sent with the email response or look for your own on [RAPID](https://rapidapi.com/juhestudio-juhestudio-default/api/exchange-rate-api1/) <br>
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

with this script you can use multiple arguments <br>
`-file=data/persons.json`: path to .json with the list of persons <br>
`-order=ASC`: ASC or DESC order the person list <br>
`-salary=100`: Filter the persons by people with more than 100 USD in salary <br>
`outputDir=./`: if included it will create a csv file with the filter result <br>

also it needs the environment variable `CurrencyExchangeApiUrl=http://localhost:8080` included in the makefile <br>
still you can test the production api with the url from the aws deployment: http://example-currencies-lb-1760636663.us-west-2.elb.amazonaws.com/

Hope You like it and left a star
