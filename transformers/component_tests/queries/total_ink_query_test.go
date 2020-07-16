package queries

import (
	"math/rand"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	query_test_helper "github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
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
		vatRepo                vat.StorageRepository
		headerRepo             datastore.HeaderRepository
		urnOne                 string
		urnTwo                 string
		blockOne, timestampOne int
		headerOne              core.Header
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		vatRepo = vat.StorageRepository{}
		vatRepo.SetDB(db)
		headerRepo = repositories.NewHeaderRepository(db)

		urnOne = test_data.RandomString(5)
		urnTwo = test_data.RandomString(5)

		blockOne = rand.Int()
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepo)
	})

	It("gets the latest ink of a single urn", func() {
		urnMetadata := query_test_helper.GetUrnMetadata(query_test_helper.FakeIlk.Hex, urnOne)

		urnSetupDataOne := query_test_helper.GetUrnSetupData()
		urnSetupDataOne[vat.UrnInk] = rand.Intn(1000000)
		query_test_helper.CreateUrn(db, urnSetupDataOne, headerOne, urnMetadata, vatRepo)

		headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
		urnSetupDataTwo := query_test_helper.GetUrnSetupData()
		urnSetupDataTwo[vat.UrnInk] = rand.Intn(1000000)
		query_test_helper.CreateUrn(db, urnSetupDataTwo, headerTwo, urnMetadata, vatRepo)

		var totalInk int
		err := db.Get(&totalInk, `SELECT * FROM api.total_ink($1)`, query_test_helper.FakeIlk.Identifier)
		Expect(err).NotTo(HaveOccurred())
		Expect(totalInk).To(Equal(urnSetupDataTwo[vat.UrnInk]))
	})

	It("sums up the latest ink of multiple urns for a given ilk", func() {
		urnOneMetadata := query_test_helper.GetUrnMetadata(query_test_helper.FakeIlk.Hex, urnOne)

		urnOneSetupData := query_test_helper.GetUrnSetupData()
		urnOneSetupData[vat.UrnInk] = rand.Intn(1000)
		query_test_helper.CreateUrn(db, urnOneSetupData, headerOne, urnOneMetadata, vatRepo)

		urnTwoMetadata := query_test_helper.GetUrnMetadata(query_test_helper.FakeIlk.Hex, urnTwo)
		urnTwoOldSetupData := query_test_helper.GetUrnSetupData()
		urnTwoOldSetupData[vat.UrnInk] = rand.Intn(1000)
		query_test_helper.CreateUrn(db, urnTwoOldSetupData, headerOne, urnTwoMetadata, vatRepo)

		headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
		urnTwoNewSetupData := query_test_helper.GetUrnSetupData()
		urnTwoNewSetupData[vat.UrnInk] = rand.Intn(1000)
		query_test_helper.CreateUrn(db, urnTwoNewSetupData, headerTwo, urnTwoMetadata, vatRepo)

		var totalInk int
		err := db.Get(&totalInk, `SELECT * FROM api.total_ink($1)`, query_test_helper.FakeIlk.Identifier)
		Expect(err).NotTo(HaveOccurred())
		Expect(totalInk).To(Equal(urnOneSetupData[vat.UrnInk].(int) + urnTwoNewSetupData[vat.UrnInk].(int)))
	})

	It("ignores ink after block number if block number is provided", func() {
		urnOneMetadata := query_test_helper.GetUrnMetadata(query_test_helper.FakeIlk.Hex, urnOne)
		urnOneSetupData := query_test_helper.GetUrnSetupData()
		urnOneSetupData[vat.UrnInk] = rand.Intn(1000)
		query_test_helper.CreateUrn(db, urnOneSetupData, headerOne, urnOneMetadata, vatRepo)

		urnTwoMetadata := query_test_helper.GetUrnMetadata(query_test_helper.FakeIlk.Hex, urnTwo)
		urnTwoOldSetupData := query_test_helper.GetUrnSetupData()
		urnTwoOldSetupData[vat.UrnInk] = rand.Intn(1000)
		query_test_helper.CreateUrn(db, urnTwoOldSetupData, headerOne, urnTwoMetadata, vatRepo)

		headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
		urnTwoNewSetupData := query_test_helper.GetUrnSetupData()
		urnTwoNewSetupData[vat.UrnInk] = rand.Intn(1000)
		query_test_helper.CreateUrn(db, urnTwoNewSetupData, headerTwo, urnTwoMetadata, vatRepo)

		var totalInk int
		err := db.Get(&totalInk, `SELECT * FROM api.total_ink($1, $2)`,
			query_test_helper.FakeIlk.Identifier, blockOne)
		Expect(err).NotTo(HaveOccurred())
		Expect(totalInk).To(Equal(urnOneSetupData[vat.UrnInk].(int) + urnTwoOldSetupData[vat.UrnInk].(int)))
	})

	It("ignores ink from urns of different ilks", func() {
		urnMetadata := query_test_helper.GetUrnMetadata(query_test_helper.FakeIlk.Hex, urnOne)
		urnSetupData := query_test_helper.GetUrnSetupData()
		urnSetupData[vat.UrnInk] = rand.Intn(1000)
		query_test_helper.CreateUrn(db, urnSetupData, headerOne, urnMetadata, vatRepo)

		irrelevantUrnMetadata := query_test_helper.GetUrnMetadata(query_test_helper.AnotherFakeIlk.Hex, urnTwo)
		irrelevantUrnSetupData := query_test_helper.GetUrnSetupData()
		query_test_helper.CreateUrn(db, irrelevantUrnSetupData, headerOne, irrelevantUrnMetadata, vatRepo)

		var totalInk int
		err := db.Get(&totalInk, `SELECT * FROM api.total_ink($1)`, query_test_helper.FakeIlk.Identifier)
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
