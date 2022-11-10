package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
	"github.com/phillipwright7/pivottechschool/marvel"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	pubKey := os.Getenv("MARVEL_PUBLIC_KEY")
	privKey := os.Getenv("MARVEL_PRIVATE_KEY")
	client := marvel.MarvelClient{
		BaseURL: "https://gateway.marvel.com:443/v1/public",
		PubKey:  pubKey,
		PrivKey: privKey,
		HttpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	limit := 5
	characters, err := client.GetCharacters(limit)
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(characters)
}
