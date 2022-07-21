package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	cfg "github.com/xxyyx/redwoods/pkg/rwconfig"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(outputanalysisCmd)
}

var outputanalysisCmd = &cobra.Command{
	Use:   "outputanalysis",
	Short: "Analyses the output after fuzzing",
	Long:  `This will go through the output and return results of the fuzzers.`,
	Run: func(cmd *cobra.Command, args []string) {
		runOutputAnalysis()

	},
}

func runOutputAnalysis() {
	fmt.Println("## Starting output analysis ###")
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	currCfg, err := cfg.ConfigRead("./redwoods-cfg.json")
	if err != nil {
		fmt.Println("Could not find a redwoods-cfg.json please run $ redwoods -c")
		return
	}
	if len(currCfg.Sourcesettings.FuzzPackages) == 0 {
		fmt.Println("No fuzzable packages found")
		return
	}

	for _, fuzzPkg := range currCfg.Sourcesettings.FuzzPackages {
		outputPath := pwd + "/fuzz/" + fuzzPkg.PackageName + "/Output/" + fuzzPkg.FuncName + "-fuzz.log"
		crashersPath := pwd + "/fuzz/" + fuzzPkg.PackageName + "/work/" + fuzzPkg.FuncName + "/crashers/"
		if _, err := os.Stat(outputPath); os.IsNotExist(err) {
			fmt.Println("package ", fuzzPkg.PackageName, fuzzPkg.FuncName, "has no output")
		} else {
			fmt.Println("package ", fuzzPkg.PackageName, fuzzPkg.FuncName, "has output:")
			parseFuzzFile(outputPath)
			files, _ := ioutil.ReadDir(crashersPath)
			if len(files) > 0 {
				fmt.Println("found ", len(files), "crashers in the work directory")
			}
		}
	}

}

func parseFuzzFile(fname string) {
	file, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	buf := make([]byte, 128)
	stat, err := os.Stat(fname)
	start := stat.Size() - 128
	_, err = file.ReadAt(buf, start)
	if err == nil {
		fmt.Printf("%s\n", buf)
	}

}
