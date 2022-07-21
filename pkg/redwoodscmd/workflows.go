package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(workflowCmd)
}

var workflowCmd = &cobra.Command{
	Use:   "workflow",
	Short: "Run predefined workflows automatically",
	Long: `
This will run any predefined workflow for you
Parameter:
	work : runs the entire workflow
	fuzz : analyses builds and runs the fuzzer
	rerun: reruns fuzz execution	
	`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "work":
			workflowWork()
		case "fuzz":
			workflowFuzz()
		case "rerun":
			workflowRerun()
		}
	},
}

func workflowWork() {
	fmt.Println("## Worflow work  Executing ###")
	initRedwoods()
	installDepdenencies()
	runAnalyzeComplete()
	runFuzzBuilder()
	runFuzzExecution()
	runOutputAnalysis()
}
func workflowRerun() {
	fmt.Println("## Worflow rerun  Executing ###")
	runFuzzExecution()
	runOutputAnalysis()
}
func workflowFuzz() {
	fmt.Println("## Worflow fuzz  Executing ###")
	runAnalyzeComplete()
	runFuzzBuilder()
	runFuzzExecution()
	runOutputAnalysis()
}
