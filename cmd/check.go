package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tylerbryy/verity-cli/pkg/client"
)

var checkCmd = &cobra.Command{
	Use:   "check [code]",
	Short: "Look up a medical code",
	Long:  "Look up a medical code (CPT, HCPCS, ICD-10, NDC) and get coverage information",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		code := args[0]

		c := client.New(getAPIKey(), getBaseURL())

		path := fmt.Sprintf("/codes/lookup?code=%s", code)

		include, _ := cmd.Flags().GetStringSlice("include")
		if len(include) > 0 {
			includeStr := ""
			for i, inc := range include {
				if i > 0 {
					includeStr += ","
				}
				includeStr += inc
			}
			path += "&include=" + includeStr
		}

		jurisdiction, _ := cmd.Flags().GetString("jurisdiction")
		if jurisdiction != "" {
			path += "&jurisdiction=" + jurisdiction
		}

		fuzzy, _ := cmd.Flags().GetBool("fuzzy")
		if !fuzzy {
			path += "&fuzzy=false"
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
			printCodeResult(result)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().StringSliceP("include", "i", []string{}, "Include additional data (rvu, policies)")
	checkCmd.Flags().StringP("jurisdiction", "j", "", "Filter by MAC jurisdiction")
	checkCmd.Flags().BoolP("fuzzy", "f", true, "Enable fuzzy matching")
}

func printCodeResult(result map[string]interface{}) {
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		fmt.Println("Invalid response format")
		return
	}

	fmt.Printf("Code: %v\n", data["code"])
	fmt.Printf("System: %v\n", data["code_system"])
	fmt.Printf("Found: %v\n", data["found"])

	if desc, ok := data["description"].(string); ok && desc != "" {
		fmt.Printf("Description: %s\n", desc)
	}

	if rvu, ok := data["rvu"].(map[string]interface{}); ok {
		fmt.Println("\nRVU Data:")
		if workRvu, ok := rvu["work_rvu"].(string); ok {
			fmt.Printf("  Work RVU: %s\n", workRvu)
		}
		if price, ok := rvu["non_facility_price"].(string); ok {
			fmt.Printf("  Non-Facility Price: $%s\n", price)
		}
		if price, ok := rvu["facility_price"].(string); ok {
			fmt.Printf("  Facility Price: $%s\n", price)
		}
	}

	if policies, ok := data["policies"].([]interface{}); ok && len(policies) > 0 {
		fmt.Println("\nPolicies:")
		for _, p := range policies {
			policy := p.(map[string]interface{})
			fmt.Printf("  - %s (%s): %s\n", policy["policy_id"], policy["policy_type"], policy["disposition"])
			fmt.Printf("    %s\n", policy["title"])
		}
	}
}
