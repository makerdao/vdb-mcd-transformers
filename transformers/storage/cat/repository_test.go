package cat_test

import (
	"github.com/ethereum/go-ethereum/common"
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/cat"
	. "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data/shared_behaviors"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/utils"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cat storage repository", func() {
	var (
		db                   *postgres.DB
		repo                 cat.CatStorageRepository
		diffID, fakeHeaderID int64
		fakeAddress          = "0x12345"
		fakeUint256          = "12345"
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repo = cat.CatStorageRepository{}
		repo.SetDB(db)
		headerRepository := repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		fakeHeaderID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
		Expect(insertHeaderErr).NotTo(HaveOccurred())
	})

	Describe("Variable", func() {
		It("panics if the metadata name is not recognized", func() {
			unrecognizedMetadata := utils.StorageValueMetadata{Name: "unrecognized"}
			repoCreate := func() {
				repo.Create(diffID, fakeHeaderID, unrecognizedMetadata, "")
			}

			Expect(repoCreate).Should(Panic())
		})

		Describe("Live", func() {
			liveMetadata := utils.GetStorageValueMetadata(cat.Live, nil, utils.Uint256)
			inputs := shared_behaviors.StorageVariableBehaviorInputs{
				ValueFieldName:   cat.Live,
				Value:            fakeUint256,
				StorageTableName: "maker.cat_live",
				Repository:       &repo,
				Metadata:         liveMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
		})

		Describe("Vat", func() {
			vatMetadata := utils.GetStorageValueMetadata(cat.Vat, nil, utils.Address)
			inputs := shared_behaviors.StorageVariableBehaviorInputs{
				ValueFieldName:   cat.Vat,
				Value:            fakeAddress,
				StorageTableName: "maker.cat_vat",
				Repository:       &repo,
				Metadata:         vatMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
		})

		Describe("Vow", func() {
			vowMetadata := utils.GetStorageValueMetadata(cat.Vow, nil, utils.Address)
			inputs := shared_behaviors.StorageVariableBehaviorInputs{
				ValueFieldName:   cat.Vow,
				Value:            fakeAddress,
				StorageTableName: "maker.cat_vow",
				Repository:       &repo,
				Metadata:         vowMetadata,
			}

			shared_behaviors.SharedStorageRepositoryVariableBehaviors(&inputs)
		})
	})

	Describe("Ilk", func() {
		BeforeEach(func() {
			fakeRawDiff := fakes.GetFakeStorageDiffForHeader(fakes.FakeHeader, common.Hash{}, common.Hash{}, common.Hash{})
			storageDiffRepo := repositories.NewStorageDiffRepository(db)
			var insertDiffErr error
			diffID, insertDiffErr = storageDiffRepo.CreateStorageDiff(fakeRawDiff)
			Expect(insertDiffErr).NotTo(HaveOccurred())
		})

		Describe("Flip", func() {
			It("writes a row", func() {
				ilkFlipMetadata := utils.GetStorageValueMetadata(cat.IlkFlip, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Address)

				err := repo.Create(diffID, fakeHeaderID, ilkFlipMetadata, fakeAddress)
				Expect(err).NotTo(HaveOccurred())

				var result MappingRes
				err = db.Get(&result, `SELECT diff_id, header_id, ilk_id AS key, flip AS value FROM maker.cat_ilk_flip`)
				Expect(err).NotTo(HaveOccurred())
				ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, diffID, fakeHeaderID, strconv.FormatInt(ilkID, 10), fakeAddress)
			})

			It("does not duplicate row", func() {
				ilkFlipMetadata := utils.GetStorageValueMetadata(cat.IlkFlip, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Address)
				insertOneErr := repo.Create(diffID, fakeHeaderID, ilkFlipMetadata, fakeAddress)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(diffID, fakeHeaderID, ilkFlipMetadata, fakeAddress)

				Expect(insertTwoErr).NotTo(HaveOccurred())
				var count int
				getCountErr := db.Get(&count, `SELECT count(*) FROM maker.cat_ilk_flip`)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})

			It("returns an error if metadata missing ilk", func() {
				malformedIlkFlipMetadata := utils.GetStorageValueMetadata(cat.IlkFlip, map[utils.Key]string{}, utils.Address)

				err := repo.Create(diffID, fakeHeaderID, malformedIlkFlipMetadata, fakeAddress)
				Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
			})

			shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
				Repository:    &repo,
				Metadata:      utils.GetStorageValueMetadata(cat.IlkFlip, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Address),
				PropertyName:  "Flip",
				PropertyValue: fakeAddress,
				TableName:     "maker.cat_ilk_flip",
			})
		})

		Describe("Chop", func() {
			It("writes a row", func() {
				ilkChopMetadata := utils.GetStorageValueMetadata(cat.IlkChop, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)

				err := repo.Create(diffID, fakeHeaderID, ilkChopMetadata, fakeUint256)
				Expect(err).NotTo(HaveOccurred())

				var result MappingRes
				err = db.Get(&result, `SELECT diff_id, header_id, ilk_id AS key, chop AS value FROM maker.cat_ilk_chop`)
				Expect(err).NotTo(HaveOccurred())
				ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, diffID, fakeHeaderID, strconv.FormatInt(ilkID, 10), fakeUint256)
			})

			It("does not duplicate row", func() {
				ilkChopMetadata := utils.GetStorageValueMetadata(cat.IlkChop, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
				insertOneErr := repo.Create(diffID, fakeHeaderID, ilkChopMetadata, fakeUint256)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(diffID, fakeHeaderID, ilkChopMetadata, fakeUint256)

				Expect(insertTwoErr).NotTo(HaveOccurred())
				var count int
				getCountErr := db.Get(&count, `SELECT count(*) FROM maker.cat_ilk_chop`)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})

			It("returns an error if metadata missing ilk", func() {
				malformedIlkChopMetadata := utils.GetStorageValueMetadata(cat.IlkChop, map[utils.Key]string{}, utils.Uint256)

				err := repo.Create(diffID, fakeHeaderID, malformedIlkChopMetadata, fakeAddress)
				Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
			})

			shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
				Repository:    &repo,
				Metadata:      utils.GetStorageValueMetadata(cat.IlkChop, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256),
				PropertyName:  "Chop",
				PropertyValue: strconv.Itoa(rand.Int()),
				TableName:     "maker.cat_ilk_chop",
			})
		})

		Describe("Lump", func() {
			It("writes a row", func() {
				ilkLumpMetadata := utils.GetStorageValueMetadata(cat.IlkLump, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)

				err := repo.Create(diffID, fakeHeaderID, ilkLumpMetadata, fakeUint256)
				Expect(err).NotTo(HaveOccurred())

				var result MappingRes
				err = db.Get(&result, `SELECT diff_id, header_id, ilk_id AS key, lump AS value FROM maker.cat_ilk_lump`)
				Expect(err).NotTo(HaveOccurred())
				ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, diffID, fakeHeaderID, strconv.FormatInt(ilkID, 10), fakeUint256)
			})

			It("does not duplicate row", func() {
				ilkLumpMetadata := utils.GetStorageValueMetadata(cat.IlkLump, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
				insertOneErr := repo.Create(diffID, fakeHeaderID, ilkLumpMetadata, fakeUint256)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(diffID, fakeHeaderID, ilkLumpMetadata, fakeUint256)

				Expect(insertTwoErr).NotTo(HaveOccurred())
				var count int
				getCountErr := db.Get(&count, `SELECT count(*) FROM maker.cat_ilk_lump`)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})

			It("returns an error if metadata missing ilk", func() {
				malformedIlkLumpMetadata := utils.GetStorageValueMetadata(cat.IlkLump, map[utils.Key]string{}, utils.Uint256)

				err := repo.Create(diffID, fakeHeaderID, malformedIlkLumpMetadata, fakeAddress)
				Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
			})

			shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
				Repository:    &repo,
				Metadata:      utils.GetStorageValueMetadata(cat.IlkLump, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256),
				PropertyName:  "Lump",
				PropertyValue: strconv.Itoa(rand.Int()),
				TableName:     "maker.cat_ilk_lump",
			})
		})
	})
})
