package cmd

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/user"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "This command will uninstall all traces of BML on your computer",
	Long: `This command will uninstall your BML installation and remove the wrapped BML command
from your bashrc`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Uninstalling BML from your system")
		// Get the directory of where the binary should be installed for this user
		usr, err := user.Current()
		if err != nil {
			fmt.Printf("FATAL ERROR Could find current user\nERROR: %s\n", err.Error())
			return
		}
		installationBinLocation := fmt.Sprintf("%s/.BML_Installation", usr.HomeDir)
		fmt.Println(installationBinLocation)

		// Remove bml wrapper command
		fmt.Printf("┳ Removing bml command from %s/.bashrc\n", usr.HomeDir)
		bmlCommandString := fmt.Sprintf(`
function bml() {
  if [ "$1" = "new" ];
  then
    %s/BML --new
  elif [ "$1" = "remove" ];
  then
    %s/BML --remove
  elif [ "$1" = "uninstall" ];
  then
    %s/BML uninstall
  elif [ "$1" = "help" ];
  then
    echo "Supported Commands:".
	echo "new: creates new bookmark entry at your current working directory"
	echo "remove: removes a bookmark from list"
	echo "uninstall: uninstalls all traces of BML"
  else
	if %s/BML ; then
		cd "$(cat /tmp/.SuperImportantTargetForBookmarks.clown)"
	fi
  fi
}`, installationBinLocation, installationBinLocation, installationBinLocation, installationBinLocation)

		bashrcInfo, err := os.Stat(usr.HomeDir + "/.bashrc")
		if err != nil {
			fmt.Printf("┗ FATAL ERROR Could not Stat %s/.bashrc\n┗ Error: %v\n", usr.HomeDir, err)
			return
		}

		bashrcContent, err := ioutil.ReadFile(usr.HomeDir + "/.bashrc")
		if err != nil {
			fmt.Printf("┗ FATAL ERROR Could not Read %s/.bashrc\n┗ Error: %v\n", usr.HomeDir, err)
			return
		}

		bashrcWithoutBML := bytes.Replace(bashrcContent, []byte(bmlCommandString), []byte(""), -1)

		if err = ioutil.WriteFile(usr.HomeDir+"/.bashrc", bashrcWithoutBML, bashrcInfo.Mode().Perm()); err != nil {
			fmt.Printf("┗ FATAL ERROR Could not Write to %s/.bashrc\n┗ Error: %v\n", usr.HomeDir, err)
			return
		} else {
			fmt.Printf("┗ bml command removed from %s/.bashrc successfully!\n", usr.HomeDir)
		}

		// Remove Install Dir
		fmt.Printf("┳ Removing Installation Directory at: %s\n", installationBinLocation)
		if err = os.RemoveAll(installationBinLocation); err != nil {
			fmt.Printf("┗ FATAL ERROR Could not remove Installation Directory: %s\n┗ Error: %v\n", installationBinLocation, err)
			return
		} else {
			fmt.Printf("┗ Installation Directory: %s removed successfully!\n", installationBinLocation)
		}

		fmt.Println("BML Uninstall Successful!")
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
