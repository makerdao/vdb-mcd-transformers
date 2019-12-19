package cat_test

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/cat"
	. "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data/shared_behaviors"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cat storage repository", func() {
	var (
		db                   = test_config.NewTestDB(test_config.NewTestNode())
		repo                 cat.CatStorageRepository
		diffID, fakeHeaderID int64
		fakeAddress          = "0x12345"
		fakeUint256          = "12345"
	)

	BeforeEach(func() {
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
			unrecognizedMetadata := storage.ValueMetadata{Name: "unrecognized"}
			repoCreate := func() {
				repo.Create(diffID, fakeHeaderID, unrecognizedMetadata, "")
			}

			Expect(repoCreate).Should(Panic())
		})

		Describe("Live", func() {
			liveMetadata := storage.GetValueMetadata(cat.Live, nil, storage.Uint256)
			inputs := shared_behaviors.StorageBehaviorInputs{
				ValueFieldName: cat.Live,
				Value:          fakeUint256,
				Schema:         constants.MakerSchema,
				TableName:      constants.CatLiveTable,
				Repository:     &repo,
				Metadata:       liveMetadata,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
		})

		Describe("Vat", func() {
			vatMetadata := storage.GetValueMetadata(cat.Vat, nil, storage.Address)
			inputs := shared_behaviors.StorageBehaviorInputs{
				ValueFieldName: cat.Vat,
				Value:          fakeAddress,
				Schema:         constants.MakerSchema,
				TableName:      constants.CatVatTable,
				Repository:     &repo,
				Metadata:       vatMetadata,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
		})

		Describe("Vow", func() {
			vowMetadata := storage.GetValueMetadata(cat.Vow, nil, storage.Address)
			inputs := shared_behaviors.StorageBehaviorInputs{
				ValueFieldName: cat.Vow,
				Value:          fakeAddress,
				Schema:         constants.MakerSchema,
				TableName:      constants.CatVowTable,
				Repository:     &repo,
				Metadata:       vowMetadata,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
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
				ilkFlipMetadata := storage.GetValueMetadata(cat.IlkFlip, map[storage.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, storage.Address)

				err := repo.Create(diffID, fakeHeaderID, ilkFlipMetadata, fakeAddress)
				Expect(err).NotTo(HaveOccurred())

				var result MappingRes
				query := fmt.Sprintf(`SELECT diff_id, header_id, ilk_id AS key, flip AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.CatIlkFlipTable))
				err = db.Get(&result, query)
				Expect(err).NotTo(HaveOccurred())
				ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, diffID, fakeHeaderID, strconv.FormatInt(ilkID, 10), fakeAddress)
			})

			It("does not duplicate row", func() {
				ilkFlipMetadata := storage.GetValueMetadata(cat.IlkFlip, map[storage.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, storage.Address)
				insertOneErr := repo.Create(diffID, fakeHeaderID, ilkFlipMetadata, fakeAddress)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(diffID, fakeHeaderID, ilkFlipMetadata, fakeAddress)

				Expect(insertTwoErr).NotTo(HaveOccurred())
				var count int
				query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.CatIlkFlipTable))
				getCountErr := db.Get(&count, query)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})

			It("returns an error if metadata missing ilk", func() {
				malformedIlkFlipMetadata := storage.GetValueMetadata(cat.IlkFlip, map[storage.Key]string{}, storage.Address)

				err := repo.Create(diffID, fakeHeaderID, malformedIlkFlipMetadata, fakeAddress)
				Expect(err).To(MatchError(storage.ErrMetadataMalformed{MissingData: constants.Ilk}))
			})

			shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
				Repository:    &repo,
				Metadata:      storage.GetValueMetadata(cat.IlkFlip, map[storage.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, storage.Address),
				PropertyName:  "Flip",
				PropertyValue: fakeAddress,
				Schema:        constants.MakerSchema,
				TableName:     constants.CatIlkFlipTable,
			})
		})

		Describe("Chop", func() {
			It("writes a row", func() {
				ilkChopMetadata := storage.GetValueMetadata(cat.IlkChop, map[storage.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, storage.Uint256)

				err := repo.Create(diffID, fakeHeaderID, ilkChopMetadata, fakeUint256)
				Expect(err).NotTo(HaveOccurred())

				var result MappingRes
				query := fmt.Sprintf(`SELECT diff_id, header_id, ilk_id AS key, chop AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.CatIlkChopTable))
				err = db.Get(&result, query)
				Expect(err).NotTo(HaveOccurred())
				ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, diffID, fakeHeaderID, strconv.FormatInt(ilkID, 10), fakeUint256)
			})

			It("does not duplicate row", func() {
				ilkChopMetadata := storage.GetValueMetadata(cat.IlkChop, map[storage.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, storage.Uint256)
				insertOneErr := repo.Create(diffID, fakeHeaderID, ilkChopMetadata, fakeUint256)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(diffID, fakeHeaderID, ilkChopMetadata, fakeUint256)

				Expect(insertTwoErr).NotTo(HaveOccurred())
				var count int
				query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.CatIlkChopTable))
				getCountErr := db.Get(&count, query)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})

			It("returns an error if metadata missing ilk", func() {
				malformedIlkChopMetadata := storage.GetValueMetadata(cat.IlkChop, map[storage.Key]string{}, storage.Uint256)

				err := repo.Create(diffID, fakeHeaderID, malformedIlkChopMetadata, fakeAddress)
				Expect(err).To(MatchError(storage.ErrMetadataMalformed{MissingData: constants.Ilk}))
			})

			shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
				Repository:    &repo,
				Metadata:      storage.GetValueMetadata(cat.IlkChop, map[storage.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, storage.Uint256),
				PropertyName:  "Chop",
				PropertyValue: strconv.Itoa(rand.Int()),
				Schema:        constants.MakerSchema,
				TableName:     constants.CatIlkChopTable,
			})
		})

		Describe("Lump", func() {
			It("writes a row", func() {
				ilkLumpMetadata := storage.GetValueMetadata(cat.IlkLump, map[storage.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, storage.Uint256)

				err := repo.Create(diffID, fakeHeaderID, ilkLumpMetadata, fakeUint256)
				Expect(err).NotTo(HaveOccurred())

				var result MappingRes
				query := fmt.Sprintf(`SELECT diff_id, header_id, ilk_id AS key, lump AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.CatIlkLumpTable))
				err = db.Get(&result, query)
				Expect(err).NotTo(HaveOccurred())
				ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, diffID, fakeHeaderID, strconv.FormatInt(ilkID, 10), fakeUint256)
			})

			It("does not duplicate row", func() {
				ilkLumpMetadata := storage.GetValueMetadata(cat.IlkLump, map[storage.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, storage.Uint256)
				insertOneErr := repo.Create(diffID, fakeHeaderID, ilkLumpMetadata, fakeUint256)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(diffID, fakeHeaderID, ilkLumpMetadata, fakeUint256)

				Expect(insertTwoErr).NotTo(HaveOccurred())
				var count int
				query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.CatIlkLumpTable))
				getCountErr := db.Get(&count, query)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})

			It("returns an error if metadata missing ilk", func() {
				malformedIlkLumpMetadata := storage.GetValueMetadata(cat.IlkLump, map[storage.Key]string{}, storage.Uint256)

				err := repo.Create(diffID, fakeHeaderID, malformedIlkLumpMetadata, fakeAddress)
				Expect(err).To(MatchError(storage.ErrMetadataMalformed{MissingData: constants.Ilk}))
			})

			shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
				Repository:    &repo,
				Metadata:      storage.GetValueMetadata(cat.IlkLump, map[storage.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, storage.Uint256),
				PropertyName:  "Lump",
				PropertyValue: strconv.Itoa(rand.Int()),
				Schema:        constants.MakerSchema,
				TableName:     constants.CatIlkLumpTable,
			})
		})
	})
})
