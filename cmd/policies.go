package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tylerbryy/verity-cli/pkg/client"
)

var policiesCmd = &cobra.Command{
	Use:   "policies",
	Short: "Search and manage policies",
	Long:  "Search and list coverage policies, or get details of a specific policy",
}

var policiesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Search and list policies",
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New(getAPIKey(), getBaseURL())

		path := "/policies?limit=50"

		query, _ := cmd.Flags().GetString("query")
		if query != "" {
			path += "&q=" + query
		}

		mode, _ := cmd.Flags().GetString("mode")
		path += "&mode=" + mode

		policyType, _ := cmd.Flags().GetString("type")
		if policyType != "" {
			path += "&policy_type=" + policyType
		}

		jurisdiction, _ := cmd.Flags().GetString("jurisdiction")
		if jurisdiction != "" {
			path += "&jurisdiction=" + jurisdiction
		}

		status, _ := cmd.Flags().GetString("status")
		path += "&status=" + status

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
			printPoliciesList(result)
		}
	},
}

var policiesGetCmd = &cobra.Command{
	Use:   "get [policy-id]",
	Short: "Get policy details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		policyID := args[0]
		c := client.New(getAPIKey(), getBaseURL())

		path := fmt.Sprintf("/policies/%s", policyID)

		include, _ := cmd.Flags().GetStringSlice("include")
		if len(include) > 0 {
			includeStr := ""
			for i, inc := range include {
				if i > 0 {
					includeStr += ","
				}
				includeStr += inc
			}
			path += "?include=" + includeStr
		}

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
			printPolicyDetail(result)
		}
	},
}

func init() {
	rootCmd.AddCommand(policiesCmd)
	policiesCmd.AddCommand(policiesListCmd)
	policiesCmd.AddCommand(policiesGetCmd)

	policiesListCmd.Flags().StringP("query", "q", "", "Search query")
	policiesListCmd.Flags().StringP("mode", "m", "keyword", "Search mode (keyword, semantic)")
	policiesListCmd.Flags().StringP("type", "t", "", "Policy type (LCD, Article, NCD)")
	policiesListCmd.Flags().StringP("jurisdiction", "j", "", "MAC jurisdiction")
	policiesListCmd.Flags().StringP("status", "s", "active", "Status (active, retired, all)")

	policiesGetCmd.Flags().StringSliceP("include", "i", []string{}, "Include additional data (criteria, codes, attachments, versions)")
}

func printPoliciesList(result map[string]interface{}) {
	data, ok := result["data"].([]interface{})
	if !ok || len(data) == 0 {
		fmt.Println("No policies found")
		return
	}

	fmt.Printf("Found %d policies:\n\n", len(data))
	for _, p := range data {
		policy := p.(map[string]interface{})
		fmt.Printf("ID: %v\n", policy["policy_id"])
		fmt.Printf("Title: %v\n", policy["title"])
		fmt.Printf("Type: %v\n", policy["policy_type"])
		if juris, ok := policy["jurisdiction"].(string); ok && juris != "" {
			fmt.Printf("Jurisdiction: %s\n", juris)
		}
		fmt.Printf("Status: %v\n", policy["status"])
		fmt.Println("---")
	}
}

func printPolicyDetail(result map[string]interface{}) {
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		fmt.Println("Invalid response format")
		return
	}

	fmt.Printf("Policy ID: %v\n", data["policy_id"])
	fmt.Printf("Title: %v\n", data["title"])
	fmt.Printf("Type: %v\n", data["policy_type"])
	fmt.Printf("Status: %v\n", data["status"])

	if juris, ok := data["jurisdiction"].(string); ok && juris != "" {
		fmt.Printf("Jurisdiction: %s\n", juris)
	}

	if date, ok := data["effective_date"].(string); ok && date != "" {
		fmt.Printf("Effective Date: %s\n", date)
	}

	if desc, ok := data["description"].(string); ok && desc != "" {
		fmt.Printf("\nDescription:\n%s\n", desc)
	}

	if summary, ok := data["summary"].(string); ok && summary != "" {
		fmt.Printf("\nSummary:\n%s\n", summary)
	}
}
