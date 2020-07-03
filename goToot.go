package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	//"net/http"
	"os"
)

// Struct for the .json configuration.
type ClientConfig struct {
	Token    string `json:"access_token"`
	Instance string `json:"instance"`
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
	baseURL := configInfo.Instance + "/api/v1/"

	fmt.Println(baseURL)
}
