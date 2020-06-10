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

	Describe("CheckNewMigrations", func() {
		It("is not an error to have no new migrations", func() {
			originalMigrations := []string{}
			newMigrations := []string{}

			err := cmd.CheckNewMigrations(originalMigrations, newMigrations)

			Expect(err).NotTo(HaveOccurred())
		})

		It("is an error if the new migrations are out of order", func() {
			originalMigrations := []string{"20200429072513_last.sql"}
			newMigrations := []string{"20200429072512_newest.sql"}

			err := cmd.CheckNewMigrations(originalMigrations, newMigrations)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp("update your timestamp"))
		})

		It("does not change the new migrations passed in", func() {
			originalMigrations := []string{"20200429072513_last.sql"}
			newMigrations := []string{"20200429072515_newest.sql", "20200429072512_newest.sql"}

			cmd.CheckNewMigrations(originalMigrations, newMigrations)

			Expect(newMigrations).To(Equal([]string{"20200429072515_newest.sql", "20200429072512_newest.sql"}))
		})

		It("does not erroneously report an error if the new Migrations come in out of order", func() {
			originalMigrations := []string{"20200429072513_last.sql"}
			newMigrations := []string{"20200429072517_newest.sql", "20200429072515_newest.sql"}

			err := cmd.CheckNewMigrations(originalMigrations, newMigrations)

			Expect(err).NotTo(HaveOccurred())
		})

		It("requires every new migration to start with a 14 digit timestamp", func() {
			originalMigrations := []string{"20200429072513_last.sql"}
			newMigrations := []string{"20200429072517_newest.sql", "2020042907251z_newest.sql"}

			err := cmd.CheckNewMigrations(originalMigrations, newMigrations)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp("does not start with a timestamp"))
		})
	})
})
