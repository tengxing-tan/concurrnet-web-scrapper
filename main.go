package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"
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

	// Anti-bot: rotate user agents
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64; rv:121.0) Gecko/20100101 Firefox/121.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15",
	}

	// Anti-bot: optional proxy via env var (HTTP_PROXY / HTTPS_PROXY)
	var transport http.Transport
	if proxy := os.Getenv("HTTPS_PROXY"); proxy != "" {
		if u, err := url.Parse(proxy); err == nil {
			transport.Proxy = http.ProxyURL(u)
		} else {
			log.Printf("Invalid HTTPS_PROXY: %v", err)
		}
	} else if proxy := os.Getenv("HTTP_PROXY"); proxy != "" {
		if u, err := url.Parse(proxy); err == nil {
			transport.Proxy = http.ProxyURL(u)
		} else {
			log.Printf("Invalid HTTP_PROXY: %v", err)
		}
	}
	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: &transport,
	}

	// Global rate limiter: 1 request every 300ms (~3.3 req/sec)
	rateTicker := time.NewTicker(300 * time.Millisecond)
	
	// Jobs channel and worker pool
	jobs := make(chan int)
	var workersWg sync.WaitGroup

	// Start fixed number of workers
	for i := 0; i < workerCount; i++ {
		workersWg.Add(1)
		go func() {
			defer workersWg.Done()
			for id := range jobs {
				// Rate limit: one fetch per tick across all workers
				<-rateTicker.C

				// Rotate UA deterministically
				ua := userAgents[id%len(userAgents)]

				post, err := fetchPost(client, id, ua)
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
		// Stop the ticker once all workers have finished
		rateTicker.Stop()
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
func fetchPost(client *http.Client, id int, userAgent string) (Post, error) {
    url := fmt.Sprintf("https://jsonplaceholder.typicode.com/posts/%d", id)

    req, err := http.NewRequest(http.MethodGet, url, nil)
    if err != nil {
        return Post{}, err
    }
    req.Header.Set("User-Agent", userAgent)
    req.Header.Set("Accept", "application/json")

    res, err := client.Do(req)
    if err != nil {
        return Post{}, err
    }
    defer res.Body.Close()

    if res.StatusCode != http.StatusOK {
        return Post{}, fmt.Errorf("unexpected status %d", res.StatusCode)
    }

    var post Post
    if err := json.NewDecoder(res.Body).Decode(&post); err != nil {
        return Post{}, err
    }
    return post, nil
}
