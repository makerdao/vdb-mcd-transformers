package queries

import (
	"math/rand"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/vat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("total ink query", func() {
	var (
		db                     *postgres.DB
		vatRepo                vat.VatStorageRepository
		headerRepo             datastore.HeaderRepository
		urnOne                 string
		urnTwo                 string
		blockOne, timestampOne int
		headerOne              core.Header
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		vatRepo = vat.VatStorageRepository{}
		vatRepo.SetDB(db)
		headerRepo = repositories.NewHeaderRepository(db)

		urnOne = test_data.RandomString(5)
		urnTwo = test_data.RandomString(5)

		blockOne = rand.Int()
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepo)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	It("gets the latest ink of a single urn", func() {
		urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, urnOne)

		urnSetupDataOne := test_helpers.GetUrnSetupData(headerOne)
		urnSetupDataOne.Ink = rand.Intn(1000000)
		test_helpers.CreateUrn(urnSetupDataOne, urnMetadata, vatRepo)

		headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
		urnSetupDataTwo := test_helpers.GetUrnSetupData(headerTwo)
		urnSetupDataTwo.Ink = rand.Intn(1000000)
		test_helpers.CreateUrn(urnSetupDataTwo, urnMetadata, vatRepo)

		var totalInk int
		err := db.Get(&totalInk, `SELECT * FROM api.total_ink($1)`, test_helpers.FakeIlk.Identifier)
		Expect(err).NotTo(HaveOccurred())
		Expect(totalInk).To(Equal(urnSetupDataTwo.Ink))
	})

	It("sums up the latest ink of multiple urns for a given ilk", func() {
		urnOneMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, urnOne)

		urnOneSetupData := test_helpers.GetUrnSetupData(headerOne)
		urnOneSetupData.Ink = rand.Intn(1000)
		test_helpers.CreateUrn(urnOneSetupData, urnOneMetadata, vatRepo)

		urnTwoMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, urnTwo)
		urnTwoOldSetupData := test_helpers.GetUrnSetupData(headerOne)
		urnTwoOldSetupData.Ink = rand.Intn(1000)
		test_helpers.CreateUrn(urnTwoOldSetupData, urnTwoMetadata, vatRepo)

		headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
		urnTwoNewSetupData := test_helpers.GetUrnSetupData(headerTwo)
		urnTwoNewSetupData.Ink = rand.Intn(1000)
		test_helpers.CreateUrn(urnTwoNewSetupData, urnTwoMetadata, vatRepo)

		var totalInk int
		err := db.Get(&totalInk, `SELECT * FROM api.total_ink($1)`, test_helpers.FakeIlk.Identifier)
		Expect(err).NotTo(HaveOccurred())
		Expect(totalInk).To(Equal(urnOneSetupData.Ink + urnTwoNewSetupData.Ink))
	})

	It("ignores ink after block number if block number is provided", func() {
		urnOneMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, urnOne)
		urnOneSetupData := test_helpers.GetUrnSetupData(headerOne)
		urnOneSetupData.Ink = rand.Intn(1000)
		test_helpers.CreateUrn(urnOneSetupData, urnOneMetadata, vatRepo)

		urnTwoMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, urnTwo)
		urnTwoOldSetupData := test_helpers.GetUrnSetupData(headerOne)
		urnTwoOldSetupData.Ink = rand.Intn(1000)
		test_helpers.CreateUrn(urnTwoOldSetupData, urnTwoMetadata, vatRepo)

		headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
		urnTwoNewSetupData := test_helpers.GetUrnSetupData(headerTwo)
		urnTwoNewSetupData.Ink = rand.Intn(1000)
		test_helpers.CreateUrn(urnTwoNewSetupData, urnTwoMetadata, vatRepo)

		var totalInk int
		err := db.Get(&totalInk, `SELECT * FROM api.total_ink($1, $2)`,
			test_helpers.FakeIlk.Identifier, blockOne)
		Expect(err).NotTo(HaveOccurred())
		Expect(totalInk).To(Equal(urnOneSetupData.Ink + urnTwoOldSetupData.Ink))
	})

	It("ignores ink from urns of different ilks", func() {
		urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, urnOne)
		urnSetupData := test_helpers.GetUrnSetupData(headerOne)
		urnSetupData.Ink = rand.Intn(1000)
		test_helpers.CreateUrn(urnSetupData, urnMetadata, vatRepo)

		irrelevantUrnMetadata := test_helpers.GetUrnMetadata(test_helpers.AnotherFakeIlk.Hex, urnTwo)
		irrelevantUrnSetupData := test_helpers.GetUrnSetupData(headerOne)
		test_helpers.CreateUrn(irrelevantUrnSetupData, irrelevantUrnMetadata, vatRepo)

		var totalInk int
		err := db.Get(&totalInk, `SELECT * FROM api.total_ink($1)`, test_helpers.FakeIlk.Identifier)
		Expect(err).NotTo(HaveOccurred())
		Expect(totalInk).To(Equal(urnSetupData.Ink))
	})
})

func createHeader(blockNumber, timestamp int, headerRepo datastore.HeaderRepository) core.Header {
	header := fakes.GetFakeHeaderWithTimestamp(int64(timestamp), int64(blockNumber))

	var insertErr error
	header.Id, insertErr = headerRepo.CreateOrUpdateHeader(header)
	Expect(insertErr).NotTo(HaveOccurred())

	return header
}
