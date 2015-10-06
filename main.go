package main

import (
	"database/sql"
	"fmt"
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
	dbInfo, port string
	tls          bool
)

func init() {

	dbInfo = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("DATABASE_SSL"))

	port = os.Getenv("API_PORT")
	if os.Getenv("TLS_ENABLED") == "1" {
		tls = true
	} else {
		tls = false
	}

	if port == "" {
		port = "80"
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
