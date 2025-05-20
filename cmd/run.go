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
	"sync"
	"time"
)

var url string
var duration time.Duration
var concurrency int

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
						resp, err := http.Get(url)
						if err == nil {
							io.Copy(io.Discard, resp.Body)
							resp.Body.Close()
						}
						reqsChan <- 1
					}
				}
			}(i)
		}

		// Count successful requests
		go func() {
			for range reqsChan {
				totalReqs++
			}
		}()

		// Wait for all workers to finish
		wg.Wait()
		close(reqsChan)

		fmt.Printf("✅ Load complete. Sent %d requests in %v\n", totalReqs, duration)
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
}
