package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type Post struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

func main() {
	start := time.Now()

	var wg sync.WaitGroup

	channel := make(chan int)
 	
	postIDs := []int{1,2,3,4,5,10,50,75,100}
	
	for _, id := range postIDs {
		wg.Add(1)
		go fetchPost(id, &wg, channel)
	}

	// Start a goroutine to Wait and Close the channel.
	go func() {
		// Block the main function here until the counter returns to 0
		wg.Wait() 
		close(channel)
	}()

	for id := range channel {
    	fmt.Println("Received:", id)
	}

	fmt.Printf("\nDone processed %d requests in %v\n", len(postIDs), time.Since(start))
}

func fetchPost(id int, wg *sync.WaitGroup, channel chan int) {
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

	channel <- post.ID

	fmt.Printf("Successfully fetched Post #%d, %s\n", post.ID, post.Title)
}
