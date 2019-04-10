package cat

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
	"github.com/vulcanize/mcd_transformers/transformers/storage/cat"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
)

var _ = Describe("Executing the transformer", func() {
	var (
		db          *postgres.DB
		ilk         string
		ilkID       int
		err         error
		mappings    = cat.CatMappings{StorageRepository: &storage2.MakerStorageRepository{}}
		repository  = cat.CatStorageRepository{}
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
		ilk = "4554480000000000000000000000000000000000000000000000000000000000"
		ilkID, err = shared.GetOrCreateIlk(ilk, db)
		Expect(err).NotTo(HaveOccurred())
		_, flipErr := db.Exec(`INSERT INTO maker.cat_nflip (nflip) VALUES ($1)`, 10)
		Expect(flipErr).NotTo(HaveOccurred())
	})

	It("reads in a Cat Live storage diff row and persists it", func() {
		blockNumber := 10501127
		catLineRow := utils.StorageDiffRow{
			Contract:     common.HexToAddress("35f3d8997ef261c7961bd7c07ddc390f5cf76bd3"),
			BlockHeight:  blockNumber,
			BlockHash:    common.HexToHash("c3c83991f0d591f66d77b77208be4ed7ec7376edfb156fa4e0bbc9be67b41834"),
			StorageKey:   common.HexToHash("0000000000000000000000000000000000000000000000000000000000000004"),
			StorageValue: common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001"),
		}
		err := transformer.Execute(catLineRow)
		Expect(err).NotTo(HaveOccurred())

		var liveResult test_helpers.VariableRes
		err = db.Get(&liveResult, `SELECT block_number, block_hash, live AS value FROM maker.cat_live`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(liveResult, blockNumber, "0xc3c83991f0d591f66d77b77208be4ed7ec7376edfb156fa4e0bbc9be67b41834", "1")
	})

	It("reads in a Cat Vat storage diff row and persists it", func() {
		blockNumber := 10501127
		catLineRow := utils.StorageDiffRow{
			Contract:     common.HexToAddress("35f3d8997ef261c7961bd7c07ddc390f5cf76bd3"),
			BlockHeight:  blockNumber,
			BlockHash:    common.HexToHash("c3c83991f0d591f66d77b77208be4ed7ec7376edfb156fa4e0bbc9be67b41834"),
			StorageKey:   common.HexToHash("0000000000000000000000000000000000000000000000000000000000000005"),
			StorageValue: common.HexToHash("00000000000000000000000067fd6c3575fc2dbe2cb596bd3bebc9edb5571fa1"),
		}
		err := transformer.Execute(catLineRow)
		Expect(err).NotTo(HaveOccurred())

		var vatResult test_helpers.VariableRes
		err = db.Get(&vatResult, `SELECT block_number, block_hash, vat AS value FROM maker.cat_vat`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(vatResult, blockNumber, "0xc3c83991f0d591f66d77b77208be4ed7ec7376edfb156fa4e0bbc9be67b41834", "0x67fd6c3575Fc2dBE2CB596bD3bEbc9EDb5571fA1")
	})

	It("reads in a Cat Vow storage diff row and persists it", func() {
		blockNumber := 10501127
		catLineRow := utils.StorageDiffRow{
			Contract:     common.HexToAddress("35f3d8997ef261c7961bd7c07ddc390f5cf76bd3"),
			BlockHeight:  blockNumber,
			BlockHash:    common.HexToHash("c3c83991f0d591f66d77b77208be4ed7ec7376edfb156fa4e0bbc9be67b41834"),
			StorageKey:   common.HexToHash("0000000000000000000000000000000000000000000000000000000000000006"),
			StorageValue: common.HexToHash("00000000000000000000000017560834075da3db54f737db74377e799c865821"),
		}
		err := transformer.Execute(catLineRow)
		Expect(err).NotTo(HaveOccurred())

		var vowResult test_helpers.VariableRes
		err = db.Get(&vowResult, `SELECT block_number, block_hash, vow AS value FROM maker.cat_vow`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(vowResult, blockNumber, "0xc3c83991f0d591f66d77b77208be4ed7ec7376edfb156fa4e0bbc9be67b41834", "0x17560834075DA3Db54f737db74377E799c865821")
	})

	It("reads in a Cat NFlip storage diff row and persists it", func() {
		blockNumber := 9377319
		catNFlipRow := utils.StorageDiffRow{
			Contract:     common.HexToAddress("2f34f22a00ee4b7a8f8bbc4eaee1658774c624e0"),
			BlockHeight:  blockNumber,
			BlockHash:    common.HexToHash("2fffcd7d71edff194bf79d773422ee5690544f811c020d8fcaa24791093383c4"),
			StorageKey:   common.HexToHash("0000000000000000000000000000000000000000000000000000000000000003"),
			StorageValue: common.HexToHash("0000000000000000000000000000000000000000000000000000000000000008"),
		}
		err := transformer.Execute(catNFlipRow)
		Expect(err).NotTo(HaveOccurred())

		var nFlipResult test_helpers.VariableRes
		err = db.Get(&nFlipResult, `SELECT block_number, block_hash, nflip AS value FROM maker.cat_nflip WHERE block_number = $1`, blockNumber)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(nFlipResult, blockNumber, "0x2fffcd7d71edff194bf79d773422ee5690544f811c020d8fcaa24791093383c4", "8")

	})

	It("reads in a Cat Ilk Flip storage diff row and persists it", func() {
		blockNumber := 10501138
		catLineRow := utils.StorageDiffRow{
			Contract:     common.HexToAddress("35f3d8997ef261c7961bd7c07ddc390f5cf76bd3"),
			BlockHeight:  blockNumber,
			BlockHash:    common.HexToHash("3f58749d3956984c2b03a84d5c02105a06efa1ad048d8aa97cf8f59aafa8f08b"),
			StorageKey:   common.HexToHash("a27f5adbce3dcb790941ebd020e02078a61e6c9748376e52ec0929d58babf97a"),
			StorageValue: common.HexToHash("000000000000000000000000d68e8045549fa4d0570c4c03dd294e6a46a121eb"),
		}
		err := transformer.Execute(catLineRow)
		Expect(err).NotTo(HaveOccurred())

		var ilkFlipResult test_helpers.MappingRes
		err = db.Get(&ilkFlipResult, `SELECT block_number, block_hash, ilk_id AS key, flip AS value FROM maker.cat_ilk_flip`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertMapping(ilkFlipResult, blockNumber, "0x3f58749d3956984c2b03a84d5c02105a06efa1ad048d8aa97cf8f59aafa8f08b", strconv.Itoa(ilkID), "0xd68e8045549Fa4d0570c4c03DD294e6a46a121Eb")
	})

	It("reads in a Cat Ilk Chop storage diff row and persists it", func() {
		blockNumber := 10501138
		catLineRow := utils.StorageDiffRow{
			Contract:     common.HexToAddress("35f3d8997ef261c7961bd7c07ddc390f5cf76bd3"),
			BlockHeight:  blockNumber,
			BlockHash:    common.HexToHash("3f58749d3956984c2b03a84d5c02105a06efa1ad048d8aa97cf8f59aafa8f08b"),
			StorageKey:   common.HexToHash("a27f5adbce3dcb790941ebd020e02078a61e6c9748376e52ec0929d58babf97b"),
			StorageValue: common.HexToHash("0000000000000000000000000000000000000000033b2e3c9fd0803ce8000000"),
		}
		err := transformer.Execute(catLineRow)
		Expect(err).NotTo(HaveOccurred())

		var ilkChopResult test_helpers.MappingRes
		err = db.Get(&ilkChopResult, `SELECT block_number, block_hash, ilk_id AS key, chop AS value FROM maker.cat_ilk_chop`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertMapping(ilkChopResult, blockNumber, "0x3f58749d3956984c2b03a84d5c02105a06efa1ad048d8aa97cf8f59aafa8f08b", strconv.Itoa(ilkID), "1000000000000000000000000000")
	})

	It("reads in a Cat Ilk Lump storage diff row and persists it", func() {
		blockNumber := 10501138
		catLineRow := utils.StorageDiffRow{
			Contract:     common.HexToAddress("35f3d8997ef261c7961bd7c07ddc390f5cf76bd3"),
			BlockHeight:  blockNumber,
			BlockHash:    common.HexToHash("3f58749d3956984c2b03a84d5c02105a06efa1ad048d8aa97cf8f59aafa8f08b"),
			StorageKey:   common.HexToHash("a27f5adbce3dcb790941ebd020e02078a61e6c9748376e52ec0929d58babf97c"),
			StorageValue: common.HexToHash("00000000000000000000000000000000000000000000021e19e0c9bab2400000"),
		}
		err := transformer.Execute(catLineRow)
		Expect(err).NotTo(HaveOccurred())

		var ilkLumpResult test_helpers.MappingRes
		err = db.Get(&ilkLumpResult, `SELECT block_number, block_hash, ilk_id AS key, lump AS value FROM maker.cat_ilk_lump`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertMapping(ilkLumpResult, blockNumber, "0x3f58749d3956984c2b03a84d5c02105a06efa1ad048d8aa97cf8f59aafa8f08b", strconv.Itoa(ilkID), "10000000000000000000000")
	})

	It("reads in a Cat Flip Ilk storage diff row and persists it", func() {
		blockNumber := 9377319
		catFlipIlkRow := utils.StorageDiffRow{
			Contract:     common.HexToAddress("2f34f22a00ee4b7a8f8bbc4eaee1658774c624e0"),
			BlockHeight:  blockNumber,
			BlockHash:    common.HexToHash("2fffcd7d71edff194bf79d773422ee5690544f811c020d8fcaa24791093383c4"),
			StorageKey:   common.HexToHash("acd8ef244210bb6898e73c48bf820ed8ecc857a3bab8d79c10e4fa92b1e9ca65"),
			StorageValue: common.HexToHash("4554480000000000000000000000000000000000000000000000000000000000"),
		}
		err := transformer.Execute(catFlipIlkRow)
		Expect(err).NotTo(HaveOccurred())

		var flipIlkResult test_helpers.MappingRes
		err = db.Get(&flipIlkResult, `SELECT block_number, block_hash, flip AS key, ilk_id AS value FROM maker.cat_flip_ilk`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertMapping(flipIlkResult, blockNumber, "0x2fffcd7d71edff194bf79d773422ee5690544f811c020d8fcaa24791093383c4", "7", strconv.Itoa(ilkID))
	})

	It("reads in a Cat Flip Urn storage diff row and persists it", func() {
		blockNumber := 9377319
		catFlipUrnRow := utils.StorageDiffRow{
			Contract:     common.HexToAddress("2f34f22a00ee4b7a8f8bbc4eaee1658774c624e0"),
			BlockHeight:  blockNumber,
			BlockHash:    common.HexToHash("2fffcd7d71edff194bf79d773422ee5690544f811c020d8fcaa24791093383c4"),
			StorageKey:   common.HexToHash("acd8ef244210bb6898e73c48bf820ed8ecc857a3bab8d79c10e4fa92b1e9ca66"),
			StorageValue: common.HexToHash("0000000000000000000000000000d8b4147eda80fec7122ae16da2479cbd7ffb"),
		}
		err := transformer.Execute(catFlipUrnRow)
		Expect(err).NotTo(HaveOccurred())

		var flipUrnResult test_helpers.MappingRes
		err = db.Get(&flipUrnResult, `SELECT block_number, block_hash, flip AS key, urn AS value FROM maker.cat_flip_urn`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertMapping(flipUrnResult, blockNumber, "0x2fffcd7d71edff194bf79d773422ee5690544f811c020d8fcaa24791093383c4", "7", "0x0000000000000000000000000000d8b4147eda80fec7122ae16da2479cbd7ffb")
	})

	It("reads in a Cat Flip Ink storage diff row and persists it", func() {
		blockNumber := 9377319
		catFlipInkRow := utils.StorageDiffRow{
			Contract:     common.HexToAddress("2f34f22a00ee4b7a8f8bbc4eaee1658774c624e0"),
			BlockHeight:  blockNumber,
			BlockHash:    common.HexToHash("2fffcd7d71edff194bf79d773422ee5690544f811c020d8fcaa24791093383c4"),
			StorageKey:   common.HexToHash("acd8ef244210bb6898e73c48bf820ed8ecc857a3bab8d79c10e4fa92b1e9ca67"),
			StorageValue: common.HexToHash("000000000000000000000000000000000000000000000004563918244f400000"),
		}
		err := transformer.Execute(catFlipInkRow)
		Expect(err).NotTo(HaveOccurred())

		var flipInkResult test_helpers.MappingRes
		err = db.Get(&flipInkResult, `SELECT block_number, block_hash, flip AS key, ink AS value FROM maker.cat_flip_ink`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertMapping(flipInkResult, blockNumber, "0x2fffcd7d71edff194bf79d773422ee5690544f811c020d8fcaa24791093383c4", "7", "80000000000000000000")
	})

	It("reads in a Cat Flip Tab storage diff row and persists it", func() {
		blockNumber := 9377319
		catFlipTabRow := utils.StorageDiffRow{
			Contract:     common.HexToAddress("2f34f22a00ee4b7a8f8bbc4eaee1658774c624e0"),
			BlockHeight:  blockNumber,
			BlockHash:    common.HexToHash("2fffcd7d71edff194bf79d773422ee5690544f811c020d8fcaa24791093383c4"),
			StorageKey:   common.HexToHash("acd8ef244210bb6898e73c48bf820ed8ecc857a3bab8d79c10e4fa92b1e9ca68"),
			StorageValue: common.HexToHash("0000000000000000000000000000000000000000000002544faa778090e00000"),
		}
		err := transformer.Execute(catFlipTabRow)
		Expect(err).NotTo(HaveOccurred())

		var flipTabResult test_helpers.MappingRes
		err = db.Get(&flipTabResult, `SELECT block_number, block_hash, flip AS key, tab AS value FROM maker.cat_flip_tab`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertMapping(flipTabResult, blockNumber, "0x2fffcd7d71edff194bf79d773422ee5690544f811c020d8fcaa24791093383c4", "7", "11000000000000000000000")
	})
})
