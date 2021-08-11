package queries

import (
	"math/rand"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
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
		headerRepo             datastore.HeaderRepository
		contractAddress        = fakes.FakeAddress.Hex()
		anotherClipAddress     = common.HexToAddress("0xabcdef123456789").Hex()
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

	Describe("all_clip_sale_events", func() {
		It("returns all clip sale events when they are all in the same block", func() {
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

			clipStorageValues := test_helpers.GetClipStorageValues(1, saleId)
			test_helpers.CreateClip(db, headerOne, clipStorageValues, test_helpers.GetClipMetadatas(strconv.Itoa(saleId)), contractAddress)

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

			clipStorageValuesBlockOne := test_helpers.GetClipStorageValues(1, saleId)
			test_helpers.CreateClip(db, headerOne, clipStorageValuesBlockOne, test_helpers.GetClipMetadatas(strconv.Itoa(saleId)), contractAddress)

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

			clipStorageValuesBlockTwo := test_helpers.GetClipStorageValues(2, saleId)
			test_helpers.CreateClip(db, headerTwo, clipStorageValuesBlockTwo, test_helpers.GetClipMetadatas(strconv.Itoa(saleId)), contractAddress)

			headerThree := createHeader(blockOne+2, timestampOne+2, headerRepo)

			clipYankLog := test_data.CreateTestLog(headerThree.Id, db)
			clipYankErr := test_helpers.CreateClipYank(test_helpers.ClipYankCreationInput{
				DB:               db,
				ContractAddress:  contractAddress,
				SaleId:           saleId,
				ClipYankHeaderId: headerThree.Id,
				ClipYankLogId:    clipYankLog.ID,
			})
			Expect(clipYankErr).NotTo(HaveOccurred())

			clipStorageValuesBlockThree := test_helpers.GetClipStorageValues(3, saleId)
			test_helpers.CreateClip(db, headerThree, clipStorageValuesBlockThree, test_helpers.GetClipMetadatas(strconv.Itoa(saleId)), contractAddress)

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

		Describe("results pagination", func() {
			var updatedClipValues map[string]interface{}

			BeforeEach(func() {
				headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)

				logID := test_data.CreateTestLog(headerTwo.Id, db).ID

				clipTakeErr := test_helpers.CreateTake(test_helpers.TakeCreationInput{
					DB:              db,
					ContractAddress: contractAddress,
					SaleId:          saleId,
					TakeHeaderId:    headerTwo.Id,
					TakeLogId:       logID,
				})
				Expect(clipTakeErr).NotTo(HaveOccurred())

				updatedClipValues = test_helpers.GetClipStorageValues(2, saleId)
				test_helpers.CreateClip(db, headerTwo, updatedClipValues, test_helpers.GetClipMetadatas(strconv.Itoa(saleId)), contractAddress)
			})

			It("limits result to latest blocks if max_results argument is provided", func() {
				maxResults := 1
				var actualSaleEvents []test_helpers.SaleEvent
				queryErr := db.Select(&actualSaleEvents, `SELECT sale_id, act FROM api.all_clip_sale_events($1)`, maxResults)
				Expect(queryErr).NotTo(HaveOccurred())

				Expect(actualSaleEvents).To(ConsistOf(
					test_helpers.SaleEvent{
						SaleId: strconv.Itoa(saleId),
						Act:    "take",
					},
				))
			})

			XIt("offsets results if offset is provided", func() {
				maxResults := 1
				resultOffset := 1
				var actualClipEvents []test_helpers.SaleEvent
				queryErr := db.Select(&actualClipEvents, `SELECT sale_id, act FROM api.all_clip_sale_events($1, $2)`, maxResults, resultOffset)
				Expect(queryErr).NotTo(HaveOccurred())

				Expect(actualClipEvents).To(ConsistOf(
					test_helpers.SaleEvent{
						SaleId: strconv.Itoa(saleId),
						Act:    "take",
					},
				))
			})
		})

		It("returns sale events from clippers that have different sale ids", func() {
			differentSaleId := rand.Int()

			clipKickLogTwo := test_data.CreateTestLog(headerOne.Id, db)

			clipKickEventTwo := test_data.ClipKickModel()
			clipKickEventTwo.ColumnValues[event.HeaderFK] = headerOne.Id
			clipKickEventTwo.ColumnValues[event.LogFK] = clipKickLogTwo.ID
			clipKickEventTwo.ColumnValues[event.AddressFK] = addressId
			clipKickEventTwo.ColumnValues[constants.SaleIDColumn] = strconv.Itoa(differentSaleId)
			clipKickErr := event.PersistModels([]event.InsertionModel{clipKickEventTwo}, db)
			Expect(clipKickErr).NotTo(HaveOccurred())

			var actualSaleEvents []test_helpers.SaleEvent
			queryErr := db.Select(&actualSaleEvents, `SELECT sale_id, act FROM api.all_clip_sale_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualSaleEvents).To(ConsistOf(
				test_helpers.SaleEvent{
					SaleId: strconv.Itoa(saleId),
					Act:    "kick"},
				test_helpers.SaleEvent{
					SaleId: clipKickEventTwo.ColumnValues["sale_id"].(string),
					Act:    "kick"},
			))
		})

		It("returns sale events from different kinds of clippers (clips with different contract addresses", func() {
			anotherAddressId, addressErr := repository.GetOrCreateAddress(db, anotherClipAddress)
			Expect(addressErr).NotTo(HaveOccurred())

			clipKickLog := test_data.CreateTestLog(headerOne.Id, db)
			clipKickEventTwo := test_data.ClipKickModel()
			clipKickEventTwo.ColumnValues[event.HeaderFK] = headerOne.Id
			clipKickEventTwo.ColumnValues[event.LogFK] = clipKickLog.ID
			clipKickEventTwo.ColumnValues[event.AddressFK] = anotherAddressId
			clipKickEventTwo.ColumnValues[constants.SaleIDColumn] = strconv.Itoa(saleId)
			clipKickErr := event.PersistModels([]event.InsertionModel{clipKickEventTwo}, db)
			Expect(clipKickErr).NotTo(HaveOccurred())

			var actualSaleEvents []test_helpers.SaleEvent
			queryErr := db.Select(&actualSaleEvents, `SELECT sale_id, act FROM api.all_clip_sale_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualSaleEvents).To(ConsistOf(
				test_helpers.SaleEvent{
					SaleId: strconv.Itoa(saleId),
					Act:    "kick"},
				test_helpers.SaleEvent{
					SaleId: clipKickEventTwo.ColumnValues["sale_id"].(string),
					Act:    "kick"},
			))
		})

		Describe("redo", func() {
			It("returns redo events from multiple blocks", func() {
				clipRedoHeaderOneLog := test_data.CreateTestLog(headerOne.Id, db)
				clipRedoHeaderOneErr := test_helpers.CreateRedo(test_helpers.RedoCreationInput{
					DB:              db,
					ContractAddress: contractAddress,
					SaleId:          saleId,
					RedoHeaderId:    headerOne.Id,
					RedoLogId:       clipRedoHeaderOneLog.ID,
				})
				Expect(clipRedoHeaderOneErr).NotTo(HaveOccurred())

				headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)

				clipRedoHeaderTwoLog := test_data.CreateTestLog(headerTwo.Id, db)
				clipRedoHeaderTwoErr := test_helpers.CreateRedo(test_helpers.RedoCreationInput{
					DB:              db,
					ContractAddress: contractAddress,
					SaleId:          saleId,
					RedoHeaderId:    headerTwo.Id,
					RedoLogId:       clipRedoHeaderTwoLog.ID,
				})
				Expect(clipRedoHeaderTwoErr).NotTo(HaveOccurred())

				var actualSaleEvents []test_helpers.SaleEvent
				queryErr := db.Select(&actualSaleEvents, `SELECT sale_id, act FROM api.all_clip_sale_events()`)
				Expect(queryErr).NotTo(HaveOccurred())

				Expect(actualSaleEvents).To(ConsistOf(
					test_helpers.SaleEvent{SaleId: strconv.Itoa(saleId), Act: "kick"},
					test_helpers.SaleEvent{SaleId: strconv.Itoa(saleId), Act: "redo"},
					test_helpers.SaleEvent{SaleId: strconv.Itoa(saleId), Act: "redo"},
				))
			})
		})

		Describe("take", func() {
			It("returns take events from multiple blocks", func() {
				clipTakeHeaderOneLog := test_data.CreateTestLog(headerOne.Id, db)
				clipTakeHeaderOneErr := test_helpers.CreateTake(test_helpers.TakeCreationInput{
					DB:              db,
					ContractAddress: contractAddress,
					SaleId:          saleId,
					TakeHeaderId:    headerOne.Id,
					TakeLogId:       clipTakeHeaderOneLog.ID,
				})
				Expect(clipTakeHeaderOneErr).NotTo(HaveOccurred())

				headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)

				clipTakeHeaderTwoLog := test_data.CreateTestLog(headerTwo.Id, db)
				clipTakeHeaderTwoErr := test_helpers.CreateTake(test_helpers.TakeCreationInput{
					DB:              db,
					ContractAddress: contractAddress,
					SaleId:          saleId,
					TakeHeaderId:    headerTwo.Id,
					TakeLogId:       clipTakeHeaderTwoLog.ID,
				})
				Expect(clipTakeHeaderTwoErr).NotTo(HaveOccurred())

				var actualSaleEvents []test_helpers.SaleEvent
				queryErr := db.Select(&actualSaleEvents, `SELECT sale_id, act FROM api.all_clip_sale_events()`)
				Expect(queryErr).NotTo(HaveOccurred())

				Expect(actualSaleEvents).To(ConsistOf(
					test_helpers.SaleEvent{SaleId: strconv.Itoa(saleId), Act: "kick"},
					test_helpers.SaleEvent{SaleId: strconv.Itoa(saleId), Act: "take"},
					test_helpers.SaleEvent{SaleId: strconv.Itoa(saleId), Act: "take"},
				))
			})
		})

		Describe("Yank", func() {
			It("returns Yank events from multiple blocks", func() {
				clipYankHeaderOneLog := test_data.CreateTestLog(headerOne.Id, db)
				clipYankHeaderOneErr := test_helpers.CreateClipYank(test_helpers.ClipYankCreationInput{
					DB:               db,
					ContractAddress:  contractAddress,
					SaleId:           saleId,
					ClipYankHeaderId: headerOne.Id,
					ClipYankLogId:    clipYankHeaderOneLog.ID,
				})
				Expect(clipYankHeaderOneErr).NotTo(HaveOccurred())

				headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)

				clipYankHeaderTwoLog := test_data.CreateTestLog(headerTwo.Id, db)
				clipYankHeaderTwoErr := test_helpers.CreateClipYank(test_helpers.ClipYankCreationInput{
					DB:               db,
					ContractAddress:  contractAddress,
					SaleId:           saleId,
					ClipYankHeaderId: headerTwo.Id,
					ClipYankLogId:    clipYankHeaderTwoLog.ID,
				})
				Expect(clipYankHeaderTwoErr).NotTo(HaveOccurred())

				var actualSaleEvents []test_helpers.SaleEvent
				queryErr := db.Select(&actualSaleEvents, `SELECT sale_id, act FROM api.all_clip_sale_events()`)
				Expect(queryErr).NotTo(HaveOccurred())

				Expect(actualSaleEvents).To(ConsistOf(
					test_helpers.SaleEvent{SaleId: strconv.Itoa(saleId), Act: "kick"},
					test_helpers.SaleEvent{SaleId: strconv.Itoa(saleId), Act: "yank"},
					test_helpers.SaleEvent{SaleId: strconv.Itoa(saleId), Act: "yank"},
				))
			})
		})
	})
})
