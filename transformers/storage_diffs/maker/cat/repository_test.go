package cat_test

import (
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage_diffs/maker/cat"
	. "github.com/vulcanize/mcd_transformers/transformers/storage_diffs/maker/test_helpers"
)

var _ = Describe("Cat storage repository", func() {
	var (
		db              *postgres.DB
		repo            cat.CatStorageRepository
		fakeBlockNumber = 123
		fakeBlockHash   = "expected_block_hash"
		fakeAddress     = "0x12345"
		fakeIlk         = "fake_ilk"
		fakeUint256     = "12345"
		fakeBytes32     = "fake_bytes32"
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repo = cat.CatStorageRepository{}
		repo.SetDB(db)
	})

	Describe("Variable", func() {
		var result VariableRes

		Describe("NFlip", func() {
			It("writes a row", func() {
				nFlipMetadata := utils.GetStorageValueMetadata(cat.NFlip, nil, utils.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, nFlipMetadata, fakeUint256)
				Expect(err).NotTo(HaveOccurred())

				err = db.Get(&result, `SELECT block_number, block_hash, nflip AS value FROM maker.cat_nflip`)
				Expect(err).NotTo(HaveOccurred())
				AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeUint256)
			})
		})

		Describe("Live", func() {
			It("writes a row", func() {
				liveMetadata := utils.GetStorageValueMetadata(cat.Live, nil, utils.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, liveMetadata, fakeUint256)
				Expect(err).NotTo(HaveOccurred())

				err = db.Get(&result, `SELECT block_number, block_hash, live AS value FROM maker.cat_live`)
				Expect(err).NotTo(HaveOccurred())
				AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeUint256)
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
		})

		Describe("Pit", func() {
			It("writes a row", func() {
				pitMetadata := utils.GetStorageValueMetadata(cat.Pit, nil, utils.Address)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, pitMetadata, fakeAddress)
				Expect(err).NotTo(HaveOccurred())

				err = db.Get(&result, `SELECT block_number, block_hash, pit AS value FROM maker.cat_pit`)
				Expect(err).NotTo(HaveOccurred())
				AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeAddress)
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
		})
	})

	Describe("Ilk", func() {
		var result MappingRes

		Describe("Flip", func() {
			It("writes a row", func() {
				ilkFlipMetadata := utils.GetStorageValueMetadata(cat.IlkFlip, map[utils.Key]string{constants.Ilk: fakeIlk}, utils.Address)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkFlipMetadata, fakeAddress)
				Expect(err).NotTo(HaveOccurred())

				err = db.Get(&result, `SELECT block_number, block_hash, ilk AS key, flip AS value FROM maker.cat_ilk_flip`)
				Expect(err).NotTo(HaveOccurred())
				ilkID, err := shared.GetOrCreateIlk(fakeIlk, db)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, fakeBlockNumber, fakeBlockHash, strconv.Itoa(ilkID), fakeAddress)
			})

			It("returns an error if metadata missing ilk", func() {
				malformedIlkFlipMetadata := utils.GetStorageValueMetadata(cat.IlkFlip, map[utils.Key]string{}, utils.Address)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkFlipMetadata, fakeAddress)
				Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
			})
		})

		Describe("Chop", func() {
			It("writes a row", func() {
				ilkChopMetadata := utils.GetStorageValueMetadata(cat.IlkChop, map[utils.Key]string{constants.Ilk: fakeIlk}, utils.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkChopMetadata, fakeUint256)
				Expect(err).NotTo(HaveOccurred())

				err = db.Get(&result, `SELECT block_number, block_hash, ilk AS key, chop AS value FROM maker.cat_ilk_chop`)
				Expect(err).NotTo(HaveOccurred())
				ilkID, err := shared.GetOrCreateIlk(fakeIlk, db)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, fakeBlockNumber, fakeBlockHash, strconv.Itoa(ilkID), fakeUint256)
			})

			It("returns an error if metadata missing ilk", func() {
				malformedIlkChopMetadata := utils.GetStorageValueMetadata(cat.IlkChop, map[utils.Key]string{}, utils.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkChopMetadata, fakeAddress)
				Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
			})
		})

		Describe("Lump", func() {
			It("writes a row", func() {
				ilkLumpMetadata := utils.GetStorageValueMetadata(cat.IlkLump, map[utils.Key]string{constants.Ilk: fakeIlk}, utils.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkLumpMetadata, fakeUint256)
				Expect(err).NotTo(HaveOccurred())

				err = db.Get(&result, `SELECT block_number, block_hash, ilk AS key, lump AS value FROM maker.cat_ilk_lump`)
				Expect(err).NotTo(HaveOccurred())
				ilkID, err := shared.GetOrCreateIlk(fakeIlk, db)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, fakeBlockNumber, fakeBlockHash, strconv.Itoa(ilkID), fakeUint256)
			})

			It("returns an error if metadata missing ilk", func() {
				malformedIlkLumpMetadata := utils.GetStorageValueMetadata(cat.IlkLump, map[utils.Key]string{}, utils.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkLumpMetadata, fakeAddress)
				Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
			})
		})
	})

	Describe("Flip", func() {
		var result MappingRes

		Describe("FlipIlk", func() {
			It("writes a row", func() {
				flipIlkMetadata := utils.GetStorageValueMetadata(cat.FlipIlk, map[utils.Key]string{constants.Flip: fakeUint256}, utils.Bytes32)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, flipIlkMetadata, fakeBytes32)
				Expect(err).NotTo(HaveOccurred())

				err = db.Get(&result, `SELECT block_number, block_hash, flip AS key, ilk AS value FROM maker.cat_flip_ilk`)
				Expect(err).NotTo(HaveOccurred())
				ilkID, err := shared.GetOrCreateIlk(fakeBytes32, db)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, fakeBlockNumber, fakeBlockHash, fakeUint256, strconv.Itoa(ilkID))
			})

			It("returns an error if metadata missing flip", func() {
				malformedFlipIlkMetadata := utils.GetStorageValueMetadata(cat.FlipIlk, map[utils.Key]string{}, utils.Bytes32)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedFlipIlkMetadata, fakeBytes32)
				Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Flip}))
			})
		})

		Describe("FlipUrn", func() {
			It("writes a row", func() {
				flipUrnMetadata := utils.GetStorageValueMetadata(cat.FlipUrn, map[utils.Key]string{constants.Flip: fakeUint256}, utils.Bytes32)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, flipUrnMetadata, fakeBytes32)
				Expect(err).NotTo(HaveOccurred())

				err = db.Get(&result, `SELECT block_number, block_hash, flip AS key, urn AS value FROM maker.cat_flip_urn`)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, fakeBlockNumber, fakeBlockHash, fakeUint256, fakeBytes32)
			})

			It("returns an error if metadata missing flip", func() {
				malformedFlipUrnMetadata := utils.GetStorageValueMetadata(cat.FlipUrn, map[utils.Key]string{}, utils.Bytes32)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedFlipUrnMetadata, fakeBytes32)
				Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Flip}))
			})
		})

		Describe("FlipInk", func() {
			It("writes a row", func() {
				flipInkMetadata := utils.GetStorageValueMetadata(cat.FlipInk, map[utils.Key]string{constants.Flip: fakeUint256}, utils.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, flipInkMetadata, fakeUint256)
				Expect(err).NotTo(HaveOccurred())

				err = db.Get(&result, `SELECT block_number, block_hash, flip AS key, ink AS value FROM maker.cat_flip_ink`)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, fakeBlockNumber, fakeBlockHash, fakeUint256, fakeUint256)
			})

			It("returns an error if metadata missing flip", func() {
				malformedFlipInkMetadata := utils.GetStorageValueMetadata(cat.FlipInk, map[utils.Key]string{}, utils.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedFlipInkMetadata, fakeUint256)
				Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Flip}))
			})
		})

		Describe("FlipTab", func() {
			It("writes a row", func() {
				flipTabMetadata := utils.GetStorageValueMetadata(cat.FlipTab, map[utils.Key]string{constants.Flip: fakeUint256}, utils.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, flipTabMetadata, fakeUint256)
				Expect(err).NotTo(HaveOccurred())

				err = db.Get(&result, `SELECT block_number, block_hash, flip AS key, tab AS value FROM maker.cat_flip_tab`)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, fakeBlockNumber, fakeBlockHash, fakeUint256, fakeUint256)
			})

			It("returns an error if metadata missing flip", func() {
				malformedFlipTabMetadata := utils.GetStorageValueMetadata(cat.FlipTab, map[utils.Key]string{}, utils.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedFlipTabMetadata, fakeUint256)
				Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Flip}))
			})
		})
	})
})
