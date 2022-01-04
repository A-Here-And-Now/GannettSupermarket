package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// A PID is a 16 digit alphanumeric product ID required in each createItems request
type Item struct {
	PID   string  `json:"pid"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var inventory []Item

// some users just want to see the inventory directly
func getInventory(w http.ResponseWriter, r *http.Request) {
	if len(inventory) == 0 {
		inventory = []Item{
			{
				PID:   "A12T-4GH7-QPL9-3N4M",
				Name:  "Lettuce",
				Price: 3.46,
			},
			{
				PID:   "E5T6-9UI3-TH15-QR88",
				Name:  "Peach",
				Price: 2.99,
			},
			{
				PID:   "YRT6-72AS-K736-L4AR",
				Name:  "Green Pepper",
				Price: 0.79,
			},
			{
				PID:   "TQ4C-VV6T-75ZX-1RMR",
				Name:  "Gala Apple",
				Price: 3.59,
			},
		}
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Function Called: getInventory()")

	w.WriteHeader(http.StatusOK) //return 200 OK
	json.NewEncoder(w).Encode(inventory)
}

// some users just want to look something up by name
// other more sophisticated users such as suppliers, exec-staff or employees
// can look up by product ID
func getItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Function Called: getItem()")

	params := mux.Vars(r)
	searchValue := params["searchValue"]

	// if our product ID format is matched, we have a PID, otherwise a name
	regex := regexp.MustCompile("^[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}$")
	isPID := regex.MatchString(searchValue)
	for _, item := range inventory {
		var itemValue = item.Name
		if isPID {
			itemValue = item.PID
		}
		//strings.ToUpper to ensure our PIDs and Names are case-insensitive
		if strings.ToUpper(itemValue) == strings.ToUpper(searchValue) {
			w.WriteHeader(http.StatusOK) //return 200 OK
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	// item not found, return a response accordingly
	log.Println("404 error - getItem(): Cannot find item: ", searchValue)
	w.WriteHeader(http.StatusNotFound) // return 404 Not Found
}

// If an array is not submitted a 400 is returned
// the 16 digit product id is received in the request to create a new item
func addItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Function Called: addItem()")

	var addItemReq Item
	err := json.NewDecoder(r.Body).Decode(&addItemReq)
	if err != nil || _checkAddItem(addItemReq) {
		// the client didn't format the JSON properly - should be a single object of Item type
		log.Println(`Could not parse the format of item received.
			Please provide a JSON object with 'price', 'name' and 'pid': `, err)
		w.WriteHeader(http.StatusBadRequest) // return 400 Bad Request
		return
	}

	// truncate the float64 provided to two decimals to ensure prices don't have more than necessary
	addItemReq.Price, err = strconv.ParseFloat(fmt.Sprintf("%.2f", addItemReq.Price), 64)
	// now we know its safe to add the items to inventory because they have been validated for format
	inventory = append(inventory, addItemReq)

	w.WriteHeader(http.StatusOK) //return 200 OK
	json.NewEncoder(w).Encode(inventory)
}

func _checkAddItem(item Item) bool {
	if item.Price == 0.00 || item.Name == "" || item.PID == "" {
		return true
	} else {
		regex := regexp.MustCompile("^[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}$")
		isValidPID := regex.MatchString(item.PID)
		if !isValidPID {
			return true
		} else {
			for _, oldItem := range inventory {
				//strings.ToUpper to ensure our PIDs and Names are case-insensitive
				if strings.ToUpper(oldItem.PID) == strings.ToUpper(item.PID) {
					log.Printf("adding PID that already exists -- %v", item.PID)
					return true
				}
			}
		}
	}
	return false
}

// If an array is not submitted a 400 is returned
// the 16 digit product id is received in the request to create a new item
func addItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Function Called: addItems()")

	var createItemsReq []Item
	err := json.NewDecoder(r.Body).Decode(&createItemsReq)
	if err != nil || _checkAddItems(createItemsReq) {
		// the client didn't format the JSON properly - should be a single object of Item type
		log.Println(`Could not parse the format of items received.
			Please provide a JSON array of objects with 'price', 'name' and 'pid': `, err)
		w.WriteHeader(http.StatusBadRequest) // return 400 Bad Request
		return
	}

	for i, item := range createItemsReq {
		// truncate the float64 provided to two decimals to ensure prices don't have more than necessary
		createItemsReq[i].Price, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", item.Price), 64)
	}

	// now we know its safe to add the items to inventory because they have been validated for format
	inventory = append(inventory, createItemsReq...)

	w.WriteHeader(http.StatusOK) //return 200 OK
	json.NewEncoder(w).Encode(inventory)
}

func _checkAddItems(items []Item) bool {
	for _, item := range items {
		if _checkAddItem(item) {
			return true
		}
	}
	return false
}

// deleting items occurs only one at a time
func deleteItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Function Called: deleteItem()")

	params := mux.Vars(r)
	pid := params["pid"]

	// if _deleteItemAt returns true, we found the item and deleted it, else 404
	if success := _deleteItemAt(pid); success {
		w.WriteHeader(http.StatusOK) //return 200 OK
		json.NewEncoder(w).Encode(inventory)
	} else {
		// item not found - return a response accordingly
		log.Println("404 error - deleteItem(): Cannot find PID: ", pid)
		w.Write([]byte("Could not find item in inventory: " + pid))
		w.WriteHeader(http.StatusNotFound) // return 404 Not Found
		return
	}
}

// return true/false so we can report a 404 Not Found from within the calling function
func _deleteItemAt(pid string) bool {
	for index, item := range inventory {
		//strings.ToUpper to ensure our PIDs are case-insensitive
		if strings.ToUpper(item.PID) == strings.ToUpper(pid) {
			// Delete item from slice
			inventory = append(inventory[:index], inventory[index+1:]...)
			return true
		}
	}
	return false
}

func handleRequests() {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/inventory", getInventory).Methods("GET")
	router.HandleFunc("/inventory/addItems", addItems).Methods("POST")
	router.HandleFunc("/inventory/addItem", addItem).Methods("POST")

	//searchValue could be a name, or it could be a product ID
	router.HandleFunc("/inventory/{searchValue}", getItem).Methods("GET")
	router.HandleFunc("/inventory/{pid}", deleteItem).Methods("DELETE")
	log.Println("Running on localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func main() {
	initialInventory := []Item{
		{
			PID:   "A12T-4GH7-QPL9-3N4M",
			Name:  "Lettuce",
			Price: 3.46,
		},
		{
			PID:   "E5T6-9UI3-TH15-QR88",
			Name:  "Peach",
			Price: 2.99,
		},
		{
			PID:   "YRT6-72AS-K736-L4AR",
			Name:  "Green Pepper",
			Price: 0.79,
		},
		{
			PID:   "TQ4C-VV6T-75ZX-1RMR",
			Name:  "Gala Apple",
			Price: 3.59,
		},
	}
	inventory = append(inventory, initialInventory...)
	handleRequests()
}
