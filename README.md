# Endpoints

### GET /inventory
Returns the current state of the grocery's inventory.

##### Body
No request body required

##### Error Codes
No errors codes at this endpoint


### POST /inventory/addItems
Adds multiple items to the inventory. It returns the inventory after adding the items.

##### Body
a JSON array of valid grocery "Item" objects

example input:<br>
[<br>
{<br>
    "PID":"A1B2-C3D4-E5F6-G7H8",<br>
    "Name": "Pear",<br>
    "Price": 1.33,<br>
},<br>
{<br>
    "PID":"Z1X2-C3V4-B5N6-M7K8",<br>
    "Name": "Orange",<br>
    "Price": 0.89,<br>
}<br>
]<br>

##### Error Codes
400 - bad json format, missing item properties, or bad PID or PID already exists


### POST /inventory/addItem
Performs exactly the same as the addItems endpoint but is only intended
to be exercised in the case that we are adding one item.
It returns the inventory after adding the item.
##### Body
The body should be a JSON formatted "Item" object

example input:<br>
{<br>
    "PID":"A1B2-C3D4-E5F6-G7H8",<br>
    "Name": "Pear",<br>
    "Price": 1.33,<br>
}<br>

##### Error Codes
400 - bad json format, missing item properties, or bad PID or PID already exists


### GET /inventory/{searchValue}
Returns the first item in the inventory that matches the searchValue, if any.
The searchValue is retrieved from the url and can be either an item name or PID.

##### Body
No request body required

##### Error Codes
404 - item not found with that PID/Name


### DELETE /inventory/{pid}
Deletes the item that matches the given pid. 
Only a PID is valid at this endpoint.
It return the inventory after deleting the item.

##### Body
No request body required

##### Error Codes
404 - item not found with that PID



# For Developers

After introducing yourself to all of the endpoints using the above documentation, the following will further assist you:

### Running the project
To run the API, navigate to the main directory of this project and run the following command: 
`go run main.go`

To run the api_test.go file, from the main directory run:
`go test`

### Code GOTCHAS and recommendations
* There are a lot of helpful comments in the code. I recommend you read through all of a function's comments if you don't understand how that function works.
* Remember to write headers before you call w.Write or your header codes won't be included in the response because w.Write returns the response immediately after it runs.
* We allow users to perform the erroneous operation of submitting a price with more than 2 digits. We will simply round to the nearest 2nd digit to conform to proper price format.
* We use the gorilla/mux library for all our router needs (as well as setting URL variables in our api_test.go file) refer to their documentation here: https://pkg.go.dev/github.com/gorilla/mux

### Potential Improvements
Here is a running list of potential improvements to be made to the project

- [ ] Add a test step in the api_test.go that validates that a 404 is received at the DELETE endpoint if the given PID is not found in the inventory
- [ ] Change "Item" object to include a Quantity value that increments when an existing item name is added (requires that PIDs be stored in an array of strings that houses each individual Apple's own PID)
- [ ] Add a PUT endpoint for suppliers to add to the Quantity of an item to represent a supplier dropping off a shipment
- [ ] Add a PUT endpoint to allow for employees to add new item names/types to the inventory
- [ ] Add the Qualities array to Item to house a list of qualities that can be attached to each item (e.g. gluten-free or grass-fed, etc.)
- [ ] Add a PUT endpoint for employees that allows them to update an item in the inventory with new Quality entries to the Qualities array 
- [ ] Add POST endpoints where transactions with customers and suppliers can be posted (This requires an array of transactions be stored in a new storage variable, much like 'inventory')
- [ ] Add a GET endpoint that processes the NET profit from the accumulated transactions
- [ ] Add timestamps to each transaction so we can try to generate daily, weekly, monthly sales/profit reports
