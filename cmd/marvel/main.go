package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/phillipwright7/pivottechschool/marvel"
)

func main() {
	var envPth string
	flag.StringVar(&envPth, "p", ".env", "Flag to find .env file")
	flag.Parse()
	if err := godotenv.Load(envPth); err != nil {
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

	characters, err := client.GetCharacters(5)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(characters)
}
