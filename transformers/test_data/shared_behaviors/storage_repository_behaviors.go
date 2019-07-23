package shared_behaviors

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/vulcanizedb/libraries/shared/repository"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"

	. "github.com/onsi/gomega"
)

type StorageVariableBehaviorInputs struct {
	FieldName        string
	Value            string
	BidId            string
	IsAMapping       bool
	StorageTableName string
	Repository       repository.StorageRepository
	Metadata         utils.StorageValueMetadata
}

func SharedStorageRepositoryVariableBehaviors(inputs *StorageVariableBehaviorInputs) {
	Describe("Create", func() {
		var (
			repository      = inputs.Repository
			fakeBlockNumber = rand.Int()
			fakeHash        = fakes.FakeHash.Hex()
			database        = test_config.NewTestDB(test_config.NewTestNode())
		)

		BeforeEach(func() {
			test_config.CleanTestDB(database)
			repository.SetDB(database)
		})

		It("persists a record", func() {
			err := repository.Create(fakeBlockNumber, fakeHash, inputs.Metadata, inputs.Value)
			Expect(err).NotTo(HaveOccurred())

			if inputs.IsAMapping == true {
				var result test_helpers.MappingRes
				query := fmt.Sprintf("SELECT block_number, block_hash, bid_id AS key, %s AS value FROM %s", inputs.FieldName, inputs.StorageTableName)
				err = database.Get(&result, query)
				Expect(err).NotTo(HaveOccurred())
				test_helpers.AssertMapping(result, fakeBlockNumber, fakeHash, inputs.BidId, inputs.Value)
			} else {
				var result test_helpers.VariableRes
				query := fmt.Sprintf("SELECT block_number, block_hash, %s AS value FROM %s", inputs.FieldName, inputs.StorageTableName)
				err = database.Get(&result, query)
				Expect(err).NotTo(HaveOccurred())
				test_helpers.AssertVariable(result, fakeBlockNumber, fakeHash, inputs.Value)
			}
		})

		It("doesn't duplicate a record", func() {
			err := repository.Create(fakeBlockNumber, fakeHash, inputs.Metadata, inputs.Value)
			Expect(err).NotTo(HaveOccurred())

			err = repository.Create(fakeBlockNumber, fakeHash, inputs.Metadata, inputs.Value)
			Expect(err).NotTo(HaveOccurred())

			var count int
			query := fmt.Sprintf("SELECT COUNT(*) FROM %s", inputs.StorageTableName)
			err = database.Get(&count, query)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})
	})
}
