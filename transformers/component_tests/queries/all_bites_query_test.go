package queries

import (
	"math/rand"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/bite"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("Bites query", func() {
	var (
		db         *postgres.DB
		biteRepo   bite.BiteRepository
		headerRepo repositories.HeaderRepository
		fakeUrn    = test_data.RandomString(5)
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		biteRepo = bite.BiteRepository{}
		biteRepo.SetDB(db)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("all_bites", func() {
		It("returns bites for an ilk", func() {
			headerOne := fakes.GetFakeHeader(1)
			headerOneId, err := headerRepo.CreateOrUpdateHeader(headerOne)
			Expect(err).NotTo(HaveOccurred())

			biteOne := generateBite(test_helpers.FakeIlk.Hex, fakeUrn)
			err = biteRepo.Create(headerOneId, []interface{}{biteOne})
			Expect(err).NotTo(HaveOccurred())

			var actualBites []test_helpers.BiteEvent
			err = db.Select(&actualBites, `SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.all_bites($1)`,
				test_helpers.FakeIlk.Identifier)
			Expect(err).NotTo(HaveOccurred())

			Expect(actualBites).To(ConsistOf(
				test_helpers.BiteEvent{IlkIdentifier: test_helpers.FakeIlk.Identifier, UrnIdentifier: fakeUrn, Ink: biteOne.Ink, Art: biteOne.Art, Tab: biteOne.Tab},
			))
		})

		It("returns bites from multiple blocks", func() {
			headerOne := fakes.GetFakeHeaderWithTimestamp(int64(111111111), 1)
			headerOneId, err := headerRepo.CreateOrUpdateHeader(headerOne)
			Expect(err).NotTo(HaveOccurred())

			biteBlockOne := generateBite(test_helpers.FakeIlk.Hex, fakeUrn)
			err = biteRepo.Create(headerOneId, []interface{}{biteBlockOne})
			Expect(err).NotTo(HaveOccurred())

			// New block
			headerTwo := fakes.GetFakeHeaderWithTimestamp(int64(222222222), 2)
			headerTwo.Hash = "anotherHash"
			headerTwoId, err := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(err).NotTo(HaveOccurred())

			biteBlockTwo := generateBite(test_helpers.FakeIlk.Hex, fakeUrn)
			err = biteRepo.Create(headerTwoId, []interface{}{biteBlockTwo})
			Expect(err).NotTo(HaveOccurred())

			var actualBites []test_helpers.BiteEvent
			err = db.Select(&actualBites, `SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.all_bites($1)`, test_helpers.FakeIlk.Identifier)
			Expect(err).NotTo(HaveOccurred())

			Expect(actualBites).To(ConsistOf(
				test_helpers.BiteEvent{IlkIdentifier: test_helpers.FakeIlk.Identifier, UrnIdentifier: fakeUrn, Ink: biteBlockTwo.Ink, Art: biteBlockTwo.Art, Tab: biteBlockTwo.Tab},
				test_helpers.BiteEvent{IlkIdentifier: test_helpers.FakeIlk.Identifier, UrnIdentifier: fakeUrn, Ink: biteBlockOne.Ink, Art: biteBlockOne.Art, Tab: biteBlockOne.Tab},
			))
		})

		It("ignores bites from irrelevant ilks", func() {
			headerOne := fakes.GetFakeHeader(1)
			headerOneId, err := headerRepo.CreateOrUpdateHeader(headerOne)
			Expect(err).NotTo(HaveOccurred())

			bite := generateBite(test_helpers.FakeIlk.Hex, fakeUrn)
			irrelevantBite := generateBite(test_helpers.AnotherFakeIlk.Hex, fakeUrn)
			irrelevantBite.TransactionIndex = bite.TransactionIndex + 1

			err = biteRepo.Create(headerOneId, []interface{}{bite, irrelevantBite})
			Expect(err).NotTo(HaveOccurred())

			var actualBites []test_helpers.BiteEvent
			err = db.Select(&actualBites, `SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.all_bites($1)`, test_helpers.FakeIlk.Identifier)
			Expect(err).NotTo(HaveOccurred())

			Expect(actualBites).To(ConsistOf(
				test_helpers.BiteEvent{IlkIdentifier: test_helpers.FakeIlk.Identifier, UrnIdentifier: fakeUrn, Ink: bite.Ink, Art: bite.Art, Tab: bite.Tab},
			))
		})

		It("fails if no argument is supplied (STRICT)", func() {
			_, err := db.Exec(`SELECT * FROM api.all_bites()`)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("function api.all_bites() does not exist"))
		})

		Describe("result pagination", func() {
			var biteBlockOne, biteBlockTwo bite.BiteModel

			BeforeEach(func() {
				headerOne := fakes.GetFakeHeaderWithTimestamp(int64(111111111), 1)
				headerOneId, headerOneErr := headerRepo.CreateOrUpdateHeader(headerOne)
				Expect(headerOneErr).NotTo(HaveOccurred())

				biteBlockOne = generateBite(test_helpers.FakeIlk.Hex, fakeUrn)
				biteOneErr := biteRepo.Create(headerOneId, []interface{}{biteBlockOne})
				Expect(biteOneErr).NotTo(HaveOccurred())

				// New block
				headerTwo := fakes.GetFakeHeaderWithTimestamp(int64(222222222), 2)
				headerTwo.Hash = "anotherHash"
				headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
				Expect(headerTwoErr).NotTo(HaveOccurred())

				biteBlockTwo = generateBite(test_helpers.FakeIlk.Hex, fakeUrn)
				biteTwoErr := biteRepo.Create(headerTwoId, []interface{}{biteBlockTwo})
				Expect(biteTwoErr).NotTo(HaveOccurred())
			})

			It("limits results to most recent blocks if max_results argument is provided", func() {
				maxResults := 1
				var actualBites []test_helpers.BiteEvent
				err := db.Select(&actualBites, `SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.all_bites($1, $2)`,
					test_helpers.FakeIlk.Identifier, maxResults)
				Expect(err).NotTo(HaveOccurred())

				Expect(actualBites).To(ConsistOf(
					test_helpers.BiteEvent{
						IlkIdentifier: test_helpers.FakeIlk.Identifier,
						UrnIdentifier: fakeUrn,
						Ink:           biteBlockTwo.Ink,
						Art:           biteBlockTwo.Art,
						Tab:           biteBlockTwo.Tab,
					},
				))
			})

			It("offsets results if offset is provided", func() {
				maxResults := 1
				resultOffset := 1
				var actualBites []test_helpers.BiteEvent
				err := db.Select(&actualBites, `SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.all_bites($1, $2, $3)`,
					test_helpers.FakeIlk.Identifier, maxResults, resultOffset)
				Expect(err).NotTo(HaveOccurred())

				Expect(actualBites).To(ConsistOf(
					test_helpers.BiteEvent{
						IlkIdentifier: test_helpers.FakeIlk.Identifier,
						UrnIdentifier: fakeUrn,
						Ink:           biteBlockOne.Ink,
						Art:           biteBlockOne.Art,
						Tab:           biteBlockOne.Tab,
					},
				))
			})
		})
	})

	Describe("urn_bites", func() {
		It("returns bites for relevant ilk/urn", func() {
			headerOne := fakes.GetFakeHeader(1)

			headerOneId, err := headerRepo.CreateOrUpdateHeader(headerOne)
			Expect(err).NotTo(HaveOccurred())

			biteOne := generateBite(test_helpers.FakeIlk.Hex, fakeUrn)
			err = biteRepo.Create(headerOneId, []interface{}{biteOne})
			Expect(err).NotTo(HaveOccurred())

			var actualBites []test_helpers.BiteEvent
			err = db.Select(&actualBites, `SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.urn_bites($1, $2)`, test_helpers.FakeIlk.Identifier, fakeUrn)
			Expect(err).NotTo(HaveOccurred())

			Expect(actualBites).To(ConsistOf(
				test_helpers.BiteEvent{IlkIdentifier: test_helpers.FakeIlk.Identifier, UrnIdentifier: fakeUrn, Ink: biteOne.Ink, Art: biteOne.Art, Tab: biteOne.Tab},
			))
		})

		It("returns bites from multiple blocks", func() {
			headerOne := fakes.GetFakeHeaderWithTimestamp(int64(111111111), 1)
			headerOneId, err := headerRepo.CreateOrUpdateHeader(headerOne)
			Expect(err).NotTo(HaveOccurred())

			biteBlockOne := generateBite(test_helpers.FakeIlk.Hex, fakeUrn)
			err = biteRepo.Create(headerOneId, []interface{}{biteBlockOne})
			Expect(err).NotTo(HaveOccurred())

			// New block
			headerTwo := fakes.GetFakeHeaderWithTimestamp(int64(222222222), 2)
			headerTwo.Hash = "anotherHash"
			headerTwoId, err := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(err).NotTo(HaveOccurred())

			biteBlockTwo := generateBite(test_helpers.FakeIlk.Hex, fakeUrn)
			err = biteRepo.Create(headerTwoId, []interface{}{biteBlockTwo})
			Expect(err).NotTo(HaveOccurred())

			var actualBites []test_helpers.BiteEvent
			err = db.Select(&actualBites, `SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.urn_bites($1, $2)`, test_helpers.FakeIlk.Identifier, fakeUrn)
			Expect(err).NotTo(HaveOccurred())

			Expect(actualBites).To(ConsistOf(
				test_helpers.BiteEvent{IlkIdentifier: test_helpers.FakeIlk.Identifier, UrnIdentifier: fakeUrn, Ink: biteBlockTwo.Ink, Art: biteBlockTwo.Art, Tab: biteBlockTwo.Tab},
				test_helpers.BiteEvent{IlkIdentifier: test_helpers.FakeIlk.Identifier, UrnIdentifier: fakeUrn, Ink: biteBlockOne.Ink, Art: biteBlockOne.Art, Tab: biteBlockOne.Tab},
			))
		})

		Describe("result pagination", func() {
			var biteBlockOne, biteBlockTwo bite.BiteModel

			BeforeEach(func() {
				headerOne := fakes.GetFakeHeaderWithTimestamp(int64(111111111), 1)
				headerOneId, err := headerRepo.CreateOrUpdateHeader(headerOne)
				Expect(err).NotTo(HaveOccurred())

				biteBlockOne = generateBite(test_helpers.FakeIlk.Hex, fakeUrn)
				err = biteRepo.Create(headerOneId, []interface{}{biteBlockOne})
				Expect(err).NotTo(HaveOccurred())

				// New block
				headerTwo := fakes.GetFakeHeaderWithTimestamp(int64(222222222), 2)
				headerTwo.Hash = "anotherHash"
				headerTwoId, err := headerRepo.CreateOrUpdateHeader(headerTwo)
				Expect(err).NotTo(HaveOccurred())

				biteBlockTwo = generateBite(test_helpers.FakeIlk.Hex, fakeUrn)
				err = biteRepo.Create(headerTwoId, []interface{}{biteBlockTwo})
				Expect(err).NotTo(HaveOccurred())
			})

			It("limits results to latest block if max_results argument is provided", func() {
				maxResults := 1
				var actualBites []test_helpers.BiteEvent
				err := db.Select(&actualBites, `SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.urn_bites($1, $2, $3)`,
					test_helpers.FakeIlk.Identifier, fakeUrn, maxResults)
				Expect(err).NotTo(HaveOccurred())

				Expect(actualBites).To(ConsistOf(
					test_helpers.BiteEvent{IlkIdentifier: test_helpers.FakeIlk.Identifier, UrnIdentifier: fakeUrn, Ink: biteBlockTwo.Ink, Art: biteBlockTwo.Art, Tab: biteBlockTwo.Tab},
				))
			})

			It("offsets results if offset is provided", func() {
				maxResults := 1
				resultOffset := 1
				var actualBites []test_helpers.BiteEvent
				err := db.Select(&actualBites, `SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.urn_bites($1, $2, $3, $4)`,
					test_helpers.FakeIlk.Identifier, fakeUrn, maxResults, resultOffset)
				Expect(err).NotTo(HaveOccurred())

				Expect(actualBites).To(ConsistOf(
					test_helpers.BiteEvent{IlkIdentifier: test_helpers.FakeIlk.Identifier, UrnIdentifier: fakeUrn, Ink: biteBlockOne.Ink, Art: biteBlockOne.Art, Tab: biteBlockOne.Tab},
				))
			})
		})

		It("ignores bites from irrelevant urns", func() {
			headerOne := fakes.GetFakeHeader(1)
			headerOneId, err := headerRepo.CreateOrUpdateHeader(headerOne)
			Expect(err).NotTo(HaveOccurred())

			bite := generateBite(test_helpers.FakeIlk.Hex, fakeUrn)
			irrelevantBite := generateBite(test_helpers.FakeIlk.Hex, "irrelevantUrn")
			irrelevantBite.TransactionIndex = bite.TransactionIndex + 1

			err = biteRepo.Create(headerOneId, []interface{}{bite, irrelevantBite})
			Expect(err).NotTo(HaveOccurred())

			var actualBites []test_helpers.BiteEvent
			err = db.Select(&actualBites, `SELECT ilk_identifier, urn_identifier, ink, art, tab FROM api.urn_bites($1, $2)`, test_helpers.FakeIlk.Identifier, fakeUrn)
			Expect(err).NotTo(HaveOccurred())

			Expect(actualBites).To(ConsistOf(
				test_helpers.BiteEvent{IlkIdentifier: test_helpers.FakeIlk.Identifier, UrnIdentifier: fakeUrn, Ink: bite.Ink, Art: bite.Art, Tab: bite.Tab},
			))
		})
	})
})

func generateBite(ilk, urn string) bite.BiteModel {
	biteEvent := test_data.BiteModel
	biteEvent.Ilk = ilk
	biteEvent.Urn = urn
	biteEvent.Ink = strconv.Itoa(rand.Int())
	biteEvent.Art = strconv.Itoa(rand.Int())
	biteEvent.Tab = strconv.Itoa(rand.Int())
	biteEvent.Id = strconv.Itoa(rand.Int())
	return biteEvent
}
