package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tylerbryy/verity-cli/pkg/client"
)

var jurisdictionsCmd = &cobra.Command{
	Use:   "jurisdictions",
	Short: "List MAC jurisdictions",
	Long:  "List all Medicare Administrative Contractor (MAC) jurisdictions",
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New(getAPIKey(), getBaseURL())

		var result map[string]interface{}
		if err := c.Get("/jurisdictions", &result); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		output := getOutput()
		if output == "json" {
			jsonData, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(jsonData))
		} else {
			printJurisdictions(result)
		}
	},
}

func init() {
	rootCmd.AddCommand(jurisdictionsCmd)
}

func printJurisdictions(result map[string]interface{}) {
	data, ok := result["data"].([]interface{})
	if !ok || len(data) == 0 {
		fmt.Println("No jurisdictions found")
		return
	}

	fmt.Printf("Found %d jurisdictions:\n\n", len(data))
	for _, j := range data {
		juris := j.(map[string]interface{})
		fmt.Printf("%-6v %v\n", juris["jurisdiction_code"], juris["mac_name"])
		if states, ok := juris["states"].([]interface{}); ok && len(states) > 0 {
			fmt.Printf("       States: ")
			for i, s := range states {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Print(s)
			}
			fmt.Println()
		}
		fmt.Println("---")
	}
}
