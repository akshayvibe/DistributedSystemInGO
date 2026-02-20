package main

import (
	"log"

	"github.com/akshayvibe/proglog/server"
)

func main() {
	srv := server.NewHttpServer(":8080")
	log.Fatal(srv.ListenAndServe())
}