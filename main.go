package main

import (
	"log"

	"github.com/askme23/hltv-api/api"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
    log.Fatal("Error loading .env file")
  }

  api.LoadMatches()
  // api.LoadTeams()
  // api.GetTeams()
}