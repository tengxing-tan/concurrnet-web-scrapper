# Concurrent Web Scrapper

A web scraper is an ideal first project because it involves both I/O-bound tasks and data processing.

Concepts Learned: Goroutines for parallel fetching, sync.WaitGroup to coordinate task completion, and channels to gather results.
Advanced Goal: Implement Rate Limiting using a time.Ticker or the golang.org/x/time/rate package to avoid overwhelming target servers

## Primary Goal

To master goroutines and channels

## What data is being fetched

Using Simple Public APIs (Recommended for Concurrency Focus)

- JSONPlaceholder: A free fake API for testing and prototyping. It returns standard JSON for posts, comments, and users, allowing you to focus entirely on your Go concurrency logic.

**Goroutine**

- Managed by Go runtime not OS, which resolve heavy threads (required pre-allocated stack)
- Reduce overhead of context switching

## How To Do

### Phase 1: Foundation & Single-Page Scraper

Before adding concurrency, establish a working baseline for fetching and parsing.

- **Initialize Project**: Create your directory and run `go mod init <name>`.
- **Select Libraries**: Install Colly for the crawling framework or **Goquery for jQuery-like HTML parsing**.
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
