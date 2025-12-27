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

	// Fixed-size worker pool
	const workerCount = 5

	// Results channel consumed by the CSV writer
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
	
	// Jobs channel and worker pool
	jobs := make(chan int)
	var workersWg sync.WaitGroup

	// Start fixed number of workers
	for i := 0; i < workerCount; i++ {
		workersWg.Add(1)
		go func() {
			defer workersWg.Done()
			for id := range jobs {
				post, err := fetchPost(id)
				if err != nil {
					log.Printf("Error fetching post %d: %v", id, err)
					continue
				}
				channel <- post
			}
		}()
	}

	// Feed jobs and close channels when done
	go func() {
		for _, id := range postIDs {
			jobs <- id
		}
		close(jobs)
		workersWg.Wait()
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

// New helper used by workers
func fetchPost(id int) (Post, error) {
    url := fmt.Sprintf("https://jsonplaceholder.typicode.com/posts/%d", id)
    res, err := http.Get(url)
    if err != nil {
        return Post{}, err
    }
    defer res.Body.Close()

    var post Post
    if err := json.NewDecoder(res.Body).Decode(&post); err != nil {
        return Post{}, err
    }
    return post, nil
}
