package cmd

import (
  "fmt"
  "github.com/spf13/cobra"
  "log"
  "os"
    "os/user"
    "path/filepath"
)


var cfgFile string


var rootCmd = &cobra.Command{
  Use:   "BML",
  Short: "BML is an application that manages your terminals working directories",
  Long: `BML is an application that manages your terminals working directories`,
  // has an action associated with it:
  	Run: func(cmd *cobra.Command, args []string) {
        // Get the directory where the binary was called from
  	    currentBinLocation, err := filepath.Abs(filepath.Dir(os.Args[0]))
        if err != nil {
            log.Fatal(err)
        }

        // Get the directory of where the binary should be installed for this user
        usr, err := user.Current()
        if err != nil {
            log.Fatal( err )
        }
        installationBinLocation := fmt.Sprintf("%s/.BML_Installation", usr.HomeDir)

        // Check if binary is in the correct location
        if currentBinLocation != installationBinLocation {
            fmt.Printf(`BML binary cannot be called as it is not in the correct install location
BML binary curently at:  %s
BML binary should be at: %s
Have you executed "BML install" command?`+"\n", currentBinLocation, installationBinLocation)
            os.Exit(1)
        }

        // Main Functionality of BML
    },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}
