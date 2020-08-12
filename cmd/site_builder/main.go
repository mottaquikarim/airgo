package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"

	airtable "github.com/mottaquikarim/go-airtable"
	log "github.com/sirupsen/logrus"
)

func main() {
	// read content from db
	var content []byte
	var err error
	if content, err = ioutil.ReadFile("data/index.json"); err != nil {
		log.Fatal(err)
	}

	// marshal content into airtable.Record
	var home = []airtable.Record{}
	if err = json.Unmarshal(content, &home); err != nil {
		log.Printf("Failed to unmarshal json %v", err)
	}

	// write file through template
	// more info here:
	// https://www.htmlgoodies.com/beyond/reference/nesting-templates-with-go-web-programming.html
	t, err := template.ParseFiles("templates/layout.html.tmpl", "templates/index.html.tmpl")
	if err != nil {
		log.Print(err)
		return
	}

	// create or truncate index.html
	f, err := os.Create("static/index.html")
	defer f.Close()
	if err != nil {
		log.Println("create file: ", err)
		return
	}

	// populate template with data
	err = t.ExecuteTemplate(f, "layout", home)
	if err != nil {
		log.Print("execute: ", err)
		return
	}

	log.Printf("File contents: %v", home)

	// start static file server for debug purposes
	// TODO: this is only for dev, on "prod" we must
	// not run this server
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	log.Println("Listening on :80...")
	err = http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal(err)
	}
}
