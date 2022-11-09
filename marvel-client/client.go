package marvel

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

type marvelClient struct {
	pubKey     string
	privKey    string
	httpClient *http.Client
}

type character struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Modified    string `json:"modified"`
	ResourceURI string `json:"resourceURI"`
	Comics      struct {
		Available     string `json:"available"`
		Returned      string `json:"returned"`
		CollectionURI string `json:"collectionURI"`
		Items         []struct {
			ResourceURI string `json:"resourceURI"`
			Name        string `json:"name"`
		} `json:"items"`
	} `json:"comics"`
}

type characterResponse struct {
	Code            string `json:"code"`
	Status          string `json:"status"`
	Copyright       string `json:"copyright"`
	AttributionText string `json:"attributionText"`
	AttributionHTML string `json:"attributionHTML"`
	Data            struct {
		Offset  string      `json:"offset"`
		Limit   string      `json:"limit"`
		Total   string      `json:"total"`
		Count   string      `json:"count"`
		Results []character `json:"results"`
	} `json:"data"`
}

func (c *marvelClient) getCharacters() ([]character, error) {
	res, err := c.httpClient.Get("https://gateway.marvel.com:443/v1/public/characters?apikey=")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var characterResponse characterResponse
	if err := json.NewDecoder(res.Body).Decode(&characterResponse); err != nil {
		return nil, err
	}
	return characterResponse.Data.Results, nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	pubKey := os.Getenv("MARVEL_PUBLIC_KEY")
	privKey := os.Getenv("MARVEL_PRIVATE_KEY")
	client := marvelClient{
		pubKey:  pubKey,
		privKey: privKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	characters, err := client.getCharacters()
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(characters)
}
