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

var (
	dbInfo, port string
)

func init() {
	dbInfo = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("DATABASE_SSL"))
	port = os.Getenv("API_PORT")
}

func main() {
	// initialize postgres backend
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Add store to resource
	p := player.PlayerResource{}

	// Register container
	wsContainer := restful.NewContainer()
	p.Register(wsContainer)

	// Setup api docs
	apidocs.Register(wsContainer, port)

	// Start server
	log.Printf("listening on localhost:%s", port)
	server := &http.Server{Addr: fmt.Sprintf(":%s", port), Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
