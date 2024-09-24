package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Person struct {
	ID        string   `json:"id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}

type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

var people []Person
var idCount = 4

func main() {
	// Seed in-memory
	people = append(people, Person{ID: "1", Firstname: "John", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}})
	people = append(people, Person{ID: "2", Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "City Z", State: "State Y"}})
	people = append(people, Person{ID: "3", Firstname: "Francis", Lastname: "Sunday"})

	router := mux.NewRouter()
	router.HandleFunc("/contact", GetPeople).Methods("GET")
	router.HandleFunc("/contact", CreatePerson).Methods("POST")
	router.HandleFunc("/contact/{id}", GetPerson).Methods("GET")
	router.HandleFunc("/contact/{id}", DeletePerson).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func GetPeople(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(people)
}

func GetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": "Person you are trying to find doesn't exist"}
	json.NewEncoder(w).Encode(response)
}

func CreatePerson(w http.ResponseWriter, r *http.Request) {
	var person Person
	err := json.NewDecoder(r.Body).Decode(&person)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{"message": "Person you are trying to find doesn't exist"}
		json.NewEncoder(w).Encode(response)
		return
	}

	person.ID = strconv.Itoa(idCount)
	people = append(people, person)
	idCount++
}

func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	for i, item := range people {
		if item.ID == params["id"] {
			people = append(people[:i], people[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": "Person you are trying to delete doesn't exist"}
	json.NewEncoder(w).Encode(response)
}
