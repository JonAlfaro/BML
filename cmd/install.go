package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "This command installs and BML and configures your bashrc",
	Long: `This command installs and BML and configures your bashrc by adding 
a wrapping command that controls the actual change in directory in your terminal session`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Installing BML onto your system")
		// Get the directory where the binary was called from
		currentBinLocation, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			fmt.Printf("Could get current BML binary location\nERROR: %s\n", err.Error())
			return
		}

		// Get the directory of where the binary should be installed for this user
		usr, err := user.Current()
		if err != nil {
			fmt.Printf("Could find current user\nERROR: %s\n", err.Error())
			return
		}
		installationBinLocation := fmt.Sprintf("%s/.BML_Installation", usr.HomeDir)

		// Check if binary is in the correct location
		if currentBinLocation == installationBinLocation {
			fmt.Printf(`BML binary is already installed
BML binary curently at:  %s
BML binary should be at: %s
If you want to uninstall execute "%s/BML uninstall"`+"\n", currentBinLocation, installationBinLocation, installationBinLocation)
			os.Exit(1)
		}

		// Create Install Directory
		fmt.Printf("┳ Creating Directory: %s\n", installationBinLocation)
		if _, err := os.Stat("installationBinLocation"); !os.IsNotExist(err) {
			fmt.Println("┗ Path already Exists, skipping this step")
		} else {
			errDir := os.MkdirAll(installationBinLocation, 0755)
			if errDir != nil {
				fmt.Println("┗ Error: " + err.Error())
				return
			} else {
				fmt.Println("┗ Directory Creation Successful")
			}
		}

		// Move Binary
		fmt.Printf("┳ Moving binary to %s\n", installationBinLocation)
		if _, err := os.Stat(installationBinLocation + "/BML"); !os.IsNotExist(err) {
			fmt.Println("┗ Warning Binary Already Exists, replacing binary file")
		}

		if _, err := os.Stat(currentBinLocation + "/BML"); os.IsNotExist(err) {
			fmt.Printf("┗ Error BML Binary does not exists at: %s/BML\n", currentBinLocation)
			return
		} else {
			if _, err := copyFile(currentBinLocation+"/BML", installationBinLocation+"/BML"); err != nil {
				fmt.Println("┗ Error: " + err.Error())
				return
			} else {
				fmt.Println("┗ BML Binary Copy Successful")
			}
		}

		// Create bookmark json
		fmt.Printf("┳ Creating %s/BML_bookmarks.json\n", installationBinLocation)
		if _, err := os.Stat(installationBinLocation + "/BML_bookmarks.json"); !os.IsNotExist(err) {
			fmt.Printf("┗ Warning JSON bookmarkfile already exists, overwriting old json file \n")
		}
		emptyJSONByte := []byte("{\n}")
		if err := ioutil.WriteFile(installationBinLocation+"/BML_bookmarks.json", emptyJSONByte, 0644); err != nil {
			fmt.Println("┗ Error: " + err.Error())
			return
		} else {
			fmt.Println("┗ Bookmark JSON Creation Successful")
		}

		// Backup bashrc
		fmt.Printf("┳ Backing up %s/.bashrc to %s/.bashrc-backup\n", usr.HomeDir, installationBinLocation)
		if _, err := os.Stat(installationBinLocation + "/.bashrc-backup"); !os.IsNotExist(err) {
			fmt.Printf("┗ Warning bashrc-backup already exists, overwriting bashrc-backup \n")
		}
		if _, err := copyFile(usr.HomeDir+"/.bashrc", installationBinLocation+"/.bashrc-backup"); err != nil {
			fmt.Println("┗ Error: " + err.Error())
			return
		} else {
			fmt.Println("┗ .bashrc backup Successful")
		}

		// Edit baschrc
		fmt.Printf("┳ Adding \"bml\" wrapper command to %s/.bashrc\n", usr.HomeDir)
		f, err := os.OpenFile(usr.HomeDir+"/.bashrc", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer f.Close()
		if err != nil {
			fmt.Println("┗ Error - Could not open " + usr.HomeDir + "/.bashrc : " + err.Error())
			return
		}

		defer f.Close()
		_, err = f.WriteString(fmt.Sprintf(`

function bml() {
  if [ "$1" = "new" ];
  then
    %s/BML --new
  else
    %s/BML
    cd "$(cat /tmp/.SuperImportantTargetForBookmarks.clown)"
  fi
}`, installationBinLocation, installationBinLocation))
		if err != nil {
			fmt.Println("┗ Error - Could write bml command to " + usr.HomeDir + "/.bashrc : " + err.Error())
			return
		} else {
			fmt.Println("┗ bml Command Appended to .bashrc Successful")
		}

		// Make binary executable: -rwxrwxr-x
		err = os.Chmod(installationBinLocation+"/BML", 0775)
		if err != nil {
			fmt.Println("-- Installation Successful with warning--")
			fmt.Printf(`Warning: Installer could make the BML binary executable, please run the command below to fix this:
chmod +x %s/BML`, installationBinLocation)
		} else {
			fmt.Println("-- Installation Successful --")
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

// Simple Copy command
func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
