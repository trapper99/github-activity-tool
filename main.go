package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Event struct {
	Type    string    `json:"type"`
	Actor   Actor     `json:"actor"`
	Repo    Repo      `json:"repo"`
	CreatedAt   time.Time  `json:"created_at"`
}

type Actor struct {
	Login string `json:"login"`
}

type Repo struct {
	Name string `json:"name"`
}

func fetchUserActivity(username string) ([]Event, error) {
// Github API endpoint for user events
url := fmt.Sprintf("https:api.github.com/users/%s/events/public", username)

//Create a new Http Client
client := &http.Client{}

req, err := http.NewRequest("GET", url, nil)
if err != nil {
	return nil, err
}

//Add headers of the API version
req.Header.Add("Accept", "application/vnd.github.v3+json")

//Adding the Github token to the header
if token := os.Getenv("GITHUB_TOKEN"); token != "" {
	req.Header.Add("Authorization", "token "+token)
}

resp, err := client.Do(req)
if err != nil {
	return nil, err
}
defer resp.Body.Close()

if resp.StatusCode != http.StatusOK {
	return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
}

var events []Event
if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
	return nil, err
}
return events, nil
}

func displayEvents(events []Event) {
	if len(events) == 0 {
		fmt.Println("No recent activity found")
		return
	}

	fmt.Println("Recent Github Activity:")
	fmt.Println("------------------------")

	for _, event := range events {
		// Formatting the time to be more readable
		timeAgo := time.Since(event.CreatedAt).Round(time.Minute)

		fmt.Printf("Type: %s\n", event.Type)
		fmt.Printf("User: %s\n", event.Actor.Login)
		fmt.Printf("Repo: %s\n", event.Repo.Name)
		fmt.Printf("When: %s ago\n", timeAgo)
		fmt.Println("----------------------")
	}
}

//TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.

func main() {
username := flag.String("user", "", "Github username")
flag.Parse()

if *username == "" {
	fmt.Println("Please provide a Github username using -user flag")
	os.Exit(1)
}

events, err := fetchUserActivity(*username)
if err != nil {
	fmt.Printf("Error fetching user activity: %v\n", err)
	os.Exit(1)
}

displayEvents(events)
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
