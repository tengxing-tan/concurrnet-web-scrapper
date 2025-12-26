package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	if err := handleMultipleUrlsInForLoop(); err != nil {
		log.Fatal(err)
	}
}

func fetchOneUrl(paramId int) error {
	url := fmt.Sprintf("https://jsonplaceholder.typicode.com/posts/%d", paramId)

	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("http request: %w", err)
	}
	defer res.Body.Close()

	// Accept any 2xx as success
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("unexpected status: %d %s", res.StatusCode, res.Status)
	}

	var post struct {
		UserID int    `json:"userId"`
		ID     int    `json:"id"`
		Title  string `json:"title"`
		Body   string `json:"body"`
	}
	if err := json.NewDecoder(res.Body).Decode(&post); err != nil {
		return fmt.Errorf("decode json: %w", err)
	}

	fmt.Printf("Post #%d Title: %s\n", post.ID, post.Title)
	return nil
}

func handleMultipleUrlsInForLoop() error {
	for i := 1; i <= 100; i++ {
        if err := fetchOneUrl(i); err != nil {
			return err
		}
    }
	return nil
}