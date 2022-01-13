package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
)

// we are hardtyping the different types of bad item submissions
// that have to do with missing properties
type BadItemNoCode struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type BadItemNoName struct {
	PID   string  `json:"pid"`
	Price float64 `json:"price"`
}

type BadItemNoPrice struct {
	PID  string `json:"pid"`
	Name string `json:"name"`
}

// this will check the response status code of a request and log unexpected values
// in: actStatus -- the actual received status code
//     expStatus -- the expected status code
//     t -- the testing.T object
// out: void
func checkStatus(actStatus int, expStatus int, t *testing.T, caller string) {
	if expStatus != actStatus {
		t.Errorf("%v -- handler returned wrong status code: actual - %v | expected - %v",
			caller, actStatus, expStatus)
	}
}

// this will check an error object passed in and go fatal if the error isn't null
// in: err -- the error object, whether null or not
//     t -- the testing.T object
// out: void
func checkError(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}

// checks a response for a json formatting error, occurs if the API is returning incorrect objects
// in: err -- the error, whether null or not
//     respRecorder -- the response recorder object (this object holds the response body)
//     expType -- whatever the calling context inputs, the type name given to the object it expects
//     t -- the testing.T object
// out: void
func checkResponseError(err error, respRecorder *httptest.ResponseRecorder, expType string, t *testing.T) {
	if err != nil {
		t.Errorf(`the response didn't format the JSON properly - should be ` + expType + ` type, but is: ` + respRecorder.Body.String())
	}
}

// checks a response for a json formatting error, occurs if the API is returning incorrect objects
// in: err -- the error, whether null or not
//     respRecorder -- the response recorder object (this object holds the response body)
//     expType -- whatever the calling context inputs, the type name given to the object it expects
//     t -- the testing.T object
// out: void
func recordResponse(action func(http.ResponseWriter, *http.Request), request *http.Request, t *testing.T) *httptest.ResponseRecorder {
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(action)
	handler.ServeHTTP(responseRecorder, request)
	return responseRecorder
}

//
func setMuxVars(req *http.Request, key string, value string) *http.Request {
	vars := map[string]string{
		key: value,
	}
	return mux.SetURLVars(req, vars)
}

// used to create and send a request for GET /inventory
// processes the response to check for errors and returns the current inventory in the API
func getInventoryReq(t *testing.T) []Item {
	req, err := http.NewRequest("GET", "/inventory", nil)
	checkError(err, t)

	respRecorder := recordResponse(getInventory, req, t)
	checkStatus(respRecorder.Code, http.StatusOK, t, "getInventoryReq")
	var Items []Item

	err = json.NewDecoder(respRecorder.Body).Decode(&Items)
	checkResponseError(err, respRecorder, "[]Item", t)

	return Items
}

// used to create and send a request for GET /inventory/{searchValue}
// checks for errors and returns the requested item from the API
func getItemReq(searchValue string, t *testing.T) Item {

	req, err := http.NewRequest("GET", "/inventory/"+searchValue, nil)
	checkError(err, t)
	req = setMuxVars(req, "searchValue", searchValue)

	respRecorder := recordResponse(getItem, req, t)
	checkStatus(respRecorder.Code, http.StatusOK, t, "getItemReq")

	var Item Item

	err = json.NewDecoder(respRecorder.Body).Decode(&Item)
	checkResponseError(err, respRecorder, "Item", t)

	return Item
}

// used to create and send a request for GET /inventory/{searchValue} with a bad searchValue
// checks for the expected 404 error
func getItemNotFoundReq(searchValue string, t *testing.T) {
	req, err := http.NewRequest("GET", "/inventory/"+searchValue, nil)
	checkError(err, t)
	req = setMuxVars(req, "searchValue", searchValue)

	respRecorder := recordResponse(getItem, req, t)

	checkStatus(respRecorder.Code, http.StatusNotFound, t, "getItemNotFoundReq")
}

// used to create and send a request for DELETE /inventory/{pid}
// checks for errors and returns the current inventory in the API
func deleteItemReq(pid string, t *testing.T) []Item {

	req, err := http.NewRequest("DELETE", "/inventory/"+pid, nil)
	checkError(err, t)
	req = setMuxVars(req, "pid", pid)

	respRecorder := recordResponse(deleteItem, req, t)
	checkStatus(respRecorder.Code, http.StatusOK, t, "deleteItemReq")

	var items []Item

	err = json.NewDecoder(respRecorder.Body).Decode(&items)
	checkResponseError(err, respRecorder, "Item", t)

	return items
}

// used to create and send a request for DELETE /inventory/{pid} with a bad pid
// checks for the expected 404 error
func deleteItemNotFoundReq(badPID string, t *testing.T) {
	req, err := http.NewRequest("GET", "/inventory/"+badPID, nil)
	checkError(err, t)
	req = setMuxVars(req, "pid", badPID)

	respRecorder := recordResponse(deleteItem, req, t)

	checkStatus(respRecorder.Code, http.StatusNotFound, t, "deleteItemNotFoundReq")
}

