package medianizer_test

import (
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/medianizer"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data/shared_behaviors"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Medianizer Storage Repository", func() {
	var (
		db          = test_config.NewTestDB(test_config.NewTestNode())
		repo        medianizer.MedianizerStorageRepository
		fakeAddress = "0x" + fakes.RandomString(20)
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		repo = medianizer.MedianizerStorageRepository{ContractAddress: test_data.MedianizerAddress()}
		repo.SetDB(db)
	})

	Describe("val", func() {
		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName: medianizer.Val,
			Value:          fakeAddress,
			Schema:         constants.MakerSchema,
			TableName:      constants.MedianizerValTable,
			Repository:     &repo,
			Metadata:       medianizer.ValMetadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})

	Describe("age", func() {
		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName: medianizer.Age,
			Value:          fakeAddress,
			Schema:         constants.MakerSchema,
			TableName:      constants.MedianizerAgeTable,
			Repository:     &repo,
			Metadata:       medianizer.AgeMetadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})

	Describe("bar", func() {
		inputs := shared_behaviors.StorageBehaviorInputs{
			ValueFieldName: medianizer.Bar,
			Value:          fakeAddress,
			Schema:         constants.MakerSchema,
			TableName:      constants.MedianizerBarTable,
			Repository:     &repo,
			Metadata:       medianizer.BarMetadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})
})
