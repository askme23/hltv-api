package main

import (
	"fmt"
	"log"

	"github.com/askme23/hltv-api/api"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
    log.Fatal("Error loading .env file")
  }

  fmt.Println(string(api.GetTeams()))
  // api.LoadMatches()
  // api.GetTeams()
}