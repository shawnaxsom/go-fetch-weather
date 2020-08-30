package main

import (
	"os"
	"log"
	"fmt"
	"github.com/joho/godotenv"
)

// use godot package to load/read the .env file and
// return the value of the key
func environmentVariable(key string) string {

  // load .env file
  err := godotenv.Load(".env")

  if err != nil {
    log.Fatalf("Error loading .env file")
  }

  return os.Getenv(key)
}

func main() {
	fmt.Println("Hello world!")

	godotenv.Load(".env")

	if key := environmentVariable("METEOSTAT_API_KEY"); key == "" {
		fmt.Println("Please enter a METEOSTAT_API_KEY in a .env file, or use an environment variable")
	} else {
		fmt.Printf("Found key: %v\n", key)
	}
}
