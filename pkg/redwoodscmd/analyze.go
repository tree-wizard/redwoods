package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/xxyyx/redwoods/pkg/analyzer"
	cfg "github.com/xxyyx/redwoods/pkg/rwconfig"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(analyzeCmd)
}

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyse the AST of the source for enhanced information",
	Long: `This tool will travel the AST of the Source to define packages, tests and
look for fuzz tests and especially their function names. Results will be stored in the /analysis folder. While traveling the AST the projects packages will be summarized into packages and available functions. Furthermore if files with fuzz or test in the name are found, it will be assume they are fuzzing / testing funntions.`,
	Run: func(cmd *cobra.Command, args []string) {
		runAnalyzeComplete()
	},
}

//runAnalyzeComplete wil create a complete ast representation of the given source project and write the found fuzzpackages to the config
func runAnalyzeComplete() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	currCfg, err := cfg.ConfigRead(pwd + "/redwoods-cfg.json")
	if err != nil {
		fmt.Println("Could not find a redwoods-cfg.json please rund redwoods -c")
		return
	}
	root := pwd + "/workspace/" + currCfg.Name
	fmt.Println("## Analyzing ", root, "###")

	//lets scan what we found
	project := analyzer.NewAnalysis(root)
	onlyfuzz := false
	project = analyzer.Analyze(project, onlyfuzz)
	analyzer.ToJson(project, pwd+"/analysis/"+currCfg.Name+".json")
	err = cfg.ConfigWrite(*currCfg, pwd+"/redwoods-cfg.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Code analysis finished, the results can be found at", pwd+"/analysis/"+currCfg.Name+".json")
}
