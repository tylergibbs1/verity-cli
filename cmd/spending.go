package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tylerbryy/verity-cli/pkg/client"
)

var spendingCmd = &cobra.Command{
	Use:   "spending [codes...]",
	Short: "Get Medicaid spending data by HCPCS code",
	Long:  "Returns aggregate Medicaid provider spending statistics per HCPCS code",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New(getAPIKey(), getBaseURL())

		var path string
		if len(args) == 1 {
			path = fmt.Sprintf("/spending/by-code?code=%s", args[0])
		} else {
			path = fmt.Sprintf("/spending/by-code?codes=%s", strings.Join(args, ","))
		}

		year, _ := cmd.Flags().GetInt("year")
		if year > 0 {
			path += fmt.Sprintf("&year=%d", year)
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
			printSpendingResult(result)
		}
	},
}

func init() {
	rootCmd.AddCommand(spendingCmd)

	spendingCmd.Flags().IntP("year", "y", 0, "Filter to a specific year")
}

func printSpendingResult(result map[string]interface{}) {
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		fmt.Println("No spending data found")
		return
	}

	for code, info := range data {
		spending, ok := info.(map[string]interface{})
		if !ok {
			continue
		}

		fmt.Printf("Code: %s\n", code)
		fmt.Printf("  Total Paid: $%v\n", spending["total_paid"])
		fmt.Printf("  Total Claims: %v\n", spending["total_claims"])
		fmt.Printf("  Unique Beneficiaries: %v\n", spending["unique_beneficiaries"])
		fmt.Printf("  Unique Providers: %v\n", spending["unique_providers"])

		if byYear, ok := spending["by_year"].([]interface{}); ok && len(byYear) > 0 {
			fmt.Println("  By Year:")
			for _, y := range byYear {
				yr := y.(map[string]interface{})
				fmt.Printf("    %v: $%v (%v claims)\n", yr["year"], yr["total_paid"], yr["total_claims"])
			}
		}
		fmt.Println("---")
	}
}
