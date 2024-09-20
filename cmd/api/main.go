package main

import (
	"fmt"
	"recommendation/internal/server"
)

func main() {

	server := server.NewServer()
	fmt.Println("Starting Server")
	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}

}
