package cmd

import (
  "fmt"
  "github.com/spf13/cobra"
  "os"
)


var cfgFile string


var rootCmd = &cobra.Command{
  Use:   "BML",
  Short: "BML is an application that manages your terminals working directories",
  Long: `BML is an application that manages your terminals working directories`,
  // Uncomment the following line if your bare application
  // has an action associated with it:
  //	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}
