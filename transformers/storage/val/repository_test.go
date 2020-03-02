package val_test

import (
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/val"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data/shared_behaviors"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"math/rand"
	"strconv"
)

var _ = Describe("Val storage repository", func() {
	var (
		db                   = test_config.NewTestDB(test_config.NewTestNode())
		repo                 val.ValStorageRepository
		diffID, fakeHeaderID int64
		fakeAddress          = "0x" + fakes.RandomString(40)
		fakeUint256          = strconv.Itoa(rand.Intn(1000000))
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		repo = val.ValStorageRepository{ContractAddress: test_data.ValAddress()}
		repo.SetDB(db)
		headerRepository := repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		fakeHeaderID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
		Expect(insertHeaderErr).NotTo(HaveOccurred())
	})

	Describe("Variable", func() {
		It("panics if the metadata name is not recognized", func() {
			unrecognizedMetadata := types.ValueMetadata{Name: "unrecognized"}
			repoCreate := func() {
				repo.Create(diffID, fakeHeaderID, unrecognizedMetadata, "")
			}

			Expect(repoCreate).Should(Panic())
		})

		Describe("Has", func() {
			hasMetadata := types.GetValueMetadata(val.Has, nil, types.Uint256)
			inputs := shared_behaviors.StorageBehaviorInputs{
				ValueFieldName: val.Has,
				Value:          fakeUint256,
				Schema:         constants.MakerSchema,
				TableName:      constants.ValHasTable,
				Repository:     &repo,
				Metadata:       hasMetadata,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
		})

		Describe("Val", func() {
			valMetadata := types.GetValueMetadata(val.Val, nil, types.Bytes32)
			inputs := shared_behaviors.StorageBehaviorInputs{
				ValueFieldName: val.Val,
				Value:          fakeAddress,
				Schema:         constants.MakerSchema,
				TableName:      constants.ValValTable,
				Repository:     &repo,
				Metadata:       valMetadata,
			}

			shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
		})
	})
})
