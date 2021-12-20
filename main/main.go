package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
)

type Quality string

//Qualities tell us what unique beneficial attributes a product has

const (
	Organic                    Quality = "organic"
	GrassFed                   Quality = "grass_fed"
	GMOFree                    Quality = "gmo_free"
	PastureRaised              Quality = "pasture_raised"
	FreeRange                  Quality = "free_range"
	GlutenFree                 Quality = "gluten_free"
	Vegetarian                 Quality = "vegetarian"
	Vegan                      Quality = "vegan"
	AntibioticFree             Quality = "antibiotic_free"
	HormoneFree                Quality = "hormone_free"
	PreservativeFree           Quality = "preservative_free"
	ArtificialFlavorFree       Quality = "artificial_flavor_free"
	ArtificialPreservativeFree Quality = "artificial_preservative_free"
	NutFree                    Quality = "nut_free"
	LocallySourced             Quality = "locally_sourced"
	FamilyOwned                Quality = "family_owned"
)

// A PID is a 16 digit alphanumeric product ID required in each createItems request
type Item struct {
	PID       string    `json:"pid"`
	Name      string    `json:"name"`
	Qualities []Quality `json:"qualities"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"quantity"`
}

// A PID is required to find the item
// A QuantityIn is required to determine how many units to add
type UpdateItemQuantityReq struct {
	PID        string `json:"pid"`
	QuantityIn int    `json:"quantity_in"`
}

// A PID is required to find the item
// A Quality array and Price is required to update the item info
type UpdateItemInfoReq struct {
	PID       string    `json:"pid"`
	Qualities []Quality `json:"qualities"`
	Price     float64   `json:"price"`
}

var inventory []Item

// some users just want to see the inventory directly
func getInventory(w http.ResponseWriter, r *http.Request) {
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
	fmt.Println("Function Called: getItemByID()")

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
}

// This is provided to other functions to grab objects before operating on the 'inventory'
// 'field' is the key/property we want to search with, 'value' is the search value
func _getItemAt(field string, value string) (bool, Item) {
	for _, item := range inventory {
		itemValue := ""
		switch field {
		case "name":
			itemValue = item.Name
			break
		case "pid":
			itemValue = item.PID
			break
		}
		//strings.ToUpper to ensure our PIDs and Names are case-insensitive
		if strings.ToUpper(itemValue) == strings.ToUpper(value) {
			return true, item
		}
	}
	return false, Item{}
}

// If an array is not submitted a 400 is returned
// the 16 digit product id is received in the request to create a new item
func createItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Function Called: createItems()")

	var createItemsReq []Item
	err := json.NewDecoder(r.Body).Decode(&createItemsReq)
	if err != nil {
		log.Println("error: ", err)
		w.WriteHeader(http.StatusBadRequest) // return 400 Bad Request
		return
	}

	for _, newItem := range createItemsReq {
		// if the item already exists we simply act as though it is
		// having its qualities, price, and quantity updated
		if success, item := _getItemAt("name", newItem.Name); success {
			_updateItem(Item{
				PID:       item.PID,
				Name:      item.Name,
				Qualities: newItem.Qualities,
				Price:     newItem.Price,
				Quantity:  newItem.Quantity + item.Quantity,
			})
		} else {
			inventory = append(inventory, newItem)
		}
	}
	w.WriteHeader(http.StatusOK) //return 200 OK
	json.NewEncoder(w).Encode(inventory)
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
		log.Println("404 error - deleteItem(): Cannot find PID: ", pid)
		w.WriteHeader(http.StatusNotFound) // return 404 Not Found
		return
	}
}

