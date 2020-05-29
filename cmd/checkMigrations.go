package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

type Set = map[string]bool

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

type GithubFile struct {
	Name string
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

	// Get the list of migration file names from staging
	var stagingMigrations []string
	for _, file := range files {
		stagingMigrations = append(stagingMigrations, file.Name)
	}

	// Get the list of files from local
	localFiles, err := ioutil.ReadDir("./db/migrations")

	if err != nil {
		return err
	}

	// Reduce to file names, only matchin .sql
	var localMigrations []string
	for _, localFile := range localFiles {
		if strings.HasSuffix(localFile.Name(), ".sql") {
			localMigrations = append(localMigrations, localFile.Name())
		}
	}

	var stagingMigrationSet = toSet(stagingMigrations)
	var localMigrationSet = toSet(localMigrations)
	var newMigrations = diff(localMigrationSet, stagingMigrationSet)

	for _, newMigration := range newMigrations {
		stagingMigrations = append(stagingMigrations, newMigration)
	}

	sort.Strings(stagingMigrations)
	sort.Strings(newMigrations)

	lastMigrations := stagingMigrations[len(stagingMigrations)-len(newMigrations) : len(stagingMigrations)]

	for idx, value := range lastMigrations {
		if newMigrations[idx] != value {
			fmt.Printf("New Migration %s is out of order. Update your timestamp! \n", newMigrations[idx])
		}
	}

	return nil
}

func toSet(list []string) Set {
	var set = make(Set)
	for _, entry := range list {
		set[entry] = true
	}
	return set
}

// Run the set difference operation (A - B)
func diff(a Set, b Set) []string {
	var diff []string

	for value := range a {
		if b[value] == false {
			diff = append(diff, value)
		}
	}

	return diff
}
