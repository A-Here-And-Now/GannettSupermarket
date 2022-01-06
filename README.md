# Endpoints

### GET /inventory
Returns the current state of the grocery's inventory.

##### Body
No response body required

##### Error Codes
No errors codes at this endpoint


### POST /inventory/addItems
Adds multiple items to the inventory. It returns the inventory after adding the items.

##### Body
a JSON array of valid grocery "Item" objects

example input: 
[
{
    "PID":"A1B2-C3D4-E5F6-G7H8",
    "Name": "Pear",
    "Price": 1.33,
},
{
    "PID":"Z1X2-C3V4-B5N6-M7K8",
    "Name": "Orange",
    "Price": 0.89,
}
]

##### Error Codes
400 - bad json format, missing item properties, or bad PID or PID already exists


### POST /inventory/addItem
Performs exactly the same as the addItems endpoint but is only intended
to be exercised in the case that we are adding one item.

##### Body
The body should be a JSON formatted "Item" object

example input: 
{
    "PID":"A1B2-C3D4-E5F6-G7H8",
    "Name": "Pear",
    "Price": 1.33,
}

##### Error Codes
400 - bad json format, missing item properties, or bad PID or PID already exists


### GET /inventory/{searchValue}
Returns the first item in the inventory that matches the searchValue, if any.
The searchValue is retrieved from the url and can be either an item name or PID.

##### Body
No response body required

##### Error Codes
404 - item not found with that PID/Name


### DELETE /inventory/{pid}
Deletes the item that matches the given pid. 
Only a PID is valid at this endpoint.

##### Body
No response body required

##### Error Codes
404 - item not found with that PID/Name



# For Developers

After introducing yourself to all of the endpoints using the above documentation, the following will further assist you:

### Running the project
To run the API, navigate to the main directory of this project and run the following command: 
`go run main.go`

To run the api_test.go file, from the main directory run:
`go test`

### Potential Improvements
Here is a running list of potential improvements to be made to the project
