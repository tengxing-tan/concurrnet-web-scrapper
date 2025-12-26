# Concurrent Web Scrapper

A web scraper is an ideal first project because it involves both I/O-bound tasks and data processing.

Concepts Learned: Goroutines for parallel fetching, sync.WaitGroup to coordinate task completion, and channels to gather results.
Advanced Goal: Implement Rate Limiting using a time.Ticker or the golang.org/x/time/rate package to avoid overwhelming target servers

## Primary Goal

To master goroutines and channels

## What data is being fetched

Using Simple Public APIs (Recommended for Concurrency Focus)

- JSONPlaceholder: A free fake API for testing and prototyping. It returns standard JSON for posts, comments, and users, allowing you to focus entirely on your Go concurrency logic.

## Things To Learn

### Goroutines: Parallel Fetching

- Managed by Go runtime not OS, which resolve heavy threads (required pre-allocated stack)
- Reduce overhead of context switching
- Purpose: Instead of waiting for one URL to finish downloading before starting the next, you can launch hundreds of fetches simultaneously.
- Note: If the main function finishes before your goroutines, the program exits immediately and the goroutines are killed. This is why synchronization is required.

### sync.WaitGroup: Coordination

- A sync.WaitGroup acts as a counter that keeps track of how many tasks are still running.
- Add(n): Call this before starting your goroutines to tell Go how many tasks it needs to wait for.
- Done(): Inside the goroutine, call defer wg.Done() as the very first line. This ensures that when the function finishes (even if it fails), the counter is decremented.
- Wait(): In your main function, call wg.Wait(). This blocks execution until the counter reaches zero, ensuring all fetches are complete before the program moves forward.

### Channels: Gathering Results

- While WaitGroups coordinate timing, Channels are the "pipes" used to move data safely between goroutines.
- Unbuffered vs. Buffered: You can use a buffered channel to collect results without blocking each worker immediately.
- Collection: As each goroutine finishes fetching a URL, it sends the result (e.g., the HTML or status code) into the channel: resultsChan <- data.
- Safe Close: A common pattern is to start a separate goroutine that calls wg.Wait() and then closes the channel. This allows you to range over the channel in your main function to process all collected data.

### Real-World Workflow Example

- Define a channel to receive results and a WaitGroup for tracking.
- Iterate over a list of URLs, calling wg.Add(1) and then launching a go func for each.
- Inside each goroutine, fetch the data and send it to the channel.
- Wait and Close: Start a helper goroutine that waits for the WaitGroup to finish and then closes the results channel.
- Process results by reading from the channel until it is empty.

## How To Do

### Phase 1: Foundation & Single-Page Scraper

Before adding concurrency, establish a working baseline for fetching and deooding JSON data.

- **Initialize Project**: Create your directory and run `go mod init <name>`.
- **Build Sequential Logic**: Write a function to fetch one URL, extract specific data (like titles or prices), and print them to the console.
- **Error Handling**: Implement basic checks for failed HTTP requests or missing HTML elements to prevent panics.

### Phase 2: Implementing Concurrency

Transition the project to parallel execution to handle multiple URLs simultaneously.

- **Define a Seed List**: Create a slice of target URLs to scrape.
- **Launch Goroutines**: Use a for loop to iterate through your URLs, launching each fetch in its own `go func()`.
- **Add Synchronization**: Use sync.WaitGroup to ensure the main function waits for all goroutines to finish before exiting.
- **Gather Results via Channels**: Create a channel to safely collect scraped data from workers back to your main thread.

### Phase 3: Scaling & Optimization (2025 Best Practices)

Refine the scraper to handle real-world challenges like rate limiting and dynamic content.

- **Implement Worker Pools**: Instead of launching a goroutine for every URL, use a fixed number of workers to prevent resource exhaustion or getting blocked by servers.
- **Add Rate Limiting**: Use `time.Sleep` or Go's rate package to introduce delays between requests and mimic human behavior.
- **Handle Dynamic Content**: For JavaScript-heavy sites, integrate Chromedp to control a headless browser.
- **Data Persistence**: Update your pipeline to save results to a structured format like a CSV file or a database.
- **Anti-Bot Strategies**: Rotate User-Agent headers and consider integrating proxy pools to avoid IP bans.
