package queries

import (
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vat"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("total ink query", func() {
	var (
		db         *postgres.DB
		vatRepo    vat.VatStorageRepository
		headerRepo repositories.HeaderRepository
		urnOne     string
		urnTwo     string
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		vatRepo = vat.VatStorageRepository{}
		vatRepo.SetDB(db)
		headerRepo = repositories.NewHeaderRepository(db)

		urnOne = test_data.RandomString(5)
		urnTwo = test_data.RandomString(5)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	It("gets the latest ink of a single urn", func() {
		urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, urnOne)

		blockNumberOne := rand.Int()
		timestampOne := int(rand.Int31())
		urnSetupDataOne := test_helpers.GetUrnSetupData(blockNumberOne, timestampOne)
		urnSetupDataOne.Ink = rand.Intn(1000000)
		test_helpers.CreateUrn(urnSetupDataOne, urnMetadata, vatRepo, headerRepo)

		blockNumberTwo := blockNumberOne + 1
		timestampTwo := timestampOne + 1
		urnSetupDataTwo := test_helpers.GetUrnSetupData(blockNumberTwo, timestampTwo)
		urnSetupDataTwo.Ink = rand.Intn(1000000)
		test_helpers.CreateUrn(urnSetupDataTwo, urnMetadata, vatRepo, headerRepo)

		var totalInk int
		err := db.Get(&totalInk, `SELECT * FROM api.total_ink($1)`, test_helpers.FakeIlk.Identifier)
		Expect(err).NotTo(HaveOccurred())
		Expect(totalInk).To(Equal(urnSetupDataTwo.Ink))
	})

	It("sums up the latest ink of multiple urns for a given ilk", func() {
		blockOne := rand.Int()
		timestampOne := int(rand.Int31())
		urnOneMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, urnOne)

		urnOneSetupData := test_helpers.GetUrnSetupData(blockOne, timestampOne)
		urnOneSetupData.Ink = rand.Intn(1000)
		test_helpers.CreateUrn(urnOneSetupData, urnOneMetadata, vatRepo, headerRepo)

		urnTwoMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, urnTwo)
		urnTwoOldSetupData := test_helpers.GetUrnSetupData(blockOne, timestampOne)
		urnTwoOldSetupData.Ink = rand.Intn(1000)
		test_helpers.CreateUrn(urnTwoOldSetupData, urnTwoMetadata, vatRepo, headerRepo)

		blockTwo := blockOne + 1
		timestampTwo := timestampOne + 1
		urnTwoNewSetupData := test_helpers.GetUrnSetupData(blockTwo, timestampTwo)
		urnTwoNewSetupData.Ink = rand.Intn(1000)
		test_helpers.CreateUrn(urnTwoNewSetupData, urnTwoMetadata, vatRepo, headerRepo)

		var totalInk int
		err := db.Get(&totalInk, `SELECT * FROM api.total_ink($1)`, test_helpers.FakeIlk.Identifier)
		Expect(err).NotTo(HaveOccurred())
		Expect(totalInk).To(Equal(urnOneSetupData.Ink + urnTwoNewSetupData.Ink))
	})

	It("ignores ink after block number if block number is provided", func() {
		blockOne := rand.Int()
		timestampOne := int(rand.Int31())

		urnOneMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, urnOne)
		urnOneSetupData := test_helpers.GetUrnSetupData(blockOne, timestampOne)
		urnOneSetupData.Ink = rand.Intn(1000)
		test_helpers.CreateUrn(urnOneSetupData, urnOneMetadata, vatRepo, headerRepo)

		urnTwoMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, urnTwo)
		urnTwoOldSetupData := test_helpers.GetUrnSetupData(blockOne, timestampOne)
		urnTwoOldSetupData.Ink = rand.Intn(1000)
		test_helpers.CreateUrn(urnTwoOldSetupData, urnTwoMetadata, vatRepo, headerRepo)

		blockTwo := blockOne + 1
		timestampTwo := timestampOne + 1
		urnTwoNewSetupData := test_helpers.GetUrnSetupData(blockTwo, timestampTwo)
		urnTwoNewSetupData.Ink = rand.Intn(1000)
		test_helpers.CreateUrn(urnTwoNewSetupData, urnTwoMetadata, vatRepo, headerRepo)

		var totalInk int
		err := db.Get(&totalInk, `SELECT * FROM api.total_ink($1, $2)`,
			test_helpers.FakeIlk.Identifier, blockOne)
		Expect(err).NotTo(HaveOccurred())
		Expect(totalInk).To(Equal(urnOneSetupData.Ink + urnTwoOldSetupData.Ink))
	})

	It("ignores ink from urns of different ilks", func() {
		blockOne := rand.Int()
		timestampOne := int(rand.Int31())

		urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, urnOne)
		urnSetupData := test_helpers.GetUrnSetupData(blockOne, timestampOne)
		urnSetupData.Ink = rand.Intn(1000)
		test_helpers.CreateUrn(urnSetupData, urnMetadata, vatRepo, headerRepo)

		irrelevantUrnMetadata := test_helpers.GetUrnMetadata(test_helpers.AnotherFakeIlk.Hex, urnTwo)
		irrelevantUrnSetupData := test_helpers.GetUrnSetupData(blockOne, timestampOne)
		test_helpers.CreateUrn(irrelevantUrnSetupData, irrelevantUrnMetadata, vatRepo, headerRepo)

		var totalInk int
		err := db.Get(&totalInk, `SELECT * FROM api.total_ink($1)`, test_helpers.FakeIlk.Identifier)
		Expect(err).NotTo(HaveOccurred())
		Expect(totalInk).To(Equal(urnSetupData.Ink))
	})
})
