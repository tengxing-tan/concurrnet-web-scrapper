# Concurrent web scrapper

A web scraper is an ideal first project because it involves both I/O-bound tasks and data processing.

Concepts Learned: Goroutines for parallel fetching, sync.WaitGroup to coordinate task completion, and channels to gather results.
Advanced Goal: Implement Rate Limiting using a time.Ticker or the golang.org/x/time/rate package to avoid overwhelming target servers
