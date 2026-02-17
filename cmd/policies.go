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

var policiesChangesCmd = &cobra.Command{
	Use:   "changes",
	Short: "Get policy change feed",
	Long:  "Track changes across all policies - new, updated, or retired",
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New(getAPIKey(), getBaseURL())

		path := "/policies/changes?limit=50"

		since, _ := cmd.Flags().GetString("since")
		if since != "" {
			path += "&since=" + since
		}

		policyID, _ := cmd.Flags().GetString("policy-id")
		if policyID != "" {
			path += "&policy_id=" + policyID
		}

		changeType, _ := cmd.Flags().GetString("change-type")
		if changeType != "" {
			path += "&change_type=" + changeType
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
			printPolicyChanges(result)
		}
	},
}

var policiesCompareCmd = &cobra.Command{
	Use:   "compare [procedure-codes...]",
	Short: "Compare policies across jurisdictions",
	Long:  "Compare coverage policies for procedures across MAC jurisdictions",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New(getAPIKey(), getBaseURL())

		reqBody := map[string]interface{}{
			"procedure_codes": args,
		}

		policyType, _ := cmd.Flags().GetString("type")
		if policyType != "" {
			reqBody["policy_type"] = policyType
		}

		jurisdictions, _ := cmd.Flags().GetStringSlice("jurisdictions")
		if len(jurisdictions) > 0 {
			reqBody["jurisdictions"] = jurisdictions
		}

		var result map[string]interface{}
		if err := c.Post("/policies/compare", reqBody, &result); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		output := getOutput()
		if output == "json" {
			jsonData, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(jsonData))
		} else {
			printPolicyComparison(result)
		}
	},
}

func init() {
	rootCmd.AddCommand(policiesCmd)
	policiesCmd.AddCommand(policiesListCmd)
	policiesCmd.AddCommand(policiesGetCmd)
	policiesCmd.AddCommand(policiesChangesCmd)
	policiesCmd.AddCommand(policiesCompareCmd)

	policiesListCmd.Flags().StringP("query", "q", "", "Search query")
	policiesListCmd.Flags().StringP("mode", "m", "keyword", "Search mode (keyword, semantic)")
	policiesListCmd.Flags().StringP("type", "t", "", "Policy type (LCD, Article, NCD)")
	policiesListCmd.Flags().StringP("jurisdiction", "j", "", "MAC jurisdiction")
	policiesListCmd.Flags().StringP("status", "s", "active", "Status (active, retired, all)")

	policiesGetCmd.Flags().StringSliceP("include", "i", []string{}, "Include additional data (criteria, codes, attachments, versions)")

	policiesChangesCmd.Flags().String("since", "", "ISO8601 timestamp - only show changes after this date")
	policiesChangesCmd.Flags().String("policy-id", "", "Filter to a specific policy")
	policiesChangesCmd.Flags().String("change-type", "", "Filter by change type (created, updated, retired)")

	policiesCompareCmd.Flags().StringP("type", "t", "", "Policy type (LCD, Article, NCD)")
	policiesCompareCmd.Flags().StringSliceP("jurisdictions", "j", []string{}, "Specific jurisdictions to compare")
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

func printPolicyChanges(result map[string]interface{}) {
	data, ok := result["data"].([]interface{})
	if !ok || len(data) == 0 {
		fmt.Println("No policy changes found")
		return
	}

	fmt.Printf("Found %d changes:\n\n", len(data))
	for _, c := range data {
		change := c.(map[string]interface{})
		fmt.Printf("Policy: %v\n", change["policy_id"])
		fmt.Printf("Type: %v\n", change["change_type"])
		if summary, ok := change["change_summary"].(string); ok && summary != "" {
			fmt.Printf("Summary: %s\n", summary)
		}
		if ts, ok := change["timestamp"].(string); ok && ts != "" {
			fmt.Printf("Date: %s\n", ts)
		}
		fmt.Println("---")
	}
}

func printPolicyComparison(result map[string]interface{}) {
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		fmt.Println("Invalid response format")
		return
	}

	if comparison, ok := data["comparison"].([]interface{}); ok {
		for _, c := range comparison {
			comp := c.(map[string]interface{})
			fmt.Printf("Jurisdiction: %v (%v)\n", comp["jurisdiction"], comp["mac_name"])
			if policies, ok := comp["policies"].([]interface{}); ok {
				fmt.Printf("  Policies: %d\n", len(policies))
				for _, p := range policies {
					policy := p.(map[string]interface{})
					fmt.Printf("    - %v: %v (%v)\n", policy["policy_id"], policy["title"], policy["disposition"])
				}
			}
			fmt.Println("---")
		}
	}
}
