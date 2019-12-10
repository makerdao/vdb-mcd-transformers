package queries

import (
	"math/rand"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/vat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("total ink query", func() {
	var (
		vatRepo                vat.VatStorageRepository
		headerRepo             datastore.HeaderRepository
		urnOne                 string
		urnTwo                 string
		blockOne, timestampOne int
		headerOne              core.Header
	)

	BeforeEach(func() {
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

	It("gets the latest ink of a single urn", func() {
		urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, urnOne)

		urnSetupDataOne := test_helpers.GetUrnSetupData()
		urnSetupDataOne[vat.UrnInk] = rand.Intn(1000000)
		test_helpers.CreateUrn(urnSetupDataOne, headerOne.Id, urnMetadata, vatRepo)

		headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
		urnSetupDataTwo := test_helpers.GetUrnSetupData()
		urnSetupDataTwo[vat.UrnInk] = rand.Intn(1000000)
		test_helpers.CreateUrn(urnSetupDataTwo, headerTwo.Id, urnMetadata, vatRepo)

		var totalInk int
		err := db.Get(&totalInk, `SELECT * FROM api.total_ink($1)`, test_helpers.FakeIlk.Identifier)
		Expect(err).NotTo(HaveOccurred())
		Expect(totalInk).To(Equal(urnSetupDataTwo[vat.UrnInk]))
	})

	It("sums up the latest ink of multiple urns for a given ilk", func() {
		urnOneMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, urnOne)

		urnOneSetupData := test_helpers.GetUrnSetupData()
		urnOneSetupData[vat.UrnInk] = rand.Intn(1000)
		test_helpers.CreateUrn(urnOneSetupData, headerOne.Id, urnOneMetadata, vatRepo)

		urnTwoMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, urnTwo)
		urnTwoOldSetupData := test_helpers.GetUrnSetupData()
		urnTwoOldSetupData[vat.UrnInk] = rand.Intn(1000)
		test_helpers.CreateUrn(urnTwoOldSetupData, headerOne.Id, urnTwoMetadata, vatRepo)

		headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
		urnTwoNewSetupData := test_helpers.GetUrnSetupData()
		urnTwoNewSetupData[vat.UrnInk] = rand.Intn(1000)
		test_helpers.CreateUrn(urnTwoNewSetupData, headerTwo.Id, urnTwoMetadata, vatRepo)

		var totalInk int
		err := db.Get(&totalInk, `SELECT * FROM api.total_ink($1)`, test_helpers.FakeIlk.Identifier)
		Expect(err).NotTo(HaveOccurred())
		Expect(totalInk).To(Equal(urnOneSetupData[vat.UrnInk] + urnTwoNewSetupData[vat.UrnInk]))
	})

	It("ignores ink after block number if block number is provided", func() {
		urnOneMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, urnOne)
		urnOneSetupData := test_helpers.GetUrnSetupData()
		urnOneSetupData[vat.UrnInk] = rand.Intn(1000)
		test_helpers.CreateUrn(urnOneSetupData, headerOne.Id, urnOneMetadata, vatRepo)

		urnTwoMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, urnTwo)
		urnTwoOldSetupData := test_helpers.GetUrnSetupData()
		urnTwoOldSetupData[vat.UrnInk] = rand.Intn(1000)
		test_helpers.CreateUrn(urnTwoOldSetupData, headerOne.Id, urnTwoMetadata, vatRepo)

		headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
		urnTwoNewSetupData := test_helpers.GetUrnSetupData()
		urnTwoNewSetupData[vat.UrnInk] = rand.Intn(1000)
		test_helpers.CreateUrn(urnTwoNewSetupData, headerTwo.Id, urnTwoMetadata, vatRepo)

		var totalInk int
		err := db.Get(&totalInk, `SELECT * FROM api.total_ink($1, $2)`,
			test_helpers.FakeIlk.Identifier, blockOne)
		Expect(err).NotTo(HaveOccurred())
		Expect(totalInk).To(Equal(urnOneSetupData[vat.UrnInk] + urnTwoOldSetupData[vat.UrnInk]))
	})

	It("ignores ink from urns of different ilks", func() {
		urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, urnOne)
		urnSetupData := test_helpers.GetUrnSetupData()
		urnSetupData[vat.UrnInk] = rand.Intn(1000)
		test_helpers.CreateUrn(urnSetupData, headerOne.Id, urnMetadata, vatRepo)

		irrelevantUrnMetadata := test_helpers.GetUrnMetadata(test_helpers.AnotherFakeIlk.Hex, urnTwo)
		irrelevantUrnSetupData := test_helpers.GetUrnSetupData()
		test_helpers.CreateUrn(irrelevantUrnSetupData, headerOne.Id, irrelevantUrnMetadata, vatRepo)

		var totalInk int
		err := db.Get(&totalInk, `SELECT * FROM api.total_ink($1)`, test_helpers.FakeIlk.Identifier)
		Expect(err).NotTo(HaveOccurred())
		Expect(totalInk).To(Equal(urnSetupData[vat.UrnInk]))
	})
})

func createHeader(blockNumber, timestamp int, headerRepo datastore.HeaderRepository) core.Header {
	header := fakes.GetFakeHeaderWithTimestamp(int64(timestamp), int64(blockNumber))

	var insertErr error
	header.Id, insertErr = headerRepo.CreateOrUpdateHeader(header)
	Expect(insertErr).NotTo(HaveOccurred())

	return header
}
