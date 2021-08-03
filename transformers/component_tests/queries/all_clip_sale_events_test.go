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

var _ = Describe("All clip sale events query", func() {
	var (
		headerRepo      datastore.HeaderRepository
		contractAddress = fakes.FakeAddress.Hex()
		addressId              int64
		saleId                 int
		blockOne, timestampOne int
		headerOne              core.Header
		clipKickEvent          event.InsertionModel
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		saleId = rand.Int()

		blockOne = rand.Int()
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepo)

		clipKickLog := test_data.CreateTestLog(headerOne.Id, db)

		var addressErr error
		addressId, addressErr = repository.GetOrCreateAddress(db, contractAddress)
		Expect(addressErr).NotTo(HaveOccurred())

		clipKickEvent = test_data.ClipKickModel()
		clipKickEvent.ColumnValues[event.HeaderFK] = headerOne.Id
		clipKickEvent.ColumnValues[event.LogFK] = clipKickLog.ID
		clipKickEvent.ColumnValues[event.AddressFK] = addressId
		clipKickEvent.ColumnValues[constants.SaleIDColumn] = strconv.Itoa(saleId)
		clipKickErr := event.PersistModels([]event.InsertionModel{clipKickEvent}, db)
		Expect(clipKickErr).NotTo(HaveOccurred())
	})

	Describe("all_clip_bid_events", func() {
		It("returns all clip bid events when they are all in the same block", func() {
			clipTakeLog := test_data.CreateTestLog(headerOne.Id, db)
			clipTakeErr := test_helpers.CreateTake(test_helpers.TakeCreationInput{
				DB:              db,
				ContractAddress: contractAddress,
				SaleId:          saleId,
				TakeHeaderId:    headerOne.Id,
				TakeLogId:       clipTakeLog.ID,
			})
			Expect(clipTakeErr).NotTo(HaveOccurred())

			clipRedoLog := test_data.CreateTestLog(headerOne.Id, db)
			clipRedoErr := test_helpers.CreateRedo(test_helpers.RedoCreationInput{
				DB:              db,
				ContractAddress: contractAddress,
				SaleId:          saleId,
				RedoHeaderId:    headerOne.Id,
				RedoLogId:       clipRedoLog.ID,
			})
			Expect(clipRedoErr).NotTo(HaveOccurred())

			clipYankLog := test_data.CreateTestLog(headerOne.Id, db)
			clipYankErr := test_helpers.CreateClipYank(test_helpers.ClipYankCreationInput{
				DB:               db,
				ContractAddress:  contractAddress,
				SaleId:           saleId,
				ClipYankHeaderId: headerOne.Id,
				ClipYankLogId:    clipYankLog.ID,
			})
			Expect(clipYankErr).NotTo(HaveOccurred())

			var actualClipSaleEvents []test_helpers.SaleEvent
			queryErr := db.Select(&actualClipSaleEvents, `SELECT sale_id, act FROM api.all_clip_sale_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualClipSaleEvents).To(ConsistOf(
				test_helpers.SaleEvent{SaleId: strconv.Itoa(saleId), Act: "kick"},
				test_helpers.SaleEvent{SaleId: strconv.Itoa(saleId), Act: "yank"},
				test_helpers.SaleEvent{SaleId: strconv.Itoa(saleId), Act: "redo"},
				test_helpers.SaleEvent{SaleId: strconv.Itoa(saleId), Act: "take"}),
			)
		})

		It("returns clip sale events across all blocks", func() {
			clipTakeLog := test_data.CreateTestLog(headerOne.Id, db)
			clipTakeErr := test_helpers.CreateTake(test_helpers.TakeCreationInput{
				DB:              db,
				ContractAddress: contractAddress,
				SaleId:          saleId,
				TakeHeaderId:    headerOne.Id,
				TakeLogId:       clipTakeLog.ID,
			})
			Expect(clipTakeErr).NotTo(HaveOccurred())

			headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)

			clipRedoLog := test_data.CreateTestLog(headerTwo.Id, db)
			clipRedoErr := test_helpers.CreateRedo(test_helpers.RedoCreationInput{
				DB:              db,
				ContractAddress: contractAddress,
				SaleId:          saleId,
				RedoHeaderId:    headerTwo.Id,
				RedoLogId:       clipRedoLog.ID,
			})
			Expect(clipRedoErr).NotTo(HaveOccurred())

			clipYankLog := test_data.CreateTestLog(headerOne.Id, db)
			clipYankErr := test_helpers.CreateClipYank(test_helpers.ClipYankCreationInput{
				DB:               db,
				ContractAddress:  contractAddress,
				SaleId:           saleId,
				ClipYankHeaderId: headerOne.Id,
				ClipYankLogId:    clipYankLog.ID,
			})
			Expect(clipYankErr).NotTo(HaveOccurred())
		})
	})
})
