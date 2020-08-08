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
	"strconv"
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
	ID                 string `json:"id"`
	ClientID           int
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

// Struct for notifications.
type Notification struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	Account   struct {
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
			Name       string      `json:"name"`
			Value      string      `json:"value"`
			VerifiedAt interface{} `json:"verified_at"`
		} `json:"fields"`
	} `json:"account"`
	Status struct {
		ID                 string `json:"id"`
		ClientID           int
		CreatedAt          time.Time   `json:"created_at"`
		InReplyToID        string      `json:"in_reply_to_id"`
		InReplyToAccountID string      `json:"in_reply_to_account_id"`
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
				Name       string      `json:"name"`
				Value      string      `json:"value"`
				VerifiedAt interface{} `json:"verified_at"`
			} `json:"fields"`
		} `json:"account"`
		MediaAttachments []interface{} `json:"media_attachments"`
		Mentions         []struct {
			ID       string `json:"id"`
			Username string `json:"username"`
			URL      string `json:"url"`
			Acct     string `json:"acct"`
		} `json:"mentions"`
		Tags   []interface{} `json:"tags"`
		Emojis []interface{} `json:"emojis"`
		Card   interface{}   `json:"card"`
		Poll   interface{}   `json:"poll"`
	} `json:"status"`
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

// Function to get the ID of a toot to boost or favorite.
func getTootID() int {
	// Prompt the user.
	fmt.Printf("\nEnter the ID.\n")
	fmt.Print("> ")

	// Get user input.
	reader := bufio.NewReader(os.Stdin)
	userInput, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		os.Exit(22)
	}
	userInput = strings.Trim(userInput, "\n")

	// Validate that it's an integer.
	inputInt, err := strconv.Atoi(userInput)
	if err != nil {
		fmt.Printf("%v is not a valid integer! You must enter a valid ID...\n", userInput)
		inputInt = 0
	}

	// Return the value.
	return inputInt
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

		// Modify the date.
		datePretty := strings.Split(allToots[i].CreatedAt.String(), ".")

		// Print the author, app, and timestamp.
		applicationName := "Web"
		if allToots[i].Application.Name != "" {
			applicationName = allToots[i].Application.Name
		}
		fmt.Printf("> %v from |%v| to |%v| at %v\n", allToots[i].Account.Acct, applicationName, allToots[i].Visibility, datePretty[0])

		// Check if there's a CW.
		if allToots[i].Sensitive {
			// Print it.
			fmt.Printf(">> CW: %v\n", allToots[i].SpoilerText)
		}

		// Print the main toot content parsed to Markdown.
		if allToots[i].Reblogged {
			fmt.Printf("\n%v\n\n", markdown)
		} else {
			fmt.Printf("\n%v\n", markdown)
		}

		// Check if there's media.
		media := allToots[i].MediaAttachments
		if len(media) > 0 {
			// Loop through it.
			for i := len(media) - 1; i >= 0; i-- {
				// Get the type of media and the URL to it.
				mediaType := media[i].(map[string]interface{})["type"]
				mediaURL := media[i].(map[string]interface{})["text_url"]
				fmt.Printf("%v: %v\n", mediaType, mediaURL)
			}
		}
		fmt.Printf("~=: ID: %v\tFavs: %v\tBoosts: %v :=~\n", allToots[i].ClientID, allToots[i].FavouritesCount, allToots[i].ReblogsCount)
		fmt.Printf("\n")
	}
}

// Function to print notifications.
func printNotifications(allNotifications []Notification) {
	// Loop through the slice backwards.
	for i := len(allNotifications) - 1; i >= 0; i-- {
		// Check the type.
		if allNotifications[i].Type == "mention" {
			// Modify the date.
			datePretty := strings.Split(allNotifications[i].Status.CreatedAt.String(), ".")

			// Get the application used.
			applicationName := "Web"
			if allNotifications[i].Status.Application.Name != "" {
				applicationName = allNotifications[i].Status.Application.Name
			}
			fmt.Printf("> Mention by %v from |%v| to |%v| at %v\n", allNotifications[i].Account.Acct, applicationName, allNotifications[i].Status.Visibility, datePretty[0])
		} else if allNotifications[i].Type == "favourite" || allNotifications[i].Type == "boost" {
			// Print the info on the fav/boost and for what toot.
			if allNotifications[i].Type == "favourite" {
				fmt.Printf("> Favorite by %v\n", allNotifications[i].Account.Acct)
			} else {
				fmt.Printf("> Boost by %v\n", allNotifications[i].Account.Acct)
			}
		} else if allNotifications[i].Type == "follow" {
			// Print information about who followed.
			markdown, err := html2text.FromString(allNotifications[i].Account.Note)
			if err != nil {
				fmt.Println(err)
				os.Exit(23)
			}
			fmt.Printf("> Followed by %v\n", allNotifications[i].Account.Acct)
			fmt.Printf(">> Has posted %v statuses, the last on %v\n", allNotifications[i].Account.StatusesCount, allNotifications[i].Account.LastStatusAt)
			fmt.Printf("%v\n", markdown)
			fmt.Printf("~=: Following: %v\tFollowers: %v :=~\n\n", allNotifications[i].Account.FollowingCount, allNotifications[i].Account.FollowersCount)
		} else {
			// Probably change this later, yeah?
			fmt.Printf("%+v\n", allNotifications[i])
			fmt.Println("Not sure what to do with a type of %v\n", allNotifications[i].Type)
		}

		// Parse the toot content and print it if there is any.
		if allNotifications[i].Status.Content != "" {
			markdown, err := html2text.FromString(allNotifications[i].Status.Content)
			if err != nil {
				fmt.Println(err)
				os.Exit(18)
			}
			fmt.Printf("\n%v\n", markdown)
			fmt.Printf("~=: ID: %v\tFavs: %v\tBoosts: %v :=~\n\n", allNotifications[i].Status.ClientID, allNotifications[i].Status.FavouritesCount, allNotifications[i].Status.ReblogsCount)
		}
	}
}

