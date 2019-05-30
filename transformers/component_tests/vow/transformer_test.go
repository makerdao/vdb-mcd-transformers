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
		vowVat := utils.StorageDiffRow{
			Contract:     common.HexToAddress("17560834075da3db54f737db74377e799c865821"),
			BlockHeight:  blockNumber,
			BlockHash:    common.HexToHash("1822bb271ce246212f0d097e59b3b04e0302819da3a2bd80e85b91e8c89fc883"),
			StorageKey:   common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001"),
			StorageValue: common.HexToHash("00000000000000000000000067fd6c3575fc2dbe2cb596bd3bebc9edb5571fa1"),
		}
		err := transformer.Execute(vowVat)
		Expect(err).NotTo(HaveOccurred())

		var vatResult test_helpers.VariableRes
		err = db.Get(&vatResult, `SELECT block_number, block_hash, vat AS value FROM maker.vow_vat`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(vatResult, blockNumber, "0x1822bb271ce246212f0d097e59b3b04e0302819da3a2bd80e85b91e8c89fc883", "0x67fd6c3575Fc2dBE2CB596bD3bEbc9EDb5571fA1")
	})

	It("reads in a Vow.cow storage diff row and persists it", func() {
		blockNumber := 10501125
		vowCow := utils.StorageDiffRow{
			Contract:     common.HexToAddress("17560834075da3db54f737db74377e799c865821"),
			BlockHeight:  blockNumber,
			BlockHash:    common.HexToHash("1822bb271ce246212f0d097e59b3b04e0302819da3a2bd80e85b91e8c89fc883"),
			StorageKey:   common.HexToHash("0000000000000000000000000000000000000000000000000000000000000002"),
			StorageValue: common.HexToHash("00000000000000000000000069803564432dee97e37ac1064e486aa52d607e03"),
		}
		err := transformer.Execute(vowCow)
		Expect(err).NotTo(HaveOccurred())

		var cowResult test_helpers.VariableRes
		err = db.Get(&cowResult, `SELECT block_number, block_hash, cow AS value FROM maker.vow_cow`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(cowResult, blockNumber, "0x1822bb271ce246212f0d097e59b3b04e0302819da3a2bd80e85b91e8c89fc883", "0x69803564432dEe97e37aC1064e486aA52D607E03")
	})

	It("reads in a Vow.row storage diff row and persists it", func() {
		blockNumber := 10501127
		vowRow := utils.StorageDiffRow{
			Contract:     common.HexToAddress("17560834075da3db54f737db74377e799c865821"),
			BlockHeight:  blockNumber,
			BlockHash:    common.HexToHash("c3c83991f0d591f66d77b77208be4ed7ec7376edfb156fa4e0bbc9be67b41834"),
			StorageKey:   common.HexToHash("0000000000000000000000000000000000000000000000000000000000000003"),
			StorageValue: common.HexToHash("000000000000000000000000187e5de065079320b50b8420692cae8aaba098a3"),
		}
		err := transformer.Execute(vowRow)
		Expect(err).NotTo(HaveOccurred())

		var rowResult test_helpers.VariableRes
		err = db.Get(&rowResult, `SELECT block_number, block_hash, row AS value FROM maker.vow_row`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(rowResult, blockNumber, "0xc3c83991f0d591f66d77b77208be4ed7ec7376edfb156fa4e0bbc9be67b41834", "0x187e5De065079320B50B8420692Cae8aABa098A3")
	})

	It("reads in a Vow.sump storage diff row and persists it", func() {
		blockNumber := 10869770
		vowSump := utils.StorageDiffRow{
			Contract:     common.HexToAddress("4afcab85f27dd2e1a5ec1008b5b294e44e487f90"),
			BlockHeight:  blockNumber,
			BlockHash:    common.HexToHash("fe124bd8436290c364692b928a59f02f4d458c642b40398cbae173252b54093c"),
			StorageKey:   common.HexToHash("0000000000000000000000000000000000000000000000000000000000000008"),
			StorageValue: common.HexToHash("00000000000000000000000000047bf19673df52e37f2410011d100000000000"),
		}
		err := transformer.Execute(vowSump)
		Expect(err).NotTo(HaveOccurred())

		var rowResult test_helpers.VariableRes
		err = db.Get(&rowResult, `SELECT block_number, block_hash, sump AS value FROM maker.vow_sump`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(rowResult, blockNumber, "0xfe124bd8436290c364692b928a59f02f4d458c642b40398cbae173252b54093c", "100000000000000000000000000000000000000000000")
	})

	It("reads in a Vow.bump storage diff row and persists it", func() {
		blockNumber := 10869768
		vowBump := utils.StorageDiffRow{
			Contract:     common.HexToAddress("4afcab85f27dd2e1a5ec1008b5b294e44e487f90"),
			BlockHeight:  blockNumber,
			BlockHash:    common.HexToHash("a750d8cf2317bb6d65b43b96ff24a179ed8c3a237f874c0e867987180b2527a8"),
			StorageKey:   common.HexToHash("0000000000000000000000000000000000000000000000000000000000000009"),
			StorageValue: common.HexToHash("00000000000000000000000000047bf19673df52e37f2410011d100000000000"),
		}
		err := transformer.Execute(vowBump)
		Expect(err).NotTo(HaveOccurred())

		var rowResult test_helpers.VariableRes
		err = db.Get(&rowResult, `SELECT block_number, block_hash, bump AS value FROM maker.vow_bump`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(rowResult, blockNumber, "0xa750d8cf2317bb6d65b43b96ff24a179ed8c3a237f874c0e867987180b2527a8", "100000000000000000000000000000000000000000000")
	})
})
