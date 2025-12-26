package main

import (
	"encoding/json"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	"os"
	"strconv"
)

type Post struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

func main() {
	start := time.Now()

	var wg sync.WaitGroup

	channel := make(chan Post, 10) // Buffered size of 10

	file, err := os.Create("posts.csv")
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close() // Best practice: Explicitly Close the File

	writer := csv.NewWriter(file)
	defer writer.Flush() // Ensures all data is physically written to the disk

	writer.Write([]string{"ID", "Title"})
 	
	postIDs := []int{1,2,3,4,5,10,50,75,100}
	
	for _, id := range postIDs {
		wg.Add(1)
		go fetchPostInChan(id, &wg, channel)
	}

	// Start a goroutine to Wait and Close the channel.
	go func() {
		// Block the main function here until the counter returns to 0
		wg.Wait() 
		close(channel)
	}()

	for post := range channel {
    	row := []string{
			strconv.Itoa(post.ID), 
			post.Title,
		}
		
		if err := writer.Write(row); err != nil {
			log.Printf("Error writing row for post %d: %v", post.ID, err)
		}
	}

	fmt.Printf("\nDone processed %d requests in %v\n", len(postIDs), time.Since(start))
}

// The 'chan<-' syntax means this function can only SEND to the channel NOT read/close
func fetchPostInChan(id int, wg *sync.WaitGroup, channel chan<- Post) {
	defer wg.Done()

	url := fmt.Sprintf("https://jsonplaceholder.typicode.com/posts/%d", id)

	res, err := http.Get(url)

	if err != nil {
		log.Printf("Error fetching post %d: %v", id, err)
		return
	}
	defer res.Body.Close()

	var post Post
	if err := json.NewDecoder(res.Body).Decode(&post); err != nil {
		log.Printf("Error decoding post %d: %v", id, err)
		return
	}

	channel <- post
}
