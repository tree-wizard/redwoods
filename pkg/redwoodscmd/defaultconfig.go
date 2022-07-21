package cmd

import (
	"fmt"

	cfg "github.com/xxyyx/redwoods/pkg/rwconfig"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(defaultconfigCmd)
}

var defaultconfigCmd = &cobra.Command{
	Use:   "defaultconfig",
	Short: "Creates the default config so it can be modified",
	Long:  `This will create the the default config in the current directory. All that its left is to configure it.`,
	Run: func(cmd *cobra.Command, args []string) {

		initdefaultconfig()
	},
}

func initdefaultconfig() {
	fmt.Println("## Creating a default Config ###")
	rwcfg := cfg.DefaultConfig()
	err := cfg.ConfigWrite(rwcfg, "./redwoods-cfg.json")
	if err != nil {
		panic(err)
	}
	fmt.Println("default config written to ./redwoods-cfg.json")
}
