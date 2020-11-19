package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	repo   string
	branch string
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
	Use:    "checkMigrations",
	PreRun: initGithubParams,
	Short:  "Check that the migrations in this repository are properly timestamped for a merge.",
	Long: `Check that any new migrations in this branch will run after all the migrations in the
target branch, prod by default`,
	Run: func(cmd *cobra.Command, args []string) {
		err := checkMigrations()

		if err != nil {
			fmt.Printf("error checking migrations %s\n", err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkMigrationsCmd)
	checkMigrationsCmd.Flags().String("repo", "makerdao/vdb-mcd-transformers", "Github Repository to check against, must be public")
	checkMigrationsCmd.Flags().String("branch", "prod", "Branch to check against")

	viper.BindPFlag("repo", checkMigrationsCmd.Flags().Lookup("repo"))
	viper.BindPFlag("branch", checkMigrationsCmd.Flags().Lookup("branch"))
}

func initGithubParams(cmd *cobra.Command, args []string) {
	repo = viper.GetString("repo")
	branch = viper.GetString("branch")
}

func checkMigrations() error {
	remoteMigrations, errGithub := getGithubFileNames()

	if errGithub != nil {
		return errGithub
	}

	localMigrations, errLocal := getLocalMigrations()

	if errLocal != nil {
		return errLocal
	}

	newMigrations := NewMigrations(localMigrations, remoteMigrations)
	logrus.Println("New Migrations are", newMigrations)

	return CheckNewMigrations(remoteMigrations, newMigrations)
}

func getGithubFileNames() ([]string, error) {
	url := fmt.Sprintf("http://api.github.com/repos/%s/contents/db/migrations?ref=%s", repo, branch)
	errorContext := "error getting migrations from Github %w"
	logrus.Println("Retrieving Migration list from", url)

	resp, httpErr := http.Get(url)

	if httpErr != nil {
		return []string{}, fmt.Errorf(errorContext, httpErr)
	}
	defer resp.Body.Close()

	body, ioErr := ioutil.ReadAll(resp.Body)

	if ioErr != nil {
		return []string{}, fmt.Errorf(errorContext, ioErr)
	}

	var files []GithubFile
	jsonErr := json.Unmarshal(body, &files)

	if jsonErr != nil {
		return []string{}, fmt.Errorf("failed to unmarshal %s, %w", string(body), jsonErr)
	}

	var namedFiles = make([]NamedFile, len(files))

	for i, file := range files {
		namedFiles[i] = file
	}

	return GetSQLFilesFromList(namedFiles), nil
}

func getLocalMigrations() ([]string, error) {
	localFiles, err := ioutil.ReadDir("./db/migrations")

	if err != nil {
		return []string{}, fmt.Errorf("error reading local migrations %w", err)
	}

	namedFiles := make([]NamedFile, len(localFiles))

	for i, file := range localFiles {
		namedFiles[i] = file.(NamedFile)
	}

	return GetSQLFilesFromList(namedFiles), nil
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

// CheckNewMigrations makes sure the new migrations are sorted correctly
func CheckNewMigrations(originalMigrations []string, newMigrations []string) error {
	err := checkAllNewMigrationsAreTimestamped(newMigrations)
	if err != nil {
		return err
	}

	return checkNewMigrationsAreAfterCurrentOnes(originalMigrations, newMigrations)
}

func checkAllNewMigrationsAreTimestamped(newMigrations []string) error {
	for _, newMigration := range newMigrations {
		matched, err := regexp.MatchString(`\d{14}_`, newMigration)

		if err != nil {
			return fmt.Errorf("error checking migration has a timestamp %w", err)
		}

		if !matched {
			return fmt.Errorf("migration %s does not start with a timestamp", newMigration)
		}
	}
	return nil
}

func checkNewMigrationsAreAfterCurrentOnes(originalMigrations []string, newMigrations []string) error {
	sortedNewMigrations := make([]string, len(newMigrations))
	copy(sortedNewMigrations, newMigrations)
	sort.Strings(sortedNewMigrations)

	finalMigrations := make([]string, len(originalMigrations))
	copy(finalMigrations, originalMigrations)
	for _, newMigration := range sortedNewMigrations {
		finalMigrations = append(finalMigrations, newMigration)
	}
	sort.Strings(finalMigrations)

	lastRunMigrations := finalMigrations[len(finalMigrations)-len(sortedNewMigrations) : len(finalMigrations)]

	for idx, value := range lastRunMigrations {
		if sortedNewMigrations[idx] != value {
			return fmt.Errorf("migration %s is out of order, update your timestamp", sortedNewMigrations[idx])
		}
	}

	return nil
}
