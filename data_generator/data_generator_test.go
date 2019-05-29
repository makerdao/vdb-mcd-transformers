package main

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"math/rand"
	"os/exec"
)

var _ = Describe("data generator", func() {
	Describe("with a given seed", func() {
		var (
			db    *postgres.DB
			state GeneratorState
			seed  int64
		)

		BeforeEach(func() {
			seed = rand.Int63()
			rand.Seed(int64(seed))
			db = test_config.NewTestDB(test_config.NewTestNode())
			test_config.CleanTestDB(db)
			state = NewGenerator(db)
		})

		// Runs twice with the same seed, dumps the DB data, repeats and compares results
		It("is repeatable", func() {
			resetIdColumnSequences(db)
			runErr := state.Run(30)
			Expect(runErr).NotTo(HaveOccurred())
			state1 := dumpTables()

			// Run two
			test_config.CleanTestDB(db)
			resetIdColumnSequences(db)
			rand.Seed(int64(seed))
			newState := NewGenerator(db)
			runErr = newState.Run(30)
			Expect(runErr).NotTo(HaveOccurred())
			state2 := dumpTables()

			Expect(state1).To(Equal(state2))
		})
	})
})

/* pg_dump is a bit erratic and sometimes moves the ordering around in the dump, so we filter out all insert-lines
representing the data, and sort alphabetically. Each line is explicit in everything including indices, so data
representation is still maintained.

Lines look like this from pg_dump --column-inserts:
INSERT INTO maker.urns (id, ilk_id, guy) VALUES (2, 1, '0x2b632c3cc964cd65822263db0d9752df3b024687');
*/
func dumpTables() []byte {
	pgDumpCmd := "pg_dump --data-only --column-inserts --schema maker --schema public vulcanize_private"
	cleanAndSortCmd := " | egrep '^INSERT' | sort"
	state, dumpErr := exec.Command("sh", "-c", pgDumpCmd+cleanAndSortCmd).Output()
	Expect(dumpErr).NotTo(HaveOccurred())
	return state
}

// Grabs all user sequences from the DB and restarts them.
// This allows the id columns to match between subsequent runs.
func resetIdColumnSequences(db *postgres.DB) {
	var sequences []string
	seqErr := db.Select(&sequences,
		`SELECT schemaname || '.' || relname from pg_statio_user_sequences
		WHERE relname NOT IN ('goose_db_version_id_seq')`)
	Expect(seqErr).NotTo(HaveOccurred())

	for _, seq := range sequences {
		_, alterErr := db.Exec(fmt.Sprintf("ALTER SEQUENCE %v RESTART", seq))
		Expect(alterErr).NotTo(HaveOccurred())
	}
}
