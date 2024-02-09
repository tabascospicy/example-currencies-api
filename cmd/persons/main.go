package main

import (
	"bytes"
	"cmp"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"sync"

	"github.com/spicyt/currencies/pkg/CustomJsonDecoder"
	"github.com/spicyt/currencies/pkg/httpClient"
)

type Person struct {
	ID         string `json:"id"`
	PersonName string `json:"personName"`
	Salary     Salary `json:"salary"`
}

type PersonsList []Person

type Salary struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}


type CurrencyExchangeRateInput struct {
	CurrencyPair string `json:"currency-pair"`
}

type CurrencyExchangeRateResponse map[string]float64

type PersonGroupedBySalary map[string][]Person



func (p PersonsList) Sort(order string) error {
	// if the order is not ASC or DESC then throw an error
	if order != "DESC" && order != "ASC" {
		
		return fmt.Errorf("Invalid order param: %s", order)
	}

	sortFunc := func(person1, person2 Person) int {

		salary1, _ := strconv.Atoi(person1.Salary.Value)
		salary2, _ := strconv.Atoi(person2.Salary.Value)
		// sort by salary ASC
		if order == "ASC" {
			return cmp.Compare(salary1, salary2)
		}

		// at this point the salary order is DESC
		return cmp.Compare(salary2, salary1)

	}
	
	slices.SortFunc(p, sortFunc)
	
	return nil
}

func (p PersonsList) GroupBySalary() PersonGroupedBySalary{
	// create an empty map to store the groups
  groups := make(PersonGroupedBySalary)

	for _, person := range p {
		// if the currency is not already in the map then create a new group
		if _, ok := groups[person.Salary.Currency]; !ok {
			groups[person.Salary.Currency] = []Person{}
		}
		// append the person to the group
		groups[person.Salary.Currency] = append(groups[person.Salary.Currency], person)
	}
 // finally return the map
	return groups
}

func (p PersonsList) FilterBySalaryInUSD(salaryMin int) PersonsList {

	groupedByCurrency := p.GroupBySalary()
 // create a wait group in order to wait for each request to complete
	wg := sync.WaitGroup{}

	ch := make(chan PersonsList)

	for currency, persons := range groupedByCurrency {
		// add one item to the wait group and run the goroutine
			wg.Add(1)
			go requestSalaryConversion(persons, currency, ch, &wg)
	}

	go func (){
		wg.Wait()
		close(ch)
		}()

	var personsList PersonsList
	// when each call is ready and the channel is closed we can iterate over the results
	for PersonsWithConvertedSalary := range ch {

		for _, newPerson := range PersonsWithConvertedSalary {
			salary,_ := strconv.Atoi(newPerson.Salary.Value)
			// filter the persons with salary greater than the min
			if salary > salaryMin {
				personsList = append(personsList, newPerson)
			}
		}
	}


	// return the filtered list
	return personsList
}


func requestSalaryConversion(persons PersonsList, currency string, ch chan PersonsList, wg *sync.WaitGroup) {


	// first we need to tell the wait group that the goroutine is done with a defer
	defer wg.Done()
	// if for some reason there is no persons then do nothing
	if(len(persons) == 0){
		return
	}
	// for USD we don't need to make any request
	if(currency == "USD"){
		ch <- persons
		return
	}


	// set the request body
	body := CurrencyExchangeRateInput{CurrencyPair: currency + "-USD"}

	json,_ := json.Marshal(body)

	apiUrl := os.Getenv("CurrencyExchangeApiUrl")

	// and the request config
	config := httpClient.RequestConfig{
		Url: apiUrl,
		Body: bytes.NewReader(json),
	}

	// make the request
  result, err :=	httpClient.Get(config)


	if err != nil {
		fmt.Printf("Error making request to %s with method %s : %v", config.Url, "GET", err)
		return
	}

	var response CurrencyExchangeRateResponse

	errorDecoding := CustomJsonDecoder.DecodeJson(result.Body, &response, "CurrencyExchangeRateResponse")

	if errorDecoding != nil {
		return 
	}

	defer result.Body.Close()

	var PersonsWithConvertedSalary = make(PersonsList, len(persons))
	// get the exchange rate that arrives as a map with the currency pair as the key like {"EUR-USD": 1.2}
	exchangeRate := response[currency + "-USD"]



	for i, person := range persons {
		// create an empty person 
		personWithNewSalary := Person{}

		salary, _ := strconv.Atoi(person.Salary.Value)

		convertedSalary := float64(salary) * exchangeRate

		// copy the person data to the new person 
		personWithNewSalary.ID = person.ID
		personWithNewSalary.PersonName = person.PersonName
		// set the new salary in USD
		personWithNewSalary.Salary.Currency = "USD"
		personWithNewSalary.Salary.Value = strconv.Itoa(int(convertedSalary))
		
		// append the new person to the list
		PersonsWithConvertedSalary[i] = personWithNewSalary
	}
 

	// send to the channel
	ch <- PersonsWithConvertedSalary
}





func ReadPersonsList(fileToReadPath string) PersonsList{
	 // Open the JSON file
	 file, err := os.Open(fileToReadPath)
	 if err != nil {
		 log.Fatalf("Error opening file: %v", err)
	 }
	 defer file.Close()
 

	var personsList PersonsList

	// assign the result of the custom json decoder to the personsList
	CustomJsonDecoder.DecodeJson(file, &personsList, "Persons")



	return personsList
}







func OperateWithPersonsList(personsList  PersonsList, order string) PersonGroupedBySalary{
	
	// sort in the required order
	personsList.Sort(order)

	fmt.Printf("Persons ordered in %s :%v"  ,order , personsList)
	 // finally group the persons by salary
	orderedBySalary := personsList.GroupBySalary()

	fmt.Printf("Persons ordered by salary :%v"  , orderedBySalary)	

	

	return orderedBySalary
}

func ReadAndFilterBySalary(personsList PersonsList, salaryMin int) PersonsList{
	// filter them by salary
	filteredPersons := personsList.FilterBySalaryInUSD(salaryMin)

	return filteredPersons
}

func main() {
// os.args[1]

	outputDirPtr := flag.String("outputDir", "", "order")
	filePtr := flag.String("file", "data/persons.json", "file to read")
	orderPtr := flag.String("order", "ASC", "order")
	salaryFilterPtr := flag.Int("salary", 0, "salary filter")

	flag.Parse()

	fmt.Println("outputPtr:", *outputDirPtr)
	fmt.Println("filePtr:", *filePtr)
	fmt.Println("orderPtr:", *orderPtr)

	args := os.Args[1:]

	fmt.Println(args)


  

	personsList := ReadPersonsList(*filePtr)
	orderedBySalary := OperateWithPersonsList(personsList, *orderPtr)

	if *salaryFilterPtr != 0 {
		personsList = ReadAndFilterBySalary(personsList, *salaryFilterPtr)
	}

	fmt.Println(orderedBySalary)



	if(*outputDirPtr == ""){
		return
	}

	file,_ := os.Create(*outputDirPtr + "persons.csv")
	headers := "ID,PersonName,SalaryValue,SalaryCurrency\n"
	file.WriteString(headers)

	for _, person := range personsList {
		file.WriteString(person.ID + "," + person.PersonName + "," + person.Salary.Value + "," + person.Salary.Currency + "\n")
	}

	file.Close()

}
