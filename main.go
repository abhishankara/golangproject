package main

import (
	"fmt"
	"log"
	"net/http"
	"server/router"
)

func main() {
	fmt.Println("Mongo DB API")

	r := router.Router()
	fmt.Println("Server is getting started...")
	log.Fatal(http.ListenAndServe(":8081", r))

}
