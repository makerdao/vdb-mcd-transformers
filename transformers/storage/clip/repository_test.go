package clip_test

import (
	"fmt"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	mcdShared "github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip"
	. "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data/shared_behaviors"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Clip storage repository", func() {
	var (
		db                   = test_config.NewTestDB(test_config.NewTestNode())
		repo                 clip.StorageRepository
		diffID, fakeHeaderID int64
		FakeAddress          = "0x" + fakes.RandomString(40)
		//fakeUint256          = strconv.Itoa(rand.Intn(1000000))
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		repo = clip.StorageRepository{ContractAddress: test_data.ClipLinkAV130Address()}
		repo.SetDB(db)
		headerRepository := repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		fakeHeaderID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
		Expect(insertHeaderErr).NotTo(HaveOccurred())
		diffID = CreateFakeDiffRecord(db)
	})

	Describe("Static Storage Variables", func() {
		It("returns an error if the metadata name is not recognized", func() {
			unrecognizedMetadata := types.ValueMetadata{Name: "unrecognized"}

			err := repo.Create(diffID, fakeHeaderID, unrecognizedMetadata, "")

			Expect(err).Should(HaveOccurred())
		})

		XIt("rolls back the record and address insertions if there's a failure", func() {
			var cuspMetadata = types.ValueMetadata{Name: storage.Cusp}
			err := repo.Create(diffID, fakeHeaderID, cuspMetadata, "")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp("pq: invalid input syntax for type numeric"))

			var addressCount int
			countErr := db.Get(&addressCount, `SELECT COUNT(*) FROM addresses`)
			Expect(countErr).NotTo(HaveOccurred())
			Expect(addressCount).To(Equal(0))
		})

		Describe("Static Variable", func() {
			Describe("ilk", func() {
				It("writes a row", func() {
					ilkMetadata := types.ValueMetadata{Name: storage.Ilk}
					insertErr := repo.Create(diffID, fakeHeaderID, ilkMetadata, FakeIlk)
					Expect(insertErr).NotTo(HaveOccurred())

					var result VariableRes
					query := fmt.Sprintf(`SELECT diff_id, header_id, ilk_id AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.ClipIlkTable))
					getErr := db.Get(&result, query)
					Expect(getErr).NotTo(HaveOccurred())
					ilkID, ilkErr := mcdShared.GetOrCreateIlk(FakeIlk, db)
					Expect(ilkErr).NotTo(HaveOccurred())

					AssertVariable(result, diffID, fakeHeaderID, strconv.FormatInt(ilkID, 10))
				})

				It("does not duplicate row", func() {
					ilkMetadata := types.ValueMetadata{Name: storage.Ilk}
					insertOneErr := repo.Create(diffID, fakeHeaderID, ilkMetadata, FakeIlk)
					Expect(insertOneErr).NotTo(HaveOccurred())

					insertTwoErr := repo.Create(diffID, fakeHeaderID, ilkMetadata, FakeIlk)
					Expect(insertTwoErr).NotTo(HaveOccurred())

					var count int
					query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.ClipIlkTable))
					getCountErr := db.Get(&count, query)
					Expect(getCountErr).NotTo(HaveOccurred())
					Expect(count).To(Equal(1))
				})
			})

			Describe("Vat", func() {
				vatMetadata := types.ValueMetadata{Name: storage.Vat}
				inputs := shared_behaviors.StorageBehaviorInputs{
					ValueFieldName: storage.Vat,
					Value:          FakeAddress,
					Schema:         constants.MakerSchema,
					TableName:      constants.ClipVatTable,
					Repository:     &repo,
					Metadata:       vatMetadata,
				}

				shared_behaviors.SharedStorageRepositoryBehaviors(&inputs)
			})
		})
	})
})
