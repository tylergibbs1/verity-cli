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

var priorAuthResearchCmd = &cobra.Command{
	Use:   "research [procedure-codes...]",
	Short: "Research prior auth requirements via AI web search",
	Long:  "Use AI-powered web research to find prior authorization requirements from payer websites",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New(getAPIKey(), getBaseURL())

		payer, _ := cmd.Flags().GetString("payer")
		state, _ := cmd.Flags().GetString("state")
		diagnosisCodes, _ := cmd.Flags().GetStringSlice("diagnosis")
		clinicalContext, _ := cmd.Flags().GetString("context")
		syncMode, _ := cmd.Flags().GetBool("sync")

		reqBody := map[string]interface{}{
			"procedure_codes": args,
			"sync":            syncMode,
		}

		if payer != "" {
			reqBody["payer"] = payer
		}
		if state != "" {
			reqBody["state"] = state
		}
		if len(diagnosisCodes) > 0 {
			reqBody["diagnosis_codes"] = diagnosisCodes
		}
		if clinicalContext != "" {
			reqBody["clinical_context"] = clinicalContext
		}

		var result map[string]interface{}
		if err := c.Post("/prior-auth/research", reqBody, &result); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		output := getOutput()
		if output == "json" {
			jsonData, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(jsonData))
		} else {
			printResearchResult(result)
		}
	},
}

var priorAuthResearchStatusCmd = &cobra.Command{
	Use:   "research-status [research-id]",
	Short: "Get prior auth research status",
	Long:  "Poll the status and results of a prior authorization research task",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New(getAPIKey(), getBaseURL())

		path := fmt.Sprintf("/prior-auth/research/%s", args[0])

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
			printResearchResult(result)
		}
	},
}

func init() {
	rootCmd.AddCommand(priorAuthCmd)
	priorAuthCmd.AddCommand(priorAuthResearchCmd)
	priorAuthCmd.AddCommand(priorAuthResearchStatusCmd)

	priorAuthCmd.Flags().StringSliceP("diagnosis", "d", []string{}, "Diagnosis codes (ICD-10)")
	priorAuthCmd.Flags().StringP("state", "s", "", "Two-letter state code")
	priorAuthCmd.Flags().StringP("payer", "p", "medicare", "Payer (medicare, aetna, uhc, all)")

	priorAuthResearchCmd.Flags().StringP("payer", "p", "", "Payer name (e.g., UnitedHealthcare, Aetna)")
	priorAuthResearchCmd.Flags().StringP("state", "s", "", "Two-letter state code")
	priorAuthResearchCmd.Flags().StringSliceP("diagnosis", "d", []string{}, "Diagnosis codes (ICD-10)")
	priorAuthResearchCmd.Flags().String("context", "", "Additional clinical context")
	priorAuthResearchCmd.Flags().Bool("sync", false, "Wait for completion instead of returning research ID")
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

func printResearchResult(result map[string]interface{}) {
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		fmt.Println("Invalid response format")
		return
	}

	fmt.Printf("Research ID: %v\n", data["research_id"])
	fmt.Printf("Status: %v\n", data["status"])

	if createdAt, ok := data["created_at"].(string); ok {
		fmt.Printf("Created: %s\n", createdAt)
	}

	if pollURL, ok := data["poll_url"].(string); ok && pollURL != "" {
		fmt.Printf("Poll URL: %s\n", pollURL)
		fmt.Println("\nUse 'verity prior-auth research-status <research-id>' to check progress")
	}

	if resResult, ok := data["result"].(map[string]interface{}); ok {
		fmt.Println("\nResults:")
		if determination, ok := resResult["determination"].(map[string]interface{}); ok {
			fmt.Printf("  PA Required: %v\n", determination["pa_required"])
			fmt.Printf("  Confidence: %v\n", determination["confidence"])
			fmt.Printf("  Reasoning: %v\n", determination["reasoning"])
		}

		if docReqs, ok := resResult["documentation_requirements"].([]interface{}); ok && len(docReqs) > 0 {
			fmt.Println("\n  Documentation Requirements:")
			for _, req := range docReqs {
				fmt.Printf("    - %v\n", req)
			}
		}

		if sources, ok := resResult["sources"].([]interface{}); ok && len(sources) > 0 {
			fmt.Println("\n  Sources:")
			for _, src := range sources {
				fmt.Printf("    - %v\n", src)
			}
		}
	}

	if errMsg, ok := data["error"].(string); ok && errMsg != "" {
		fmt.Printf("\nError: %s\n", errMsg)
	}
}
