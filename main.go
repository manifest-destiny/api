package main

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/emicklei/go-restful"
	"github.com/manifest-destiny/api/apidocs"
	"github.com/manifest-destiny/api/player"
)

const (
	tlsCert = "keys/tls_cert.crt"
	tlsKey  = "keys/tls_key.key"
)

var (
	dbInfo, port, cert, key string
	tls                     bool
)

func init() {

	dbInfo = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("DATABASE_SSL"))

	port = os.Getenv("API_PORT")
	if port == "" {
		port = "80"
	}

	// write TLS cert and key to file
	os.Remove(tlsKey)
	os.Remove(tlsCert)

	if os.Getenv("TLS_KEY_BASE64") == "" || os.Getenv("TLS_CERTIFICATE_BASE64") == "" {
		tls = false
	} else {
		tls = true
		tlsKeyB, err := base64.StdEncoding.DecodeString(os.Getenv("TLS_KEY_BASE64"))
		fatal(err)
		err = ioutil.WriteFile(tlsKey, tlsKeyB, 0400)
		fatal(err)
		tlsCertB, err := base64.StdEncoding.DecodeString(os.Getenv("TLS_CERTIFICATE_BASE64"))
		fatal(err)
		err = ioutil.WriteFile(tlsCert, tlsCertB, 0400)
		fatal(err)
	}
}

func main() {
	// initialize postgres backend
	db, err := sql.Open("postgres", dbInfo)
	fatal(err)

	defer db.Close()

	// Add store to resource
	p := player.PlayerResource{}

	// Register container
	wsContainer := restful.NewContainer()
	p.Register(wsContainer)

	// Setup api docs
	apidocs.Register(wsContainer, port, tls)

	// Start server
	log.Printf("listening on localhost:%s", port)
	server := &http.Server{Addr: fmt.Sprintf(":%s", port), Handler: wsContainer}
	if tls {
		log.Fatal(server.ListenAndServeTLS(tlsCert, tlsKey))
	} else {
		log.Fatal(server.ListenAndServe())
	}
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
