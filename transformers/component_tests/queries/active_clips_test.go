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
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip"
	storage_test_helpers "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Active clip sales query", func() {
	var (
		headerOne              core.Header
		headerRepository       datastore.HeaderRepository
		contractAddress        = test_data.ClipLinkAV130Address()
		addressId              int64
		fakeSaleId             int
		blockOne, timestampOne int
		dogBarkUrnAddress      = common.HexToAddress(test_data.DogBarkEventLog.Log.Topics[2].Hex()).Hex()
		transformer            storage.Transformer
		clipStorageValues      map[string]interface{}
	)

	BeforeEach(func() {
		fakeSaleId = 50
		blockOne = rand.Int()
		timestampOne = int(rand.Int31())

		test_config.CleanTestDB(db)

		headerRepository = repositories.NewHeaderRepository(db)
		headerOne = createHeader(blockOne, timestampOne, headerRepository)

		storageKeysLookup := storage.NewKeysLookup(clip.NewKeysLoader(&mcdStorage.MakerStorageRepository{}, contractAddress))
		clipRepository := clip.StorageRepository{ContractAddress: contractAddress}
		transformer = storage.Transformer{
			Address:           common.HexToAddress(contractAddress),
			StorageKeysLookup: storageKeysLookup,
			Repository:        &clipRepository,
		}
		transformer.NewTransformer(db)

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

		clipStorageValues = test_helpers.GetClipStorageValues(1, fakeSaleId)
		test_helpers.CreateClip(db, headerOne, clipStorageValues, test_helpers.GetClipMetadatas(strconv.Itoa(fakeSaleId)), contractAddress)
	})

	Describe("active clips", func() {
		It("returns clips that are currently active", func() {
			key := common.HexToHash("000000000000000000000000000000000000000000000000000000000000000b")
			value := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001")
			clipActiveLengthDiff := storage_test_helpers.CreateDiffRecord(db, headerOne, common.HexToAddress(contractAddress), key, value)

			err := transformer.Execute(clipActiveLengthDiff)
			Expect(err).NotTo(HaveOccurred())

			key2 := common.HexToHash("0175b7a638427703f0dbe7bb9bbf987a2551717b34e79f33b5b1008d1fa01db9")
			value2 := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000032")
			clipActiveSaleIDDiff := storage_test_helpers.CreateDiffRecord(db, headerOne, common.HexToAddress(contractAddress), key2, value2)

			barkUrnId, urnErr := shared.GetOrCreateUrn(dogBarkUrnAddress, test_helpers.FakeIlk.Hex, db)
			Expect(urnErr).NotTo(HaveOccurred())

			ilkId, ilkErr := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
			Expect(ilkErr).NotTo(HaveOccurred())

			activeErr := transformer.Execute(clipActiveSaleIDDiff)
			Expect(activeErr).NotTo(HaveOccurred())
			clipSaleOne := test_helpers.ClipSale{
				BlockHeight: strconv.Itoa(blockOne),
				SaleId:      strconv.Itoa(fakeSaleId),
				IlkId:       strconv.Itoa(int(ilkId)),
				UrnId:       strconv.Itoa(int(barkUrnId)),
				Pos:         clipStorageValues[clip.SalePos].(string),
				Tab:         clipStorageValues[clip.SaleTab].(string),
				Lot:         clipStorageValues[clip.SaleLot].(string),
				Usr:         clipStorageValues[clip.Packed].(map[int]string)[0],
				Tic:         clipStorageValues[clip.Packed].(map[int]string)[1],
				Top:         clipStorageValues[clip.SaleTop].(string),
				Created:     sql.NullString{String: time.Unix(int64(timestampOne), 0).UTC().Format(time.RFC3339), Valid: true},
				Updated:     sql.NullString{String: time.Unix(int64(timestampOne), 0).UTC().Format(time.RFC3339), Valid: true},
				ClipAddress: contractAddress,
			}

			var actualSales []test_helpers.ClipSale
			saleQueryErr := db.Select(&actualSales, `SELECT block_height, sale_id, ilk_id, urn_id, pos, tab, lot, usr, tic, "top", created, updated, clip_address from api.active_clips($1)`, test_helpers.FakeIlk.Identifier)
			Expect(saleQueryErr).NotTo(HaveOccurred())
			Expect(actualSales).To(ContainElement(clipSaleOne))
		})
	})
})
