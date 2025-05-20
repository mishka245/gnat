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
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
	"time"
)

var url string
var duration time.Duration

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run gnat load generator",
	Run: func(cmd *cobra.Command, args []string) {
		if url == "" {
			fmt.Fprintln(os.Stderr, "❌ Please provide a URL with --url")
			os.Exit(1)
		}
		endTime := time.Now().Add(duration)
		count := 0

		for time.Now().Before(endTime) {
			resp, err := http.Get(url)
			if err != nil {
				fmt.Fprintf(os.Stderr, "❌ Request %d failed: %v\n", count, err)
				continue
			}

			io.Copy(io.Discard, resp.Body) // Discard response for now
			resp.Body.Close()

			count++
		}

		fmt.Printf("✅ Finished. Sent %d requests in %v\n", count, duration)
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
}
