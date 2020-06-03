package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// NamedFile interface represents a local file object, with a Name function
type NamedFile interface {
	Name() string
}

type Set = map[string]bool

type githubJSON struct {
	Name string
}

type GithubFile struct {
	githubJSON
}

func (file *GithubFile) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &file.githubJSON)
}

func (file GithubFile) Name() string {
	return file.githubJSON.Name
}

var checkMigrationsCmd = &cobra.Command{
	Use:   "checkMigrations",
	Short: "Check that the migrations in this repository are properly timestamped for a merge.",
	Long: `Check that any new migrations in this branch will run after all the migrations in the
staging branch`,
	Run: func(cmd *cobra.Command, args []string) {
		err := checkMigrations()

		if err != nil {
			fmt.Printf("Error checking migrations %s\n", err.Error())
			return
		}

		return
	},
}

func init() {
	rootCmd.AddCommand(checkMigrationsCmd)
}

func checkMigrations() error {
	// Connect to Github (URL should be a param) (so should branch)
	resp, err := http.Get("http://api.github.com/repos/makerdao/vdb-mcd-transformers/contents/db/migrations?ref=staging")

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	// Convert JSON to Struct
	var files []GithubFile
	err = json.Unmarshal(body, &files)

	if err != nil {
		return err
	}

	stagingMigrations := getGithubFileNames(files)

	localFiles, err := ioutil.ReadDir("./db/migrations")

	if err != nil {
		return err
	}

	localMigrations := getLocalFileNamesFrom(localFiles)
	newMigrations := NewMigrations(localMigrations, stagingMigrations)

	// Retrieve what the latest migrations will be
	for _, newMigration := range newMigrations {
		stagingMigrations = append(stagingMigrations, newMigration)
	}

	sort.Strings(stagingMigrations)
	sort.Strings(newMigrations)

	lastMigrations := stagingMigrations[len(stagingMigrations)-len(newMigrations) : len(stagingMigrations)]

	// Check - and return an error!
	for idx, value := range lastMigrations {
		if newMigrations[idx] != value {
			fmt.Printf("New Migration %s is out of order. Update your timestamp! \n", newMigrations[idx])
		}
	}

	return nil
}

// GetGithubFileNames returns the list of file names from a list of
// GithubFile objects
func getGithubFileNames(files []GithubFile) []string {
	var namedFiles = make([]NamedFile, len(files))

	for i, file := range files {
		namedFiles[i] = file
	}

	return GetSQLFilesFromList(namedFiles)
}

// GetLocalFileNamesFrom returns the list of file names from a list of
// local file objects
func getLocalFileNamesFrom(files []os.FileInfo) []string {
	var namedFiles = make([]NamedFile, len(files))

	for i, file := range files {
		namedFiles[i] = file.(NamedFile)
	}

	return GetSQLFilesFromList(namedFiles)
}

// GetSQLFilesFromList returns the list of file names
// from a list of, local file objects but only SQL files
func GetSQLFilesFromList(files []NamedFile) []string {
	var fileNames []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			fileNames = append(fileNames, file.Name())
		}
	}

	return fileNames
}

// NewMigrations gets the list of brand new migrations
func NewMigrations(newList []string, oldList []string) []string {
	oldSet := toSet(oldList)
	newSet := toSet(newList)
	return diff(newSet, oldSet)
}

func toSet(list []string) Set {
	set := make(Set)
	for _, entry := range list {
		set[entry] = true
	}
	return set
}

func diff(a Set, b Set) []string {
	diff := []string{}

	for value := range a {
		if b[value] == false {
			diff = append(diff, value)
		}
	}

	return diff
}
