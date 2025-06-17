/*
Package cmd
Copyright © 2025 Mikheil Lomidze

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var url string
var duration time.Duration
var concurrency int
var method string
var body string

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run gnat load generator",
	Run: func(cmd *cobra.Command, args []string) {
		if url == "" {
			fmt.Fprintln(os.Stderr, "❌ Please provide a URL with --url")
			os.Exit(1)
		}
		fmt.Printf("🚀 Sending requests to %s for %v with %d workers\n", url, duration, concurrency)

		ctx, cancel := context.WithTimeout(context.Background(), duration)
		defer cancel()

		var wg sync.WaitGroup
		var totalReqs int64
		responseStatusCodes := make(map[int]int)
		reqsChan := make(chan int)

		// Worker function
		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for {
					select {
					case <-ctx.Done():
						return
					default:
						var resp *http.Response
						var err error

						if method == "POST" {
							req, err := http.NewRequest("POST", url, strings.NewReader(body))
							if err == nil {
								req.Header.Set("Content-Type", "application/json")
								resp, err = http.DefaultClient.Do(req)
							}
						} else {
							resp, err = http.Get(url)
						}
						var statusCode int
						if err == nil && resp != nil {
							io.Copy(io.Discard, resp.Body)
							statusCode = resp.StatusCode
							resp.Body.Close()
						}
						reqsChan <- statusCode
					}
				}
			}(i)
		}

		// Count successful requests
		go func() {
			for sc := range reqsChan {
				totalReqs++
				responseStatusCodes[sc]++
			}
		}()

		// Wait for all workers to finish
		wg.Wait()
		close(reqsChan)

		fmt.Printf("✅ Load complete. Sent %d requests in %v\n", totalReqs, duration)
		if totalReqs > 0 {
			fmt.Printf("Average requests per second: %.2f\n", float64(totalReqs)/duration.Seconds())
			fmt.Println("Response status codes:")
			for code, count := range responseStatusCodes {
				if count > 0 {
					fmt.Printf("  %d: %d times (%.2f%%)\n", code, count, float64(count)/float64(totalReqs)*100)
				}
			}
		} else {
			fmt.Println("No requests were sent.")

		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	runCmd.Flags().StringVarP(&url, "url", "u", "", "Url to send requests to")
	runCmd.Flags().DurationVarP(&duration, "duration", "d", 5*time.Second, "How long to run the GET requests (e.g. 10s, 1m)")
	runCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 1, "Number of concurrent workers")
	runCmd.Flags().StringVarP(&method, "method", "m", "GET", "HTTP method to use (GET or POST)")
	runCmd.Flags().StringVar(&body, "body", "", "Raw JSON body for POST request")
}
