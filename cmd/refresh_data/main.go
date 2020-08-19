package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	airtable "github.com/mottaquikarim/go-airtable"
	flag "github.com/namsral/flag"
	log "github.com/sirupsen/logrus"
)

// these are the key variables that store the input
// required to make this service run
var (
	fs     = flag.NewFlagSetWithEnvPrefix(os.Args[0], "STREAK", 0)
	apiKey = fs.String("api-key", "XXXXX", "Airtable API Key")
	baseId = fs.String("base-id", "XXXXX", "Airtable Base Id")
)

// helper function to hint to the user
// what the required key variables are
// (in this case, apiKey and baseId)
func usage() {
	fmt.Println("Usage: ./service.go [flags]")
	fs.PrintDefaults()
}

// this func kicks off the main implementation
// of the code
func main() {
	// first, parse the input and extract
	// apiKey and baseId variables
	fs.Usage = usage
	fs.Parse(os.Args[1:])

	// we create an account struct
	// that is passed in to airtable
	// module
	account := airtable.Account{
		ApiKey: *apiKey,
		BaseId: *baseId,
	}

	// TODO: we extract the `config.json` from
	// the /templates folder and use it to
	// fetch data from airtable

	// for now, just read the hardocded "Home" table
	var nameContent []airtable.Record
	var err error
	var opts = airtable.Options{
		MaxRecords: 100,
		View:       "Content",
	}
	mainSite := airtable.NewTable("Name", account)
	for account := 1; account < 11; account++ {
		fmt.Printf("%d Name", account)

	}
	if nameContent, err = mainSite.List(opts); err != nil {
		log.Printf("Error! %v", err)
	}
	// write to file now
	message, err := json.Marshal(nameContent)
	err = ioutil.WriteFile("data/index.json", message, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("data: %v", nameContent)

	log.Printf("Successfully completed")
	os.Exit(0) // redundant?

	
	

 }



