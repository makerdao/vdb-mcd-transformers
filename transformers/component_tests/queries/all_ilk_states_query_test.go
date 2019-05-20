package queries

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/storage/cat"
	"github.com/vulcanize/mcd_transformers/transformers/storage/jug"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vat"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
)

var _ = Describe("Ilk State History Query", func() {
	var (
		db                       *postgres.DB
		blockOne                 = rand.Int()
		blockTwo                 = blockOne + 1
		vatRepository            vat.VatStorageRepository
		catRepository            cat.CatStorageRepository
		jugRepository            jug.JugStorageRepository
		blockOneIlkValues        map[string]string
		blockTwoIlkValues        map[string]string
		expectedBlockOneIlkState test_helpers.IlkState
		expectedBlockTwoIlkState test_helpers.IlkState
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		vatRepository.SetDB(db)
		catRepository.SetDB(db)
		jugRepository.SetDB(db)
		headerRepository = repositories.NewHeaderRepository(db)

		blockOneHeader = fakes.GetFakeHeader(int64(blockOne))
		_, err := headerRepository.CreateOrUpdateHeader(blockOneHeader)
		Expect(err).NotTo(HaveOccurred())

		blockTwoHeader = fakes.GetFakeHeader(int64(blockTwo))
		blockTwoHeader.Timestamp = blockTwoHeader.Timestamp + "2"
		blockTwoHeader.Hash = "block2Hash"
		_, err = headerRepository.CreateOrUpdateHeader(blockTwoHeader)
		Expect(err).NotTo(HaveOccurred())

		blockOneIlkValues = test_helpers.GetIlkValues(0)
		test_helpers.CreateVatRecords(blockOneHeader, blockOneIlkValues, test_helpers.FakeIlkVatMetadatas, vatRepository)
		test_helpers.CreateCatRecords(blockOneHeader, blockOneIlkValues, test_helpers.FakeIlkCatMetadatas, catRepository)
		test_helpers.CreateJugRecords(blockOneHeader, blockOneIlkValues, test_helpers.FakeIlkJugMetadatas, jugRepository)
		expectedBlockOneIlkState = test_helpers.IlkStateFromValues(test_helpers.FakeIlk.Hex, blockOneHeader.Timestamp, blockOneHeader.Timestamp, blockOneIlkValues)

		blockTwoIlkValues = test_helpers.GetIlkValues(1)
		test_helpers.CreateVatRecords(blockTwoHeader, blockTwoIlkValues, test_helpers.FakeIlkVatMetadatas, vatRepository)
		test_helpers.CreateCatRecords(blockTwoHeader, blockTwoIlkValues, test_helpers.FakeIlkCatMetadatas, catRepository)
		test_helpers.CreateJugRecords(blockTwoHeader, blockTwoIlkValues, test_helpers.FakeIlkJugMetadatas, jugRepository)
		expectedBlockTwoIlkState = test_helpers.IlkStateFromValues(test_helpers.FakeIlk.Hex, blockTwoHeader.Timestamp, blockOneHeader.Timestamp, blockTwoIlkValues)
	})

	It("returns the history of an ilk from the given block height", func() {
		var dbResult []test_helpers.IlkState
		err := db.Select(&dbResult,
			`SELECT ilk_name, rate, art, spot, line, dust, chop, lump, flip, rho, duty, created, updated from api.all_ilk_states($1, $2)`,
			test_helpers.FakeIlk.Name, blockTwo)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(dbResult)).To(Equal(2))
		Expect(dbResult).To(ConsistOf([]test_helpers.IlkState{
			expectedBlockOneIlkState,
			expectedBlockTwoIlkState,
		}))
	})

	It("can handle multiple ilks in the db", func() {
		blockOneAnotherFakeIlkValues := test_helpers.GetIlkValues(3)

		test_helpers.CreateVatRecords(blockOneHeader, blockOneAnotherFakeIlkValues, test_helpers.AnotherFakeIlkVatMetadatas, vatRepository)
		test_helpers.CreateCatRecords(blockOneHeader, blockOneAnotherFakeIlkValues, test_helpers.AnotherFakeIlkCatMetadatas, catRepository)
		test_helpers.CreateJugRecords(blockOneHeader, blockOneAnotherFakeIlkValues, test_helpers.AnotherFakeIlkJugMetadatas, jugRepository)

		var dbResult []test_helpers.IlkState
		err := db.Select(&dbResult,
			`SELECT ilk_name, rate, art, spot, line, dust, chop, lump, flip, rho, duty, created, updated from api.all_ilk_states($1, $2)`,
			test_helpers.AnotherFakeIlk.Name, blockTwo)
		Expect(err).NotTo(HaveOccurred())
		expectedBlockOneAnotherIlkState := test_helpers.IlkStateFromValues(test_helpers.AnotherFakeIlk.Hex,
			blockOneHeader.Timestamp, blockOneHeader.Timestamp, blockOneAnotherFakeIlkValues)
		//does not include fake ilk's results
		Expect(len(dbResult)).To(Equal(1))
		Expect(dbResult).To(ConsistOf([]test_helpers.IlkState{
			expectedBlockOneAnotherIlkState,
		}))
	})

	It("handles a query with a block height before the ilk is in the db", func() {
		blockZero := blockOne - 1

		blockZeroHeader := fakes.GetFakeHeader(int64(blockZero))
		_, err := headerRepository.CreateOrUpdateHeader(blockZeroHeader)
		Expect(err).NotTo(HaveOccurred())

		var dbResult []test_helpers.IlkState
		err = db.Select(&dbResult,
			`SELECT ilk_name, rate, art, spot, line, dust, chop, lump, flip, rho, duty, created, updated from api.all_ilk_states($1, $2)`,
			test_helpers.FakeIlk.Name, blockZero)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult).To(BeEmpty())
	})

	It("handles when there have been no recent updates to the ilk", func() {
		blockOneHundred := int64(blockOne + 100)
		blockOneHundredHeader := fakes.GetFakeHeader(blockOneHundred)
		_, err := headerRepository.CreateOrUpdateHeader(blockOneHundredHeader)
		Expect(err).NotTo(HaveOccurred())

		test_helpers.CreateVatRecords(blockOneHundredHeader, test_helpers.GetIlkValues(100), test_helpers.AnotherFakeIlkVatMetadatas, vatRepository)

		var dbResult []test_helpers.IlkState
		err = db.Select(&dbResult,
			`SELECT ilk_name, rate, art, spot, line, dust, chop, lump, flip, rho, duty, created, updated from api.all_ilk_states($1, $2)`,
			test_helpers.FakeIlk.Name, blockOneHundred)
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
		err := db.Select(&dbResult, `SELECT block_height FROM api.all_ilk_states($1)`, test_helpers.FakeIlk.Name)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbResult).To(Equal([]int{blockTwo, blockOne}))
	})
})
