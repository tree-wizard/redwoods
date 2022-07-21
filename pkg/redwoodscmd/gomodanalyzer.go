package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	cfg "github.com/xxyyx/redwoods/pkg/rwconfig"
	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"
)

func init() {
	rootCmd.AddCommand(gomodanalyzerCmd)
}

var gomodanalyzerCmd = &cobra.Command{
	Use:   "gomodanalyzer",
	Short: "Analyse the go.mod file for all packages",
	Long:  `Gives you and overview of all modules and replacements used in the project as well as their versions`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("## Analyzing go.mod ###")
		pwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		currCfg, err := cfg.ConfigRead(pwd + "/redwoods-cfg.json")
		if err != nil {
			fmt.Println("Could not finde a redwoods-cfg.json please rund redwoods -c")
			return
		}
		outputPath := pwd + "/analysis/packages-" + currCfg.Name + ".json"
		analysis := GoModAnalysis{
			Project:      currCfg.Name,
			Path:         pwd + "/workspace/" + currCfg.Name + "/go.mod",
			Modules:      make(map[string]string),
			Replacements: make(map[string]string),
		}
		file, err := os.Open(analysis.Path)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		fileinfo, err := file.Stat()
		if err != nil {
			fmt.Println(err)
			return
		}

		filesize := fileinfo.Size()
		buffer := make([]byte, filesize)

		_, err = file.Read(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}
		f, err := modfile.Parse("go.mod", buffer, nil)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("required packages")
		for _, c := range f.Require {
			analysis.Modules[c.Mod.Path] = c.Mod.Version
			fmt.Println(c.Mod.Path, " ", c.Mod.Version)
		}
		fmt.Println("replaced packages")
		for _, c := range f.Replace {
			analysis.Replacements[c.Old.Path] = c.New.Path
			fmt.Println(c.Old.Path, " ", c.New.Path)
		}

		json, err := json.MarshalIndent(analysis, "", " ")
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(outputPath, json, 0644)
		if err != nil {
			panic(err)
		}
	},
}

type GoModAnalysis struct {
	Project      string            `json:"project"`
	Path         string            `json:"path"`
	Modules      map[string]string `json:"modules"`
	Replacements map[string]string `json:"replacements"`
}
