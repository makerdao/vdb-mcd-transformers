package medianizer_test

import (
	"fmt"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/medianizer"
	. "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
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

var _ = Describe("Medianizer Storage Repository", func() {
	var (
		db                   = test_config.NewTestDB(test_config.NewTestNode())
		repo                 medianizer.MedianizerStorageRepository
		fakeUint256          = strconv.Itoa(rand.Intn(1000000))
		blockNumber          int64
		diffID, fakeHeaderID int64
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		repo = medianizer.MedianizerStorageRepository{ContractAddress: test_data.MedianizerAddress()}
		repo.SetDB(db)
		blockNumber = rand.Int63()
		var insertHeaderErr error
		headerRepository := repositories.NewHeaderRepository(db)
		fakeHeaderID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(blockNumber))
		Expect(insertHeaderErr).NotTo(HaveOccurred())
		diffID = CreateFakeDiffRecord(db)
	})

	Describe("val and age", func() {
		packedNames := make(map[int]string)
		packedNames[0] = medianizer.Val
		packedNames[1] = medianizer.Age
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
			valQuery := fmt.Sprintf(`SELECT diff_id, header_id, val AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.MedianizerValTable))
			err = db.Get(&valResult, valQuery)
			Expect(err).NotTo(HaveOccurred())
			AssertVariable(valResult, diffID, fakeHeaderID, fakeVal)

			var ageResult VariableRes
			ageQuery := fmt.Sprintf(`SELECT diff_id, header_id, age AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.MedianizerAgeTable))
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
			ValueFieldName: medianizer.Bar,
			Value:          fakeUint256,
			Schema:         constants.MakerSchema,
			TableName:      constants.MedianizerBarTable,
			Repository:     &repo,
			Metadata:       medianizer.BarMetadata,
		}

		shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
	})
})
