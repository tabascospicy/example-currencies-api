package main

import (
	"cmp"
	"encoding/json"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
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

type DecodeError struct {
	Err     error
	Context string
}

func (d *DecodeError) Error() string {
	return "Error decoding data: " + d.Err.Error() + " with context " + d.Context
}

func (p PersonsList) Sort(order string) {
	if order != "DESC" && order != "ASC" {
		return
	}
	sortFunc := func(person1, person2 Person) int {

		salary1, _ := strconv.Atoi(person1.Salary.Value)
		salary2, _ := strconv.Atoi(person2.Salary.Value)

		if order == "ASC" {
			return cmp.Compare(salary1, salary2)
		}

		return cmp.Compare(salary2, salary1)

	}

	slices.SortFunc(p, sortFunc)

}

func DecodeJson(response io.Reader, data interface{}, name string) error {

	val, err := io.ReadAll(response)

	if err != nil {
		return &DecodeError{Err: err, Context: name}
	}

	errJson := json.Unmarshal(val, &data)

	if errJson != nil {
		return &DecodeError{Err: err, Context: name}
	}

	return nil
}

func main() {
	jsonFile, err := os.OpenFile("data/persons.json")

	if err != nil {
		log.Fatal("Error reading file: " + err.Error())
	}
	var personsList PersonsList
	personsList.Sort("ASC")
	personsList.Sort("DESC")

	DecodeJson(jsonFile, personsList, "Persons")
}
