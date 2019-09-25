package cat_test

import (
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage/cat"
	. "github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
)

var _ = Describe("Cat storage repository", func() {
	var (
		db              *postgres.DB
		repo            cat.CatStorageRepository
		fakeBlockNumber = 123
		fakeBlockHash   = "expected_block_hash"
		fakeAddress     = "0x12345"
		fakeUint256     = "12345"
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repo = cat.CatStorageRepository{}
		repo.SetDB(db)
	})

	Describe("Variable", func() {
		var result VariableRes

		Describe("Live", func() {
			It("writes a row", func() {
				liveMetadata := utils.GetStorageValueMetadata(cat.Live, nil, utils.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, liveMetadata, fakeUint256)
				Expect(err).NotTo(HaveOccurred())

				err = db.Get(&result, `SELECT block_number, block_hash, live AS value FROM maker.cat_live`)
				Expect(err).NotTo(HaveOccurred())
				AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeUint256)
			})

			It("does not duplicate row", func() {
				liveMetadata := utils.GetStorageValueMetadata(cat.Live, nil, utils.Uint256)
				insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, liveMetadata, fakeUint256)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, liveMetadata, fakeUint256)

				Expect(insertTwoErr).NotTo(HaveOccurred())
				var count int
				getCountErr := db.Get(&count, `SELECT count(*) FROM maker.cat_live`)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})
		})

		Describe("Vat", func() {
			It("writes a row", func() {
				vatMetadata := utils.GetStorageValueMetadata(cat.Vat, nil, utils.Address)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, vatMetadata, fakeAddress)
				Expect(err).NotTo(HaveOccurred())

				err = db.Get(&result, `SELECT block_number, block_hash, vat AS value FROM maker.cat_vat`)
				Expect(err).NotTo(HaveOccurred())
				AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeAddress)
			})

			It("does not duplicate row", func() {
				vatMetadata := utils.GetStorageValueMetadata(cat.Vat, nil, utils.Address)
				insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, vatMetadata, fakeAddress)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, vatMetadata, fakeAddress)

				Expect(insertTwoErr).NotTo(HaveOccurred())
				var count int
				getCountErr := db.Get(&count, `SELECT count(*) FROM maker.cat_vat`)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})
		})

		Describe("Vow", func() {
			It("writes a row", func() {
				vowMetadata := utils.GetStorageValueMetadata(cat.Vow, nil, utils.Address)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, vowMetadata, fakeAddress)
				Expect(err).NotTo(HaveOccurred())

				err = db.Get(&result, `SELECT block_number, block_hash, vow AS value FROM maker.cat_vow`)
				Expect(err).NotTo(HaveOccurred())
				AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeAddress)
			})

			It("does not duplicate row", func() {
				vowMetadata := utils.GetStorageValueMetadata(cat.Vow, nil, utils.Address)
				insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, vowMetadata, fakeAddress)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, vowMetadata, fakeAddress)

				Expect(insertTwoErr).NotTo(HaveOccurred())
				var count int
				getCountErr := db.Get(&count, `SELECT count(*) FROM maker.cat_vow`)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})
		})
	})

	Describe("Ilk", func() {
		var result MappingRes

		Describe("Flip", func() {
			It("writes a row", func() {
				ilkFlipMetadata := utils.GetStorageValueMetadata(cat.IlkFlip, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Address)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkFlipMetadata, fakeAddress)
				Expect(err).NotTo(HaveOccurred())

				err = db.Get(&result, `SELECT block_number, block_hash, ilk_id AS key, flip AS value FROM maker.cat_ilk_flip`)
				Expect(err).NotTo(HaveOccurred())
				ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, fakeBlockNumber, fakeBlockHash, strconv.FormatInt(ilkID, 10), fakeAddress)
			})

			It("does not duplicate row", func() {
				ilkFlipMetadata := utils.GetStorageValueMetadata(cat.IlkFlip, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Address)
				insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkFlipMetadata, fakeAddress)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkFlipMetadata, fakeAddress)

				Expect(insertTwoErr).NotTo(HaveOccurred())
				var count int
				getCountErr := db.Get(&count, `SELECT count(*) FROM maker.cat_ilk_flip`)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})

			It("returns an error if metadata missing ilk", func() {
				malformedIlkFlipMetadata := utils.GetStorageValueMetadata(cat.IlkFlip, map[utils.Key]string{}, utils.Address)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkFlipMetadata, fakeAddress)
				Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
			})
		})

		Describe("Chop", func() {
			It("writes a row", func() {
				ilkChopMetadata := utils.GetStorageValueMetadata(cat.IlkChop, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkChopMetadata, fakeUint256)
				Expect(err).NotTo(HaveOccurred())

				err = db.Get(&result, `SELECT block_number, block_hash, ilk_id AS key, chop AS value FROM maker.cat_ilk_chop`)
				Expect(err).NotTo(HaveOccurred())
				ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, fakeBlockNumber, fakeBlockHash, strconv.FormatInt(ilkID, 10), fakeUint256)
			})

			It("does not duplicate row", func() {
				ilkChopMetadata := utils.GetStorageValueMetadata(cat.IlkChop, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
				insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkChopMetadata, fakeUint256)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkChopMetadata, fakeUint256)

				Expect(insertTwoErr).NotTo(HaveOccurred())
				var count int
				getCountErr := db.Get(&count, `SELECT count(*) FROM maker.cat_ilk_chop`)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})

			It("returns an error if metadata missing ilk", func() {
				malformedIlkChopMetadata := utils.GetStorageValueMetadata(cat.IlkChop, map[utils.Key]string{}, utils.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkChopMetadata, fakeAddress)
				Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
			})
		})

		Describe("Lump", func() {
			It("writes a row", func() {
				ilkLumpMetadata := utils.GetStorageValueMetadata(cat.IlkLump, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkLumpMetadata, fakeUint256)
				Expect(err).NotTo(HaveOccurred())

				err = db.Get(&result, `SELECT block_number, block_hash, ilk_id AS key, lump AS value FROM maker.cat_ilk_lump`)
				Expect(err).NotTo(HaveOccurred())
				ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, fakeBlockNumber, fakeBlockHash, strconv.FormatInt(ilkID, 10), fakeUint256)
			})

			It("does not duplicate row", func() {
				ilkLumpMetadata := utils.GetStorageValueMetadata(cat.IlkLump, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
				insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkLumpMetadata, fakeUint256)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkLumpMetadata, fakeUint256)

				Expect(insertTwoErr).NotTo(HaveOccurred())
				var count int
				getCountErr := db.Get(&count, `SELECT count(*) FROM maker.cat_ilk_lump`)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})

			It("returns an error if metadata missing ilk", func() {
				malformedIlkLumpMetadata := utils.GetStorageValueMetadata(cat.IlkLump, map[utils.Key]string{}, utils.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkLumpMetadata, fakeAddress)
				Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
			})
		})
	})
})
