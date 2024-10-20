package main

import (
	"log"
	"os"

	"github.com/mdw-smarty/calc-apps/handlers"
)

func main() {
	handler := handlers.NewCSVHandler(
		os.Stdin,
		os.Stdout,
		log.New(os.Stderr, "", 0),
	)
	err := handler.Handle()
	if err != nil {
		log.Fatalln(err)
	}
}
