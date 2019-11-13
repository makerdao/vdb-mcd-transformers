package shared_behaviors

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	. "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
	"reflect"
	"strconv"
)

type StorageVariableBehaviorInputs struct {
	KeyFieldName     string
	ValueFieldName   string
	Key              string
	Value            string
	IsAMapping       bool
	StorageTableName string
	Repository       storage.Repository
	Metadata         utils.StorageValueMetadata
}

func SharedStorageRepositoryVariableBehaviors(inputs *StorageVariableBehaviorInputs) {
	Describe("Create", func() {
		var (
			repo            = inputs.Repository
			fakeBlockNumber = rand.Int()
			fakeHash        = fakes.FakeHash.Hex()
			database        = test_config.NewTestDB(test_config.NewTestNode())
		)

		BeforeEach(func() {
			test_config.CleanTestDB(database)
			repo.SetDB(database)
		})

		It("persists a record", func() {
			err := repo.Create(fakeBlockNumber, fakeHash, inputs.Metadata, inputs.Value)
			Expect(err).NotTo(HaveOccurred())

			if inputs.IsAMapping == true {
				var result MappingRes
				query := fmt.Sprintf("SELECT block_number, block_hash, %s AS key, %s AS value FROM %s",
					inputs.KeyFieldName, inputs.ValueFieldName, inputs.StorageTableName)
				err = database.Get(&result, query)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, fakeBlockNumber, fakeHash, inputs.Key, inputs.Value)
			} else {
				var result VariableRes
				query := fmt.Sprintf("SELECT block_number, block_hash, %s AS value FROM %s", inputs.ValueFieldName, inputs.StorageTableName)
				err = database.Get(&result, query)
				Expect(err).NotTo(HaveOccurred())
				AssertVariable(result, fakeBlockNumber, fakeHash, inputs.Value)
			}
		})

		It("doesn't duplicate a record", func() {
			err := repo.Create(fakeBlockNumber, fakeHash, inputs.Metadata, inputs.Value)
			Expect(err).NotTo(HaveOccurred())

			err = repo.Create(fakeBlockNumber, fakeHash, inputs.Metadata, inputs.Value)
			Expect(err).NotTo(HaveOccurred())

			var count int
			query := fmt.Sprintf("SELECT COUNT(*) FROM %s", inputs.StorageTableName)
			err = database.Get(&count, query)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})
	})
}

type IlkTriggerTestInput struct {
	Repository    storage.Repository
	Metadata      utils.StorageValueMetadata
	PropertyName  string
	PropertyValue string
}

