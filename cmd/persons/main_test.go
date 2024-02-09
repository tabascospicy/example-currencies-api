package main

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strconv"
	s "strings"
	"testing"
)

const TestPersonsListJson = `
[
    { "id": "1", "personName": "Cadanaut Ruben", "salary": { "value": "200000", "currency": "NPR" } },
    { "id": "2", "personName": "Cadanaut Juan", "salary": { "value": "40000", "currency": "NZD" } },
    { "id": "3", "personName": "Cadanaut Albert", "salary": { "value": "1300", "currency": "USD" } },
    { "id": "4", "personName": "Cadanaut Max", "salary": { "value": "4200", "currency": "EUR" } },
    { "id": "5", "personName": "Cadanaut G", "salary": { "value": "5000", "currency": "USD" } },
    { "id": "6", "personName": "Cadanaut A", "salary": { "value": "2010", "currency": "JPY" } },
    { "id": "7", "personName": "Cadanaut M", "salary": { "value": "1100", "currency": "JPY" } },
    { "id": "8", "personName": "Cadanaut Z", "salary": { "value": "1200", "currency": "USD" } },
    { "id": "9", "personName": "Cadanaut Y", "salary": { "value": "140000", "currency": "NPR" } },
    { "id": "10", "personName": "Cadanaut X", "salary": { "value": "3010", "currency": "USD" } }
]
`

type testServer struct {
	name             string
	server           *httptest.Server
	expectedErr      error
}

func TestReadPersonsListAndOrder(t *testing.T) {
	// TestReadPersonsListAndOrder tests the function ReadPersonsListAndOrder
	// of the main.go file


	t.Run("Test Read Person Lists and order them in ASC", func(t *testing.T) {
		// Test the function ReadPersonsListAndOrder
		// of the main.go file

   expectedLength := 10

	 personsList := PersonsList{}

	 json.Unmarshal([]byte(TestPersonsListJson), &personsList)

	 personsList.Sort("ASC")



	 if(len(personsList) != expectedLength){
		t.Errorf("Error reading json file")
	 }

	 lastSalary := 0

	 for person := range personsList {

		salary, _ := strconv.Atoi(personsList[person].Salary.Value)
		if salary < lastSalary {
			t.Errorf("Error ordering by salary")
		}
		lastSalary = salary
	 }

	})
	t.Run("Test Read Person Lists and passing an unexpected order param", func(t *testing.T) {
		// Test the function ReadPersonsListAndOrder
		// of the main.go file


		personsList1 := PersonsList{}

		personList2 := PersonsList{}
 
		json.Unmarshal([]byte(TestPersonsListJson), &personsList1)
		json.Unmarshal([]byte(TestPersonsListJson), &personList2)
		err := personsList1.Sort("Something")

		if err == nil {
			t.Errorf("Error displaying error when passing invalid sort json file")
		 }
		
		

	})
	t.Run("Test Read Person Lists and order them in DESC", func(t *testing.T) {
		// Test the function ReadPersonsListAndOrder
		// of the main.go file
		expectedLength := 10


		personsList := PersonsList{}
 
		json.Unmarshal([]byte(TestPersonsListJson), &personsList)

		personsList.Sort("DESC")

		if(len(personsList) != expectedLength){
			t.Errorf("Error reading json file")
		 }
	
		 lastSalary := 0
	
		 for person := range personsList {
	
			salary, _ := strconv.Atoi(personsList[person].Salary.Value)
			if salary > lastSalary && lastSalary != 0{
				t.Errorf("Error ordering by salary DESC")
			}
			lastSalary = salary
		 }
 

	 if(len(personsList) != expectedLength){
		t.Errorf("Error reading json file by DESC")
	 }

	})
	t.Run("Test Read Person Lists and order them in ASC", func(t *testing.T) {
		// Test the function ReadPersonsListAndOrder
		// of the main.go file
		expectedLength := 10


		personsList := PersonsList{}
 
		json.Unmarshal([]byte(TestPersonsListJson), &personsList)

		personsList.Sort("ASC")

		if(len(personsList) != expectedLength){
			t.Errorf("Error reading json file")
		 }
	
		 lastSalary := 0
	
		 for person := range personsList {
	
			salary, _ := strconv.Atoi(personsList[person].Salary.Value)
			if salary < lastSalary && lastSalary != 0 {
				t.Errorf("Error ordering by salary ASC")
			}
			lastSalary = salary
		 }
 

	 if(len(personsList) != expectedLength){
		t.Errorf("Error reading json file by ASC")
	 }

	})

	t.Run("Test Read Person Lists and orderby salary", func(t *testing.T) {
		// Test the function ReadPersonsListAndOrder
		// of the main.go file
		personsList := PersonsList{}

		currenciesList := "NPR USD EUR JPY NZD"
  
		json.Unmarshal([]byte(TestPersonsListJson), &personsList)

		result :=	OperateWithPersonsList(personsList, "ASC")
		
		for currency, _ := range result {


			if !s.Contains(currenciesList, currency) {
				t.Errorf("Error ordering by salary by currency")
			}
			
		}

	fmt.Println(result)

	})

}