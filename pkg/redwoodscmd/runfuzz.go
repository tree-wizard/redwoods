package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	cfg "github.com/xxyyx/redwoods/pkg/rwconfig"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runfuzzCmd)
}

var runfuzzCmd = &cobra.Command{
	Use:   "runfuzz",
	Short: "Build the fuzz-build archives",
	Long:  `Takes the proposed fuzzpackages and runs them through go-fuzz`,
	Run: func(cmd *cobra.Command, args []string) {
		runFuzzExecution()

	},
}

func runFuzzExecution() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	currCfg, err := cfg.ConfigRead("./redwoods-cfg.json")
	if err != nil {
		fmt.Println("Could not find a redwoods-cfg.json please run $ redwoods -c")
		return
	}
	if len(currCfg.Sourcesettings.FuzzPackages) > 0 {
		fmt.Println("Could not finde a previous code analysis!")
		runAnalyze()
		currCfg, err = cfg.ConfigRead("./redwoods-cfg.json")
		if err != nil {
			fmt.Println("Tried to update redwoods-cfg.json please but failed run $ redwoods -c")
			return
		}
	}
	exeGofuzz(pwd, currCfg)
}

func exeGofuzz(pwd string, currCfg *cfg.RWConfiguration) {
	fmt.Println("## Running Fuzz Execution ###")
	if len(currCfg.Sourcesettings.FuzzPackages) == 0 {
		fmt.Println("No fuzzable packages found")
		return
	}
	for _, fuzzPkg := range currCfg.Sourcesettings.FuzzPackages {
		if _, err := os.Stat(fuzzPkg.BuildArchive); os.IsNotExist(err) {
			fmt.Println("fuzz build package for ", fuzzPkg.BuildArchive, "does not exist, skipping")
		} else {

			dir, _ := filepath.Split(fuzzPkg.BuildArchive)
			os.MkdirAll(dir+"/Output", os.ModePerm)
			//COMMAND="go-fuzz -bin=$FUZZ_BASE_WORKDIRECTORY/$pkg/$NAME-fuzz.zip -workdir=$WORKDIR/Work -procs=$NOOFPROCESSES -timeout=$TIMEOUT_RUN_TIME "
			fmt.Println("running fuzz for", fuzzPkg.PackageName, "with function", fuzzPkg.FuncName, "for", currCfg.Fuzzsettings.FUZZ_TEST_RUN_TIME, "seconds")
			e := exec.Command("go-fuzz", "-bin="+fuzzPkg.BuildArchive, "-workdir="+pwd+"/fuzz/"+fuzzPkg.PackageName+"/work/"+fuzzPkg.FuncName, "-timeout="+strconv.FormatInt(currCfg.Fuzzsettings.FUZZ_TEST_TIMEOUT, 10), "-procs="+strconv.FormatInt(currCfg.Fuzzsettings.FUZZ_NUM_PROCESSES, 10))
			var stdBuffer bytes.Buffer
			mw := io.MultiWriter(os.Stdout, &stdBuffer)
			e.Stdout = mw
			e.Stderr = mw
			err := e.Start()
			if err != nil {
				log.Fatal(err)
			}
			done := make(chan error)
			go func() { done <- e.Wait() }()

			// Start a timer
			timeout := time.After(time.Duration(currCfg.Fuzzsettings.FUZZ_TEST_RUN_TIME) * time.Second)

			select {
			case <-timeout:
				e.Process.Kill()
				fmt.Println("Execution time passed. Killing proccess after", currCfg.Fuzzsettings.FUZZ_TEST_RUN_TIME, "seconds")
			case err := <-done:
				fmt.Println("Execution finished before the timout of", currCfg.Fuzzsettings.FUZZ_TEST_RUN_TIME, "seconds")
				if err != nil {
					fmt.Println("Non-zero exit code:", err)
				}
			}
			err = ioutil.WriteFile(dir+"Output/"+fuzzPkg.FuncName+"-fuzz.log", stdBuffer.Bytes(), 0777)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Output written to", dir+"Output/"+fuzzPkg.FuncName+"-fuzz.log")
		}

	}
}
