package main

import (
	"encrypted-chat-relay/handlers"
	"fmt"
	"net/http"
)

func main() {
	db := connect()
	fmt.Println("Connected to DB!")

	handler := handlers.NewHandler(db)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.HandleHome)
	mux.HandleFunc("/register", handler.HandleRegister)

	corsHandler := handlers.CORSMiddleware(mux)

	fmt.Println("Starting http server on :8080!")
	http.ListenAndServe(":8080", corsHandler)
}
