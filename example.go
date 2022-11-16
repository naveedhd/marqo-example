package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const serverPort = 8882

// Http helpers
func url(endpoint string) string {
	return fmt.Sprintf("http://localhost:%d/%s", serverPort, endpoint)
}

func Post(url string, body []byte) (string, error) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(responseData), nil
}

func Get(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(responseData), nil
}

// Marqo endpoints
func health() (string, error) {
	requestUrl := url("health")
	return Get(requestUrl)
}

func createDocument(index string, requestBody []byte) (string, error) {
	requestUrl := url(fmt.Sprintf("indexes/%s/documents", index))
	return Post(requestUrl, requestBody)
}

func search(index string, requestBody []byte) (string, error) {
	requestUrl := url(fmt.Sprintf("indexes/%s/search", index))
	return Post(requestUrl, requestBody)
}

func main() {
	healthResponse, err := health()
	if err != nil {
		fmt.Printf("Error making request: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println(healthResponse)

	// create documents
	index := "my-first-index"
	requestBody := []byte(`[
		{
			"Title": "The Travels of Marco Polo",
			"Description": "A 13th-century travelogue describing Polo's travels"
		},
		{
			"Title": "Extravehicular Mobility Unit (EMU)",
			"Description": "The EMU is a spacesuit that provides environmental protection, mobility, life support, and communications for astronauts",
			"_id": "article_591"
		}
	]`)
	documentsResponse, err := createDocument(index, requestBody)
	if err != nil {
		fmt.Printf("Error making request: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println(documentsResponse)

	// search
	requestBody = []byte(`{
		"q":"What is the best outfit to wear on the moon?"
	  }`)
	searchResponse, err := search(index, requestBody)
	if err != nil {
		fmt.Printf("Error making request: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println(searchResponse)
}