// updating item info is more of an employee use case
// you can only update one item at a time in this use case
// this reuses the quantity and name/pid, but updates the qualities/price
func updateItemInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Function Called: updateItemInfo()")

	params := mux.Vars(r)
	pid := params["pid"]

	var itemReq UpdateItemInfoReq
	err := json.NewDecoder(r.Body).Decode(&itemReq)
	if err != nil {
		log.Println("error: ", err)
		w.WriteHeader(http.StatusBadRequest) // return 400 Bad Request
		return
	}

	if success, item := _getItemAt("pid", pid); success {
		newItem := Item{
			PID:       pid,
			Name:      item.Name,
			Qualities: itemReq.Qualities,
			Price:     itemReq.Price,
			Quantity:  item.Quantity,
		}
		// if _updateItem returns true, we found the item and deleted it, else 404
		if success := _updateItem(newItem); success {
			w.WriteHeader(http.StatusOK) //return 200 OK
			json.NewEncoder(w).Encode(inventory)
		} else {
			//some 500 error, we know its not a 404
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		log.Println("404 error - updateItemInfo(): Cannot find PID: ", pid)
		w.WriteHeader(http.StatusNotFound) // return 404 Not Found
		return
	}
}

// this is more of a supplier/stocking employee use case
// an array is used to allow us to update all the items in a delivery at once
func increaseInventory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Function Called: increaseInventory()")

	var quantityReq []UpdateItemQuantityReq
	err := json.NewDecoder(r.Body).Decode(&quantityReq)
	if err != nil {
		log.Println("error: ", err)
		w.WriteHeader(http.StatusBadRequest) // return 400 Bad Request
		return
	}

	newItems := make([]Item, 0)

	// we make sure all the items exist before operating on the inventory
	for _, q := range quantityReq {
		if success, item := _getItemAt("pid", q.PID); success {
			newItems = append(newItems, Item{
				PID:       q.PID,
				Name:      item.Name,
				Qualities: item.Qualities,
				Price:     item.Price,
				Quantity:  item.Quantity + q.QuantityIn,
			})
		} else {
			log.Println("404 error - increaseInventory(): Cannot find PID: ", q.PID)
			w.WriteHeader(http.StatusNotFound) // return 404 Not Found
			return
		}
	}

	//we know all the items exist
	for _, newItem := range newItems {
		if success := _updateItem(newItem); !success {
			//some 500 error, we know its not a 404
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK) //return 200 OK
	json.NewEncoder(w).Encode(inventory)
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

// return true/false so we can report a 404 Not Found from within the calling function
func _updateItem(newItem Item) bool {
	for index, item := range inventory {
		//strings.ToUpper to ensure our PIDs are case-insensitive
		if strings.ToUpper(item.PID) == strings.ToUpper(newItem.PID) {
			// Re-assign item in slice
			inventory[index] = newItem
			return true
		}
	}
	return false
}

func handleRequests() {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/inventory", getInventory).Methods("GET")
	router.HandleFunc("/inventory", createItems).Methods("POST")
	router.HandleFunc("/inventory", increaseInventory).Methods("PUT")

	//searchValue could be a name, or it could be a product ID
	router.HandleFunc("/inventory/{searchValue}", getItem).Methods("GET")
	router.HandleFunc("/inventory/{pid}", deleteItem).Methods("DELETE")
	router.HandleFunc("/inventory/{pid}", updateItemInfo).Methods("PUT")
	log.Println("Running on localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func main() {
	initialInventory := []Item{
		{
			PID:  "A12T-4GH7-QPL9-3N4M",
			Name: "Lettuce",
			Qualities: []Quality{
				LocallySourced,
				GMOFree,
				Organic,
			},
			Price: 1.99,
		},
		{
			PID:  "E5T6-9UI3-TH15-QR88",
			Name: "Peach",
			Qualities: []Quality{
				Organic,
				GMOFree,
			},
			Price: 1.79,
		},
		{
			PID:  "YRT6-72AS-K736-L4AR",
			Name: "Green Pepper",
			Qualities: []Quality{
				LocallySourced,
				GMOFree,
			},
			Price: 3.50,
		},
		{
			PID:  "TQ4C-VV6T-75ZX-1RMR",
			Name: "Gala Apple",
			Qualities: []Quality{
				FamilyOwned,
				Organic,
			},
			Price: 1.49,
		},
	}
	inventory = append(inventory, initialInventory...)
	handleRequests()
}
