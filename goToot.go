package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"jaytaylorcom/html2text"
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

// Struct for a single toot. Used in response when posting.
type SingleToot struct {
	ID                 string      `json:"id"`
	CreatedAt          time.Time   `json:"created_at"`
	InReplyToID        interface{} `json:"in_reply_to_id"`
	InReplyToAccountID interface{} `json:"in_reply_to_account_id"`
	Sensitive          bool        `json:"sensitive"`
	SpoilerText        string      `json:"spoiler_text"`
	Visibility         string      `json:"visibility"`
	Language           string      `json:"language"`
	URI                string      `json:"uri"`
	URL                string      `json:"url"`
	RepliesCount       int         `json:"replies_count"`
	ReblogsCount       int         `json:"reblogs_count"`
	FavouritesCount    int         `json:"favourites_count"`
	Favourited         bool        `json:"favourited"`
	Reblogged          bool        `json:"reblogged"`
	Muted              bool        `json:"muted"`
	Bookmarked         bool        `json:"bookmarked"`
	Pinned             bool        `json:"pinned"`
	Content            string      `json:"content"`
	Reblog             interface{} `json:"reblog"`
	Application        struct {
		Name    string `json:"name"`
		Website string `json:"website"`
	} `json:"application"`
	Account struct {
		ID             string        `json:"id"`
		Username       string        `json:"username"`
		Acct           string        `json:"acct"`
		DisplayName    string        `json:"display_name"`
		Locked         bool          `json:"locked"`
		Bot            bool          `json:"bot"`
		Discoverable   bool          `json:"discoverable"`
		Group          bool          `json:"group"`
		CreatedAt      time.Time     `json:"created_at"`
		Note           string        `json:"note"`
		URL            string        `json:"url"`
		Avatar         string        `json:"avatar"`
		AvatarStatic   string        `json:"avatar_static"`
		Header         string        `json:"header"`
		HeaderStatic   string        `json:"header_static"`
		FollowersCount int           `json:"followers_count"`
		FollowingCount int           `json:"following_count"`
		StatusesCount  int           `json:"statuses_count"`
		LastStatusAt   string        `json:"last_status_at"`
		Emojis         []interface{} `json:"emojis"`
		Fields         []struct {
			Name       string    `json:"name"`
			Value      string    `json:"value"`
			VerifiedAt time.Time `json:"verified_at"`
		} `json:"fields"`
	} `json:"account"`
	MediaAttachments []interface{} `json:"media_attachments"`
	Mentions         []interface{} `json:"mentions"`
	Tags             []interface{} `json:"tags"`
	Emojis           []interface{} `json:"emojis"`
	Card             interface{}   `json:"card"`
	Poll             interface{}   `json:"poll"`
}

// Function to query data from Mastodon.
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

// Function to push content to Mastodon.
func postToMasto(bearer string, url string, content string, replyId string, sensitive bool, spoiler string) string {
	// Create the url.
	url = fmt.Sprintf("%v/statuses", url)

	// Create the map for the form data.
	formData := make(map[string]string)
	formData["status"] = content
	if replyId != "" {
		formData["in_reply_to_id"] = replyId
	}
	if sensitive {
		formData["sensitive"] = "true"
		formData["spoiler_text"] = spoiler
	}

	// Put together the form body.
	reqBody, err := json.Marshal(formData)
	if err != nil {
		fmt.Println(err)
		os.Exit(10)
	}

	// Put together the client.
	client := &http.Client{}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println(err)
		os.Exit(11)
	}
	request.Header.Set("Authorization", bearer)
	request.Header.Set("Content-Type", "application/json")

	// Make the request.
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		os.Exit(12)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(13)
	}

	// Parse the toot to a struct and return the ID.
	var postedToot SingleToot
	err = json.Unmarshal(body, &postedToot)
	if err != nil {
		fmt.Println(err)
		os.Exit(14)
	}

	return postedToot.ID
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

// Function to get toot content.
func getTootContent() string {
	var text string
	var err error
	reader := bufio.NewReader(os.Stdin)
	shortEnough := false
	// Prompt the user for their text.
	fmt.Printf("\nEnter your toot.\n")
	for !shortEnough {
		fmt.Print("> ")
		text, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(9)
		}

		// Trim the return.
		text = strings.Trim(text, "\n")

		// Verify we're within the length limit.
		if len(text) > 500 {
			fmt.Println("That toot is too long! Try again...")
			continue
		} else {
			shortEnough = true
		}
	}
	return text
}

// Function to print the toots in a timeline.
func printToots(allToots []SingleToot) {
	// Loop through the slice backwards.
	for i := len(allToots) - 1; i >= 0; i-- {
		// Parse the HTML of the post to Markdown-like text.
		markdown, err := html2text.FromString(allToots[i].Content)
		if err != nil {
			fmt.Println(err)
			os.Exit(15)
		}
		fmt.Printf("%v at %v\n\n", allToots[i].Account.Acct, allToots[i].CreatedAt)
		fmt.Println(markdown)
		fmt.Printf("Favs: %v\tBoosts: %v\n", allToots[i].FavouritesCount, allToots[i].ReblogsCount)
		fmt.Printf("\n\n")
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
	if !verifyToken(bearerHeader, baseURL) {
		fmt.Println("Token is invalid!")
		os.Exit(6)
	}

	// Verify the user information.
	currentUser := verifyUserCreds(bearerHeader, baseURL)
	fmt.Printf("Logged in as: %v\n", currentUser.Acct)
	fmt.Printf("%v statuses, last one posted on %v\n\n", currentUser.StatusesCount, currentUser.LastStatusAt)

	// Start the main loop to see what the user would like to do.
	var userChoice string
	var cwText string
	var currentPost string
	var currentTimeline []byte
	var currentTLParsed []SingleToot
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
			// Get the byte slice for the timeline.
			currentTimeline = queryMasto(bearerHeader, fmt.Sprintf("%v/timelines/home?limit=2", baseURL))
			err = json.Unmarshal(currentTimeline, &currentTLParsed)
			if err != nil {
				fmt.Println(err)
				os.Exit(15)
			}
			//fmt.Printf("%+v\n", currentTLParsed)
			printToots(currentTLParsed)
		case "local":
			fmt.Println("Display 'Local' timeline.")
		case "note":
			fmt.Println("Display 'Notification' feed.")
		case "toot":
			// Prompt the user for their text.
			text := getTootContent()

			// Pass to the function.
			currentPost = postToMasto(bearerHeader, baseURL, text, "", false, "")
			fmt.Printf("Successfully posted toot: %v\n\n", currentPost)
		case "cwtoot":
			// Prompt the user for their spoiler text.
			fmt.Printf("\nEnter your spoiler text.\n")
			fmt.Print("> ")
			cwText, err = reader.ReadString('\n')

			// Prompt the user for their text.
			text := getTootContent()
			// Pass to the function.
			currentPost = postToMasto(bearerHeader, baseURL, text, "", true, cwText)
			fmt.Printf("Successfully posted toot: %v\n\n", currentPost)
		default:
			continue
		}
	}
}
