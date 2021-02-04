package queries

import (
	"database/sql"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const getMinMaxTransformedDiffsQuery = `SELECT * FROM api.get_storage_transformations_for_status($1)`

var _ = Describe("api.storage_transformations_for_status($1)", func() {

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("returns two rows, one for the max value and one for min value", func() {
		header := fakes.GetFakeHeader(rand.Int63())
		persistedDiff := test_helpers.CreateDiffRecord(db, header, fakes.FakeAddress, fakes.FakeHash, fakes.AnotherFakeHash)

		var stateDiffMetaData []metaData
		Expect(db.Select(&stateDiffMetaData, getMinMaxTransformedDiffsQuery, "new")).To(Succeed())

		Expect(len(stateDiffMetaData)).To(Equal(2))
		elementMap := convertMinMaxRowsToHash(stateDiffMetaData)

		Expect(elementMap["min"].MinOrMax).To(Equal("min"))
		Expect(elementMap["min"].Address).To(Equal(fakes.FakeAddress.Bytes()))
		Expect(elementMap["min"].BlockHash).To(Equal(persistedDiff.RawDiff.BlockHash))
		Expect(elementMap["min"].BlockHeight).To(Equal(header.BlockNumber))
		Expect(elementMap["min"].FromBackfill).To(Equal(false))
		Expect(elementMap["min"].Status).To(Equal("new"))
		Expect(elementMap["min"].StorageKey).To(Equal(fakes.FakeHash.Bytes()))
		Expect(elementMap["min"].StorageValue).To(Equal(fakes.AnotherFakeHash.Bytes()))

		Expect(elementMap["max"].MinOrMax).To(Equal("max"))
		Expect(elementMap["max"].Address).To(Equal(fakes.FakeAddress.Bytes()))
		Expect(elementMap["max"].BlockHash).To(Equal(persistedDiff.RawDiff.BlockHash))
		Expect(elementMap["max"].BlockHeight).To(Equal(header.BlockNumber))
		Expect(elementMap["max"].FromBackfill).To(Equal(false))
		Expect(elementMap["max"].Status).To(Equal("new"))
		Expect(elementMap["max"].StorageKey).To(Equal(fakes.FakeHash.Bytes()))
		Expect(elementMap["max"].StorageValue).To(Equal(fakes.AnotherFakeHash.Bytes()))
	})

	It("returns the max block height for the max row, min block height for the min row", func() {
		headerOne := fakes.GetFakeHeader(rand.Int63())
		test_helpers.CreateDiffRecord(db, headerOne, fakes.FakeAddress, fakes.FakeHash, fakes.AnotherFakeHash)

		headerTwo := fakes.GetFakeHeader(headerOne.BlockNumber + 1)
		test_helpers.CreateDiffRecord(db, headerTwo, fakes.FakeAddress, fakes.FakeHash, fakes.AnotherFakeHash)

		var stateDiffMetaData []metaData
		Expect(db.Select(&stateDiffMetaData, getMinMaxTransformedDiffsQuery, "new")).To(Succeed())

		elementMap := convertMinMaxRowsToHash(stateDiffMetaData)

		Expect(elementMap["min"].MinOrMax).To(Equal("min"))
		Expect(elementMap["min"].BlockHeight).To(Equal(headerOne.BlockNumber))

		Expect(elementMap["max"].MinOrMax).To(Equal("max"))
		Expect(elementMap["max"].BlockHeight).To(Equal(headerTwo.BlockNumber))
	})

	It("respects the passed in status", func() {
		headerOne := fakes.GetFakeHeader(rand.Int63())
		persistedDiff := test_helpers.CreateDiffRecord(db, headerOne, fakes.FakeAddress, fakes.FakeHash, fakes.AnotherFakeHash)
		diffRepo := storage.NewDiffRepository(db)
		diffRepo.MarkTransformed(persistedDiff.ID)

		var stateDiffMetaData []metaData
		Expect(db.Select(&stateDiffMetaData, getMinMaxTransformedDiffsQuery, "new")).To(Succeed())

		Expect(stateDiffMetaData).To(BeEmpty())
	})

	It("does not require a case-sensitive status to be passed in", func() {
		headerOne := fakes.GetFakeHeader(rand.Int63())
		test_helpers.CreateDiffRecord(db, headerOne, fakes.FakeAddress, fakes.FakeHash, fakes.AnotherFakeHash)

		var stateDiffMetaData []metaData
		Expect(db.Select(&stateDiffMetaData, getMinMaxTransformedDiffsQuery, "NEW")).To(Succeed())

		Expect(len(stateDiffMetaData)).To(Equal(2))
	})

})

func convertMinMaxRowsToHash(rows []metaData) map[string]metaData {
	elementMap := make(map[string]metaData)
	for _, data := range rows {
		elementMap[data.MinOrMax] = data
	}
	return elementMap
}

type metaData struct {
	MinOrMax     string         `db:"min_or_max"`
	Address      []byte         `db:"address"`
	BlockHash    common.Hash    `db:"block_hash"`
	BlockHeight  int64          `db:"block_height"`
	FromBackfill bool           `db:"from_backfill"`
	Status       string         `db:"status"`
	StorageKey   []byte         `db:"storage_key"`
	StorageValue []byte         `db:"storage_value"`
	Created      sql.NullString `db:"created"`
	Updated      sql.NullString `db:"updated"`
}
