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
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
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
		//fakeAddress          = "0x" + fakes.RandomString(40)
		fakeUint256 = strconv.Itoa(rand.Intn(1000000))
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
	})
})
