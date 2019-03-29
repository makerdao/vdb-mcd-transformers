package queries

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_frob"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
	"strconv"
)

var _ = Describe("Frobs query", func() {
	var (
		db         *postgres.DB
		frobRepo   vat_frob.VatFrobRepository
		headerRepo repositories.HeaderRepository
		fakeIlk    = test_data.RandomString(5)
		fakeUrn    = test_data.RandomString(5)
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		frobRepo = vat_frob.VatFrobRepository{}
		frobRepo.SetDB(db)
	})

	It("returns frobs for relevant ilk/urn", func() {
		headerOne := fakes.GetFakeHeader(1)

		headerOneId, err := headerRepo.CreateOrUpdateHeader(headerOne)
		Expect(err).NotTo(HaveOccurred())

		frobBlockOne := test_data.VatFrobModel
		frobBlockOne.Ilk = fakeIlk
		frobBlockOne.Urn = fakeUrn
		frobBlockOne.Dink = strconv.Itoa(rand.Int())
		frobBlockOne.Dart = strconv.Itoa(rand.Int())

		irrelevantFrob := test_data.VatFrobModel
		irrelevantFrob.Ilk = "wrong ilk"
		irrelevantFrob.Urn = fakeUrn
		irrelevantFrob.Dink = strconv.Itoa(rand.Int())
		irrelevantFrob.Dart = strconv.Itoa(rand.Int())
		irrelevantFrob.TransactionIndex = frobBlockOne.TransactionIndex + 1

		err = frobRepo.Create(headerOneId, []interface{}{frobBlockOne, irrelevantFrob})
		Expect(err).NotTo(HaveOccurred())

		// New block
		headerTwo := fakes.GetFakeHeader(2)
		headerTwo.Hash = "anotherHash"
		headerTwoId, err := headerRepo.CreateOrUpdateHeader(headerTwo)
		Expect(err).NotTo(HaveOccurred())

		frobBlockTwo := test_data.VatFrobModel
		frobBlockTwo.Ilk = fakeIlk
		frobBlockTwo.Urn = fakeUrn
		frobBlockTwo.Dink = strconv.Itoa(rand.Int())
		frobBlockTwo.Dart = strconv.Itoa(rand.Int())

		err = frobRepo.Create(headerTwoId, []interface{}{frobBlockTwo})
		Expect(err).NotTo(HaveOccurred())

		var actualFrobs []FrobEvent
		err = db.Select(&actualFrobs, `SELECT ilkId, urnId, dink, dart FROM maker.frobs($1, $2)`, fakeIlk, fakeUrn)
		Expect(err).NotTo(HaveOccurred())

		Expect(actualFrobs).To(ConsistOf(
			FrobEvent{fakeIlk, fakeUrn, frobBlockOne.Dink, frobBlockOne.Dart},
			FrobEvent{fakeIlk, fakeUrn, frobBlockTwo.Dink, frobBlockTwo.Dart},
		))
	})
})

type FrobEvent struct {
	IlkId string
	UrnId string
	Dink  string
	Dart  string
}
