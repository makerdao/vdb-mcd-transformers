package queries

import (
	"math/rand"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/vat"
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
		dogBarkUrnAddress      = common.HexToAddress(test_data.DogBarkEventLog.Log.Topics[2].Hex()).Hex()
	)

	BeforeEach(func() {
		fakeSaleId = rand.Int()
		blockOne = rand.Int()
		timestampOne = int(rand.Int31())

		test_config.CleanTestDB(db)

		headerRepository = repositories.NewHeaderRepository(db)
		headerOne = createHeader(blockOne, timestampOne, headerRepository)
		dogBarkLogOne := test_data.CreateTestLog(headerOne.Id, db)

		_, _ = shared.GetOrCreateUrn(dogBarkUrnAddress, test_helpers.FakeIlk.Hex, db)
		ilkID, _ := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
		dogBarkEventOne := test_data.DogBarkModel()
		dogBarkEventOne.ColumnValues[event.HeaderFK] = headerOne.Id
		dogBarkEventOne.ColumnValues[event.LogFK] = dogBarkLogOne.ID
		dogBarkEventOne.ColumnValues[constants.SaleIDColumn] = strconv.Itoa(fakeSaleId)
		dogBarkEventOne.ColumnValues[constants.IlkColumn] = strconv.Itoa(int(ilkID))
		test_data.AssignIlkID(dogBarkEventOne, test_helpers.FakeIlk.Identifier, db)
		test_data.AssignUrnID(dogBarkEventOne, db)
		test_data.AssignAddressID(test_data.DogBarkEventLog, dogBarkEventOne, db)
		test_data.AssignClip(contractAddress, dogBarkEventOne, db)

		dogBarkErr := event.PersistModels([]event.InsertionModel{dogBarkEventOne}, db)
		Expect(dogBarkErr).NotTo(HaveOccurred())

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

	Describe("clip_sale_snapshot_ilk", func() {
		It("returns ilk_snapshot for a clip_sale_snapshot", func() {
			ilkValues := test_helpers.GetIlkValues(0)
			test_helpers.CreateIlk(db, headerOne, ilkValues, test_helpers.FakeIlkVatMetadatas, test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas, test_helpers.FakeIlkSpotMetadatas)

			expectedIlk := test_helpers.IlkSnapshotFromValues(test_helpers.FakeIlk.Hex, headerOne.Timestamp, headerOne.Timestamp, ilkValues)

			var result test_helpers.IlkSnapshot
			getIlkErr := db.Get(&result, `
				SELECT ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, dunk, created, updated
				FROM api.clip_sale_snapshot_ilk(
					(SELECT (block_height, sale_id, ilk_id, urn_id, pos, tab, lot, usr, tic, "top", created, updated, clip_address)::api.clip_sale_snapshot
					 FROM api.get_clip_with_address($1, $2, $3))
			)`, fakeSaleId, contractAddress, blockOne)

			Expect(getIlkErr).NotTo(HaveOccurred())
			Expect(result).To(Equal(expectedIlk))
		})
	})

	Describe("clip_sale_snapshot_urn", func() {
		It("returns urn_snapshot for a clip_sale_snapshot", func() {
			urnSetupData := test_helpers.GetUrnSetupData()
			urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, dogBarkUrnAddress)
			vatRepository := vat.StorageRepository{}
			vatRepository.SetDB(db)
			test_helpers.CreateUrn(db, urnSetupData, headerOne, urnMetadata, vatRepository)

			var actualUrn test_helpers.UrnState
			getUrnErr := db.Get(&actualUrn, `
				SELECT urn_identifier, ilk_identifier
				FROM api.clip_sale_snapshot_urn(
					(SELECT (block_height, sale_id, ilk_id, urn_id, pos, tab, lot, usr, tic, "top", created, updated, clip_address)::api.clip_sale_snapshot
					FROM api.get_clip_with_address($1, $2, $3))
			)`, fakeSaleId, contractAddress, blockOne)

			Expect(getUrnErr).NotTo(HaveOccurred())

			expectedUrn := test_helpers.UrnState{
				UrnIdentifier: dogBarkUrnAddress,
				IlkIdentifier: test_helpers.FakeIlk.Identifier,
			}

			test_helpers.AssertUrn(actualUrn, expectedUrn)
		})
	})

})
