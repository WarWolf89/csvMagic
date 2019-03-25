# CSV Magic
This App is built to handle huge csv files, validate them given certain custom validators that you define on the structs themselves, and then save them to mongo along with the file meta.
## Getting started
* fire up a mongodb instance on local, this can be done via a docker container, or installing directly. Default values are set for the APP. which means 127.0.0.1:27017 and DB = "local"
* run `go get ./...` to install all dependencies for the app from the root directory
* `go run main.go` to start the mux server 
## Useful commands 
If you want to upload full CSVs, you can find samples in the resources\csv folder.
`curl http://localhost:8080/csvUpload -F "file=@ENTERYOURPATHHERE/csvMagic/resources/csv/fullTest.csv" -vvv`
All other endpoints are self-explanatory, you can call them from Postman or curl. 
## API 
### {specifiedPort}/csvUpload POST
Reads the csv, parses the csv, validates entries, tries to fix, generates metadata, saves all values to mongo.
* request body:
the csv file you're uploading
* sample response (response as a JSON file): 
` HTTP/1.1 201 Created
  Content-Disposition: attachment; filename=fullTest
  Content-Type: application/json
  Date: Mon, 25 Mar 2019 14:45:54 GMT
  Content-Length: 181
 {
  "file_id": "5c98e9a18261caa977dc2079",
  "Name": "fullTest",
  "stats": {
   "fixed": 4,
   "processed": 1000,
   "unfixable": 796,
   "valid": 200
  },
  "execution_time": 0.062052487
 }` 
 ### {specifiedPort}/valSingle/{num} GET
 Validates a single string for any phonenumber enveloped in the param.
 * request
 Enter number you would like to validate in the pathparam(there's a very specific reason why this is not a POST request)
* sample response
`{
     "IsValid": false,
     "ValErr": "there was no valid number in the corrupted field"
 }`
 ### {specifiedPort}/filedata/{fid} GET
 Fetches metadata for any previously uploaded file. 
 * request
 {fid} has to be an existing objectId for a file that has been previously uploaded. If an incorrect value is entered either Hex error for malformed string is thrown or a mongo error for no document found
 * sample response 
 `{
      "file_id": "5c98e9a18261caa977dc2079",
      "Name": "fullTest",
      "stats": {
          "fixed": 4,
          "processed": 1000,
          "unfixable": 796,
          "valid": 200
      },
      "execution_time": 0.062052487
  }`
  ## Testing
  Run `go test ./... -v` from root directory to run all tests. Validators, fixers and reader is fully tested. Test csv file is in the resources folder.
  ## Limitations and faults that I know of
  These are the thing that I will correct in the future, I've delayed submitting this as long as it is so I don't want to postpone it any further:
  * if non-existing fileID is in pathparam, 404 should be returned instead of 500
  * The fixer could be refactored heavily
  * the "//D" regex is most likely not necessary
  * endpoints should be unit tested
  * pathparams should be validated in the handlers
  * endpoints should be created in a separate go file like App.go and reused in tests
  * more fixing logic could be added for area codes in regex group etc.
  * godep should be used and everything should be shipped with vendor dir
  * app could be moved to container with mongo, set up a docker-compose and just fire everything up
  * unique key could be set up for the csv rows in mongo for either ID or the number itself so that no duplicate entries are inserted
  * better error handling
  * better logging
  