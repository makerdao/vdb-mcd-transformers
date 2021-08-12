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

var _ = Describe("All clips view", func() {
	var (
		headerRepo                   datastore.HeaderRepository
		contractAddress              = fakes.FakeAddress.Hex()
		anotherContractAddress       = fakes.AnotherFakeAddress.Hex()
		blockOne, timestampOne       int
		headerOne                    core.Header
		fakeSaleIdOne, fakeSaleIdTwo int
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)

		blockOne = rand.Int()
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepo)
		fakeSaleIdOne = rand.Int()
		fakeSaleIdTwo = fakeSaleIdOne + 1
	})

	It("gets the latest state of every sale on the clipper", func() {
		clipKickLogOne := test_data.CreateTestLog(headerOne.Id, db)

		var addressErr error
		addressId, addressErr := repository.GetOrCreateAddress(db, contractAddress)
		Expect(addressErr).NotTo(HaveOccurred())

		clipKickEventOne := test_data.ClipKickModel()
		clipKickEventOne.ColumnValues[event.HeaderFK] = headerOne.Id
		clipKickEventOne.ColumnValues[event.LogFK] = clipKickLogOne.ID
		clipKickEventOne.ColumnValues[event.AddressFK] = addressId
		clipKickEventOne.ColumnValues[constants.SaleIDColumn] = strconv.Itoa(fakeSaleIdOne)
		clipKickErr := event.PersistModels([]event.InsertionModel{clipKickEventOne}, db)
		Expect(clipKickErr).NotTo(HaveOccurred())

		clipTakeLogOne := test_data.CreateTestLog(headerOne.Id, db)
		clipTakeOneErr := test_helpers.CreateTake(test_helpers.TakeCreationInput{
			DB:              db,
			ContractAddress: contractAddress,
			SaleId:          fakeSaleIdOne,
			TakeHeaderId:    headerOne.Id,
			TakeLogId:       clipTakeLogOne.ID,
		})
		Expect(clipTakeOneErr).NotTo(HaveOccurred())

		clipStorageValuesOne := test_helpers.GetClipStorageValues(1, fakeSaleIdOne)
		test_helpers.CreateClip(db, headerOne, clipStorageValuesOne, test_helpers.GetClipMetadatas(strconv.Itoa(fakeSaleIdOne)), contractAddress)

		clipKickLogTwo := test_data.CreateTestLog(headerOne.Id, db)

		clipKickEventTwo := test_data.ClipKickModel()
		clipKickEventTwo.ColumnValues[event.HeaderFK] = headerOne.Id
		clipKickEventTwo.ColumnValues[event.LogFK] = clipKickLogTwo.ID
		clipKickEventTwo.ColumnValues[event.AddressFK] = addressId
		clipKickEventTwo.ColumnValues[constants.SaleIDColumn] = strconv.Itoa(fakeSaleIdTwo)
		clipKickErr = event.PersistModels([]event.InsertionModel{clipKickEventTwo}, db)
		Expect(clipKickErr).NotTo(HaveOccurred())

		clipTakeLogTwo := test_data.CreateTestLog(headerOne.Id, db)
		clipTakeTwoErr := test_helpers.CreateTake(test_helpers.TakeCreationInput{
			DB:              db,
			ContractAddress: contractAddress,
			SaleId:          fakeSaleIdTwo,
			TakeHeaderId:    headerOne.Id,
			TakeLogId:       clipTakeLogTwo.ID,
		})
		Expect(clipTakeTwoErr).NotTo(HaveOccurred())

		clipStorageValuesTwo := test_helpers.GetClipStorageValues(2, fakeSaleIdTwo)
		test_helpers.CreateClip(db, headerOne, clipStorageValuesTwo, test_helpers.GetClipMetadatas(strconv.Itoa(fakeSaleIdTwo)), contractAddress)

		var saleCount int
		countQueryErr := db.Get(&saleCount, `SELECT COUNT(*) FROM api.all_clips($1)`, contractAddress)
		Expect(countQueryErr).NotTo(HaveOccurred())
		Expect(saleCount).To(Equal(2))
	})

	It("ignores sales from other contracts", func() {
		clipKickLogOne := test_data.CreateTestLog(headerOne.Id, db)

		var addressErr error
		addressId, addressErr := repository.GetOrCreateAddress(db, contractAddress)
		Expect(addressErr).NotTo(HaveOccurred())

		clipKickEventOne := test_data.ClipKickModel()
		clipKickEventOne.ColumnValues[event.HeaderFK] = headerOne.Id
		clipKickEventOne.ColumnValues[event.LogFK] = clipKickLogOne.ID
		clipKickEventOne.ColumnValues[event.AddressFK] = addressId
		clipKickEventOne.ColumnValues[constants.SaleIDColumn] = strconv.Itoa(fakeSaleIdOne)
		clipKickErr := event.PersistModels([]event.InsertionModel{clipKickEventOne}, db)
		Expect(clipKickErr).NotTo(HaveOccurred())

		clipTakeLogOne := test_data.CreateTestLog(headerOne.Id, db)
		clipTakeOneErr := test_helpers.CreateTake(test_helpers.TakeCreationInput{
			DB:              db,
			ContractAddress: contractAddress,
			SaleId:          fakeSaleIdOne,
			TakeHeaderId:    headerOne.Id,
			TakeLogId:       clipTakeLogOne.ID,
		})
		Expect(clipTakeOneErr).NotTo(HaveOccurred())

		clipStorageValuesOne := test_helpers.GetClipStorageValues(1, fakeSaleIdOne)
		test_helpers.CreateClip(db, headerOne, clipStorageValuesOne, test_helpers.GetClipMetadatas(strconv.Itoa(fakeSaleIdOne)), contractAddress)

		addressIdTwo, addressErr := repository.GetOrCreateAddress(db, anotherContractAddress)
		Expect(addressErr).NotTo(HaveOccurred())

		clipKickLogTwo := test_data.CreateTestLog(headerOne.Id, db)

		clipKickEventTwo := test_data.ClipKickModel()
		clipKickEventTwo.ColumnValues[event.HeaderFK] = headerOne.Id
		clipKickEventTwo.ColumnValues[event.LogFK] = clipKickLogTwo.ID
		clipKickEventTwo.ColumnValues[event.AddressFK] = addressIdTwo
		clipKickEventTwo.ColumnValues[constants.SaleIDColumn] = strconv.Itoa(fakeSaleIdTwo)
		clipKickErr = event.PersistModels([]event.InsertionModel{clipKickEventTwo}, db)
		Expect(clipKickErr).NotTo(HaveOccurred())

		clipTakeLogTwo := test_data.CreateTestLog(headerOne.Id, db)
		clipTakeTwoErr := test_helpers.CreateTake(test_helpers.TakeCreationInput{
			DB:              db,
			ContractAddress: anotherContractAddress,
			SaleId:          fakeSaleIdTwo,
			TakeHeaderId:    headerOne.Id,
			TakeLogId:       clipTakeLogTwo.ID,
		})
		Expect(clipTakeTwoErr).NotTo(HaveOccurred())

		clipStorageValuesTwo := test_helpers.GetClipStorageValues(2, fakeSaleIdTwo)
		test_helpers.CreateClip(db, headerOne, clipStorageValuesTwo, test_helpers.GetClipMetadatas(strconv.Itoa(fakeSaleIdTwo)), anotherContractAddress)

		var saleCount int
		countQueryErr := db.Get(&saleCount, `SELECT COUNT(*) FROM api.all_clips($1)`, contractAddress)
		Expect(countQueryErr).NotTo(HaveOccurred())
		Expect(saleCount).To(Equal(1))
	})

	It("gets the right sales when there are the same ids on different contracts", func() {
		clipKickLogOne := test_data.CreateTestLog(headerOne.Id, db)

		var addressErr error
		addressId, addressErr := repository.GetOrCreateAddress(db, contractAddress)
		Expect(addressErr).NotTo(HaveOccurred())

		clipKickEventOne := test_data.ClipKickModel()
		clipKickEventOne.ColumnValues[event.HeaderFK] = headerOne.Id
		clipKickEventOne.ColumnValues[event.LogFK] = clipKickLogOne.ID
		clipKickEventOne.ColumnValues[event.AddressFK] = addressId
		clipKickEventOne.ColumnValues[constants.SaleIDColumn] = strconv.Itoa(fakeSaleIdOne)
		clipKickErr := event.PersistModels([]event.InsertionModel{clipKickEventOne}, db)
		Expect(clipKickErr).NotTo(HaveOccurred())

		clipTakeLogOne := test_data.CreateTestLog(headerOne.Id, db)
		clipTakeOneErr := test_helpers.CreateTake(test_helpers.TakeCreationInput{
			DB:              db,
			ContractAddress: contractAddress,
			SaleId:          fakeSaleIdOne,
			TakeHeaderId:    headerOne.Id,
			TakeLogId:       clipTakeLogOne.ID,
		})
		Expect(clipTakeOneErr).NotTo(HaveOccurred())

		clipStorageValuesOne := test_helpers.GetClipStorageValues(1, fakeSaleIdOne)
		test_helpers.CreateClip(db, headerOne, clipStorageValuesOne, test_helpers.GetClipMetadatas(strconv.Itoa(fakeSaleIdOne)), contractAddress)

		clipKickLogTwo := test_data.CreateTestLog(headerOne.Id, db)

		addressIdTwo, addressErr := repository.GetOrCreateAddress(db, anotherContractAddress)
		Expect(addressErr).NotTo(HaveOccurred())

		clipKickEventTwo := test_data.ClipKickModel()
		clipKickEventTwo.ColumnValues[event.HeaderFK] = headerOne.Id
		clipKickEventTwo.ColumnValues[event.LogFK] = clipKickLogTwo.ID
		clipKickEventTwo.ColumnValues[event.AddressFK] = addressIdTwo
		clipKickEventTwo.ColumnValues[constants.SaleIDColumn] = strconv.Itoa(fakeSaleIdOne)
		clipKickErr = event.PersistModels([]event.InsertionModel{clipKickEventTwo}, db)
		Expect(clipKickErr).NotTo(HaveOccurred())

		clipTakeLogTwo := test_data.CreateTestLog(headerOne.Id, db)
		clipTakeTwoErr := test_helpers.CreateTake(test_helpers.TakeCreationInput{
			DB:              db,
			ContractAddress: anotherContractAddress,
			SaleId:          fakeSaleIdOne,
			TakeHeaderId:    headerOne.Id,
			TakeLogId:       clipTakeLogTwo.ID,
		})
		Expect(clipTakeTwoErr).NotTo(HaveOccurred())

		clipStorageValuesTwo := test_helpers.GetClipStorageValues(1, fakeSaleIdOne)
		test_helpers.CreateClip(db, headerOne, clipStorageValuesTwo, test_helpers.GetClipMetadatas(strconv.Itoa(fakeSaleIdTwo)), anotherContractAddress)

		var saleCount int
		countQueryErr := db.Get(&saleCount, `SELECT COUNT(*) FROM api.all_clips($1)`, contractAddress)
		Expect(countQueryErr).NotTo(HaveOccurred())
		Expect(saleCount).To(Equal(1))
	})
})