// Function to assign indexes to all toots for reference.
func assignIndexToots(allToots []SingleToot, indexStart int) ([]SingleToot, int) {
	for i := len(allToots) - 1; i >= 0; i-- {
		// Increment the counter.
		indexStart++

		// Assign the ID.
		allToots[i].ClientID = indexStart
	}

	// Return the updated array and the new index.
	return allToots, indexStart
}

// Function to assign an index to notifications.
func assignIndexNotes(allNotes []Notification, indexStart int) ([]Notification, int) {
	for i := len(allNotes) - 1; i >= 0; i-- {
		// Increment the counter.
		indexStart++

		// Assign the ID.
		allNotes[i].Status.ClientID = indexStart
	}

	// Return the updated slice and the new index.
	return allNotes, indexStart
}

// Function to favorite a toot.
func favOrBoostToot(bearer string, url string, tootID string, tootBody string, updateType string) {
	// Parse the appropriate URL.
	if updateType == "boost" {
		url = fmt.Sprintf("%v/statuses/%v/reblog", url, tootID)
	} else {
		url = fmt.Sprintf("%v/statuses/%v/favourite", url, tootID)
	}

	// Put together the client.
	client := &http.Client{}
	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(19)
	}
	request.Header.Set("Authorization", bearer)
	request.Header.Set("Content-Type", "application/json")

	// Make the request.
	_, err = client.Do(request)
	if err != nil {
		fmt.Println(err)
		os.Exit(20)
	} else {
		// Parse the toot content to plaintext.
		markdown, err := html2text.FromString(tootBody)
		if err != nil {
			fmt.Println(err)
			os.Exit(21)
		}
		if updateType == "boost" {
			fmt.Printf("Successfully boosted: %v\n", markdown)
		} else {
			fmt.Printf("Successfully favorited: %v\n", markdown)
		}
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

	// Initialize the counter for toot IDs.
	tootCounter := 0

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
	var currentNotifications []byte
	var currentTLParsed []SingleToot
	var currentNotesParsed []Notification
	var lastTootsReceived string
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

			// Assign each toot an index for this app.
			currentTLParsed, tootCounter = assignIndexToots(currentTLParsed, tootCounter)

			//fmt.Printf("%+v\n", currentTLParsed)
			printToots(currentTLParsed)

			// Set where we got toots from.
			lastTootsReceived = "tl"
		case "local":
			// Get the byte slice.
			currentTimeline = queryMasto(bearerHeader, fmt.Sprintf("%v/timelines/public?local=true&limit=2", baseURL))
			err = json.Unmarshal(currentTimeline, &currentTLParsed)
			if err != nil {
				fmt.Println(err)
				os.Exit(16)
			}
			currentTLParsed, tootCounter = assignIndexToots(currentTLParsed, tootCounter)
			printToots(currentTLParsed)
			lastTootsReceived = "tl"
		case "note", "notes":
			// Get the byte slice.
			currentNotifications = queryMasto(bearerHeader, fmt.Sprintf("%v/notifications?limit=2", baseURL))

			// Parse to a struct.
			err = json.Unmarshal(currentNotifications, &currentNotesParsed)
			if err != nil {
				fmt.Println(err)
				os.Exit(17)
			}
			currentNotesParsed, tootCounter = assignIndexNotes(currentNotesParsed, tootCounter)

			// Print the notifications.
			printNotifications(currentNotesParsed)
			lastTootsReceived = "notes"
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
		case "fav":
			// Prompt the user for the ID of the toot to fav.
			tootSelection := getTootID()

			// Don't do anything if it was 0.
			if tootSelection != 0 {
				if lastTootsReceived == "" {
					fmt.Println("No toots in the local database to fav!")
				} else if lastTootsReceived == "tl" {
					for _, toot := range currentTLParsed {
						if toot.ClientID == tootSelection {
							// Submit the fav and break.
							favOrBoostToot(bearerHeader, baseURL, toot.ID, toot.Content, "fav")
							break
						}
					}
				} else {
					for _, note := range currentNotesParsed {
						if note.Type == "mention" {
							if note.Status.ClientID == tootSelection {
								favOrBoostToot(bearerHeader, baseURL, note.Status.ID, note.Status.Content, "fav")
								break
							}
						}
					}
				}
			}
		case "exit":
			// Just reset the userChoice variable to quit.
			userChoice = "quit"
		default:
			continue
		}
	}
}
