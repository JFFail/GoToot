package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// Struct for the .json configuration.
type ClientConfig struct {
	Token    string `json:"access_token"`
	Instance string `json:"instance"`
}

// Struct for the user's account.
type CurrentUser struct {
	ID             string    `json:"id"`
	Username       string    `json:"username"`
	Acct           string    `json:"acct"`
	DisplayName    string    `json:"display_name"`
	Locked         bool      `json:"locked"`
	Bot            bool      `json:"bot"`
	Discoverable   bool      `json:"discoverable"`
	Group          bool      `json:"group"`
	CreatedAt      time.Time `json:"created_at"`
	Note           string    `json:"note"`
	URL            string    `json:"url"`
	Avatar         string    `json:"avatar"`
	AvatarStatic   string    `json:"avatar_static"`
	Header         string    `json:"header"`
	HeaderStatic   string    `json:"header_static"`
	FollowersCount int       `json:"followers_count"`
	FollowingCount int       `json:"following_count"`
	StatusesCount  int       `json:"statuses_count"`
	LastStatusAt   string    `json:"last_status_at"`
	Source         struct {
		Privacy   string      `json:"privacy"`
		Sensitive bool        `json:"sensitive"`
		Language  interface{} `json:"language"`
		Note      string      `json:"note"`
		Fields    []struct {
			Name       string    `json:"name"`
			Value      string    `json:"value"`
			VerifiedAt time.Time `json:"verified_at"`
		} `json:"fields"`
		FollowRequestsCount int `json:"follow_requests_count"`
	} `json:"source"`
	Emojis []interface{} `json:"emojis"`
	Fields []struct {
		Name       string    `json:"name"`
		Value      string    `json:"value"`
		VerifiedAt time.Time `json:"verified_at"`
	} `json:"fields"`
}

func queryMasto(bearer string, url string) []byte {
	// Create an HTTP client
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
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

	// Return the data.
	return respData
}

func verifyToken(bearer string, url string) bool {
	// Append the endpoint to the URL.
	fullUrl := fmt.Sprintf("%v/apps/verify_credentials", url)

	// Query.
	queryResult := queryMasto(bearer, fullUrl)

	// Ensure it's not null.
	if queryResult != nil {
		return true
	} else {
		return false
	}
}

// Function to verify the user.
func verifyUserCreds(bearer string, url string) CurrentUser {
	// Append the endpoint to the URL.
	fullUrl := fmt.Sprintf("%v/accounts/verify_credentials", url)

	// Query.
	queryResult := queryMasto(bearer, fullUrl)

	// Create the user.
	var thisUser CurrentUser
	err := json.Unmarshal(queryResult, &thisUser)
	if err != nil {
		fmt.Println(err)
		os.Exit(7)
	}

	return thisUser
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
	if !verifyToken(bearerHeader, baseURL) {
		fmt.Println("Token is invalid!")
		os.Exit(6)
	}

	// Verify the user information.
	currentUser := verifyUserCreds(bearerHeader, baseURL)
	fmt.Printf("Logged in as: %v\n", currentUser.Acct)
	fmt.Printf("%v statuses, last one posted on %v\n\n", currentUser.StatusesCount, currentUser.LastStatusAt)

	// Start the main loop to see what the user would like to do.
	userChoice := ""
	userPrompt := fmt.Sprintf("[%v]: ", currentUser.Acct)
	reader := bufio.NewReader(os.Stdin)
	for userChoice != "quit" {
		fmt.Print(userPrompt)

		// Get the user's input.
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(8)
		}
		userChoice = strings.ToLower(strings.Trim(text, "\n"))

		// Figure out what action to take based on user input.
		switch userChoice {
		case "home":
			fmt.Println("Display 'Home' timeline.")
		case "local":
			fmt.Println("Display 'Local' timeline.")
		case "note":
			fmt.Println("Display 'Notification' feed.")
		case "toot":
			fmt.Println("Prompt the user to post.")
		default:
			continue
		}
	}
}
