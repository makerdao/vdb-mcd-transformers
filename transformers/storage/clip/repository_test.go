package clip_test

import (
	"math/rand"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/clip"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
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
		//fakeAddress          = "0x" + fakes.RandomString(40)
		//fakeUint256          = strconv.Itoa(rand.Intn(1000000))
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		repo = clip.StorageRepository{ContractAddress: test_data.Clip130Address()}
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
	})
})