package clip_test

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip"
	. "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data/shared_behaviors"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Clip storage repository", func() {
	var (
		db                   = test_config.NewTestDB(test_config.NewTestNode())
		repo                 clip.StorageRepository
		diffID, fakeHeaderID int64
		fakeUint256          = strconv.Itoa(rand.Intn(1000000))
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		repo = clip.StorageRepository{ContractAddress: test_data.ClipLinkAV130Address()}
		repo.SetDB(db)
		headerRepository := repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		fakeHeaderID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
		Expect(insertHeaderErr).NotTo(HaveOccurred())
		diffID = rand.Int63()
	})

	Describe("Wards mapping", func() {
		BeforeEach(func() {
			fakeRawDiff := GetFakeStorageDiffForHeader(fakes.FakeHeader, common.Address{}, common.Hash{}, common.Hash{})
			storageDiffRepo := vdbStorage.NewDiffRepository(db)
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
			contractAddressID, contractAddressErr := repository.GetOrCreateAddress(db, repo.ContractAddress)
			Expect(contractAddressErr).NotTo(HaveOccurred())
			userAddressID, userAddressErr := repository.GetOrCreateAddress(db, fakeUserAddress)
			Expect(userAddressErr).NotTo(HaveOccurred())
			AssertMappingWithAddress(result, diffID, fakeHeaderID, contractAddressID, strconv.FormatInt(userAddressID, 10), fakeUint256)
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

	Describe("Static Storage Variables", func() {
		It("returns an error if the metadata name is not recognized", func() {
			unrecognizedMetadata := types.ValueMetadata{Name: "unrecognized"}

			err := repo.Create(diffID, fakeHeaderID, unrecognizedMetadata, "")

			Expect(err).Should(HaveOccurred())
		})

		Describe("clip dog", func() {
			BeforeEach(func() {
				diffID = CreateFakeDiffRecord(db)
			})

			It("writes a row", func() {
				dogMetadata := types.ValueMetadata{Name: clip.Dog}
				insertErr := repo.Create(diffID, fakeHeaderID, dogMetadata, FakeAddress)
				Expect(insertErr).NotTo(HaveOccurred())

				var result VariableRes
				query := fmt.Sprintf(`SELECT diff_id, header_id, dog AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.ClipDogTable))
				getErr := db.Get(&result, query)
				Expect(getErr).NotTo(HaveOccurred())
				addressID, addressErr := repository.GetOrCreateAddress(db, FakeAddress)
				Expect(addressErr).NotTo(HaveOccurred())
				AssertVariable(result, diffID, fakeHeaderID, strconv.FormatInt(addressID, 10))
			})

			It("does not duplicate a row", func() {
				dogMetadata := types.ValueMetadata{Name: clip.Dog}
				insertOneErr := repo.Create(diffID, fakeHeaderID, dogMetadata, FakeAddress)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(diffID, fakeHeaderID, dogMetadata, FakeAddress)
				Expect(insertTwoErr).NotTo(HaveOccurred())
				var count int
				query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.ClipDogTable))
				getCountErr := db.Get(&count, query)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))

			})
		})

		Describe("clip vow", func() {
			BeforeEach(func() {
				diffID = CreateFakeDiffRecord(db)
			})

			It("writes a row", func() {
				vowMetadata := types.ValueMetadata{Name: clip.Vow}
				insertErr := repo.Create(diffID, fakeHeaderID, vowMetadata, FakeAddress)
				Expect(insertErr).NotTo(HaveOccurred())

				var result VariableRes
				query := fmt.Sprintf(`SELECT diff_id, header_id, vow AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.ClipVowTable))
				getErr := db.Get(&result, query)
				Expect(getErr).NotTo(HaveOccurred())
				addressID, addressErr := repository.GetOrCreateAddress(db, FakeAddress)
				Expect(addressErr).NotTo(HaveOccurred())
				AssertVariable(result, diffID, fakeHeaderID, strconv.FormatInt(addressID, 10))
			})

			It("does not duplicate a row", func() {
				vowMetadata := types.ValueMetadata{Name: clip.Vow}
				insertOneErr := repo.Create(diffID, fakeHeaderID, vowMetadata, FakeAddress)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(diffID, fakeHeaderID, vowMetadata, FakeAddress)
				Expect(insertTwoErr).NotTo(HaveOccurred())
				var count int
				query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.ClipVowTable))
				getCountErr := db.Get(&count, query)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))

			})
		})

		Describe("clip spotter", func() {
			BeforeEach(func() {
				diffID = CreateFakeDiffRecord(db)
			})

			It("writes a row", func() {
				spotterMetadata := types.ValueMetadata{Name: clip.Spotter}
				insertErr := repo.Create(diffID, fakeHeaderID, spotterMetadata, FakeAddress)
				Expect(insertErr).NotTo(HaveOccurred())

				var result VariableRes
				query := fmt.Sprintf(`SELECT diff_id, header_id, spotter AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.ClipSpotterTable))
				getErr := db.Get(&result, query)
				Expect(getErr).NotTo(HaveOccurred())
				addressID, addressErr := repository.GetOrCreateAddress(db, FakeAddress)
				Expect(addressErr).NotTo(HaveOccurred())
				AssertVariable(result, diffID, fakeHeaderID, strconv.FormatInt(addressID, 10))
			})

			It("does not duplicate a row", func() {
				spotterMetadata := types.ValueMetadata{Name: clip.Spotter}
				insertOneErr := repo.Create(diffID, fakeHeaderID, spotterMetadata, FakeAddress)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(diffID, fakeHeaderID, spotterMetadata, FakeAddress)
				Expect(insertTwoErr).NotTo(HaveOccurred())
				var count int
				query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.ClipSpotterTable))
				getCountErr := db.Get(&count, query)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))

			})
		})

		Describe("clip calc", func() {
			BeforeEach(func() {
				diffID = CreateFakeDiffRecord(db)
			})

			It("writes a row", func() {
				calcMetadata := types.ValueMetadata{Name: clip.Calc}
				insertErr := repo.Create(diffID, fakeHeaderID, calcMetadata, FakeAddress)
				Expect(insertErr).NotTo(HaveOccurred())

				var result VariableRes
				query := fmt.Sprintf(`SELECT diff_id, header_id, calc AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.ClipCalcTable))
				getErr := db.Get(&result, query)
				Expect(getErr).NotTo(HaveOccurred())
				addressID, addressErr := repository.GetOrCreateAddress(db, FakeAddress)
				Expect(addressErr).NotTo(HaveOccurred())
				AssertVariable(result, diffID, fakeHeaderID, strconv.FormatInt(addressID, 10))
			})

			It("does not duplicate a row", func() {
				calcMetadata := types.ValueMetadata{Name: clip.Calc}
				insertOneErr := repo.Create(diffID, fakeHeaderID, calcMetadata, FakeAddress)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(diffID, fakeHeaderID, calcMetadata, FakeAddress)
				Expect(insertTwoErr).NotTo(HaveOccurred())
				var count int
				query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.ClipCalcTable))
				getCountErr := db.Get(&count, query)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))

			})
		})

		Describe("clip buf", func() {
			metadata := types.GetValueMetadata(clip.Buf, nil, types.Uint256)
			inputs := shared_behaviors.StorageBehaviorInputs{
				ValueFieldName:  clip.Buf,
				Value:           fakeUint256,
				Schema:          constants.MakerSchema,
				TableName:       constants.ClipBufTable,
				ContractAddress: test_data.ClipLinkAV130Address(),
				Repository:      &repo,
				Metadata:        metadata,
				IsAMapping:      false,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
		})

		Describe("clip tail", func() {
			metadata := types.GetValueMetadata(clip.Tail, nil, types.Uint256)
			inputs := shared_behaviors.StorageBehaviorInputs{
				ValueFieldName:  clip.Tail,
				Value:           fakeUint256,
				Schema:          constants.MakerSchema,
				TableName:       constants.ClipTailTable,
				ContractAddress: test_data.ClipLinkAV130Address(),
				Repository:      &repo,
				Metadata:        metadata,
				IsAMapping:      false,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
		})

		Describe("clip cusp", func() {
			metadata := types.GetValueMetadata(clip.Cusp, nil, types.Uint256)
			inputs := shared_behaviors.StorageBehaviorInputs{
				ValueFieldName:  clip.Cusp,
				Value:           fakeUint256,
				Schema:          constants.MakerSchema,
				TableName:       constants.ClipCuspTable,
				ContractAddress: test_data.ClipLinkAV130Address(),
				Repository:      &repo,
				Metadata:        metadata,
				IsAMapping:      false,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
		})

		Describe("clip chip and tip", func() {
			BeforeEach(func() {
				diffID = CreateFakeDiffRecord(db)
			})

			packedNames := make(map[int]string)
			packedNames[0] = clip.Chip
			packedNames[1] = clip.Tip
			var chipAndTipMetadata = types.ValueMetadata{
				Name:        clip.Packed,
				PackedNames: packedNames,
			}

			var fakeChip = strconv.Itoa(rand.Intn(100))
			var fakeTip = strconv.Itoa(rand.Intn(100))
			values := make(map[int]string)
			values[0] = fakeChip
			values[1] = fakeTip

			It("persists chip and tip records", func() {
				err := repo.Create(diffID, fakeHeaderID, chipAndTipMetadata, values)
				Expect(err).NotTo(HaveOccurred())

				var chipResult VariableRes
				chipQuery := fmt.Sprintf(`SELECT diff_id, header_id, chip AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.ClipChipTable))
				err = db.Get(&chipResult, chipQuery)
				Expect(err).NotTo(HaveOccurred())
				AssertVariable(chipResult, diffID, fakeHeaderID, fakeChip)

				var tipResult VariableRes
				tipQuery := fmt.Sprintf(`SELECT diff_id, header_id, tip AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.ClipTipTable))
				err = db.Get(&tipResult, tipQuery)
				Expect(err).NotTo(HaveOccurred())
				AssertVariable(tipResult, diffID, fakeHeaderID, fakeTip)

			})

			It("panics if the packed name is not recognized", func() {
				packedNames := make(map[int]string)
				packedNames[0] = "notRecognized"

				var badMetadata = types.ValueMetadata{
					Name:        clip.Packed,
					PackedNames: packedNames,
				}

				err := repo.Create(diffID, fakeHeaderID, badMetadata, values)
				Expect(err).Should(HaveOccurred())
			})

			It("returns an error if inserting fails", func() {
				badValues := make(map[int]string)
				badValues[0] = ""
				err := repo.Create(diffID, fakeHeaderID, chipAndTipMetadata, badValues)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("pq: invalid input syntax"))
			})
		})

		Describe("clip chost", func() {
			metadata := types.GetValueMetadata(clip.Chost, nil, types.Uint256)
			inputs := shared_behaviors.StorageBehaviorInputs{
				ValueFieldName:  clip.Chost,
				Value:           fakeUint256,
				Schema:          constants.MakerSchema,
				TableName:       constants.ClipChostTable,
				ContractAddress: test_data.ClipLinkAV130Address(),
				Repository:      &repo,
				Metadata:        metadata,
				IsAMapping:      false,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
		})

		Describe("clip kicks", func() {
			metadata := types.GetValueMetadata(clip.Kicks, nil, types.Uint256)
			inputs := shared_behaviors.StorageBehaviorInputs{
				ValueFieldName:  clip.Kicks,
				Value:           fakeUint256,
				Schema:          constants.MakerSchema,
				TableName:       constants.ClipKicksTable,
				ContractAddress: test_data.ClipLinkAV130Address(),
				Repository:      &repo,
				Metadata:        metadata,
				IsAMapping:      false,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
		})
	})

	Describe("Dynamic Sales storage field", func() {
		var (
			fakeUint256 = strconv.Itoa(rand.Intn(100))
			fakeSaleID  = strconv.Itoa(rand.Intn(100))
		)
		BeforeEach(func() {
			fakeRawDiff := GetFakeStorageDiffForHeader(fakes.FakeHeader, common.Address{}, common.Hash{}, common.Hash{})
			storageDiffRepo := vdbStorage.NewDiffRepository(db)
			var insertDiffErr error
			diffID, insertDiffErr = storageDiffRepo.CreateStorageDiff(fakeRawDiff)
			Expect(insertDiffErr).NotTo(HaveOccurred())
		})

		Describe("Pos", func() {
			It("writes a row", func() {
				salePosMetadata := types.GetValueMetadata(clip.SalePos, map[types.Key]string{constants.SaleId: fakeSaleID}, types.Uint256)

				err := repo.Create(diffID, fakeHeaderID, salePosMetadata, fakeUint256)
				Expect(err).NotTo(HaveOccurred())

				var result MappingResWithAddress
				query := fmt.Sprintf(`SELECT diff_id, header_id, address_id, sale_id AS key, pos AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.ClipSalePosTable))
				err = db.Get(&result, query)
				Expect(err).NotTo(HaveOccurred())
				contractAddressID, contractAddressErr := repository.GetOrCreateAddress(db, repo.ContractAddress)
				Expect(contractAddressErr).NotTo(HaveOccurred())
				AssertMappingWithAddress(result, diffID, fakeHeaderID, contractAddressID, fakeSaleID, fakeUint256)
			})

			It("does not duplicate row", func() {
				salePosMetadata := types.GetValueMetadata(clip.SalePos, map[types.Key]string{constants.SaleId: fakeSaleID}, types.Uint256)
				insertOneErr := repo.Create(diffID, fakeHeaderID, salePosMetadata, fakeUint256)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(diffID, fakeHeaderID, salePosMetadata, fakeUint256)

				Expect(insertTwoErr).NotTo(HaveOccurred())
				var count int
				query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.ClipSalePosTable))
				getCountErr := db.Get(&count, query)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})

			It("returns an error if metadata missing ilk", func() {
				malformedSalePosMetadata := types.GetValueMetadata(clip.SalePos, map[types.Key]string{}, types.Uint256)

				err := repo.Create(diffID, fakeHeaderID, malformedSalePosMetadata, fakeUint256)
				Expect(err).To(MatchError(types.ErrMetadataMalformed{MissingData: constants.SaleId}))
			})
		})

		Describe("Tab", func() {
			It("writes a row", func() {
				saleTabMetadata := types.GetValueMetadata(clip.SaleTab, map[types.Key]string{constants.SaleId: fakeSaleID}, types.Uint256)

				err := repo.Create(diffID, fakeHeaderID, saleTabMetadata, fakeUint256)
				Expect(err).NotTo(HaveOccurred())

				var result MappingResWithAddress
				query := fmt.Sprintf(`SELECT diff_id, header_id, address_id, sale_id AS key, tab AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.ClipSaleTabTable))
				err = db.Get(&result, query)
				Expect(err).NotTo(HaveOccurred())
				contractAddressID, contractAddressErr := repository.GetOrCreateAddress(db, repo.ContractAddress)
				Expect(contractAddressErr).NotTo(HaveOccurred())
				AssertMappingWithAddress(result, diffID, fakeHeaderID, contractAddressID, fakeSaleID, fakeUint256)
			})
		})
	})
})
