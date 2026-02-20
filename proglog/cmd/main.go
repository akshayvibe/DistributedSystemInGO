package main

import (
	"log"

	"github.com/akshayvibe/proglog/server"
)

func main() {
	srv := server.NewHttpServer(":8077")
	log.Fatal(srv.ListenAndServe())
}