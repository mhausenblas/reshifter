package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays the ReShifter version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", releaseVersion)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
