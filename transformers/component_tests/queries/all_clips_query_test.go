package queries

import (
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
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
		headerRepo                 datastore.HeaderRepository
		anotherContractAddress     = fakes.AnotherFakeAddress.Hex()
		ilkBatIdentifier           = "BAT-A"
		hexBatIlk                  = "0x4241542d41"
		ilkEthIdentifier           = "ETH-A"
		hexEthIlk                  = "0x4554482d41"
		blockOne, blockTwo         int
		timestampOne, timestampTwo int
		headerOne                  core.Header
		fakeSaleIdOne              int
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)

		blockOne = rand.Int()
		blockTwo = blockOne + 1
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepo)
		fakeSaleIdOne = rand.Int()

		dogBarkLogOne := test_data.CreateTestLog(headerOne.Id, db)
		_, _ = shared.GetOrCreateUrn(test_data.UrnAddress, hexEthIlk, db)
		dogBarkEventOne := test_data.DogBarkModel()
		dogBarkEventOne.ColumnValues[event.HeaderFK] = headerOne.Id
		dogBarkEventOne.ColumnValues[event.LogFK] = dogBarkLogOne.ID
		dogBarkEventOne.ColumnValues[constants.SaleIDColumn] = strconv.Itoa(fakeSaleIdOne)
		test_data.AssignIlkID(dogBarkEventOne, ilkEthIdentifier, db)
		test_data.AssignUrnID(dogBarkEventOne, db)
		test_data.AssignAddressID(test_data.DogBarkEventLog, dogBarkEventOne, db)
		test_data.AssignClip(test_data.ClipAddress, dogBarkEventOne, db)

		dogBarkErr := event.PersistModels([]event.InsertionModel{dogBarkEventOne}, db)
		Expect(dogBarkErr).NotTo(HaveOccurred())
	})

	It("gets the state of a single sale on the clipper", func() {
		clipKickLogOne := test_data.CreateTestLog(headerOne.Id, db)

		var addressErr error
		addressId, addressErr := repository.GetOrCreateAddress(db, test_data.ClipAddress)
		Expect(addressErr).NotTo(HaveOccurred())

		clipKickEventOne := test_data.ClipKickModel()
		clipKickEventOne.ColumnValues[event.HeaderFK] = headerOne.Id
		clipKickEventOne.ColumnValues[event.LogFK] = clipKickLogOne.ID
		clipKickEventOne.ColumnValues[event.AddressFK] = addressId
		clipKickEventOne.ColumnValues[constants.SaleIDColumn] = strconv.Itoa(fakeSaleIdOne)
		clipKickErr := event.PersistModels([]event.InsertionModel{clipKickEventOne}, db)
		Expect(clipKickErr).NotTo(HaveOccurred())

		clipStorageValuesOne := test_helpers.GetClipStorageValues(1, fakeSaleIdOne)
		test_helpers.CreateClip(db, headerOne, clipStorageValuesOne, test_helpers.GetClipMetadatas(strconv.Itoa(fakeSaleIdOne)), test_data.ClipAddress)

		var saleCount int
		countQueryErr := db.Get(&saleCount, `SELECT COUNT(*) FROM api.all_clips($1)`, ilkEthIdentifier)
		Expect(countQueryErr).NotTo(HaveOccurred())
		Expect(saleCount).To(Equal(1))
	})

	It("gets the latest state of multiple sales on the clipper", func() {
		clipKickLogOne := test_data.CreateTestLog(headerOne.Id, db)

		var addressErr error
		addressId, addressErr := repository.GetOrCreateAddress(db, test_data.ClipAddress)
		Expect(addressErr).NotTo(HaveOccurred())

		clipKickEventOne := test_data.ClipKickModel()
		clipKickEventOne.ColumnValues[event.HeaderFK] = headerOne.Id
		clipKickEventOne.ColumnValues[event.LogFK] = clipKickLogOne.ID
		clipKickEventOne.ColumnValues[event.AddressFK] = addressId
		clipKickEventOne.ColumnValues[constants.SaleIDColumn] = strconv.Itoa(fakeSaleIdOne)
		clipKickErr := event.PersistModels([]event.InsertionModel{clipKickEventOne}, db)
		Expect(clipKickErr).NotTo(HaveOccurred())

		clipStorageValuesOne := test_helpers.GetClipStorageValues(1, fakeSaleIdOne)
		test_helpers.CreateClip(db, headerOne, clipStorageValuesOne, test_helpers.GetClipMetadatas(strconv.Itoa(fakeSaleIdOne)), test_data.ClipAddress)

		headerTwo := createHeader(blockTwo, timestampTwo, headerRepo)
		clipTakeLogOne := test_data.CreateTestLog(headerTwo.Id, db)
		clipTakeOneErr := test_helpers.CreateTake(test_helpers.TakeCreationInput{
			DB:              db,
			ContractAddress: test_data.ClipAddress,
			SaleId:          fakeSaleIdOne,
			TakeHeaderId:    headerTwo.Id,
			TakeLogId:       clipTakeLogOne.ID,
		})
		Expect(clipTakeOneErr).NotTo(HaveOccurred())

		clipStorageValuesTwo := test_helpers.GetClipStorageValues(2, fakeSaleIdOne)
		test_helpers.CreateClip(db, headerTwo, clipStorageValuesTwo, test_helpers.GetClipMetadatas(strconv.Itoa(fakeSaleIdOne)), test_data.ClipAddress)

		var saleCount, blockHeight int
		countQueryErr := db.Get(&saleCount, `SELECT COUNT(*) FROM api.all_clips($1)`, ilkEthIdentifier)
		Expect(countQueryErr).NotTo(HaveOccurred())
		Expect(saleCount).To(Equal(1))

		blockHeightQueryErr := db.Get(&blockHeight, `SELECT block_height FROM api.all_clips($1)`, ilkEthIdentifier)
		Expect(blockHeightQueryErr).NotTo(HaveOccurred())
		Expect(blockHeight).To(Equal(blockTwo))
	})

	XIt("gets the right sale when there are the same ids on different contracts/ilks", func() {
		clipKickLogOne := test_data.CreateTestLog(headerOne.Id, db)

		var addressErr error
		addressId, addressErr := repository.GetOrCreateAddress(db, test_data.ClipAddress)
		Expect(addressErr).NotTo(HaveOccurred())

		clipKickEventOne := test_data.ClipKickModel()
		clipKickEventOne.ColumnValues[event.HeaderFK] = headerOne.Id
		clipKickEventOne.ColumnValues[event.LogFK] = clipKickLogOne.ID
		clipKickEventOne.ColumnValues[event.AddressFK] = addressId
		clipKickEventOne.ColumnValues[constants.SaleIDColumn] = strconv.Itoa(fakeSaleIdOne)
		clipKickErr := event.PersistModels([]event.InsertionModel{clipKickEventOne}, db)
		Expect(clipKickErr).NotTo(HaveOccurred())

		clipStorageValuesOne := test_helpers.GetClipStorageValues(1, fakeSaleIdOne)
		test_helpers.CreateClip(db, headerOne, clipStorageValuesOne, test_helpers.GetClipMetadatas(strconv.Itoa(fakeSaleIdOne)), test_data.ClipAddress)

		dogBarkLogTwo := test_data.CreateTestLog(headerOne.Id, db)
		_, _ = shared.GetOrCreateUrn(test_data.UrnAddress, hexBatIlk, db)
		dogBarkEventTwo := test_data.DogBarkModel()
		dogBarkEventTwo.ColumnValues[event.HeaderFK] = headerOne.Id
		dogBarkEventTwo.ColumnValues[event.LogFK] = dogBarkLogTwo.ID
		test_data.AssignIlkID(dogBarkEventTwo, ilkBatIdentifier, db)
		test_data.AssignUrnID(dogBarkEventTwo, db)
		test_data.AssignAddressID(test_data.DogBarkEventLog, dogBarkEventTwo, db)
		test_data.AssignClip(anotherContractAddress, dogBarkEventTwo, db)

		dogBarkErr := event.PersistModels([]event.InsertionModel{dogBarkEventTwo}, db)
		Expect(dogBarkErr).NotTo(HaveOccurred())
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

		clipStorageValuesTwo := test_helpers.GetClipStorageValues(1, fakeSaleIdOne)
		test_helpers.CreateClip(db, headerOne, clipStorageValuesTwo, test_helpers.GetClipMetadatas(strconv.Itoa(fakeSaleIdOne)), anotherContractAddress)

		var saleCount int
		countQueryErr := db.Get(&saleCount, `SELECT COUNT(*) FROM api.all_clips($1)`, ilkEthIdentifier)
		Expect(countQueryErr).NotTo(HaveOccurred())
		Expect(saleCount).To(Equal(1))
	})
})