func SharedIlkTriggerTests(input IlkTriggerTestInput) {
	Describe("updating historical_ilk_state trigger table", func() {
		var (
			blockOne,
			blockTwo int
			headerOne,
			headerTwo core.Header
			repo             = input.Repository
			database         = test_config.NewTestDB(test_config.NewTestNode())
			hashOne          = common.BytesToHash([]byte{1, 2, 3, 4, 5})
			hashTwo          = common.BytesToHash([]byte{5, 4, 3, 2, 1})
			getStateQuery    = `SELECT ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, updated FROM api.historical_ilk_state ORDER BY block_number`
			getFieldQuery    = fmt.Sprintf(`SELECT %s FROM api.historical_ilk_state ORDER BY block_number`, input.Metadata.Name)
			insertFieldQuery = fmt.Sprintf(`INSERT INTO api.historical_ilk_state (ilk_identifier, block_number, %s) VALUES ($1, $2, $3)`, input.Metadata.Name)
		)

		BeforeEach(func() {
			test_config.CleanTestDB(database)
			repo.SetDB(database)
			blockOne = rand.Int()
			blockTwo = blockOne + 1
			rawTimestampOne := int64(rand.Int31())
			rawTimestampTwo := rawTimestampOne + 1
			headerOne = CreateHeaderWithHash(hashOne.String(), rawTimestampOne, blockOne, database)
			headerTwo = CreateHeaderWithHash(hashTwo.String(), rawTimestampTwo, blockTwo, database)
		})

		It("inserts a row for new ilk-block", func() {
			initialIlkValues := test_helpers.GetIlkValues(0)
			test_helpers.CreateIlk(database, headerOne, initialIlkValues, test_helpers.FakeIlkVatMetadatas,
				test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas, test_helpers.FakeIlkSpotMetadatas)

			err := repo.Create(blockTwo, hashTwo.String(), input.Metadata, input.PropertyValue)
			Expect(err).NotTo(HaveOccurred())

			var ilkStates []test_helpers.IlkState
			queryErr := database.Select(&ilkStates, getStateQuery)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(len(ilkStates)).To(Equal(2))
			initialIlkValues[input.Metadata.Name] = input.PropertyValue
			expectedIlk := test_helpers.IlkStateFromValues(test_helpers.FakeIlk.Hex,
				headerTwo.Timestamp, headerOne.Timestamp, initialIlkValues)
			assertIlk(ilkStates[1], expectedIlk, headerTwo.BlockNumber)
		})

		It("updates row if ilk-block combination already exists in table", func() {
			initialIlkValues := test_helpers.GetIlkValues(0)
			test_helpers.CreateIlk(database, headerOne, initialIlkValues, test_helpers.FakeIlkVatMetadatas,
				test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas, test_helpers.FakeIlkSpotMetadatas)

			err := repo.Create(blockOne, hashOne.String(), input.Metadata, input.PropertyValue)
			Expect(err).NotTo(HaveOccurred())

			var ilkStates []test_helpers.IlkState
			queryErr := database.Select(&ilkStates, getStateQuery)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(len(ilkStates)).To(Equal(1))
			initialIlkValues[input.Metadata.Name] = input.PropertyValue
			expectedIlk := test_helpers.IlkStateFromValues(test_helpers.FakeIlk.Hex,
				headerOne.Timestamp, headerOne.Timestamp, initialIlkValues)
			assertIlk(ilkStates[0], expectedIlk, headerOne.BlockNumber)
		})

		It("updates field in subsequent blocks", func() {
			initialIlkValues := test_helpers.GetIlkValues(0)
			_, setupErr := database.Exec(insertFieldQuery,
				test_helpers.FakeIlk.Identifier, headerTwo.BlockNumber, initialIlkValues[input.Metadata.Name])
			Expect(setupErr).NotTo(HaveOccurred())

			err := repo.Create(blockOne, hashOne.String(), input.Metadata, input.PropertyValue)
			Expect(err).NotTo(HaveOccurred())

			var ilkStates []test_helpers.IlkState
			queryErr := database.Select(&ilkStates, getFieldQuery)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(len(ilkStates)).To(Equal(2))
			Expect(getIlkProperty(ilkStates[1], input.PropertyName)).To(Equal(input.PropertyValue))
		})

		It("ignores rows from blocks after the next time the field is updated", func() {
			initialIlkValues := test_helpers.GetIlkValues(0)
			setupErr := repo.Create(blockTwo, hashTwo.String(), input.Metadata, initialIlkValues[input.Metadata.Name])
			Expect(setupErr).NotTo(HaveOccurred())

			err := repo.Create(blockOne, hashOne.String(), input.Metadata, input.PropertyValue)
			Expect(err).NotTo(HaveOccurred())

			var ilkStates []test_helpers.IlkState
			queryErr := database.Select(&ilkStates, getFieldQuery)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(len(ilkStates)).To(Equal(2))
			Expect(getIlkProperty(ilkStates[1], input.PropertyName)).To(Equal(initialIlkValues[input.Metadata.Name]))
		})

		It("ignores rows from different ilk", func() {
			initialIlkValues := test_helpers.GetIlkValues(0)
			_, setupErr := database.Exec(insertFieldQuery,
				test_helpers.AnotherFakeIlk.Identifier, headerTwo.BlockNumber, initialIlkValues[input.Metadata.Name])
			Expect(setupErr).NotTo(HaveOccurred())

			err := repo.Create(blockOne, hashOne.String(), input.Metadata, input.PropertyValue)
			Expect(err).NotTo(HaveOccurred())

			var ilkStates []test_helpers.IlkState
			queryErr := database.Select(&ilkStates, getFieldQuery)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(len(ilkStates)).To(Equal(2))
			Expect(getIlkProperty(ilkStates[1], input.PropertyName)).To(Equal(initialIlkValues[input.Metadata.Name]))
		})

		It("ignores rows from earlier blocks", func() {
			initialIlkValues := test_helpers.GetIlkValues(0)
			_, setupErr := database.Exec(insertFieldQuery,
				test_helpers.FakeIlk.Identifier, headerOne.BlockNumber, initialIlkValues[input.Metadata.Name])
			Expect(setupErr).NotTo(HaveOccurred())

			err := repo.Create(blockTwo, hashTwo.String(), input.Metadata, input.PropertyValue)
			Expect(err).NotTo(HaveOccurred())

			var ilkStates []test_helpers.IlkState
			queryErr := database.Select(&ilkStates, getFieldQuery)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(len(ilkStates)).To(Equal(2))
			Expect(getIlkProperty(ilkStates[0], input.PropertyName)).To(Equal(initialIlkValues[input.Metadata.Name]))
		})
	})
}

func getIlkProperty(ilk test_helpers.IlkState, fieldName string) string {
	r := reflect.ValueOf(ilk)
	property := reflect.Indirect(r).FieldByName(fieldName)
	return property.String()
}

func assertIlk(actualIlk test_helpers.IlkState, expectedIlk test_helpers.IlkState, expectedBlockNumber int64) {
	Expect(actualIlk.IlkIdentifier).To(Equal(expectedIlk.IlkIdentifier))
	Expect(actualIlk.BlockNumber).To(Equal(strconv.FormatInt(expectedBlockNumber, 10)))
	Expect(actualIlk.Rate).To(Equal(expectedIlk.Rate))
	Expect(actualIlk.Art).To(Equal(expectedIlk.Art))
	Expect(actualIlk.Spot).To(Equal(expectedIlk.Spot))
	Expect(actualIlk.Line).To(Equal(expectedIlk.Line))
	Expect(actualIlk.Dust).To(Equal(expectedIlk.Dust))
	Expect(actualIlk.Chop).To(Equal(expectedIlk.Chop))
	Expect(actualIlk.Lump).To(Equal(expectedIlk.Lump))
	Expect(actualIlk.Flip).To(Equal(expectedIlk.Flip))
	Expect(actualIlk.Rho).To(Equal(expectedIlk.Rho))
	Expect(actualIlk.Duty).To(Equal(expectedIlk.Duty))
	Expect(actualIlk.Pip).To(Equal(expectedIlk.Pip))
	Expect(actualIlk.Mat).To(Equal(expectedIlk.Mat))
	Expect(actualIlk.Updated).To(Equal(expectedIlk.Updated))
}
