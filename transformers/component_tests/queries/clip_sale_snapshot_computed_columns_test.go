package queries

import (
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("clip_sale_snapshot computed columns", func() {
	var (
		headerOne              core.Header
		headerRepository       datastore.HeaderRepository
		contractAddress        = fakes.FakeAddress.Hex()
		addressId              int64
		fakeSaleId             int
		blockOne, timestampOne int
	)

	BeforeEach(func() {
		fakeSaleId = rand.Int()
		blockOne = rand.Int()
		timestampOne = int(rand.Int31())

		test_config.CleanTestDB(db)

		headerRepository = repositories.NewHeaderRepository(db)
		headerOne = createHeader(blockOne, timestampOne, headerRepository)

		clipKickLog := test_data.CreateTestLog(headerOne.Id, db)

		var addressErr error
		addressId, addressErr = repository.GetOrCreateAddress(db, contractAddress)
		Expect(addressErr).NotTo(HaveOccurred())

		clipKickEvent := test_data.ClipKickModel()
		clipKickEvent.ColumnValues[event.HeaderFK] = headerOne.Id
		clipKickEvent.ColumnValues[event.LogFK] = clipKickLog.ID
		clipKickEvent.ColumnValues[event.AddressFK] = addressId
		clipKickEvent.ColumnValues[constants.SaleIDColumn] = strconv.Itoa(fakeSaleId)
		clipKickErr := event.PersistModels([]event.InsertionModel{clipKickEvent}, db)
		Expect(clipKickErr).NotTo(HaveOccurred())

		clipStorageValues := test_helpers.GetClipStorageValues(1, fakeSaleId)
		test_helpers.CreateClip(db, headerOne, clipStorageValues, test_helpers.GetClipMetadatas(strconv.Itoa(fakeSaleId)), contractAddress)
	})

	Describe("clip sale snapshot sale events", func() {
		It("returns the sale events for a clip", func() {
			expectedClipKickEvent := test_helpers.SaleEvent{
				SaleId:          strconv.Itoa(fakeSaleId),
				Act:             "kick",
				ContractAddress: contractAddress,
			}

			var actualSaleEvents []test_helpers.SaleEvent
			queryErr := db.Select(&actualSaleEvents,
				`SELECT sale_id, act, contract_address FROM api.clip_sale_snapshot_sale_events(
    					(SELECT (block_height, sale_id, ilk_id, urn_id, pos, tab, lot, usr, tic, "top", created, updated, clip_address)::api.clip_sale_snapshot
    					 FROM api.get_clip_with_address($1, $2, $3)))`, fakeSaleId, contractAddress, blockOne)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(actualSaleEvents).To(ConsistOf(expectedClipKickEvent))
		})
	})
})
