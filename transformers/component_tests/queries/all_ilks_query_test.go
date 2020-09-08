package queries

import (
	"math/rand"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	ilkHelpers "github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	headerHelpers "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("All Ilks query", func() {
	var (
		headerOne, headerTwo core.Header

		ethIlkValuesOne, ethIlkValuesTwo, batIlkValuesOne, batIlkValuesTwo map[string]interface{}
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)

		blockOne := rand.Int()
		blockTwo := blockOne + 1
		timestampOne := int64(rand.Int31())
		timestampTwo := timestampOne + 1
		headerOne = headerHelpers.CreateHeader(timestampOne, blockOne, db)
		headerTwo = headerHelpers.CreateHeader(timestampTwo, blockTwo, db)

		// set up 2 different ilks in 2 different blocks
		ethIlkValuesOne = ilkHelpers.GetIlkValues(0)
		ethIlkValuesTwo = ilkHelpers.GetIlkValues(1)
		batIlkValuesOne = ilkHelpers.GetIlkValues(2)
		batIlkValuesTwo = ilkHelpers.GetIlkValues(3)
		ilkHelpers.CreateIlk(db, headerOne, ethIlkValuesOne, ilkHelpers.FakeIlkVatMetadatas, ilkHelpers.FakeIlkCatMetadatas, ilkHelpers.FakeIlkJugMetadatas, ilkHelpers.FakeIlkSpotMetadatas)
		ilkHelpers.CreateIlk(db, headerTwo, ethIlkValuesTwo, ilkHelpers.FakeIlkVatMetadatas, ilkHelpers.FakeIlkCatMetadatas, ilkHelpers.FakeIlkJugMetadatas, ilkHelpers.FakeIlkSpotMetadatas)
		ilkHelpers.CreateIlk(db, headerOne, batIlkValuesOne, ilkHelpers.AnotherFakeIlkVatMetadatas, ilkHelpers.AnotherFakeIlkCatMetadatas, ilkHelpers.AnotherFakeIlkJugMetadatas, ilkHelpers.AnotherFakeIlkSpotMetadatas)
		ilkHelpers.CreateIlk(db, headerTwo, batIlkValuesTwo, ilkHelpers.AnotherFakeIlkVatMetadatas, ilkHelpers.AnotherFakeIlkCatMetadatas, ilkHelpers.AnotherFakeIlkJugMetadatas, ilkHelpers.AnotherFakeIlkSpotMetadatas)
	})

	Context("when called with no arguments", func() {
		It("returns the most recent state of each ilk", func() {
			expectedEthIlk := ilkHelpers.IlkSnapshotFromValues(ilkHelpers.FakeIlk.Hex, headerTwo.Timestamp,
				headerOne.Timestamp, ethIlkValuesTwo)
			expectedBatIlk := ilkHelpers.IlkSnapshotFromValues(ilkHelpers.AnotherFakeIlk.Hex, headerTwo.Timestamp,
				headerOne.Timestamp, batIlkValuesTwo)

			var actualIlks []ilkHelpers.IlkSnapshot
			err := db.Select(&actualIlks, `SELECT ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, dunk, created, updated FROM api.all_ilks()`)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualIlks).To(ConsistOf(expectedEthIlk, expectedBatIlk))
		})
	})

	Context("when called with block number", func() {
		It("returns the most recent state of each ilk as of the given block number", func() {
			expectedEthIlk := ilkHelpers.IlkSnapshotFromValues(ilkHelpers.FakeIlk.Hex, headerOne.Timestamp,
				headerOne.Timestamp, ethIlkValuesOne)
			expectedBatIlk := ilkHelpers.IlkSnapshotFromValues(ilkHelpers.AnotherFakeIlk.Hex, headerOne.Timestamp,
				headerOne.Timestamp, batIlkValuesOne)

			var actualIlks []ilkHelpers.IlkSnapshot
			err := db.Select(&actualIlks, `SELECT ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, dunk, created, updated FROM api.all_ilks($1)`,
				headerOne.BlockNumber)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualIlks).To(ConsistOf(expectedEthIlk, expectedBatIlk))
		})
	})
})
