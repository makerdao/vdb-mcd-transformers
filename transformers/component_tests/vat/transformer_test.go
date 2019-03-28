package vat

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_tune"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	storage2 "github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vat"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("Executing the transformer", func() {
	var (
		db          *postgres.DB
		mappings    = vat.VatMappings{StorageRepository: &storage2.MakerStorageRepository{}}
		repository  = vat.VatStorageRepository{}
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
		ilk := "4554480000000000000000000000000000000000000000000000000000000000"
		_, err := shared.GetOrCreateIlk(ilk, db)
		Expect(err).NotTo(HaveOccurred())

		h := repositories.NewHeaderRepository(db)
		i, _ := h.CreateOrUpdateHeader(fakes.FakeHeader)
		v := vat_tune.VatTuneRepository{}
		v.SetDB(db)

		m := vat_tune.VatTuneModel{
			Ilk:              "4554480000000000000000000000000000000000000000000000000000000000",
			Urn:              "84271a423a68d9a3904fe8107185d9ff58a64974000000000000000000000037",
			V:                "v",
			W:                "w",
			Dink:             test_data.VatTuneModel.Dink,
			Dart:             test_data.VatTuneModel.Dart,
			TransactionIndex: 0,
			LogIndex:         0,
			Raw:              test_data.VatTuneModel.Raw,
		}

		v.Create(i, []interface{}{m})
	})

	It("reads in a Vat.rate storage diff row and persists it", func() {
		blockNumber := 10501138
		vatLineRow := utils.StorageDiffRow{
			Contract:     common.HexToAddress("67fd6c3575fc2dbe2cb596bd3bebc9edb5571fa1"),
			BlockHeight:  blockNumber,
			BlockHash:    common.HexToHash("3f58749d3956984c2b03a84d5c02105a06efa1ad048d8aa97cf8f59aafa8f08b"),
			StorageKey:   common.HexToHash("d9402e47e68a5154d8f9cdc9d4ffcf4a300546026372c3b04224fbe62656c602"),
			StorageValue: common.HexToHash("0000000000000000000000000000000000000000033b2e3c9fd0803ce8000000"),
		}
		err := transformer.Execute(vatLineRow)
		Expect(err).NotTo(HaveOccurred())

		var spotResult test_helpers.VariableRes
		err = db.Get(&spotResult, `SELECT block_number, block_hash, rate AS value FROM maker.vat_ilk_rate`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(spotResult, blockNumber, "0x3f58749d3956984c2b03a84d5c02105a06efa1ad048d8aa97cf8f59aafa8f08b", "1000000000000000000000000000")
	})

	It("reads in a Vat.spot storage diff row and persists it", func() {
		blockNumber := 10501138
		vatLineRow := utils.StorageDiffRow{
			Contract:     common.HexToAddress("67fd6c3575fc2dbe2cb596bd3bebc9edb5571fa1"),
			BlockHeight:  blockNumber,
			BlockHash:    common.HexToHash("3f58749d3956984c2b03a84d5c02105a06efa1ad048d8aa97cf8f59aafa8f08b"),
			StorageKey:   common.HexToHash("d9402e47e68a5154d8f9cdc9d4ffcf4a300546026372c3b04224fbe62656c603"),
			StorageValue: common.HexToHash("0000000000000000000000000000000000000001ba9f5611e5769eabb9000000"),
		}
		err := transformer.Execute(vatLineRow)
		Expect(err).NotTo(HaveOccurred())

		var spotResult test_helpers.VariableRes
		err = db.Get(&spotResult, `SELECT block_number, block_hash, spot AS value FROM maker.vat_ilk_spot`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(spotResult, blockNumber, "0x3f58749d3956984c2b03a84d5c02105a06efa1ad048d8aa97cf8f59aafa8f08b", "136985000000000000000000000000")
	})

	It("reads in a Vat.live storage diff row and persists it", func() {
		blockNumber := 10501122
		vatLineRow := utils.StorageDiffRow{
			Contract:     common.HexToAddress("67fd6c3575fc2dbe2cb596bd3bebc9edb5571fa1"),
			BlockHeight:  blockNumber,
			BlockHash:    common.HexToHash("1622e1531ade0154465dd99a9d25e3b4e4b8b9338edae51b71961446158f177b"),
			StorageKey:   common.HexToHash("000000000000000000000000000000000000000000000000000000000000000a"),
			StorageValue: common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001"),
		}
		err := transformer.Execute(vatLineRow)
		Expect(err).NotTo(HaveOccurred())

		var liveResult test_helpers.VariableRes
		err = db.Get(&liveResult, `SELECT block_number, block_hash, live AS value FROM maker.vat_live`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(liveResult, blockNumber, "0x1622e1531ade0154465dd99a9d25e3b4e4b8b9338edae51b71961446158f177b", "1")
	})
})
