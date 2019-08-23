package queries

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/storage/cdp_manager"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
	"strconv"
)

var _ = Describe("Get managed CDP by ID query", func() {
	var (
		db         *postgres.DB
		headerRepo repositories.HeaderRepository
		repo       cdp_manager.CdpManagerStorageRepository
		fakeCdpi   = rand.Int()
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		repo = cdp_manager.CdpManagerStorageRepository{}
		repo.SetDB(db)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	It("only gets requested CDP", func() {
		fakeIlk := test_helpers.FakeIlk.Hex
		fakeUrn := test_data.FakeUrn
		headerBlock := rand.Int()

		header := fakes.GetFakeHeader(int64(headerBlock))
		_, headerErr := headerRepo.CreateOrUpdateHeader(header)
		Expect(headerErr).NotTo(HaveOccurred())

		cdpManagerStorageValues := test_helpers.GetCdpManagerStorageValues(1, fakeIlk, fakeUrn, fakeCdpi)
		cdpErr1 := test_helpers.CreateManagedCdp(db, header, cdpManagerStorageValues,
			test_helpers.GetCdpManagerMetadatas(strconv.Itoa(fakeCdpi)))
		Expect(cdpErr1).NotTo(HaveOccurred())

		irrelevantCdpi := fakeCdpi + 1
		irrelevantStorageValues := test_helpers.GetCdpManagerStorageValues(2, fakeIlk, fakeUrn, irrelevantCdpi)
		cdpErr2 := test_helpers.CreateManagedCdp(db, header, irrelevantStorageValues,
			test_helpers.GetCdpManagerMetadatas(strconv.Itoa(irrelevantCdpi)))
		Expect(cdpErr2).NotTo(HaveOccurred())

		expectedCdp := test_helpers.ManagedCdpFromValues(
			test_helpers.FakeIlk.Identifier, header.Timestamp, cdpManagerStorageValues)

		var actualCdp test_helpers.ManagedCdp
		queryErr := db.Get(&actualCdp, `SELECT usr, id, urn_identifier, ilk_identifier, created FROM api.get_managed_cdp($1)`, fakeCdpi)
		Expect(queryErr).NotTo(HaveOccurred())

		Expect(expectedCdp).To(Equal(actualCdp))
	})

	It("gets the latest owner of the CDP", func() {
		fakeIlk := test_helpers.FakeIlk.Hex
		fakeUrn := test_data.FakeUrn

		headerOneBlock := rand.Int()
		headerOneTimestamp := int(rand.Int31())
		headerOne := fakes.GetFakeHeaderWithTimestamp(int64(headerOneTimestamp), int64(headerOneBlock))
		_, headerOneErr := headerRepo.CreateOrUpdateHeader(headerOne)
		Expect(headerOneErr).NotTo(HaveOccurred())

		headerTwoBlock := headerOneBlock + 1
		headerTwoTimestamp := headerOneTimestamp + 1000
		headerTwo := fakes.GetFakeHeaderWithTimestamp(int64(headerTwoTimestamp), int64(headerTwoBlock))
		_, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
		Expect(headerTwoErr).NotTo(HaveOccurred())

		cdpManagerStorageValues := test_helpers.GetCdpManagerStorageValues(1, fakeIlk, fakeUrn, fakeCdpi)
		cdpErr := test_helpers.CreateManagedCdp(db, headerOne, cdpManagerStorageValues,
			test_helpers.GetCdpManagerMetadatas(strconv.Itoa(fakeCdpi)))
		Expect(cdpErr).NotTo(HaveOccurred())

		newOwner := "0x16Fb96a5fa0427Af0C8F7cF1eB4870231c8154B6"
		_, ownsErr := db.Exec(cdp_manager.InsertOwnsQuery, headerTwo.BlockNumber, headerTwo.Hash, fakeCdpi, newOwner)
		Expect(ownsErr).NotTo(HaveOccurred())

		cdpManagerStorageValues[cdp_manager.CdpManagerOwns] = newOwner
		expectedCdp := test_helpers.ManagedCdpFromValues(
			test_helpers.FakeIlk.Identifier, headerOne.Timestamp, cdpManagerStorageValues)

		var actualCdp test_helpers.ManagedCdp
		queryErr := db.Get(&actualCdp, `SELECT usr, id, urn_identifier, ilk_identifier, created FROM api.get_managed_cdp($1)`,
			fakeCdpi)
		Expect(queryErr).NotTo(HaveOccurred())

		Expect(expectedCdp).To(Equal(actualCdp))
	})

	It("gets time created based on when cdpi changed", func() {
		fakeIlk := test_helpers.FakeIlk.Hex
		fakeUrn := test_data.FakeUrn

		headerOneBlock := rand.Int()
		headerOneTimestamp := int(rand.Int31())
		headerOne := fakes.GetFakeHeaderWithTimestamp(int64(headerOneTimestamp), int64(headerOneBlock))
		_, headerOneErr := headerRepo.CreateOrUpdateHeader(headerOne)
		Expect(headerOneErr).NotTo(HaveOccurred())

		headerTwoBlock := headerOneBlock + 1
		headerTwoTimestamp := headerOneTimestamp + 1000
		headerTwo := fakes.GetFakeHeaderWithTimestamp(int64(headerTwoTimestamp), int64(headerTwoBlock))
		_, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
		Expect(headerTwoErr).NotTo(HaveOccurred())

		_, cdpiErr := db.Exec(cdp_manager.InsertCdpiQuery, headerOne.BlockNumber, headerOne.Hash, fakeCdpi)
		Expect(cdpiErr).NotTo(HaveOccurred())

		cdpManagerStorageValues := test_helpers.GetCdpManagerStorageValues(1, fakeIlk, fakeUrn, fakeCdpi)
		cdpErr := test_helpers.CreateManagedCdp(db, headerTwo, cdpManagerStorageValues,
			test_helpers.GetCdpManagerMetadatas(strconv.Itoa(fakeCdpi)))
		Expect(cdpErr).NotTo(HaveOccurred())

		expectedCdp := test_helpers.ManagedCdpFromValues(
			test_helpers.FakeIlk.Identifier, headerOne.Timestamp, cdpManagerStorageValues)

		var actualCdp test_helpers.ManagedCdp
		queryErr := db.Get(&actualCdp, `SELECT usr, id, urn_identifier, ilk_identifier, created FROM api.get_managed_cdp($1)`,
			fakeCdpi)
		Expect(queryErr).NotTo(HaveOccurred())

		Expect(expectedCdp).To(Equal(actualCdp))
	})
})
