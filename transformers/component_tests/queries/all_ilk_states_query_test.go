package queries

import (
	storage_helper "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"math/rand"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/cat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/jug"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/spot"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/vat"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ilk State History Query", func() {
	var (
		vatRepository            vat.VatStorageRepository
		catRepository            cat.CatStorageRepository
		jugRepository            jug.JugStorageRepository
		spotRepository           spot.SpotStorageRepository
		blockOneIlkValues        map[string]string
		blockTwoIlkValues        map[string]string
		expectedBlockOneIlkState test_helpers.IlkState
		expectedBlockTwoIlkState test_helpers.IlkState
		headerRepository         repositories.HeaderRepository
		blockOne, timestampOne   int
		headerOne, headerTwo     core.Header
		diffID int64
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		vatRepository.SetDB(db)
		catRepository.SetDB(db)
		jugRepository.SetDB(db)
		spotRepository.SetDB(db)
		headerRepository = repositories.NewHeaderRepository(db)

		blockOne = rand.Int()
		blockTwo := blockOne + 1
		timestampOne = int(rand.Int31())
		timestampTwo := timestampOne + 1
		headerOne = createHeader(blockOne, timestampOne, headerRepository)
		headerTwo = createHeader(blockTwo, timestampTwo, headerRepository)

		diffID = storage_helper.CreateDiffRecord(db)

		blockOneIlkValues = test_helpers.GetIlkValues(0)
		test_helpers.CreateVatRecords(diffID, headerOne, blockOneIlkValues, test_helpers.FakeIlkVatMetadatas, vatRepository)
		test_helpers.CreateCatRecords(diffID, headerOne.Id, blockOneIlkValues, test_helpers.FakeIlkCatMetadatas, catRepository)
		test_helpers.CreateJugRecords(diffID, headerOne, blockOneIlkValues, test_helpers.FakeIlkJugMetadatas, jugRepository)
		test_helpers.CreateSpotRecords(diffID, headerOne, blockOneIlkValues, test_helpers.FakeIlkSpotMetadatas, spotRepository)
		expectedBlockOneIlkState = test_helpers.IlkStateFromValues(test_helpers.FakeIlk.Hex, headerOne.Timestamp, headerOne.Timestamp, blockOneIlkValues)

		blockTwoIlkValues = test_helpers.GetIlkValues(1)
		test_helpers.CreateVatRecords(diffID, headerTwo, blockTwoIlkValues, test_helpers.FakeIlkVatMetadatas, vatRepository)
		test_helpers.CreateCatRecords(diffID, headerTwo.Id, blockTwoIlkValues, test_helpers.FakeIlkCatMetadatas, catRepository)
		test_helpers.CreateJugRecords(diffID, headerTwo, blockTwoIlkValues, test_helpers.FakeIlkJugMetadatas, jugRepository)
		test_helpers.CreateSpotRecords(diffID, headerTwo, blockTwoIlkValues, test_helpers.FakeIlkSpotMetadatas, spotRepository)
		expectedBlockTwoIlkState = test_helpers.IlkStateFromValues(test_helpers.FakeIlk.Hex, headerTwo.Timestamp, headerOne.Timestamp, blockTwoIlkValues)
	})

	It("returns the history of an ilk from the given block height", func() {
		var dbResult []test_helpers.IlkState
		err := db.Select(&dbResult,
			`SELECT ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated from api.all_ilk_states($1, $2)`,
			test_helpers.FakeIlk.Identifier, headerTwo.BlockNumber)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(2))
		Expect(dbResult).To(ConsistOf([]test_helpers.IlkState{
			expectedBlockOneIlkState,
			expectedBlockTwoIlkState,
		}))
	})

	It("can handle multiple ilks in the db", func() {
		blockOneAnotherFakeIlkValues := test_helpers.GetIlkValues(3)

		test_helpers.CreateVatRecords(diffID, headerOne, blockOneAnotherFakeIlkValues, test_helpers.AnotherFakeIlkVatMetadatas, vatRepository)
		test_helpers.CreateCatRecords(diffID, headerOne.Id, blockOneAnotherFakeIlkValues, test_helpers.AnotherFakeIlkCatMetadatas, catRepository)
		test_helpers.CreateJugRecords(diffID, headerOne, blockOneAnotherFakeIlkValues, test_helpers.AnotherFakeIlkJugMetadatas, jugRepository)
		test_helpers.CreateSpotRecords(diffID, headerOne, blockOneAnotherFakeIlkValues, test_helpers.AnotherFakeIlkSpotMetadatas, spotRepository)

		var dbResult []test_helpers.IlkState
		err := db.Select(&dbResult,
			`SELECT ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated from api.all_ilk_states($1, $2)`,
			test_helpers.AnotherFakeIlk.Identifier, headerTwo.BlockNumber)
		Expect(err).NotTo(HaveOccurred())
		expectedBlockOneAnotherIlkState := test_helpers.IlkStateFromValues(test_helpers.AnotherFakeIlk.Hex,
			headerOne.Timestamp, headerOne.Timestamp, blockOneAnotherFakeIlkValues)
		//does not include fake ilk's results
		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult).To(ConsistOf([]test_helpers.IlkState{
			expectedBlockOneAnotherIlkState,
		}))
	})

	It("handles a query with a block height before the ilk is in the db", func() {
		blockZero := blockOne - 1
		_ = createHeader(blockZero, timestampOne-1, headerRepository)

		var dbResult []test_helpers.IlkState
		err := db.Select(&dbResult,
			`SELECT ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated from api.all_ilk_states($1, $2)`,
			test_helpers.FakeIlk.Identifier, blockZero)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult).To(BeEmpty())
	})

	It("handles when there have been no recent updates to the ilk", func() {
		blockOneHundred := blockOne + 100
		blockOneHundredHeader := createHeader(blockOneHundred, timestampOne+100, headerRepository)

		test_helpers.CreateVatRecords(diffID, blockOneHundredHeader, test_helpers.GetIlkValues(100), test_helpers.AnotherFakeIlkVatMetadatas, vatRepository)

		var dbResult []test_helpers.IlkState
		err := db.Select(&dbResult,
			`SELECT ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated from api.all_ilk_states($1, $2)`,
			test_helpers.FakeIlk.Identifier, blockOneHundred)
		Expect(err).NotTo(HaveOccurred())

		Expect(dbResult).To(ConsistOf(expectedBlockOneIlkState, expectedBlockTwoIlkState))
	})

	It("fails if no argument is supplied (STRICT)", func() {
		_, err := db.Exec(`SELECT * FROM api.all_ilk_states()`)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("function api.all_ilk_states() does not exist"))
	})

	It("uses default value for blockHeight if not supplied", func() {
		var dbResult []int
		err := db.Select(&dbResult, `SELECT block_height FROM api.all_ilk_states($1)`, test_helpers.FakeIlk.Identifier)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult).To(Equal([]int{int(headerTwo.BlockNumber), blockOne}))
	})

	Describe("result pagination", func() {
		It("limits the results to the most recent [limit] ilk states when a limit argument is provided", func() {
			maxResults := 1
			var dbResult []test_helpers.IlkState
			err := db.Select(&dbResult,
				`SELECT ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated from api.all_ilk_states($1, $2, $3)`,
				test_helpers.FakeIlk.Identifier, headerTwo.BlockNumber, maxResults)
			Expect(err).NotTo(HaveOccurred())

			Expect(dbResult).To(ConsistOf(expectedBlockTwoIlkState))
		})

		It("offsets results if offset is provided", func() {
			maxResults := 1
			resultOffset := 1
			var dbResult []test_helpers.IlkState
			err := db.Select(&dbResult,
				`SELECT ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated from api.all_ilk_states($1, $2, $3, $4)`,
				test_helpers.FakeIlk.Identifier, headerTwo.BlockNumber, maxResults, resultOffset)
			Expect(err).NotTo(HaveOccurred())

			Expect(dbResult).To(ConsistOf(expectedBlockOneIlkState))
		})
	})
})
