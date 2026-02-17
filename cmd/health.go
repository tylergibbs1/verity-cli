package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tylerbryy/verity-cli/pkg/client"
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check API health status",
	Long:  "Check the health status of the Verity API including database and Redis checks",
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New(getAPIKey(), getBaseURL())

		var result map[string]interface{}
		if err := c.Get("/health", &result); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		output := getOutput()
		if output == "json" {
			jsonData, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(jsonData))
		} else {
			printHealthResult(result)
		}
	},
}

func init() {
	rootCmd.AddCommand(healthCmd)
}

func printHealthResult(result map[string]interface{}) {
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		fmt.Println("Invalid response format")
		return
	}

	fmt.Printf("Status: %v\n", data["status"])
	fmt.Printf("Version: %v\n", data["version"])
	fmt.Printf("Timestamp: %v\n", data["timestamp"])

	if checks, ok := data["checks"].(map[string]interface{}); ok {
		fmt.Println("\nChecks:")
		for name, check := range checks {
			if checkMap, ok := check.(map[string]interface{}); ok {
				fmt.Printf("  %s: %v\n", name, checkMap["status"])
			}
		}
	}
}
