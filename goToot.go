package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Struct for the .json configuration.
type ClientConfig struct {
	Token    string `json:"access_token"`
	Instance string `json:"instance"`
}

func verifyToken(bearer string, url string) bool {
	// Append the endpoint to the URL.
	fullUrl := fmt.Sprintf("%v/apps/verify_credentials", url)

	// Create an HTTP client.
	client := &http.Client{}
	request, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	// Set the header.
	request.Header.Set("Authorization", bearer)

	// Make the request.
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		os.Exit(4)
	}

	defer response.Body.Close()

	// Read the data.
	respData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(5)
	}

	// Ensure it's not null.
	if respData != nil {
		return true
	} else {
		return false
	}
}

// Main function.
func main() {
	// Import the file with the config.
	config, err := ioutil.ReadFile("./client.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create the struct with the config.
	var configInfo ClientConfig
	err = json.Unmarshal(config, &configInfo)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// Create the base URL for our instance.
	baseURL := fmt.Sprintf("%v/api/v1", configInfo.Instance)

	// Create the header we'll use for authorization.
	bearerHeader := fmt.Sprintf("Bearer %v", configInfo.Token)

	// Verify the token is valid.
	if verifyToken(bearerHeader, baseURL) {
		fmt.Println("Token is valid!")
	} else {
		fmt.Println("Bad token!")
	}
}
