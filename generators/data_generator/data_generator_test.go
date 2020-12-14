package main

import (
	"fmt"
	"math/rand"
	"os/exec"

	"github.com/jmoiron/sqlx"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vulcanizedb/pkg/config"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
			var newDBerr error
			db, newDBerr = newTestDB()
			Expect(newDBerr).NotTo(HaveOccurred())

			test_config.CleanTestDB(db)
			cleanEthNodesErr := deleteEthNodes(db)
			Expect(cleanEthNodesErr).NotTo(HaveOccurred())

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

/* This special NewDB creaets a db WITHOUT a node. It does this because the code under*/
/* test creates a node, and so using the standard postgres.NewDB here can cause false*/
/* positives. */
func newTestDB() (*postgres.DB, error) {
	connectString := config.DbConnectionString(test_config.DBConfig)
	db, connectErr := sqlx.Connect("postgres", connectString)
	if connectErr != nil {
		return &postgres.DB{}, postgres.ErrDBConnectionFailed(connectErr)
	}
	pg := postgres.DB{DB: db}
	return &pg, nil
}

func deleteEthNodes(db *postgres.DB) error {
	_, err := db.Exec("TRUNCATE TABLE eth_nodes CASCADE;")
	return err
}

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
