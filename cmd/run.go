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
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var url string

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run gnat load generator",
	Run: func(cmd *cobra.Command, args []string) {
		if url == "" {
			fmt.Fprintln(os.Stderr, "❌ Please provide a URL with --url")
			os.Exit(1)
		}
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Request failed: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to read response body: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Response from %s:\n\n%s\n", url, string(body))
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
}
