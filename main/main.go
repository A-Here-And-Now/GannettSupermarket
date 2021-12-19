package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Item struct {
	UID   string  `json:"uid"`
	Name  string  `json:"name"`
	Desc  string  `json:"desc"`
	Price float64 `json:"price"`
}

var inventory []Item

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Endpoint called: homePage()")
}

func getInventory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Function Called: getInventory()")

	json.NewEncoder(w).Encode(inventory)
}

func createItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var items []Item
	err := json.NewDecoder(r.Body).Decode(&items)
	if err != nil {
		log.Println("error: ", err)
		os.Exit(1)
	}
	inventory = append(inventory, items...)
	json.NewEncoder(w).Encode(items)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	_deleteItemAtUid(params["uid"])

	json.NewEncoder(w).Encode(inventory)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	var item Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		log.Println("error: ", err)
		os.Exit(1)
	}

	_updateItemAtUid(params["uid"], item)

	json.NewEncoder(w).Encode(inventory)
}

func _deleteItemAtUid(uid string) {
	for index, item := range inventory {
		if item.UID == uid {
			// Delete item from slice
			inventory = append(inventory[:index], inventory[index+1:]...)
			break
		}
	}
}

func _updateItemAtUid(uid string, thing Item) {
	for index, item := range inventory {
		if item.UID == uid {
			// Re-assign item in slice
			inventory[index] = thing
			break
		}
	}
}

func handleRequests() {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homePage).Methods("GET")
	router.HandleFunc("/inventory", getInventory).Methods("GET")
	router.HandleFunc("/inventory", createItem).Methods("POST")
	router.HandleFunc("/inventory/{uid}", deleteItem).Methods("DELETE")
	router.HandleFunc("/inventory/{uid}", updateItem).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func main() {
	inventory = append(inventory, Item{
		UID:   "0",
		Name:  "Cheese",
		Desc:  "A fine block of cheese",
		Price: 4.99,
	})
	inventory = append(inventory, Item{
		UID:   "1",
		Name:  "Milk",
		Desc:  "Grass-fed milk",
		Price: 5.50,
	})
	handleRequests()
}
