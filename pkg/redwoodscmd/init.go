package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	cfg "github.com/xxyyx/redwoods/pkg/rwconfig"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates Workfolders and downloads the project form git",
	Long:  `This will create the outputt,work and fuzz directories as well as to clone the project into them. If a commit is given, that will be checked out.`,
	Run: func(cmd *cobra.Command, args []string) {

		initRedwoods()
	},
}

func initRedwoods() {
	currCfg, err := cfg.ConfigRead("./redwoods-cfg.json")
	if err != nil {
		fmt.Println("Could not find a redwoods-cfg.json please run $ redwoods -c")
		return
	}
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("## Creating workspaces ###")
	//create the fuzz folder
	os.MkdirAll(pwd+"/fuzz", os.ModePerm)
	//create the workdir folder
	os.MkdirAll(pwd+"/workspace", os.ModePerm)

	//create the analysis folder
	os.MkdirAll(pwd+"/analysis", os.ModePerm)
	fmt.Println("## Cloning Repository ###")
	runGitclone(pwd, *currCfg)
	if currCfg.Sourcesettings.Commit != "latest" && currCfg.Sourcesettings.Commit != "" {
		runGitCheckout(pwd, *currCfg)
	}
}

func runGitclone(pwd string, currCfg cfg.RWConfiguration) {
	fmt.Println("Cloning ", currCfg.Repository)
	if currCfg.Sourcesettings.Branch != "main" && currCfg.Sourcesettings.Branch != "" {
		e := exec.Command("git", "clone", currCfg.Repository, "-b"+currCfg.Sourcesettings.Branch, pwd+"/workspace/"+currCfg.Name)
		e.Stdout = os.Stdout
		e.Stderr = os.Stderr
		_ = e.Run()
	} else {
		e := exec.Command("git", "clone", currCfg.Repository, pwd+"/workspace/"+currCfg.Name)
		e.Stdout = os.Stdout
		e.Stderr = os.Stderr
		_ = e.Run()
	}
	//e.Path = pwd + "/workspace"

	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func runGitCheckout(pwd string, currCfg cfg.RWConfiguration) {
	fmt.Println("Checking out Commit ", currCfg.Sourcesettings.Commit)
	//git checkout -b $(SOURCE_BRANCH) $(SOURCE_COMMIT)
	e := exec.Command("git", "checkout", currCfg.Sourcesettings.Commit)
	e.Dir = pwd + "/workspace/" + currCfg.Name
	e.Stdout = os.Stdout
	e.Stderr = os.Stderr
	_ = e.Run()
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
