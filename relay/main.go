package main

import (
	"encrypted-chat-relay/handlers"
	"fmt"
	"net/http"
)

func main() {
	_ = connect()
	fmt.Println("Connected to DB!")

	http.HandleFunc("/", handlers.Home)

	fmt.Println("Starting http server on :8080!")
	http.ListenAndServe(":8080", nil)
}
