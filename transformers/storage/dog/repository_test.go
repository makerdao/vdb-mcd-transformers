package dog_test

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	mcdShared "github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/dog"
	. "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data/shared_behaviors"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dog storage repository", func() {
	var (
		db                   = test_config.NewTestDB(test_config.NewTestNode())
		repo                 dog.StorageRepository
		diffID, fakeHeaderID int64
		fakeAddress          = "0x" + fakes.RandomString(40)
		fakeUint256          = strconv.Itoa(rand.Intn(1000000))
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		repo = dog.StorageRepository{ContractAddress: test_data.Dog130Address()}
		repo.SetDB(db)
		headerRepository := repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		fakeHeaderID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
		Expect(insertHeaderErr).NotTo(HaveOccurred())
		diffID = rand.Int63()
	})

	Describe("Static Storage Variables", func() {
		It("returns an error if the metadata name is not recognized", func() {
			unrecognizedMetadata := types.ValueMetadata{Name: "unrecognized"}

			err := repo.Create(diffID, fakeHeaderID, unrecognizedMetadata, "")

			Expect(err).Should(HaveOccurred())
		})

		Describe("dirt", func() {
			metadata := types.GetValueMetadata(dog.Dirt, nil, types.Uint256)
			inputs := shared_behaviors.StorageBehaviorInputs{
				ValueFieldName:  dog.Dirt,
				Value:           fakeUint256,
				Schema:          constants.MakerSchema,
				TableName:       constants.DogDirtTable,
				ContractAddress: test_data.Dog130Address(),
				Repository:      &repo,
				Metadata:        metadata,
				IsAMapping:      false,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
		})

		Describe("hole", func() {
			metadata := types.GetValueMetadata(dog.Hole, nil, types.Uint256)
			inputs := shared_behaviors.StorageBehaviorInputs{
				ValueFieldName:  dog.Hole,
				Value:           fakeUint256,
				Schema:          constants.MakerSchema,
				TableName:       constants.DogHoleTable,
				ContractAddress: test_data.Dog130Address(),
				Repository:      &repo,
				Metadata:        metadata,
				IsAMapping:      false,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
		})

		Describe("live", func() {
			metadata := types.GetValueMetadata(dog.Live, nil, types.Uint256)
			inputs := shared_behaviors.StorageBehaviorInputs{
				ValueFieldName:  dog.Live,
				Value:           fakeUint256,
				Schema:          constants.MakerSchema,
				TableName:       constants.DogLiveTable,
				ContractAddress: test_data.Dog130Address(),
				Repository:      &repo,
				Metadata:        metadata,
				IsAMapping:      false,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
		})

		Describe("vat", func() {
			var (
				diffID, headerID int64
				vatMetadata      = types.GetValueMetadata(dog.Vat, nil, types.Address)
			)

			BeforeEach(func() {
				headerRepository := repositories.NewHeaderRepository(db)
				var insertHeaderErr error
				headerID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
				Expect(insertHeaderErr).NotTo(HaveOccurred())

				diffID = CreateFakeDiffRecord(db)
			})

			It("persists a record", func() {
				err := repo.Create(diffID, headerID, vatMetadata, fakeAddress)
				Expect(err).NotTo(HaveOccurred())

				contractAddressID, contractAddressErr := repository.GetOrCreateAddress(db, test_data.Dog130Address())
				Expect(contractAddressErr).NotTo(HaveOccurred())

				vatAddressID, vatAddressErr := repository.GetOrCreateAddress(db, fakeAddress)
				Expect(vatAddressErr).NotTo(HaveOccurred())
				vatAddressIDString := strconv.FormatInt(vatAddressID, 10)

				var result VariableResWithAddress
				query := fmt.Sprintf("SELECT diff_id, header_id, address_id, %s AS value FROM %s", "vat", "maker.dog_vat")
				err = db.Get(&result, query)
				Expect(err).NotTo(HaveOccurred())

				AssertVariableWithAddress(result, diffID, headerID, contractAddressID, vatAddressIDString)
			})

			It("doesn't duplicate a record", func() {
				createErr := repo.Create(diffID, headerID, vatMetadata, fakeAddress)
				Expect(createErr).NotTo(HaveOccurred())

				createErr2 := repo.Create(diffID, headerID, vatMetadata, fakeAddress)
				Expect(createErr2).NotTo(HaveOccurred())

				var count int
				query := fmt.Sprintf("SELECT COUNT(*) FROM %s", "maker.dog_vat")
				countErr := db.Get(&count, query)
				Expect(countErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})
		})
	})

	Describe("vow", func() {
		var (
			diffID, headerID int64
			vowMetadata      = types.GetValueMetadata(dog.Vow, nil, types.Address)
		)

		BeforeEach(func() {
			headerRepository := repositories.NewHeaderRepository(db)
			var insertHeaderErr error
			headerID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(insertHeaderErr).NotTo(HaveOccurred())

			diffID = CreateFakeDiffRecord(db)
		})

		It("persists a record", func() {
			err := repo.Create(diffID, headerID, vowMetadata, fakeAddress)
			Expect(err).NotTo(HaveOccurred())

			contractAddressID, contractAddressErr := repository.GetOrCreateAddress(db, test_data.Dog130Address())
			Expect(contractAddressErr).NotTo(HaveOccurred())

			vowAddressID, vowAddressErr := repository.GetOrCreateAddress(db, fakeAddress)
			Expect(vowAddressErr).NotTo(HaveOccurred())
			vatAddressIDString := strconv.FormatInt(vowAddressID, 10)

			var result VariableResWithAddress
			query := fmt.Sprintf("SELECT diff_id, header_id, address_id, %s AS value FROM %s", "vow", "maker.dog_vow")
			err = db.Get(&result, query)
			Expect(err).NotTo(HaveOccurred())

			AssertVariableWithAddress(result, diffID, headerID, contractAddressID, vatAddressIDString)
		})

		It("doesn't duplicate a record", func() {
			createErr := repo.Create(diffID, headerID, vowMetadata, fakeAddress)
			Expect(createErr).NotTo(HaveOccurred())

			createErr2 := repo.Create(diffID, headerID, vowMetadata, fakeAddress)
			Expect(createErr2).NotTo(HaveOccurred())

			var count int
			query := fmt.Sprintf("SELECT COUNT(*) FROM %s", "maker.dog_vow")
			countErr := db.Get(&count, query)
			Expect(countErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})
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

	Describe("Ilks", func() {
		BeforeEach(func() {
			fakeRawDiff := GetFakeStorageDiffForHeader(fakes.FakeHeader, common.Address{}, common.Hash{}, common.Hash{})
			storageDiffRepo := storage.NewDiffRepository(db)
			var insertDiffErr error
			diffID, insertDiffErr = storageDiffRepo.CreateStorageDiff(fakeRawDiff)
			Expect(insertDiffErr).NotTo(HaveOccurred())
		})

		Describe("Clip", func() {
			It("writes a row", func() {
				ilkClipMetadata := types.GetValueMetadata(dog.IlkClip, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Address)

				err := repo.Create(diffID, fakeHeaderID, ilkClipMetadata, fakeAddress)
				Expect(err).NotTo(HaveOccurred())

				var result MappingResWithAddress
				query := fmt.Sprintf(`SELECT diff_id, header_id, address_id, ilk_id AS key, clip AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.DogIlkClipTable))
				err = db.Get(&result, query)
				Expect(err).NotTo(HaveOccurred())
				ilkID, err := mcdShared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
				Expect(err).NotTo(HaveOccurred())
				contractAddressID, contractAddressErr := repository.GetOrCreateAddress(db, repo.ContractAddress)
				Expect(contractAddressErr).NotTo(HaveOccurred())
				AssertMappingWithAddress(result, diffID, fakeHeaderID, contractAddressID, strconv.FormatInt(ilkID, 10), fakeAddress)
			})

			It("does not duplicate row", func() {
				ilkClipMetadata := types.GetValueMetadata(dog.IlkClip, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256)
				insertOneErr := repo.Create(diffID, fakeHeaderID, ilkClipMetadata, fakeUint256)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(diffID, fakeHeaderID, ilkClipMetadata, fakeUint256)

				Expect(insertTwoErr).NotTo(HaveOccurred())
				var count int
				query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.DogIlkClipTable))
				getCountErr := db.Get(&count, query)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})

			It("returns an error if metadata missing ilk", func() {
				malformedIlkClipMetadata := types.GetValueMetadata(dog.IlkClip, map[types.Key]string{}, types.Address)

				err := repo.Create(diffID, fakeHeaderID, malformedIlkClipMetadata, fakeAddress)
				Expect(err).To(MatchError(types.ErrMetadataMalformed{MissingData: constants.Ilk}))
			})
			//TODO: Add trigger table tests when refactoring snapshots
			//shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
			//	Repository:    &repo,
			//	Metadata:      types.GetValueMetadata(dog.IlkClip, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Address),
			//	PropertyName:  "Clip",
			//	PropertyValue: fakeAddress,
			//	Schema:        constants.MakerSchema,
			//	TableName:     constants.DogIlkClipTable,
			//})
		})

		Describe("Chop", func() {
			It("writes a row", func() {
				ilkChopMetadata := types.GetValueMetadata(dog.IlkChop, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256)

				err := repo.Create(diffID, fakeHeaderID, ilkChopMetadata, fakeUint256)
				Expect(err).NotTo(HaveOccurred())

				var result MappingResWithAddress
				query := fmt.Sprintf(`SELECT diff_id, header_id, address_id, ilk_id AS key, chop AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.DogIlkChopTable))
				err = db.Get(&result, query)
				Expect(err).NotTo(HaveOccurred())
				ilkID, err := mcdShared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
				Expect(err).NotTo(HaveOccurred())
				contractAddressID, contractAddressErr := repository.GetOrCreateAddress(db, repo.ContractAddress)
				Expect(contractAddressErr).NotTo(HaveOccurred())
				AssertMappingWithAddress(result, diffID, fakeHeaderID, contractAddressID, strconv.FormatInt(ilkID, 10), fakeUint256)
			})

			It("does not duplicate row", func() {
				ilkChopMetadata := types.GetValueMetadata(dog.IlkChop, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256)
				insertOneErr := repo.Create(diffID, fakeHeaderID, ilkChopMetadata, fakeUint256)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(diffID, fakeHeaderID, ilkChopMetadata, fakeUint256)

				Expect(insertTwoErr).NotTo(HaveOccurred())
				var count int
				query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.DogIlkChopTable))
				getCountErr := db.Get(&count, query)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})

			//TODO: Add trigger table tests when refactoring snapshots
			//shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
			//	Repository:    &repo,
			//	Metadata:      types.GetValueMetadata(dog.IlkChop, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256),
			//	PropertyName:  "Chop",
			//	PropertyValue: fakeAddress,
			//	Schema:        constants.MakerSchema,
			//	TableName:     constants.DogIlkChopTable,
			//})
		})
	})
})
