package main

import (
	_ "database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/emicklei/go-restful"
	_ "github.com/lib/pq"

	"github.com/manifest-destiny/api"
	"github.com/manifest-destiny/api/apidocs"
	"github.com/manifest-destiny/api/user"
)

const (
	tlsCert = "keys/tls_cert.crt"
	tlsKey  = "keys/tls_key.key"
)

var (
	dbInfo, port, webClientID string
	tls                       bool
)

func init() {
	dbInfo = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("DATABASE_SSL"))

	webClientID = os.Getenv("WEB_CLIENT_ID")

	port = os.Getenv("API_PORT")
	if port == "" {
		port = "80"
	}

	tls = os.Getenv("TLS_ENABLED") == "1"
}

func main() {
	// initialize postgres backend
	conn, err := api.NewDB("postgres", dbInfo)
	fatal(err)
	defer conn.Close()

	// initialize token validator
	v, err := api.GoogleTokenValidator(webClientID)
	fatal(err)

	// Add db and token validator to user resource
	userResource := &user.Resource{conn, v}

	// Register container
	wsContainer := restful.NewContainer()
	user.RegisterContainer(wsContainer, userResource)

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
