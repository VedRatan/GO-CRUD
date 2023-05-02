package main

import (
	"flipr_assignment/routers"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main(){
	router := mux.NewRouter()
	routers.RegisterNoteRoutes(router)
	fmt.Println("Listening on port 5000...")
	log.Fatal(http.ListenAndServe(":5000", (router)))
}