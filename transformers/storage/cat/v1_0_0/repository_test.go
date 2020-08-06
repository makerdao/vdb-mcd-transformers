package v1_0_0_test

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/cat/v1_0_0"
	. "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data/shared_behaviors"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cat storage repository", func() {
	var (
		db                   = test_config.NewTestDB(test_config.NewTestNode())
		repo                 v1_0_0.StorageRepository
		diffID, fakeHeaderID int64
		fakeAddress          = "0x" + fakes.RandomString(40)
		fakeUint256          = strconv.Itoa(rand.Intn(1000000))
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		repo = v1_0_0.StorageRepository{ContractAddress: test_data.CatAddress()}
		repo.SetDB(db)
		headerRepository := repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		fakeHeaderID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
		Expect(insertHeaderErr).NotTo(HaveOccurred())
	})

	Describe("Variable", func() {
		It("panics if the metadata name is not recognized", func() {
			unrecognizedMetadata := types.ValueMetadata{Name: "unrecognized"}
			repoCreate := func() {
				repo.Create(diffID, fakeHeaderID, unrecognizedMetadata, "")
			}

			Expect(repoCreate).Should(Panic())
		})

		Describe("Live", func() {
			liveMetadata := types.GetValueMetadata(v1_0_0.Live, nil, types.Uint256)
			inputs := shared_behaviors.StorageBehaviorInputs{
				ValueFieldName: v1_0_0.Live,
				Value:          fakeUint256,
				Schema:         constants.MakerSchema,
				TableName:      constants.CatLiveTable,
				Repository:     &repo,
				Metadata:       liveMetadata,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
		})

		Describe("Vat", func() {
			vatMetadata := types.GetValueMetadata(v1_0_0.Vat, nil, types.Address)
			inputs := shared_behaviors.StorageBehaviorInputs{
				ValueFieldName: v1_0_0.Vat,
				Value:          fakeAddress,
				Schema:         constants.MakerSchema,
				TableName:      constants.CatVatTable,
				Repository:     &repo,
				Metadata:       vatMetadata,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
		})

		Describe("Vow", func() {
			vowMetadata := types.GetValueMetadata(v1_0_0.Vow, nil, types.Address)
			inputs := shared_behaviors.StorageBehaviorInputs{
				ValueFieldName: v1_0_0.Vow,
				Value:          fakeAddress,
				Schema:         constants.MakerSchema,
				TableName:      constants.CatVowTable,
				Repository:     &repo,
				Metadata:       vowMetadata,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
		})
	})

	Describe("Wards mapping", func() {
		BeforeEach(func() {
			fakeRawDiff := GetFakeStorageDiffForHeader(fakes.FakeHeader, common.Hash{}, common.Hash{}, common.Hash{})
			storageDiffRepo := storage.NewDiffRepository(db)
			var insertDiffErr error
			diffID, insertDiffErr = storageDiffRepo.CreateStorageDiff(fakeRawDiff)
			Expect(insertDiffErr).NotTo(HaveOccurred())
		})

		It("writes a row", func() {
			fakeUserAddress := "0x" + fakes.RandomString(40)
			wardsMetadata := types.GetValueMetadata(wards.Wards, map[types.Key]string{constants.User: fakeUserAddress}, types.Uint256)

			setupErr := repo.Create(diffID, fakeHeaderID, wardsMetadata, fakeUint256)
			Expect(setupErr).NotTo(HaveOccurred())

			var result MappingResWithAddress
			query := fmt.Sprintf(`SELECT diff_id, header_id, address_id, usr AS key, wards AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.WardsTable))
			err := db.Get(&result, query)
			Expect(err).NotTo(HaveOccurred())
			contractAddressID, contractAddressErr := shared.GetOrCreateAddress(repo.ContractAddress, db)
			Expect(contractAddressErr).NotTo(HaveOccurred())
			userAddressID, userAddressErr := shared.GetOrCreateAddress(fakeUserAddress, db)
			Expect(userAddressErr).NotTo(HaveOccurred())
			Expect(result.AddressID).To(Equal(strconv.FormatInt(contractAddressID, 10)))
			AssertMapping(result.MappingRes, diffID, fakeHeaderID, strconv.FormatInt(userAddressID, 10), fakeUint256)
		})

		It("does not duplicate row", func() {
			fakeUserAddress := "0x" + fakes.RandomString(40)
			wardsMetadata := types.GetValueMetadata(wards.Wards, map[types.Key]string{constants.User: fakeUserAddress}, types.Uint256)
			insertOneErr := repo.Create(diffID, fakeHeaderID, wardsMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(diffID, fakeHeaderID, wardsMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.WardsTable))
			getCountErr := db.Get(&count, query)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns an error if metadata missing user", func() {
			malformedWardsMetadata := types.GetValueMetadata(wards.Wards, map[types.Key]string{}, types.Uint256)

			err := repo.Create(diffID, fakeHeaderID, malformedWardsMetadata, fakeUint256)
			Expect(err).To(MatchError(types.ErrMetadataMalformed{MissingData: constants.User}))
		})
	})

	Describe("Ilk", func() {
		BeforeEach(func() {
			fakeRawDiff := GetFakeStorageDiffForHeader(fakes.FakeHeader, common.Hash{}, common.Hash{}, common.Hash{})
			storageDiffRepo := storage.NewDiffRepository(db)
			var insertDiffErr error
			diffID, insertDiffErr = storageDiffRepo.CreateStorageDiff(fakeRawDiff)
			Expect(insertDiffErr).NotTo(HaveOccurred())
		})

		Describe("Flip", func() {
			It("writes a row", func() {
				ilkFlipMetadata := types.GetValueMetadata(v1_0_0.IlkFlip, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Address)

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
				ilkFlipMetadata := types.GetValueMetadata(v1_0_0.IlkFlip, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Address)
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
				malformedIlkFlipMetadata := types.GetValueMetadata(v1_0_0.IlkFlip, map[types.Key]string{}, types.Address)

				err := repo.Create(diffID, fakeHeaderID, malformedIlkFlipMetadata, fakeAddress)
				Expect(err).To(MatchError(types.ErrMetadataMalformed{MissingData: constants.Ilk}))
			})

			shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
				Repository:    &repo,
				Metadata:      types.GetValueMetadata(v1_0_0.IlkFlip, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Address),
				PropertyName:  "Flip",
				PropertyValue: fakeAddress,
				Schema:        constants.MakerSchema,
				TableName:     constants.CatIlkFlipTable,
			})
		})

		Describe("Chop", func() {
			It("writes a row", func() {
				ilkChopMetadata := types.GetValueMetadata(v1_0_0.IlkChop, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256)

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
				ilkChopMetadata := types.GetValueMetadata(v1_0_0.IlkChop, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256)
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
				malformedIlkChopMetadata := types.GetValueMetadata(v1_0_0.IlkChop, map[types.Key]string{}, types.Uint256)

				err := repo.Create(diffID, fakeHeaderID, malformedIlkChopMetadata, fakeAddress)
				Expect(err).To(MatchError(types.ErrMetadataMalformed{MissingData: constants.Ilk}))
			})

			shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
				Repository:    &repo,
				Metadata:      types.GetValueMetadata(v1_0_0.IlkChop, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256),
				PropertyName:  "Chop",
				PropertyValue: strconv.Itoa(rand.Int()),
				Schema:        constants.MakerSchema,
				TableName:     constants.CatIlkChopTable,
			})
		})

		Describe("Lump", func() {
			It("writes a row", func() {
				ilkLumpMetadata := types.GetValueMetadata(v1_0_0.IlkLump, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256)

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
				ilkLumpMetadata := types.GetValueMetadata(v1_0_0.IlkLump, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256)
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
				malformedIlkLumpMetadata := types.GetValueMetadata(v1_0_0.IlkLump, map[types.Key]string{}, types.Uint256)

				err := repo.Create(diffID, fakeHeaderID, malformedIlkLumpMetadata, fakeAddress)
				Expect(err).To(MatchError(types.ErrMetadataMalformed{MissingData: constants.Ilk}))
			})

			shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
				Repository:    &repo,
				Metadata:      types.GetValueMetadata(v1_0_0.IlkLump, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256),
				PropertyName:  "Lump",
				PropertyValue: strconv.Itoa(rand.Int()),
				Schema:        constants.MakerSchema,
				TableName:     constants.CatIlkLumpTable,
			})
		})
	})
})
