package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/xxyyx/redwoods/pkg/clicommon"
	cfg "github.com/xxyyx/redwoods/pkg/rwconfig"
	"github.com/xxyyx/redwoods/pkg/version"
	"github.com/spf13/cobra"
)

var (
	vers       bool
	showConfig bool
	setConfig  bool
)

func init() {
	cobra.OnInitialize()
	rootCmd.Flags().BoolVarP(&vers, "version", "v", false, "show current version of CLI")
	rootCmd.Flags().BoolVarP(&showConfig, "show config", "s", false, "show current configuration")
	rootCmd.Flags().BoolVarP(&setConfig, "set config", "c", false, "set new configuration")
}

var rootCmd = &cobra.Command{
	Use:   "redwoods",
	Short: "TODO",
	Long: `
###############################################################
Welcome to the Redwoods fuzzing Suite
###############################################################
The Idea behind this project is to have a base for automated fuzzing of your code packages.
Fuzzing is a great way to find critical errors in your code that otherwise would have remained hidden.

How to use:
In order to use redwoods you first need to create a config. You can do that either through:
$ redwoods -c
which will take you through a wizzard or use the 
$ redwoods defaultconfig 
to create a redwoods-cfg.json in you working directory.

NOTE: This program will create subfolders like "fuzz", "analyze" and "workspace" so its best to run it
in a directory of your choice apart from a project.

You find addionnal help with:
$ redwoods [command] -h
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !cfg.ConfigExists("./redwoods-cfg.json") {
			rwcfg := cfg.DefaultConfig()
			err := cfg.ConfigWrite(rwcfg, "./redwoods-cfg.json")
			if err != nil {
				return err
			}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if vers {
			return printVersion()
		}
		if showConfig {
			return showCurrentConfig()
		}
		if setConfig {
			return setNewConfig()
		}
		return cmd.Help()
	},
}

func printVersion() error {
	log.Printf(
		"Starting the service...\ncommit: %s, build time: %s, release: %s",
		version.Commit, version.BuildTime, version.Release,
	)
	return nil
}

func showCurrentConfig() error {
	currCfg, err := cfg.ConfigRead("./redwoods-cfg.json")
	if err != nil {
		return err
	}
	fmt.Println("current configuration:")
	empJSON, err := json.MarshalIndent(currCfg, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("MarshalIndent funnction output %s\n", string(empJSON))
	return nil
}

func setNewConfig() error {
	currCfg, err := promptSetNewConfig()
	if err != nil {
		return err
	}
	sources, _ := clicommon.PromptConfirm("Would you like to configure the sources?")
	if err != nil {
		return err
	}
	if strings.ToLower(sources) == "y" {
		currCfg, err = promptAddNewSource(currCfg)
		if err != nil {
			return err
		}
	}
	dependencies, _ := clicommon.PromptConfirm("Would you like to configure the dependencies?")
	if err != nil {
		return err
	}
	if strings.ToLower(dependencies) == "y" {
		currCfg, err = promptAddDependencies(currCfg)
		if err != nil {
			return err
		}
	}
	fuzzsettings, _ := clicommon.PromptConfirm("Would you like to configure the fuzzing settings?")
	if err != nil {
		return err
	}
	if strings.ToLower(fuzzsettings) == "y" {
		currCfg, err = promptAddFuzzesttings(currCfg)
		if err != nil {
			return err
		}
	}
	golangsettings, _ := clicommon.PromptConfirm("Would you like to configure the golang settings?")
	if err != nil {
		return err
	}
	if strings.ToLower(golangsettings) == "y" {
		currCfg, err = promptAddGolangsettings(currCfg)
		if err != nil {
			return err
		}
	}

	err = cfg.ConfigWrite(*currCfg, "./redwoods-cfg.json")
	if err != nil {
		return err
	}
	fmt.Println("Override new configuration with details:")

	return nil
}

func promptSetNewConfig() (*cfg.RWConfiguration, error) {
	currCfg, err := cfg.ConfigRead("./redwoods-cfg.json")
	if err != nil {
		return nil, err
	}
	fmt.Println("## Welcome to the redwoods fuzzing suite setup! ##")
	reporUrl, err := clicommon.PromptNotEmptyString("Git Repository URL", currCfg.Repository)
	if err != nil {
		return nil, err
	}
	currCfg.Repository = reporUrl
	name, err := clicommon.PromptNotEmptyString("Project Name", currCfg.Name)
	if err != nil {
		return nil, err
	}
	currCfg.Name = name
	containerrepo, err := clicommon.PromptNotEmptyString("Container Repositor URL", currCfg.ContainerRepository)
	if err != nil {
		return nil, err
	}
	currCfg.ContainerRepository = containerrepo

	useDocker, err := clicommon.PromptConfirm("Can this machine use Docker?")
	if err != nil {
		return nil, err
	}
	if strings.ToLower(useDocker) == "y" {
		currCfg.CanUseDocker = true
	} else {
		currCfg.CanUseDocker = false
	}

	return currCfg, nil
}

//promptAddNewSource confgiures the source part of th config
func promptAddNewSource(currCfg *cfg.RWConfiguration) (*cfg.RWConfiguration, error) {
	if currCfg == nil {
		currCfg, err := cfg.ConfigRead("./redwoods-cfg.json")
		fmt.Println("loaded config", currCfg.Name)
		if err != nil {
			return nil, err
		}
	}

	branch, err := clicommon.PromptString("Branch (blank for main)", currCfg.Sourcesettings.Branch)
	if err != nil {
		return nil, err
	}
	if branch == "" {
		currCfg.Sourcesettings.Branch = "main"
	} else if branch != "" {
		currCfg.Sourcesettings.Branch = branch
	}

	commit, err := clicommon.PromptString("Commit (blank for latest)", currCfg.Sourcesettings.Commit)
	if err != nil {
		return nil, err
	}
	if commit == "" {
		currCfg.Sourcesettings.Commit = "latest"
	} else if commit != "" {
		currCfg.Sourcesettings.Commit = commit
	}

	return currCfg, nil
}

//promptAddNewSource confgiures the source part of th config
func promptAddDependencies(currCfg *cfg.RWConfiguration) (*cfg.RWConfiguration, error) {
	if currCfg == nil {
		currCfg, err := cfg.ConfigRead("./redwoods-cfg.json")
		fmt.Println("loaded config", currCfg.Name)
		if err != nil {
			return nil, err
		}
	}

	gOFUZZ_COMMIT, err := clicommon.PromptString("Commit for go-fuzz (blank for latest)", currCfg.Dependencysettings.GOFUZZ_COMMIT)
	if err != nil {
		return nil, err
	}
	if gOFUZZ_COMMIT == "" {
		currCfg.Dependencysettings.GOFUZZ_COMMIT = "latest"
	} else if gOFUZZ_COMMIT != "" {
		currCfg.Dependencysettings.GOFUZZ_COMMIT = gOFUZZ_COMMIT
	}

	gOFUZZ_BUILD_PKG, err := clicommon.PromptString("Repo for go-fuzz-build (blank for github.com/dvyukov/go-fuzz/go-fuzz-build)", currCfg.Dependencysettings.GOFUZZ_BUILD_PKG)
	if err != nil {
		return nil, err
	}
	if gOFUZZ_BUILD_PKG == "" {
		currCfg.Dependencysettings.GOFUZZ_BUILD_PKG = "github.com/dvyukov/go-fuzz/go-fuzz-build"
	} else if gOFUZZ_BUILD_PKG != "" {
		currCfg.Dependencysettings.GOFUZZ_BUILD_PKG = gOFUZZ_BUILD_PKG
	}

	gOFUZZ_PKG, err := clicommon.PromptString("Repo for go-fuzz (blank for github.com/dvyukov/go-fuzz/go-fuzz) current: ", currCfg.Dependencysettings.GOFUZZ_PKG)
	if err != nil {
		return nil, err
	}
	if gOFUZZ_PKG == "" {
		currCfg.Dependencysettings.GOFUZZ_PKG = "github.com/dvyukov/go-fuzz/go-fuzz"
	} else if gOFUZZ_PKG != "" {
		currCfg.Dependencysettings.GOFUZZ_PKG = gOFUZZ_PKG
	}
	gOFUZZ_DEP_PKG, err := clicommon.PromptString("Repo for go-fuzz-dep (blank forgithub.com/dvyukov/go-fuzz/go-fuzz-dep) current: ", currCfg.Dependencysettings.GOFUZZ_DEP_PKG)
	if err != nil {
		return nil, err
	}
	if gOFUZZ_DEP_PKG == "" {
		currCfg.Dependencysettings.GOFUZZ_DEP_PKG = "github.com/dvyukov/go-fuzz/go-fuzz-dep"
	} else if gOFUZZ_DEP_PKG != "" {
		currCfg.Dependencysettings.GOFUZZ_DEP_PKG = gOFUZZ_DEP_PKG
	}

	return currCfg, nil
}

//promptAddFuzzesttings confgiures the source part of th config
func promptAddFuzzesttings(currCfg *cfg.RWConfiguration) (*cfg.RWConfiguration, error) {
	if currCfg == nil {
		currCfg, err := cfg.ConfigRead("./redwoods-cfg.json")
		fmt.Println("loaded config", currCfg.Name)
		if err != nil {
			return nil, err
		}
	}

	fUZZ_TEST_RUN_TIME, err := clicommon.PromptInteger("Fuzz Runtime in Seconds current: " + string(currCfg.Fuzzsettings.FUZZ_TEST_RUN_TIME))
	if err != nil {
		return nil, err
	}
	currCfg.Fuzzsettings.FUZZ_TEST_RUN_TIME = fUZZ_TEST_RUN_TIME

	fUZZ_TEST_TIMEOUT, err := clicommon.PromptInteger("Fuzz Timeout in Seconds current: " + string(currCfg.Fuzzsettings.FUZZ_TEST_TIMEOUT))
	if err != nil {
		return nil, err
	}
	currCfg.Fuzzsettings.FUZZ_TEST_TIMEOUT = fUZZ_TEST_TIMEOUT
	fUZZ_NUM_PROCESSES, err := clicommon.PromptInteger("Fuzz Number of parralell processes current: " + string(currCfg.Fuzzsettings.FUZZ_NUM_PROCESSES))
	if err != nil {
		return nil, err
	}
	currCfg.Fuzzsettings.FUZZ_NUM_PROCESSES = fUZZ_NUM_PROCESSES

	return currCfg, nil
}

func promptAddGolangsettings(currCfg *cfg.RWConfiguration) (*cfg.RWConfiguration, error) {
	if currCfg == nil {
		currCfg, err := cfg.ConfigRead("./redwoods-cfg.json")
		fmt.Println("loaded config", currCfg.Name)
		if err != nil {
			return nil, err
		}
	}

	gO111MODULE, err := clicommon.PromptConfirm("Use GO111MODULES? current: " + currCfg.Gosettings.GO111MODULE)
	if err != nil && gO111MODULE == "" {
		fmt.Println("Caught: ", err)
	}
	if strings.ToLower(gO111MODULE) == "y" {
		currCfg.Gosettings.GO111MODULE = "on"
	} else {
		currCfg.Gosettings.GO111MODULE = "off"
	}

	gOPROXY, err := clicommon.PromptString("Use GOPROXY? current: ", currCfg.Gosettings.GOPROXY)
	if err != nil {
		return nil, err
	}
	currCfg.Gosettings.GOPROXY = gOPROXY

	gONOSUMDB, err := clicommon.PromptString("Use GONOSUMDB? current: ", currCfg.Gosettings.GONOSUMDB)
	if err != nil {
		return nil, err
	}
	currCfg.Gosettings.GONOSUMDB = gONOSUMDB

	gOOS, err := clicommon.PromptGOOS("Which GOOS? current: " + currCfg.Gosettings.GOOS)
	if err != nil {
		return nil, err
	}
	currCfg.Gosettings.GOOS = gOOS

	gOARCH, err := clicommon.PromptGOARCH("Which GOARCH? current: " + currCfg.Gosettings.GOARCH)
	if err != nil {
		return nil, err
	}
	currCfg.Gosettings.GOARCH = gOARCH

	gO_BIN, err := clicommon.PromptGOOS("Would you like to set a Go bininary path? (blank to get from environment) current: " + currCfg.Gosettings.GO_BIN)
	if err != nil {
		return nil, err
	}
	if gO_BIN == "" {
		currCfg.Gosettings.GO_BIN = os.Getenv("GOPATH") + "/bin"
	} else if gO_BIN != "" {
		currCfg.Gosettings.GO_BIN = gO_BIN
	}

	return currCfg, nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
