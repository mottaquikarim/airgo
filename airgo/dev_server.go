package airgo

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func DevServer(dirname string, port int) {
	// start static file server for debug purposes
	// TODO: this is only for dev, on "prod" we must
	// not run this server
	fs := http.FileServer(http.Dir(fmt.Sprintf("./%s", dirname)))
	http.Handle("/", fs)

	log.Printf("Listening on :%d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
