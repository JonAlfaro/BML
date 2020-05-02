package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "This command will uninstall all traces of BML on your computer",
	Long: `This command will uninstall your BML installation and remove the wrapped BML command
from your bashrc`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("uninstall called")
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
