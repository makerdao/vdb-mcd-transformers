package shared_behaviors

import (
	"fmt"
	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/storage"
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
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
				var result test_helpers.MappingRes
				query := fmt.Sprintf("SELECT block_number, block_hash, %s AS key, %s AS value FROM %s",
					inputs.KeyFieldName, inputs.ValueFieldName, inputs.StorageTableName)
				err = database.Get(&result, query)
				Expect(err).NotTo(HaveOccurred())
				test_helpers.AssertMapping(result, fakeBlockNumber, fakeHash, inputs.Key, inputs.Value)
			} else {
				var result test_helpers.VariableRes
				query := fmt.Sprintf("SELECT block_number, block_hash, %s AS value FROM %s", inputs.ValueFieldName, inputs.StorageTableName)
				err = database.Get(&result, query)
				Expect(err).NotTo(HaveOccurred())
				test_helpers.AssertVariable(result, fakeBlockNumber, fakeHash, inputs.Value)
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
