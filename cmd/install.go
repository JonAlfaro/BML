package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "This command installs and BML and configures your bashrc",
	Long: `This command installs and BML and configures your bashrc by adding 
a wrapping command that controls the actual change in directory in your terminal session`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("install called")
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
