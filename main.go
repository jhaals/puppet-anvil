package main

import (
	"log"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	modulePath := os.Getenv("MODULEPATH")

	if len(port) == 0 {
		log.Fatal("Missing PORT environment variable")
	}
	if len(modulePath) == 0 {
		log.Fatal("Missing MODULEPATH environment variable")
	}

	svc := NewAnvilService(port, modulePath)

	if err := svc.Run(); err != nil {
		log.Print(err)
	}
}
