package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Return the version of cpsp",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v0.1")

	},
}
