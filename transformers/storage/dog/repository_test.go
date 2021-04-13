package dog_test

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/dog"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
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
		repo = dog.StorageRepository{ContractAddress: test_data.Dog1xxAddress()}
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
				ContractAddress: test_data.Dog1xxAddress(),
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
				ContractAddress: test_data.Dog1xxAddress(),
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
				ContractAddress: test_data.Dog1xxAddress(),
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

				diffID = test_helpers.CreateFakeDiffRecord(db)
			})

			It("persists a record", func() {
				err := repo.Create(diffID, headerID, vatMetadata, fakeAddress)
				Expect(err).NotTo(HaveOccurred())

				contractAddressID, contractAddressErr := repository.GetOrCreateAddress(db, test_data.Dog1xxAddress())
				Expect(contractAddressErr).NotTo(HaveOccurred())

				vatAddressID, vatAddressErr := repository.GetOrCreateAddress(db, fakeAddress)
				Expect(vatAddressErr).NotTo(HaveOccurred())
				vatAddressIDString := strconv.FormatInt(vatAddressID, 10)

				var result test_helpers.VariableResWithAddress
				query := fmt.Sprintf("SELECT diff_id, header_id, address_id, %s AS value FROM %s", "vat", "maker.dog_vat")
				err = db.Get(&result, query)
				Expect(err).NotTo(HaveOccurred())

				test_helpers.AssertVariableWithAddress(result, diffID, headerID, contractAddressID, vatAddressIDString)
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

			diffID = test_helpers.CreateFakeDiffRecord(db)
		})

		It("persists a record", func() {
			err := repo.Create(diffID, headerID, vowMetadata, fakeAddress)
			Expect(err).NotTo(HaveOccurred())

			contractAddressID, contractAddressErr := repository.GetOrCreateAddress(db, test_data.Dog1xxAddress())
			Expect(contractAddressErr).NotTo(HaveOccurred())

			vowAddressID, vowAddressErr := repository.GetOrCreateAddress(db, fakeAddress)
			Expect(vowAddressErr).NotTo(HaveOccurred())
			vatAddressIDString := strconv.FormatInt(vowAddressID, 10)

			var result test_helpers.VariableResWithAddress
			query := fmt.Sprintf("SELECT diff_id, header_id, address_id, %s AS value FROM %s", "vow", "maker.dog_vow")
			err = db.Get(&result, query)
			Expect(err).NotTo(HaveOccurred())

			test_helpers.AssertVariableWithAddress(result, diffID, headerID, contractAddressID, vatAddressIDString)
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
})
