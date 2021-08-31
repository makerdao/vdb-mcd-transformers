package queries

import (
	"database/sql"
	"math/rand"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
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

var _ = Describe("Single clip view", func() {
	var (
		blockOne, timestampOne int
		addressId              int64
		urnID                  int64
		ilkEthIdentifier       = "ETH-A"
		hexEthIlk              = "0x4554482d41"
		fakeSaleId             = rand.Int()
		headerOne              core.Header
		headerRepo             datastore.HeaderRepository
		dogBarkUrnAddress      = common.HexToAddress(test_data.DogBarkEventLog.Log.Topics[2].Hex()).Hex()
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)

		blockOne = rand.Int()
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepo)

		clipKickLog := test_data.CreateTestLog(headerOne.Id, db)
		dogBarkLogOne := test_data.CreateTestLog(headerOne.Id, db)

		urnID, _ = shared.GetOrCreateUrn(dogBarkUrnAddress, hexEthIlk, db)

		dogBarkEventOne := test_data.DogBarkModel()
		dogBarkEventOne.ColumnValues[event.HeaderFK] = headerOne.Id
		dogBarkEventOne.ColumnValues[event.LogFK] = dogBarkLogOne.ID
		dogBarkEventOne.ColumnValues[constants.SaleIDColumn] = strconv.Itoa(fakeSaleId)
		test_data.AssignIlkID(dogBarkEventOne, ilkEthIdentifier, db)
		test_data.AssignUrnID(dogBarkEventOne, db)
		test_data.AssignAddressID(test_data.DogBarkEventLog, dogBarkEventOne, db)
		test_data.AssignClip(test_data.ClipAddress, dogBarkEventOne, db)

		dogBarkErr := event.PersistModels([]event.InsertionModel{dogBarkEventOne}, db)
		Expect(dogBarkErr).NotTo(HaveOccurred())

		var addressErr error
		addressId, addressErr = repository.GetOrCreateAddress(db, test_data.ClipAddress)
		Expect(addressErr).NotTo(HaveOccurred())

		clipKickEvent := test_data.ClipKickModel()
		clipKickEvent.ColumnValues[event.HeaderFK] = headerOne.Id
		clipKickEvent.ColumnValues[event.LogFK] = clipKickLog.ID
		clipKickEvent.ColumnValues[event.AddressFK] = addressId
		clipKickEvent.ColumnValues[constants.SaleIDColumn] = strconv.Itoa(fakeSaleId)
		clipKickErr := event.PersistModels([]event.InsertionModel{clipKickEvent}, db)
		Expect(clipKickErr).NotTo(HaveOccurred())
	})

	Context("get_clip_with_address", func() {
		It("gets only the specified clip", func() {
			clipStorageValuesOne := test_helpers.GetClipStorageValues(1, fakeSaleId)
			test_helpers.CreateClip(db, headerOne, clipStorageValuesOne, test_helpers.GetClipMetadatas(strconv.Itoa(fakeSaleId)), test_data.ClipAddress)

			var actualSale test_helpers.ClipSale
			queryErr := db.Get(&actualSale, `SELECT block_height, sale_id, ilk_id, urn_id, pos, tab, lot, usr, tic, "top", created, updated FROM api.get_clip_with_address($1, $2, $3)`,
				fakeSaleId, test_data.ClipAddress, blockOne)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(actualSale.BlockHeight).To(Equal(strconv.Itoa(blockOne)))
			Expect(actualSale.SaleId).To(Equal(strconv.Itoa(fakeSaleId)))
			Expect(actualSale.UrnId).To(Equal(strconv.Itoa(int(urnID))))
			Expect(actualSale.Created).To(Equal(sql.NullString{String: time.Unix(int64(timestampOne), 0).UTC().Format(time.RFC3339), Valid: true}))
			Expect(actualSale.Updated).To(Equal(sql.NullString{String: time.Unix(int64(timestampOne), 0).UTC().Format(time.RFC3339), Valid: true}))
		})
		It("gets the right clip when the salesIDs are the same for different clippers in the same block", func() {
			anotherClipAddress := "0x" + fakes.RandomString(38)
			addressId, addressErr := repository.GetOrCreateAddress(db, anotherClipAddress)
			Expect(addressErr).NotTo(HaveOccurred())

			clipStorageValuesOne := test_helpers.GetClipStorageValues(1, fakeSaleId)
			test_helpers.CreateClip(db, headerOne, clipStorageValuesOne, test_helpers.GetClipMetadatas(strconv.Itoa(fakeSaleId)), test_data.ClipAddress)
			test_helpers.CreateClip(db, headerOne, clipStorageValuesOne, test_helpers.GetClipMetadatas(strconv.Itoa(fakeSaleId)), anotherClipAddress)

			clipKickLogTwo := test_data.CreateTestLog(headerOne.Id, db)
			clipKickEventTwo := test_data.ClipKickModel()
			clipKickEventTwo.ColumnValues[event.HeaderFK] = headerOne.Id
			clipKickEventTwo.ColumnValues[event.LogFK] = clipKickLogTwo.ID
			clipKickEventTwo.ColumnValues[event.AddressFK] = addressId
			clipKickEventTwo.ColumnValues[constants.SaleIDColumn] = strconv.Itoa(fakeSaleId)
			clipKickTwoErr := event.PersistModels([]event.InsertionModel{clipKickEventTwo}, db)
			Expect(clipKickTwoErr).NotTo(HaveOccurred())

			var actualSale test_helpers.ClipSale
			queryErr := db.Get(&actualSale, `SELECT block_height, sale_id, ilk_id, urn_id, pos, tab, lot, usr, tic, "top", created, updated FROM api.get_clip_with_address($1, $2, $3)`,
				fakeSaleId, test_data.ClipAddress, blockOne)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(actualSale.BlockHeight).To(Equal(strconv.Itoa(blockOne)))
			Expect(actualSale.SaleId).To(Equal(strconv.Itoa(fakeSaleId)))
			Expect(actualSale.UrnId).To(Equal(strconv.Itoa(int(urnID))))
			Expect(actualSale.Created).To(Equal(sql.NullString{String: time.Unix(int64(timestampOne), 0).UTC().Format(time.RFC3339), Valid: true}))
			Expect(actualSale.Updated).To(Equal(sql.NullString{String: time.Unix(int64(timestampOne), 0).UTC().Format(time.RFC3339), Valid: true}))
		})
	})
})
