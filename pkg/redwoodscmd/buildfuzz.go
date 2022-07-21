package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	cfg "github.com/xxyyx/redwoods/pkg/rwconfig"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(buildfuzzCmd)
}

var buildfuzzCmd = &cobra.Command{
	Use:   "buildfuzz",
	Short: "Build the fuzz-build archives",
	Long:  `After the analyser has found all fuzzable packages, we shall build them for fuzzing. This tool calls go-fuzz-build.`,
	Run: func(cmd *cobra.Command, args []string) {
		runFuzzBuilder()
	},
}

func runFuzzBuilder() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	currCfg, err := cfg.ConfigRead(pwd + "/redwoods-cfg.json")
	if err != nil {
		fmt.Println("Could not find a redwoods-cfg.json please run $ redwoods -c")
		return
	}
	if len(currCfg.Sourcesettings.FuzzPackages) == 0 {
		fmt.Println("Could not finde a previous code analysis!")
		runAnalyze()
		currCfg, err = cfg.ConfigRead(pwd + "/redwoods-cfg.json")
		if err != nil {
			fmt.Println("Tried to update redwoods-cfg.json please but failed run $ redwoods -c")
			return
		}
	}
	root := pwd + "/workspace/" + currCfg.Name
	fmt.Println("Starting to Build the Fuzz packages for", root)
	if len(currCfg.Sourcesettings.FuzzPackages) == 0 {
		fmt.Println("No fuzzable packages found")
		return
	}
	for _, fuzzPkg := range currCfg.Sourcesettings.FuzzPackages {
		if _, err := os.Stat(fuzzPkg.BuildArchive); os.IsNotExist(err) {
			fmt.Println("package ", fuzzPkg.PackageName, "has not been built yet, building")
			fuzzBuild(fuzzPkg.PackageName, fuzzPkg.FuncName, fuzzPkg.BuildArchive, pwd, *currCfg)
		} else {
			fmt.Println("package ", fuzzPkg.PackageName, "has been built already, skipping.")
		}

	}

}

//fuzzBuild uses go-fuzz-build to create the packages for the go-fuzzer
func fuzzBuild(pkgName, funcName, buildarchivePath, pwd string, currCfg cfg.RWConfiguration) {
	// go-fuzz-build -func "Fuzz_$NAME" -o "$FUZZ_BASE_WORKDIRECTORY/$pkg/$NAME-fuzz.zip" "$FUZZ_WORKSPACE/$pkg"
	fmt.Println("go-fuzz-build", "-func=\""+funcName+"\"", "-o="+buildarchivePath, pwd+"/workspace/"+currCfg.Name+"/"+pkgName)
	dir, _ := filepath.Split(buildarchivePath)
	os.MkdirAll(dir, os.ModePerm)
	e := exec.Command("go-fuzz-build", "-func="+funcName, "-o="+buildarchivePath, pwd+"/workspace/"+currCfg.Name+"/"+pkgName)
	e.Dir = pwd + "/workspace/" + currCfg.Name
	e.Stdout = os.Stdout
	e.Stderr = os.Stderr
	err := e.Run()
	if err != nil {
		log.Fatal(err)
	}
}
