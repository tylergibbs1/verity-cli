package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tylerbryy/verity-cli/pkg/client"
)

var coverageCmd = &cobra.Command{
	Use:   "coverage",
	Short: "Coverage criteria commands",
	Long:  "Search and manage coverage criteria",
}

var coverageSearchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search coverage criteria",
	Long:  "Search coverage criteria text across all policies",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New(getAPIKey(), getBaseURL())

		path := fmt.Sprintf("/coverage/criteria?q=%s", args[0])

		section, _ := cmd.Flags().GetString("section")
		if section != "" {
			path += "&section=" + section
		}

		policyType, _ := cmd.Flags().GetString("type")
		if policyType != "" {
			path += "&policy_type=" + policyType
		}

		jurisdiction, _ := cmd.Flags().GetString("jurisdiction")
		if jurisdiction != "" {
			path += "&jurisdiction=" + jurisdiction
		}

		limit, _ := cmd.Flags().GetInt("limit")
		path += fmt.Sprintf("&limit=%d", limit)

		var result map[string]interface{}
		if err := c.Get(path, &result); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		output := getOutput()
		if output == "json" {
			jsonData, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(jsonData))
		} else {
			printCriteriaResults(result)
		}
	},
}

func init() {
	rootCmd.AddCommand(coverageCmd)
	coverageCmd.AddCommand(coverageSearchCmd)

	coverageSearchCmd.Flags().StringP("section", "s", "", "Filter by section (indications, limitations, documentation)")
	coverageSearchCmd.Flags().StringP("type", "t", "", "Policy type (LCD, Article, NCD)")
	coverageSearchCmd.Flags().StringP("jurisdiction", "j", "", "MAC jurisdiction")
	coverageSearchCmd.Flags().IntP("limit", "l", 50, "Results per page (1-100)")
}

func printCriteriaResults(result map[string]interface{}) {
	data, ok := result["data"].([]interface{})
	if !ok || len(data) == 0 {
		fmt.Println("No criteria found")
		return
	}

	fmt.Printf("Found %d criteria blocks:\n\n", len(data))
	for _, c := range data {
		criteria := c.(map[string]interface{})
		if policyId, ok := criteria["policy_id"].(string); ok {
			fmt.Printf("Policy: %s", policyId)
			if title, ok := criteria["policy_title"].(string); ok {
				fmt.Printf(" - %s", title)
			}
			fmt.Println()
		}
		if section, ok := criteria["section"].(string); ok {
			fmt.Printf("Section: %s\n", section)
		}
		if text, ok := criteria["text"].(string); ok {
			if len(text) > 200 {
				text = text[:200] + "..."
			}
			fmt.Printf("  %s\n", text)
		}
		fmt.Println("---")
	}
}
