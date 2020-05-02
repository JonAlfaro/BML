package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"github.com/manifoldco/promptui"
)

var bookmarkList = make([]Bookmark, 0)
var newFlag = false

const (
	BOOKMARK_JSON = "/BML_bookmarks.json"
)

type Bookmark struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

var rootCmd = &cobra.Command{
	Use:   "BML",
	Short: "BML is an application that manages your terminals working directories",
	Long:  `BML is an application that manages your terminals working directories`,
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
			log.Fatal(err)
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

		// Load in bookmarks
		bookmarkFile, err := os.Open(installationBinLocation + BOOKMARK_JSON)
		if err != nil {
			fmt.Printf("Could not open bookmark json file: %s \nError: %v\n", installationBinLocation+BOOKMARK_JSON, err)
			return
		}

		// defer the closing of our jsonFile so that we can parse it later on
		defer bookmarkFile.Close()

		byteValue, err := ioutil.ReadAll(bookmarkFile)
		if err != nil {
			fmt.Printf("Could not read bookmark json file: %s \nError: %v\n", installationBinLocation+BOOKMARK_JSON, err)
			return
		}

		bookmarkFile.Close()

		json.Unmarshal(byteValue, &bookmarkList)

		if newFlag {
			validateName := func(input string) error {
				input = strings.TrimSpace(input)
				if input == "" {
					return errors.New("Bookmark cannot be empty")
				}

				for _, bookmark := range bookmarkList {
					if input == bookmark.Name {
						return errors.New("Bookmark name already exists")
					}
				}
				return nil
			}

			promptName := promptui.Prompt{
				Label:    "Bookmark Name",
				Validate: validateName,
			}

			newBookmarkName, err := promptName.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			newBookmarkPath, err := os.Getwd()
			if err != nil {
				fmt.Printf("Could not get current working directory %v\n", err)
				return
			}

			bookmarkList = append(bookmarkList, Bookmark{
				Name: newBookmarkName,
				Path: newBookmarkPath,
			})

			// Backup old bookmark file
			err = ioutil.WriteFile(installationBinLocation+BOOKMARK_JSON+"-old", byteValue, 0644)
			if err != nil {
				panic(err)
			}

			// TODO: Append to file, instead of overwriting file
			bookmarkFile, err = os.Create(installationBinLocation + BOOKMARK_JSON)
			if err != nil {
				// I dont know how this could happen unless there is a permission problem, or someone is tampering
				// with bookmark json
				fmt.Printf("Could not create new bookmark file   %v\n", err)
				return
			}
			defer bookmarkFile.Close()

			bookmarkFileListData, _ := json.MarshalIndent(bookmarkList, "", " ")
			if _, err := bookmarkFile.Write(bookmarkFileListData); err != nil {
				fmt.Printf("Could not to bookmark file  %v\n", err)
				return
			}

			fmt.Printf("Successfully Created new bookmark: %s PATH: %s\n", newBookmarkName, newBookmarkPath)

		} else {
			searcher := func(input string, index int) bool {
				bookmark := bookmarkList[index]
				name := strings.Replace(strings.ToLower(bookmark.Name), " ", "", -1)
				input = strings.Replace(strings.ToLower(input), " ", "", -1)

				return strings.Contains(name, input)
			}

			template := &promptui.SelectTemplates{
				Active:   "{{.Name | cyan}}",
				Inactive: "{{.Name}}",
				Selected: "{{.Name}}",
				Details: `
--------- Details ----------
{{ "Name:" | faint }}	{{ .Name }}
{{ "Path:" | faint }}	{{ .Path }}`,
			}

			prompt := promptui.Select{
				Label:             "Bookmark List",
				Items:             bookmarkList,
				Size:              10,
				Searcher:          searcher,
				StartInSearchMode: true,
				Templates:         template,
			}

			i, _, err := prompt.Run()

			if err != nil {
				fmt.Printf("Failed in prompt selection %v\n", err)
				return
			}

			fmt.Printf("Changing Working Directory to [%s] PATH: %s\n", bookmarkList[i].Name, bookmarkList[i].Path)
			if fi, err := os.Stat(bookmarkList[i].Path); os.IsNotExist(err) {
				// Path does not exist
				fmt.Printf("Failed changing working directory: Path does not exist\n")
				return
			} else if !fi.IsDir() {
				// Path is not a Dir
				fmt.Printf("Failed changing working directory: Path is not a directory\n")
				return
			}

			// Write to file, this is just a file where target will be written to.
			// It is done this way to the wrapper bash command can cd cat the file
			file, err := os.Create("/tmp/.SuperImportantTargetForBookmarks.clown")
			if err != nil {
				fmt.Printf("Failed changing working directory: Could not create target\n")
				return
			}
			defer file.Close()

			if _, err := file.WriteString(bookmarkList[i].Path); err != nil {
				fmt.Printf("Failed changing working directory: Could not write target\n")
				return
			}
		}
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

func init() {
	rootCmd.Flags().BoolVar(&newFlag, "new", false, "create new bookmark")
}
