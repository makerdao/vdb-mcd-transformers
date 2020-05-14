package median_test

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/median"
	. "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data/shared_behaviors"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Median Storage Repository", func() {
	var (
		db                   = test_config.NewTestDB(test_config.NewTestNode())
		repo                 median.MedianStorageRepository
		fakeUint256          = strconv.Itoa(rand.Intn(1000000))
		blockNumber          int64
		diffID, fakeHeaderID int64
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		repo = median.MedianStorageRepository{ContractAddress: test_data.MedianEthAddress()}
		repo.SetDB(db)
		blockNumber = rand.Int63()
		var insertHeaderErr error
		headerRepository := repositories.NewHeaderRepository(db)
		fakeHeaderID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(blockNumber))
		Expect(insertHeaderErr).NotTo(HaveOccurred())
		diffID = CreateFakeDiffRecord(db)
	})

	Describe("Wards mapping", func() {
		var fakeUint256 = strconv.Itoa(rand.Intn(1000000))

		It("writes a row", func() {
			fakeUserAddress := "0x" + fakes.RandomString(40)
			wardsMetadata := types.GetValueMetadata(wards.Wards, map[types.Key]string{constants.User: fakeUserAddress}, types.Uint256)

			writeErr := repo.Create(diffID, fakeHeaderID, wardsMetadata, fakeUint256)
			Expect(writeErr).NotTo(HaveOccurred())

			var result MappingResWithAddress
			query := fmt.Sprintf(`SELECT diff_id, header_id, address_id, usr AS key, wards AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.WardsTable))
			readErr := db.Get(&result, query)
			Expect(readErr).NotTo(HaveOccurred())
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

	Describe("val and age", func() {
		packedNames := make(map[int]string)
		packedNames[0] = median.Val
		packedNames[1] = median.Age
		var valAndAgeMetadata = types.ValueMetadata{
			Name:        storage.Packed,
			PackedNames: packedNames,
		}

		var fakeVal = strconv.Itoa(rand.Intn(100))
		var fakeAge = strconv.Itoa(rand.Intn(100))
		values := make(map[int]string)
		values[0] = fakeVal
		values[1] = fakeAge

		It("persists val and age records", func() {
			err := repo.Create(diffID, fakeHeaderID, valAndAgeMetadata, values)
			Expect(err).NotTo(HaveOccurred())

			var valResult VariableRes
			valQuery := fmt.Sprintf(`SELECT diff_id, header_id, val AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.MedianValTable))
			err = db.Get(&valResult, valQuery)
			Expect(err).NotTo(HaveOccurred())
			AssertVariable(valResult, diffID, fakeHeaderID, fakeVal)

			var ageResult VariableRes
			ageQuery := fmt.Sprintf(`SELECT diff_id, header_id, age AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.MedianAgeTable))
			err = db.Get(&ageResult, ageQuery)
			Expect(err).NotTo(HaveOccurred())
			AssertVariable(ageResult, diffID, fakeHeaderID, fakeAge)
		})

		It("panics if the packed name is not recognized", func() {
			packedNames := make(map[int]string)
			packedNames[0] = "notRecognized"

			var badMetadata = types.ValueMetadata{
				Name:        storage.Packed,
				PackedNames: packedNames,
			}

			createFunc := func() {
				repo.Create(diffID, fakeHeaderID, badMetadata, values)
			}
			Expect(createFunc).To(Panic())
		})

		It("returns an error if inserting fails", func() {
			badValues := make(map[int]string)
			badValues[0] = ""
			err := repo.Create(diffID, fakeHeaderID, valAndAgeMetadata, badValues)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("pq: invalid input syntax"))
		})
	})

	Describe("bar", func() {
		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName: median.Bar,
			Value:          fakeUint256,
			Schema:         constants.MakerSchema,
			TableName:      constants.MedianBarTable,
			Repository:     &repo,
			Metadata:       median.BarMetadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})

	Describe("bud mapping", func() {
		var fakeUint256 = strconv.Itoa(rand.Intn(1000000))

		It("writes a row", func() {
			fakeBudAddress := "0x" + fakes.RandomString(40)
			budMetadata := types.GetValueMetadata(median.Bud, map[types.Key]string{constants.A: fakeBudAddress}, types.Uint256)

			writeErr := repo.Create(diffID, fakeHeaderID, budMetadata, fakeUint256)
			Expect(writeErr).NotTo(HaveOccurred())

			var result MappingResWithAddress
			query := fmt.Sprintf(`SELECT diff_id, header_id, address_id, a AS key, bud AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.MedianBudTable))
			readErr := db.Get(&result, query)
			Expect(readErr).NotTo(HaveOccurred())
			contractAddressID, contractAddressErr := shared.GetOrCreateAddress(repo.ContractAddress, db)
			Expect(contractAddressErr).NotTo(HaveOccurred())
			budAddressID, budAddressErr := shared.GetOrCreateAddress(fakeBudAddress, db)
			Expect(budAddressErr).NotTo(HaveOccurred())
			Expect(result.AddressID).To(Equal(strconv.FormatInt(contractAddressID, 10)))
			AssertMapping(result.MappingRes, diffID, fakeHeaderID, strconv.FormatInt(budAddressID, 10), fakeUint256)
		})

		It("does not duplicate row", func() {
			fakeBudAddress := "0x" + fakes.RandomString(40)
			budMetadata := types.GetValueMetadata(median.Bud, map[types.Key]string{constants.A: fakeBudAddress}, types.Uint256)
			insertOneErr := repo.Create(diffID, fakeHeaderID, budMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(diffID, fakeHeaderID, budMetadata, fakeUint256)
			Expect(insertTwoErr).NotTo(HaveOccurred())

			var count int
			query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.MedianBudTable))
			getCountErr := db.Get(&count, query)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns an error if metadata missing 'a' address", func() {
			malformedBudMetadata := types.GetValueMetadata(median.Bud, map[types.Key]string{}, types.Uint256)

			err := repo.Create(diffID, fakeHeaderID, malformedBudMetadata, fakeUint256)
			Expect(err).To(MatchError(types.ErrMetadataMalformed{MissingData: constants.A}))
		})
	})

	Describe("orcl mapping", func() {
		var fakeUint256 = strconv.Itoa(rand.Intn(1000000))
		It("writes a row", func() {
			fakeOrclAddress := "0x" + fakes.RandomString(40)
			orclMetadata := types.GetValueMetadata(median.Orcl, map[types.Key]string{constants.Address: fakeOrclAddress}, types.Uint256)

			writeErr := repo.Create(diffID, fakeHeaderID, orclMetadata, fakeUint256)
			Expect(writeErr).NotTo(HaveOccurred())

			var result MappingResWithAddress
			query := fmt.Sprintf(`SELECT diff_id, header_id, address_id, a AS key, orcl AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.MedianOrclTable))
			readErr := db.Get(&result, query)
			Expect(readErr).NotTo(HaveOccurred())
			contractAddressID, contractAddressErr := shared.GetOrCreateAddress(repo.ContractAddress, db)
			Expect(contractAddressErr).NotTo(HaveOccurred())
			orclAddressID, orclAddressErr := shared.GetOrCreateAddress(fakeOrclAddress, db)
			Expect(orclAddressErr).NotTo(HaveOccurred())
			Expect(result.AddressID).To(Equal(strconv.FormatInt(contractAddressID, 10)))
			AssertMapping(result.MappingRes, diffID, fakeHeaderID, strconv.FormatInt(orclAddressID, 10), fakeUint256)
		})

		It("returns an error if metadata missing 'a' address", func() {
			malformedOrclMetadata := types.GetValueMetadata(median.Orcl, map[types.Key]string{}, types.Uint256)

			err := repo.Create(diffID, fakeHeaderID, malformedOrclMetadata, fakeUint256)
			Expect(err).To(MatchError(types.ErrMetadataMalformed{MissingData: constants.A}))
		})

	})
})
