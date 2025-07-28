package main

import (
	"fmt"
	"go-postgress/middleware"
	"go-postgress/routes"
	"log"
	"net/http"
)

func main() {

	r := routes.Router()
	middleware.InitDB()

	fmt.Println("server starting on port 4000")

	log.Fatal(http.ListenAndServe(":4000", r))

}
