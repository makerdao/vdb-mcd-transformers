package spot

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	storage2 "github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/spot"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
)

var _ = Describe("Executing the transformer", func() {
	var (
		db          *postgres.DB
		err         error
		ilkID       int64
		mappings    = spot.SpotMappings{StorageRepository: &storage2.MakerStorageRepository{}}
		repository  = spot.SpotStorageRepository{}
		transformer = storage.Transformer{
			Address:    common.Address{},
			Mappings:   &mappings,
			Repository: &repository,
		}
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		transformer.NewTransformer(db)
		ilk := "0x434f4c352d410000000000000000000000000000000000000000000000000000"
		ilkID, err = shared.GetOrCreateIlk(ilk, db)
		Expect(err).NotTo(HaveOccurred())
	})

	It("reads in a Spot Vat storage diff row and persists it", func() {
		blockNumber := 11257169
		spotVatRow := utils.StorageDiffRow{
			Contract:     common.HexToAddress("25a008bf942ce6d5b362f91ed7ae3e4104286a12"),
			BlockHeight:  blockNumber,
			BlockHash:    common.HexToHash("d39fe1598fad020726983eeb76bdca943d2757dc3be91864ab00f2cb0931628a"),
			StorageKey:   common.HexToHash("0000000000000000000000000000000000000000000000000000000000000002"),
			StorageValue: common.HexToHash("00000000000000000000000057aa8b02f5d3e28371fedcf672c8668869f9aac7"),
		}
		err := transformer.Execute(spotVatRow)
		Expect(err).NotTo(HaveOccurred())

		var vatResult test_helpers.VariableRes
		err = db.Get(&vatResult, `SELECT block_number, block_hash, vat AS value FROM maker.spot_vat`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(vatResult, blockNumber, "0xd39fe1598fad020726983eeb76bdca943d2757dc3be91864ab00f2cb0931628a", "0x57aA8B02F5D3E28371FEdCf672C8668869f9AAC7")
	})

	It("reads in a Spot Par storage diff row and persists it", func() {
		blockNumber := 11257169
		spotParRow := utils.StorageDiffRow{
			Contract:     common.HexToAddress("25a008bf942ce6d5b362f91ed7ae3e4104286a12"),
			BlockHeight:  blockNumber,
			BlockHash:    common.HexToHash("d39fe1598fad020726983eeb76bdca943d2757dc3be91864ab00f2cb0931628a"),
			StorageKey:   common.HexToHash("0000000000000000000000000000000000000000000000000000000000000003"),
			StorageValue: common.HexToHash("0000000000000000000000000000000000000000033b2e3c9fd0803ce8000000"),
		}
		err := transformer.Execute(spotParRow)
		Expect(err).NotTo(HaveOccurred())

		var parResult test_helpers.VariableRes
		err = db.Get(&parResult, `SELECT block_number, block_hash, par AS value FROM maker.spot_par`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(parResult, blockNumber, "0xd39fe1598fad020726983eeb76bdca943d2757dc3be91864ab00f2cb0931628a", "1000000000000000000000000000")
	})

	It("reads in a Spot Ilk Pip storage diff row and persists it", func() {
		blockNumber := 11257255
		spotIlkPipRow := utils.StorageDiffRow{
			Contract:     common.HexToAddress("a57d4123c8a80ac410e924df9d5e47765ffd1375"),
			BlockHeight:  blockNumber,
			BlockHash:    common.HexToHash("8c57727c0e057bd603e27304762c7144df161fc96990a573fddf23916b64c7df"),
			StorageKey:   common.HexToHash("1730ac98111482efebd8acadb14d7fa301298e0d95bf3c34c3378ef524670bc6"),
			StorageValue: common.HexToHash("000000000000000000000000a53e6efb4cbed841eace02220498860905e94998"),
		}
		err := transformer.Execute(spotIlkPipRow)
		Expect(err).NotTo(HaveOccurred())

		var ilkPipResult test_helpers.MappingRes
		err = db.Get(&ilkPipResult, `SELECT block_number, block_hash, ilk_id AS key, pip AS value FROM maker.spot_ilk_pip`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertMapping(ilkPipResult, blockNumber, "0x8c57727c0e057bd603e27304762c7144df161fc96990a573fddf23916b64c7df", strconv.FormatInt(ilkID, 10), "0xA53e6EFB4cBeD841Eace02220498860905E94998")
	})

	It("reads in a Spot Ilk Mat storage diff row and persists it", func() {
		blockNumber := 11257407
		spotIlkMatRow := utils.StorageDiffRow{
			Contract:     common.HexToAddress("a57d4123c8a80ac410e924df9d5e47765ffd1375"),
			BlockHeight:  blockNumber,
			BlockHash:    common.HexToHash("d95e007739c8451f5e7b73fa1139b450aae37a6bf7735bcdb1f858cd32873726"),
			StorageKey:   common.HexToHash("1730ac98111482efebd8acadb14d7fa301298e0d95bf3c34c3378ef524670bc7"),
			StorageValue: common.HexToHash("000000000000000000000000000000000000000006765c793fa10079d0000000"),
		}
		err := transformer.Execute(spotIlkMatRow)
		Expect(err).NotTo(HaveOccurred())

		var ilkRhoResult test_helpers.MappingRes
		err = db.Get(&ilkRhoResult, `SELECT block_number, block_hash, ilk_id AS key, mat AS value FROM maker.spot_ilk_mat`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertMapping(ilkRhoResult, blockNumber, "0xd95e007739c8451f5e7b73fa1139b450aae37a6bf7735bcdb1f858cd32873726", strconv.FormatInt(ilkID, 10), "2000000000000000000000000000")
	})
})
