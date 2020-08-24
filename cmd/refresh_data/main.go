package main

import (
	"fmt"
	"os"

	flag "github.com/namsral/flag"
	log "github.com/sirupsen/logrus"

	airgo "github.com/mottaquikarim/airgo/airgo"
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

	if err := airgo.RefreshData(airgo.RefreshOpts{
		ApiKey: *apiKey,
		BaseId: *baseId,
		Conf:   "config.json",
		Dest:   "data",
	}); err != nil {
		log.Fatal(err)
	}

	os.Exit(0) // redundant?
}
