package cmd_test

import (
	"github.com/makerdao/vdb-mcd-transformers/cmd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type MockFile struct {
	name string
}

func (f MockFile) Name() string {
	return f.name
}

var _ = Describe("Check Migrations", func() {

	Describe("GetSQLFilesfromList", func() {
		It("maps a list of named files to their actual file names", func() {
			files := []MockFile{{name: "one.sql"}, {name: "two.sql"}}
			var namedFiles = make([]cmd.NamedFile, len(files))
			for i, file := range files {
				namedFiles[i] = file
			}

			fileNames := cmd.GetSQLFilesFromList(namedFiles)

			Expect(fileNames).To(Equal([]string{"one.sql", "two.sql"}))
		})

		It("skips anything that isn't a sql file", func() {
			files := []MockFile{{name: "one.sql"}, {name: "two.txt"}}
			var namedFiles = make([]cmd.NamedFile, len(files))
			for i, file := range files {
				namedFiles[i] = file
			}

			fileNames := cmd.GetSQLFilesFromList(namedFiles)

			Expect(fileNames).To(Equal([]string{"one.sql"}))
		})
	})

	Describe("NewMigrations", func() {
		It("returns no new migrations when the new migrations list is empty", func() {
			newMigrationList := []string{}
			oldMigrationList := []string{}

			Expect(cmd.NewMigrations(newMigrationList, oldMigrationList)).To(Equal([]string{}))
		})

		It("returns no new migrations when the migration lists are identical", func() {
			newMigrationList := []string{"one.sql"}
			oldMigrationList := []string{"one.sql"}

			Expect(cmd.NewMigrations(newMigrationList, oldMigrationList)).To(Equal([]string{}))
		})

		It("returns any migrations that are in the new, but not the old", func() {
			newMigrationList := []string{"one.sql"}
			oldMigrationList := []string{}

			Expect(cmd.NewMigrations(newMigrationList, oldMigrationList)).To(Equal([]string{"one.sql"}))
		})

		It("ignores oldMigrations", func() {
			newMigrationList := []string{"one.sql"}
			oldMigrationList := []string{"two.sql"}

			Expect(cmd.NewMigrations(newMigrationList, oldMigrationList)).To(Equal([]string{"one.sql"}))
		})

		It("keeps the list to unique entries", func() {
			newMigrationList := []string{"one.sql", "one.sql"}
			oldMigrationList := []string{"two.sql", "two.sql"}

			Expect(cmd.NewMigrations(newMigrationList, oldMigrationList)).To(Equal([]string{"one.sql"}))
		})
	})
})