func addItemReq(item Item, t *testing.T) []Item {
	// we must marshall the provided golang object into a json object
	body, err := json.Marshal(item)
	checkError(err, t)

	// we have to convert body to bytes.NewBuffer before passing into the request creator
	req, err := http.NewRequest("GET", "/inventory/addItem", bytes.NewBuffer(body))
	checkError(err, t)

	// at this point our golang object and our request didn't have any runtime errors
	// we want to have the recording response available because it has the .Code and .Body
	respRecorder := recordResponse(addItem, req, t)
	checkStatus(respRecorder.Code, http.StatusOK, t, "addItemReq")
	var Items []Item

	//decode the response into the expected golang object (or array of objects)
	err = json.NewDecoder(respRecorder.Body).Decode(&Items)
	checkResponseError(err, respRecorder, "[]Item", t)

	//return the inventory to the calling context so it can keep track of its own testing progress
	return Items
}

func addItemsReq(items []Item, t *testing.T) []Item {

	body, err := json.Marshal(items)
	checkError(err, t)

	req, err := http.NewRequest("GET", "/inventory/addItems", bytes.NewBuffer(body))
	checkError(err, t)

	respRecorder := recordResponse(addItems, req, t)
	checkStatus(respRecorder.Code, http.StatusOK, t, "addItemsReq")
	var Items []Item

	err = json.NewDecoder(respRecorder.Body).Decode(&Items)
	checkResponseError(err, respRecorder, "[]Item", t)

	return Items
}

func addBadItemAllCasesReq(item Item, t *testing.T) {
	badCases := []string{"No Name", "No Code", "No Price", "Bad PID format", "Not Unique PID"}

	for _, badCase := range badCases {
		body, err := json.Marshal(_buildBadItem(badCase, item))
		checkError(err, t)

		req, err := http.NewRequest("GET", "/inventory/addItem", bytes.NewBuffer(body))
		checkError(err, t)

		respRecorder := recordResponse(addItem, req, t)
		checkStatus(respRecorder.Code, http.StatusBadRequest, t, "addBadItemAllCasesReq")
	}
}

func _buildBadItem(badCase string, item Item) interface{} {

	switch badCase {
	case "No Name":
		return BadItemNoName{
			PID:   item.PID,
			Price: item.Price,
		}
	case "No Code":
		return BadItemNoCode{
			Name:  item.Name,
			Price: item.Price,
		}
	case "No Price":
		return BadItemNoPrice{
			PID:  item.PID,
			Name: item.Name,
		}
	case "Bad PID format":
		item.PID += "-"
		return item
	case "Not Unique PID":
		item.PID = "E5T6-9UI3-TH15-QR88" // this is the Peach PID
		return item
	}
	// should never reach this next line of code, but what the compiler wants, the compiler gets
	return Item{}
}

// we frequently evaluate expected inventory with actual, when relevant
// after we check and log any inconcsistencies, we return actual so it can be
// set to expInventory so we can proceed with the rest of the tests regardless
func compareActualWithExpected(actual []Item, expected []Item, t *testing.T, checkpointNumber string) []Item {
	if len(actual) != len(expected) {
		t.Errorf("%v -- actual and expected inventory sizes are not equal: len(actual)- %v | len(expected) - %v",
			checkpointNumber, len(actual), len(expected))
	} else {
		for i := range actual {
			if actual[i] != expected[i] {
				t.Errorf("%v -- actual and expected inventories differ at i == %v: actual[i] - %v | expected[i] - %v",
					checkpointNumber, i, actual[i], expected[i])
			}
		}
	}
	return actual
}

