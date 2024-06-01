package main

import (
	"GraphQLPostComments/internal/server"
	"log"
	"net/http"
)

func main() {
	srv := server.NewServer()

	http.Handle("/graphql", srv)
	log.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
