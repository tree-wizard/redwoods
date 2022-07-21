package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/xxyyx/redwoods/pkg/analyzer"
	cfg "github.com/xxyyx/redwoods/pkg/rwconfig"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(analyzeFuzzCmd)
}

var analyzeFuzzCmd = &cobra.Command{
	Use:   "analyzefuzz",
	Short: "Find the fuzzable packages in the given repository reduced to fuzzing packages",
	Long: `
	This tool will travel the AST of the Source to define packages, tests and 
	look for fuzz tests and especially their function names. Results will be stored in the /analysis folder.
	
	While traveling the AST the projects packages will be summarized into packages and available functions.
	Furthermore if files with fuzz or test in the name are found, it will be assume they are fuzzing / testing funntions.

	The result will be reduced to fuzzing packages
	`,
	Run: func(cmd *cobra.Command, args []string) {
		runAnalyze()
	},
}

func runAnalyze() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	currCfg, err := cfg.ConfigRead(pwd + "/redwoods-cfg.json")
	if err != nil {
		fmt.Println("Could not finde a redwoods-cfg.json please rund redwoods -c")
		return
	}
	root := pwd + "/workspace/" + currCfg.Name
	fmt.Println("## Analyzing Fuzz for ", root, "###")
	project := analyzer.Project{}
	onlyfuzz := true

	if _, err := os.Stat(pwd + "/analysis/" + currCfg.Name + ".json"); os.IsNotExist(err) {
		project = analyzer.NewAnalysis(root)
		project = analyzer.Analyze(project, onlyfuzz)
	} else {
		fmt.Println("using previous analysis in ", pwd+"/analysis/"+currCfg.Name+".json")
		project, err = analyzer.FromJson(pwd + "/analysis/" + currCfg.Name + ".json")
		if err != nil {
			fmt.Println("Problem parsing the analysis")
			return
		}
	}
	//reset the fuzzpackages because they could have changed
	currCfg.Sourcesettings.FuzzPackages = make([]cfg.Fuzzpackage, 0)
	for _, pkg := range project.FuzzPackages {
		if pkg.HasGoFuzzTests {
			for _, pkgFunc := range pkg.Functions {
				if pkgFunc.HasGoFuzzTests {
					fuzzPkg := cfg.Fuzzpackage{
						FuncName:     pkgFunc.Name,
						PackageName:  strings.Replace(pkg.Location, root+"/", "", -1),
						BuildArchive: pwd + "/fuzz/" + strings.Replace(pkg.Location, root+"/", "", -1) + "/" + pkg.Name + "-" + pkgFunc.Name + ".zip",
					}
					currCfg.Sourcesettings.FuzzPackages = append(currCfg.Sourcesettings.FuzzPackages, fuzzPkg)
				}
			}
		}
	}

	err = cfg.ConfigWrite(*currCfg, pwd+"/redwoods-cfg.json")
	if err != nil {
		fmt.Println(err)
		return
	}
}
