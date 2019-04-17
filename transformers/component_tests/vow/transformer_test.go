package vow

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	storage2 "github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vow"
	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

var _ = Describe("Executing the transformer", func() {
	var (
		db          *postgres.DB
		mappings    = vow.VowMappings{StorageRepository: &storage2.MakerStorageRepository{}}
		repository  = vow.VowStorageRepository{}
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
	})

	It("reads in a Vow.vat storage diff row and persists it", func() {
		blockNumber := 10501125
		vowLineRow := utils.StorageDiffRow{
			Contract:     common.HexToAddress("17560834075da3db54f737db74377e799c865821"),
			BlockHeight:  blockNumber,
			BlockHash:    common.HexToHash("1822bb271ce246212f0d097e59b3b04e0302819da3a2bd80e85b91e8c89fc883"),
			StorageKey:   common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001"),
			StorageValue: common.HexToHash("00000000000000000000000067fd6c3575fc2dbe2cb596bd3bebc9edb5571fa1"),
		}
		err := transformer.Execute(vowLineRow)
		Expect(err).NotTo(HaveOccurred())

		var vatResult test_helpers.VariableRes
		err = db.Get(&vatResult, `SELECT block_number, block_hash, vat AS value FROM maker.vow_vat`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(vatResult, blockNumber, "0x1822bb271ce246212f0d097e59b3b04e0302819da3a2bd80e85b91e8c89fc883", "0x67fd6c3575Fc2dBE2CB596bD3bEbc9EDb5571fA1")
	})

	It("reads in a Vow.cow storage diff row and persists it", func() {
		blockNumber := 10501125
		vowLineRow := utils.StorageDiffRow{
			Contract:     common.HexToAddress("17560834075da3db54f737db74377e799c865821"),
			BlockHeight:  blockNumber,
			BlockHash:    common.HexToHash("1822bb271ce246212f0d097e59b3b04e0302819da3a2bd80e85b91e8c89fc883"),
			StorageKey:   common.HexToHash("0000000000000000000000000000000000000000000000000000000000000002"),
			StorageValue: common.HexToHash("00000000000000000000000069803564432dee97e37ac1064e486aa52d607e03"),
		}
		err := transformer.Execute(vowLineRow)
		Expect(err).NotTo(HaveOccurred())

		var cowResult test_helpers.VariableRes
		err = db.Get(&cowResult, `SELECT block_number, block_hash, cow AS value FROM maker.vow_cow`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(cowResult, blockNumber, "0x1822bb271ce246212f0d097e59b3b04e0302819da3a2bd80e85b91e8c89fc883", "0x69803564432dEe97e37aC1064e486aA52D607E03")
	})

	It("reads in a Vow.row storage diff row and persists it", func() {
		blockNumber := 10501127
		vowLineRow := utils.StorageDiffRow{
			Contract:     common.HexToAddress("17560834075da3db54f737db74377e799c865821"),
			BlockHeight:  blockNumber,
			BlockHash:    common.HexToHash("c3c83991f0d591f66d77b77208be4ed7ec7376edfb156fa4e0bbc9be67b41834"),
			StorageKey:   common.HexToHash("0000000000000000000000000000000000000000000000000000000000000003"),
			StorageValue: common.HexToHash("000000000000000000000000187e5de065079320b50b8420692cae8aaba098a3"),
		}
		err := transformer.Execute(vowLineRow)
		Expect(err).NotTo(HaveOccurred())

		var rowResult test_helpers.VariableRes
		err = db.Get(&rowResult, `SELECT block_number, block_hash, row AS value FROM maker.vow_row`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(rowResult, blockNumber, "0xc3c83991f0d591f66d77b77208be4ed7ec7376edfb156fa4e0bbc9be67b41834", "0x187e5De065079320B50B8420692Cae8aABa098A3")
	})
})
