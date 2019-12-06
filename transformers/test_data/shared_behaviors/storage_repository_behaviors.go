package shared_behaviors

import (
	"database/sql"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	. "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/utils"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
			repo             = inputs.Repository
			database         = test_config.NewTestDB(test_config.NewTestNode())
			diffID, headerID int64
		)

		BeforeEach(func() {
			test_config.CleanTestDB(database)
			repo.SetDB(database)
			headerRepository := repositories.NewHeaderRepository(database)
			var insertHeaderErr error
			headerID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(insertHeaderErr).NotTo(HaveOccurred())

			fakeRawDiff := fakes.GetFakeStorageDiffForHeader(fakes.FakeHeader, common.Hash{}, common.Hash{}, common.Hash{})
			storageDiffRepo := repositories.NewStorageDiffRepository(database)
			var insertDiffErr error
			diffID, insertDiffErr = storageDiffRepo.CreateStorageDiff(fakeRawDiff)
			Expect(insertDiffErr).NotTo(HaveOccurred())
		})

		It("persists a record", func() {
			err := repo.Create(diffID, headerID, inputs.Metadata, inputs.Value)
			Expect(err).NotTo(HaveOccurred())

			if inputs.IsAMapping == true {
				var result MappingRes
				query := fmt.Sprintf("SELECT header_id, %s AS key, %s AS value FROM %s",
					inputs.KeyFieldName, inputs.ValueFieldName, inputs.StorageTableName)
				err = database.Get(&result, query)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, diffID, headerID, inputs.Key, inputs.Value)
			} else {
				var result VariableRes
				query := fmt.Sprintf("SELECT diff_id, header_id, %s AS value FROM %s", inputs.ValueFieldName, inputs.StorageTableName)
				err = database.Get(&result, query)
				Expect(err).NotTo(HaveOccurred())
				AssertVariable(result, diffID, headerID, inputs.Value)
			}
		})

		It("doesn't duplicate a record", func() {
			err := repo.Create(diffID, headerID, inputs.Metadata, inputs.Value)
			Expect(err).NotTo(HaveOccurred())

			err = repo.Create(diffID, headerID, inputs.Metadata, inputs.Value)
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
	TableName     string
	PropertyName  string
	PropertyValue string
}

func SharedIlkTriggerTests(input IlkTriggerTestInput) {
	Describe("updating historical_ilk_state trigger table", func() {
		var (
			blockOne,
			blockTwo,
			blockThree int
			headerOne,
			headerTwo core.Header
			diffID           int64
			repo             = input.Repository
			database         = test_config.NewTestDB(test_config.NewTestNode())
			hashOne          = common.BytesToHash([]byte{1, 2, 3, 4, 5})
			hashTwo          = common.BytesToHash([]byte{5, 4, 3, 2, 1})
			hashThree        = common.BytesToHash([]byte{6, 7, 8, 9, 0})
			getStateQuery    = `SELECT ilk_identifier, block_number, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, updated FROM api.historical_ilk_state ORDER BY block_number`
			getFieldQuery    = fmt.Sprintf(`SELECT %s FROM api.historical_ilk_state ORDER BY block_number`, input.Metadata.Name)
			insertFieldQuery = fmt.Sprintf(`INSERT INTO api.historical_ilk_state (ilk_identifier, block_number, %s) VALUES ($1, $2, $3)`, input.Metadata.Name)
			deleteRowQuery   = fmt.Sprintf(`DELETE FROM %s WHERE header_id = $1`, input.TableName)
		)

		BeforeEach(func() {
			test_config.CleanTestDB(database)
			repo.SetDB(database)
			blockOne = rand.Int()
			blockTwo = blockOne + 1
			blockThree = blockTwo + 1
			rawTimestampOne := int64(rand.Int31())
			rawTimestampTwo := rawTimestampOne + 1
			rawTimestampThree := rawTimestampTwo + 1
			headerOne = CreateHeaderWithHash(hashOne.String(), rawTimestampOne, blockOne, database)
			headerTwo = CreateHeaderWithHash(hashTwo.String(), rawTimestampTwo, blockTwo, database)
			CreateHeaderWithHash(hashThree.String(), rawTimestampThree, blockThree, database)

			fakeRawDiff := fakes.GetFakeStorageDiffForHeader(fakes.FakeHeader, common.Hash{}, common.Hash{}, common.Hash{})
			storageDiffRepo := repositories.NewStorageDiffRepository(database)
			var insertDiffErr error
			diffID, insertDiffErr = storageDiffRepo.CreateStorageDiff(fakeRawDiff)
			Expect(insertDiffErr).NotTo(HaveOccurred())
		})

		It("inserts a row for new ilk-block", func() {
			initialIlkValues := test_helpers.GetIlkValues(0)
			test_helpers.CreateIlk(database, diffID, headerOne, initialIlkValues, test_helpers.FakeIlkVatMetadatas, test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas, test_helpers.FakeIlkSpotMetadatas)

			err := repo.Create(diffID, headerTwo.Id, input.Metadata, input.PropertyValue)
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
			test_helpers.CreateIlk(database, diffID, headerOne, initialIlkValues, test_helpers.FakeIlkVatMetadatas, test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas, test_helpers.FakeIlkSpotMetadatas)

			err := repo.Create(diffID, headerOne.Id, input.Metadata, input.PropertyValue)
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

			err := repo.Create(diffID, headerOne.Id, input.Metadata, input.PropertyValue)
			Expect(err).NotTo(HaveOccurred())

			var ilkStates []test_helpers.IlkState
			queryErr := database.Select(&ilkStates, getFieldQuery)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(len(ilkStates)).To(Equal(2))
			Expect(getIlkProperty(ilkStates[1], input.PropertyName)).To(Equal(input.PropertyValue))
		})

		It("ignores rows from blocks after the next time the field is updated", func() {
			initialIlkValues := test_helpers.GetIlkValues(0)
			setupErr := repo.Create(diffID, headerTwo.Id, input.Metadata, initialIlkValues[input.Metadata.Name])
			Expect(setupErr).NotTo(HaveOccurred())

			err := repo.Create(diffID, headerOne.Id, input.Metadata, input.PropertyValue)
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

			err := repo.Create(diffID, headerOne.Id, input.Metadata, input.PropertyValue)
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

			err := repo.Create(diffID, headerTwo.Id, input.Metadata, input.PropertyValue)
			Expect(err).NotTo(HaveOccurred())

			var ilkStates []test_helpers.IlkState
			queryErr := database.Select(&ilkStates, getFieldQuery)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(len(ilkStates)).To(Equal(2))
			Expect(getIlkProperty(ilkStates[0], input.PropertyName)).To(Equal(initialIlkValues[input.Metadata.Name]))
		})

		Describe("when diff is deleted", func() {
			It("updates field to previous value in subsequent rows", func() {
				initialIlkValues := test_helpers.GetIlkValues(0)
				setupErrOne := repo.Create(diffID, headerOne.Id, input.Metadata, initialIlkValues[input.Metadata.Name])
				Expect(setupErrOne).NotTo(HaveOccurred())

				subsequentIlkValues := test_helpers.GetIlkValues(1)
				setupErrTwo := repo.Create(diffID, headerTwo.Id, input.Metadata, subsequentIlkValues[input.Metadata.Name])
				Expect(setupErrTwo).NotTo(HaveOccurred())
				_, setupErrThree := database.Exec(insertFieldQuery,
					test_helpers.FakeIlk.Identifier, blockThree, subsequentIlkValues[input.Metadata.Name])
				Expect(setupErrThree).NotTo(HaveOccurred())

				_, deleteErr := database.Exec(deleteRowQuery, headerTwo.Id)
				Expect(deleteErr).NotTo(HaveOccurred())

				var ilkStates []test_helpers.IlkState
				queryErr := database.Select(&ilkStates, getFieldQuery)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(getIlkProperty(ilkStates[1], input.PropertyName)).To(Equal(initialIlkValues[input.Metadata.Name]))
			})

			It("sets field in subsequent rows to null if no previous diff exists", func() {
				initialIlkValues := test_helpers.GetIlkValues(0)
				setupErrOne := repo.Create(diffID, headerOne.Id, input.Metadata, initialIlkValues[input.Metadata.Name])
				Expect(setupErrOne).NotTo(HaveOccurred())
				_, setupErrTwo := database.Exec(insertFieldQuery,
					test_helpers.FakeIlk.Identifier, blockTwo, initialIlkValues[input.Metadata.Name])
				Expect(setupErrTwo).NotTo(HaveOccurred())

				_, deleteErr := database.Exec(deleteRowQuery, headerOne.Id)
				Expect(deleteErr).NotTo(HaveOccurred())

				var fieldValues []sql.NullString
				queryErr := database.Select(&fieldValues, getFieldQuery)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(fieldValues[0].Valid).To(BeFalse())
			})

			It("deletes ilk state associated with diff if identical to previous state", func() {
				initialIlkValues := test_helpers.GetIlkValues(0)
				setupErrOne := repo.Create(diffID, headerOne.Id, input.Metadata, initialIlkValues[input.Metadata.Name])
				Expect(setupErrOne).NotTo(HaveOccurred())

				subsequentIlkValues := test_helpers.GetIlkValues(1)
				setupErrTwo := repo.Create(diffID, headerTwo.Id, input.Metadata, subsequentIlkValues[input.Metadata.Name])
				Expect(setupErrTwo).NotTo(HaveOccurred())
				_, setupErrThree := database.Exec(insertFieldQuery,
					test_helpers.FakeIlk.Identifier, blockThree, subsequentIlkValues[input.Metadata.Name])
				Expect(setupErrThree).NotTo(HaveOccurred())

				_, deleteErr := database.Exec(deleteRowQuery, headerTwo.Id)
				Expect(deleteErr).NotTo(HaveOccurred())

				var ilkStates []test_helpers.IlkState
				queryErr := database.Select(&ilkStates, getFieldQuery)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(len(ilkStates)).To(Equal(2))
			})

			It("deletes ilk state associated with diff if it's the earliest state in the table", func() {
				initialIlkValues := test_helpers.GetIlkValues(0)
				setupErrOne := repo.Create(diffID, headerOne.Id, input.Metadata, initialIlkValues[input.Metadata.Name])
				Expect(setupErrOne).NotTo(HaveOccurred())

				_, deleteErr := database.Exec(deleteRowQuery, headerOne.Id)
				Expect(deleteErr).NotTo(HaveOccurred())

				var ilkStates []test_helpers.IlkState
				queryErr := database.Select(&ilkStates, getFieldQuery)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(len(ilkStates)).To(Equal(0))
			})
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
