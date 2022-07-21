package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	cfg "github.com/xxyyx/redwoods/pkg/rwconfig"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(dependenciesCmd)
}

var dependenciesCmd = &cobra.Command{
	Use:   "dependencies",
	Short: "Download all dependencies needed for this project",
	Long:  `This will run the installation packages for all go related tools`,
	Run: func(cmd *cobra.Command, args []string) {
		installDepdenencies()
	},
}

func installDepdenencies() {
	currCfg, err := cfg.ConfigRead("./redwoods-cfg.json")
	if err != nil {
		fmt.Println("Could not find a redwoods-cfg.json please run $ redwoods -c")
		return
	}
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("## Installing dependencies ###")
	if currCfg.Dependencysettings.GOFUZZ_BUILD_PKG_VERSION != "latest" && currCfg.Dependencysettings.GOFUZZ_BUILD_PKG_VERSION != "" {
		depInstall(pwd, *currCfg, currCfg.Dependencysettings.GOFUZZ_BUILD_PKG+"@"+currCfg.Dependencysettings.GOFUZZ_BUILD_PKG_VERSION)
	} else {
		depInstall(pwd, *currCfg, currCfg.Dependencysettings.GOFUZZ_BUILD_PKG+"@latest")
	}
	if currCfg.Dependencysettings.GOFUZZ_COMMIT != "latest" && currCfg.Dependencysettings.GOFUZZ_COMMIT != "" {
		depInstall(pwd, *currCfg, currCfg.Dependencysettings.GOFUZZ_PKG+"@"+currCfg.Dependencysettings.GOFUZZ_COMMIT)
	} else {
		depInstall(pwd, *currCfg, currCfg.Dependencysettings.GOFUZZ_PKG+"@latest")
	}
	if currCfg.Dependencysettings.GOCYCLO_PKG_VERSION != "latest" && currCfg.Dependencysettings.GOCYCLO_PKG_VERSION != "" {
		depInstall(pwd, *currCfg, currCfg.Dependencysettings.GOCYCLO_PKG+"@"+currCfg.Dependencysettings.GOCYCLO_PKG_VERSION)
	} else {
		depInstall(pwd, *currCfg, currCfg.Dependencysettings.GOCYCLO_PKG+"@latest")
	}
	if currCfg.Dependencysettings.GOFUZZ_DEP_PKG_VERSION != "latest" && currCfg.Dependencysettings.GOFUZZ_DEP_PKG_VERSION != "" {
		depGetWorkspace(pwd, *currCfg, currCfg.Dependencysettings.GOFUZZ_DEP_PKG+"@"+currCfg.Dependencysettings.GOFUZZ_DEP_PKG_VERSION)
	} else {
		depGetWorkspace(pwd, *currCfg, currCfg.Dependencysettings.GOFUZZ_DEP_PKG+"@latest")
	}

}

func depGet(pwd string, currCfg cfg.RWConfiguration, pkg string) {
	fmt.Println("Installing package: ", pkg)
	//git checkout -b $(SOURCE_BRANCH) $(SOURCE_COMMIT)
	e := exec.Command("go", "get", pkg)
	e.Env = os.Environ()
	e.Env = append(e.Env, "GO111MODULE="+currCfg.Gosettings.GO111MODULE)
	e.Env = append(e.Env, "GOPROXY="+currCfg.Gosettings.GOPROXY)
	e.Env = append(e.Env, "GONOSUMDB="+currCfg.Gosettings.GONOSUMDB)
	e.Stdout = os.Stdout
	e.Stderr = os.Stderr
	_ = e.Run()

}

func depInstall(pwd string, currCfg cfg.RWConfiguration, pkg string) {
	fmt.Println("Installing package: ", pkg)
	//git checkout -b $(SOURCE_BRANCH) $(SOURCE_COMMIT)
	e := exec.Command("go", "install", pkg)
	e.Env = os.Environ()
	e.Env = append(e.Env, "GO111MODULE="+currCfg.Gosettings.GO111MODULE)
	e.Env = append(e.Env, "GOPROXY="+currCfg.Gosettings.GOPROXY)
	e.Env = append(e.Env, "GONOSUMDB="+currCfg.Gosettings.GONOSUMDB)
	e.Stdout = os.Stdout
	e.Stderr = os.Stderr
	_ = e.Run()

}

func depGetWorkspace(pwd string, currCfg cfg.RWConfiguration, pkg string) {
	fmt.Println("Installing package in workspace: ", pkg, "in", pwd+"/workspace/"+currCfg.Name)
	e := exec.Command("go", "get", pkg)
	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)
	e.Env = os.Environ()
	e.Env = append(e.Env, "GO111MODULE="+currCfg.Gosettings.GO111MODULE)
	e.Env = append(e.Env, "GOPROXY="+currCfg.Gosettings.GOPROXY)
	e.Env = append(e.Env, "GONOSUMDB="+currCfg.Gosettings.GONOSUMDB)

	e.Dir = pwd + "/workspace/" + currCfg.Name + "/"
	e.Stdout = mw
	e.Stderr = mw
	_ = e.Run()

}
