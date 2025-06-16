package main

import (
	"encrypted-chat-relay/handlers"
	"fmt"
	"log"
	"net/http"
)

func main() {
	db := connect()
	fmt.Println("Connected to DB!")

	httpHandler := handlers.NewHandler(db)
	wsHandler := handlers.NewWSHandler(db)

	mux := http.NewServeMux()
	mux.HandleFunc("/", httpHandler.HandleHome)
	mux.HandleFunc("/register", httpHandler.HandleRegister)
	mux.HandleFunc("/ws", wsHandler.HandleWS)

	corsHandler := handlers.CORSMiddleware(mux)

	fmt.Println("Starting http server on :8080!")
	if err := http.ListenAndServe(":8080", corsHandler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
