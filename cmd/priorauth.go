package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tylerbryy/verity-cli/pkg/client"
)

var priorAuthCmd = &cobra.Command{
	Use:   "prior-auth [procedure-codes...]",
	Short: "Check prior authorization requirements",
	Long:  "Check if procedures require prior authorization based on codes and state",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New(getAPIKey(), getBaseURL())

		diagnosisCodes, _ := cmd.Flags().GetStringSlice("diagnosis")
		state, _ := cmd.Flags().GetString("state")
		payer, _ := cmd.Flags().GetString("payer")

		reqBody := map[string]interface{}{
			"procedure_codes":    args,
			"payer":              payer,
			"criteria_page":      1,
			"criteria_per_page":  25,
		}

		if len(diagnosisCodes) > 0 {
			reqBody["diagnosis_codes"] = diagnosisCodes
		}

		if state != "" {
			reqBody["state"] = state
		}

		var result map[string]interface{}
		if err := c.Post("/prior-auth/check", reqBody, &result); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		output := getOutput()
		if output == "json" {
			jsonData, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(jsonData))
		} else {
			printPriorAuthResult(result)
		}
	},
}

func init() {
	rootCmd.AddCommand(priorAuthCmd)

	priorAuthCmd.Flags().StringSliceP("diagnosis", "d", []string{}, "Diagnosis codes (ICD-10)")
	priorAuthCmd.Flags().StringP("state", "s", "", "Two-letter state code")
	priorAuthCmd.Flags().StringP("payer", "p", "medicare", "Payer (medicare, aetna, uhc, all)")
}

func printPriorAuthResult(result map[string]interface{}) {
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		fmt.Println("Invalid response format")
		return
	}

	fmt.Printf("Prior Authorization Required: %v\n", data["pa_required"])
	fmt.Printf("Confidence: %v\n", data["confidence"])
	fmt.Printf("Reason: %v\n\n", data["reason"])

	if policies, ok := data["matched_policies"].([]interface{}); ok && len(policies) > 0 {
		fmt.Println("Matched Policies:")
		for _, p := range policies {
			policy := p.(map[string]interface{})
			fmt.Printf("  - %s: %s\n", policy["policy_id"], policy["title"])
		}
		fmt.Println()
	}

	if checklist, ok := data["documentation_checklist"].([]interface{}); ok && len(checklist) > 0 {
		fmt.Println("Documentation Checklist:")
		for _, item := range checklist {
			fmt.Printf("  - %v\n", item)
		}
	}
}