func TestAPI(t *testing.T) {
	// 0. expInventory ===================================================================================================
	// expInventory is initialized identically to actual in the main.go file``
	expInventory := []Item{
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


	// 1. getInventory ===================================================================================================
	t.Log("1. getInventory")

	// check that the getInventory endpoint is working correctly after
	// we have compared the actual with expected, we've already logged
	// discrepencies so we can set expected to actual and continue
	actInventory := getInventoryReq(t)
	expInventory = compareActualWithExpected(actInventory, expInventory, t, "1")



	// 2. delete Lettuce =================================================================================================
	t.Log("2. delete Lettuce")

	// we will frequently throw in various upper/lower cases to test case insensitivity.
	// steps 2-5 will alter the expInventory to match what we expected actInventory
	// to look like when the next response comes back, since we deleted lettuce and
	// lettuce is the first item of the inventory array, we remove that from expIntentory
	expInventory = expInventory[1:]
	actInventory = deleteItemReq("a12T-4Gh7-QPl9-3n4M", t) // also tests case insensitivity
	expInventory = compareActualWithExpected(actInventory, expInventory, t, "2")



	// 3. add Tomato =====================================================================================================
	t.Log("3. add Tomato")

	// we can test the 2 decimal requirement by sending in 3 decimals and expecting to only get 2 back
	tomato_with_3_decimals := Item{
		PID:   "M4N5-F0C3-F4gk-si00",
		Name:  "Tomato",
		Price: 3.355,
	}
	tomato_expected := Item{
		PID:   "M4N5-F0C3-F4gk-si00",
		Name:  "Tomato",
		Price: 3.35,
	}
	
	expInventory = append(expInventory, tomato_expected)
	actInventory = addItemReq(tomato_with_3_decimals, t)
	expInventory = compareActualWithExpected(actInventory, expInventory, t, "3")



	// 4. add Pickle, Broccoli, Chicken Breast ============================================================================
	t.Log("4. add Pickle, Broccoli, Chicken Breast")

	// we can just initialize the expected array (2_decimal) and the submitted (3_decimal) array side by side
	items_with_3_decimals := []Item{
		{
			PID:   "F4J6-D4M2-J0G5-G3E5",
			Name:  "Pickle",
			Price: 1.299,
		},
		{
			PID:   "0g44-gm33-4jf9-FGM4",
			Name:  "Broccoli",
			Price: 2.208,
		},
		{
			PID:   "1A2S-3F5G-6HJ7-4R6V",
			Name:  "Chicken Breast",
			Price: 6.493,
		},
	}
	items_expected := []Item{
		{
			PID:   "F4J6-D4M2-J0G5-G3E5",
			Name:  "Pickle",
			Price: 1.3,
		},
		{
			PID:   "0g44-gm33-4jf9-FGM4",
			Name:  "Broccoli",
			Price: 2.21,
		},
		{
			PID:   "1A2S-3F5G-6HJ7-4R6V",
			Name:  "Chicken Breast",
			Price: 6.49,
		},
	}
	expInventory = append(expInventory, items_expected...)
	actInventory = addItemsReq(items_with_3_decimals, t)
	expInventory = compareActualWithExpected(actInventory, expInventory, t, "4")



	// 5. delete Broccoli, Gala Apple, Pepper ================================================================================
	t.Log("5. delete Broccoli, Gala Apple, Pepper")

	// at this point the expected inventory is ['peach', 'pepper', 'apple', 'tomato', 'pickle', 'broccoli', 'chicken']
	PIDs := []string{"0g44-gm33-4jf9-FGM4", "YRT6-72AS-K736-L4AR", "TQ4C-VV6T-75ZX-1RMR"}
	expectedIndexes := []int{5, 1, 1} // expected index of broccoli is 5, then pepper 1, and apple 1 since we remove pepper
	for i, pid := range PIDs {
		expInventory = append(expInventory[:expectedIndexes[i]], expInventory[expectedIndexes[i]+1:]...)
		actInventory = deleteItemReq(pid, t)
		expInventory = compareActualWithExpected(actInventory, expInventory, t, "5"+strconv.Itoa(i))
	}



	// 6. Using addItem... add BAD Potato due to -- (No Name), then (No Code), ===============================================
	t.Log("6. Using addItem... add BAD Potato due to -- (No Name), then (No Code)")
	//     then (No Price), then (Bad PID format),
	//     then (Not Unique PID)

	// we make a good potato and give it to the badItemReq func, which will
	// use our potato to create all the different kinds of bad potatoes
	good_potato := Item{
		PID:   "b6N3-C5X3-Z0F6-2K0J",
		Name:  "potato",
		Price: 0.49,
	}
	// we don't need to compare actual and expected cus we just expect to get 400's back
	// on all 5 of the different requests that addBadItemAllCasesReq will send out
	addBadItemAllCasesReq(good_potato, t)



	// 7. good get Item by name ===============================================================================================
	t.Log("7. good get Item by name")

	tomato_actual := getItemReq("ToMaTo", t)
	tomato_expected = expInventory[1]
	// just initialize arrays using the singular tomato item we get back from getItemReq
	// so we can reuse compareActualWithExpected's logic w/o writing more unnecessarily
	compareActualWithExpected([]Item{tomato_actual}, []Item{tomato_expected}, t, "6")

	// 8. good get Item by PID ================================================================================================
	t.Log("8. good get Item by PID")

	pickle_actual := getItemReq("F4j6-D4m2-j0G5-G3e5", t)
	pickle_expected := items_expected[0]
	// repeat use of compareActualWithExpected as was done above
	compareActualWithExpected([]Item{pickle_actual}, []Item{pickle_expected}, t, "7")



	// 9. bad get Item by name, 404 ===========================================================================================
	t.Log("9. bad get Item by name, 404")
	getItemNotFoundReq("tomatoe", t) // wrong spelling of tomato supplied, expect a 404



	// 10. bad get Item by PID, 404 ===========================================================================================
	t.Log("10. bad get Item by name, 404")
	getItemNotFoundReq(pickle_expected.PID+"-", t) // malformed pickle PID, expect a 404



	// 11. bad get Item that was deleted already, 404 =========================================================================
	t.Log("11. bad get Item that was deleted already, 404")
	getItemNotFoundReq(items_expected[1].PID, t) // broccoli doesn't exist anymore, expect a 404



	// 12. bad delete Item using PID that does not exist =========================================================================
	t.Log("12. bad delete Item using PID that does not exist")
	deleteItemNotFoundReq("Th1s-P1Dd-N0t3-X1ST", t) // valid PID format but it doesn't exist
}
