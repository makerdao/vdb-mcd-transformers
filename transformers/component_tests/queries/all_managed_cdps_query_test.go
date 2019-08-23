package queries

import (
	"math/rand"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/storage/cdp_manager"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("All managed CDPS", func() {
	var (
		db         *postgres.DB
		headerRepo repositories.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	It("returns a CDP for every unique CDPI (same block)", func() {
		fakeCdpiOne := rand.Int()
		fakeCdpiTwo := fakeCdpiOne + 1
		blockNumber := rand.Int()

		header := fakes.GetFakeHeader(int64(blockNumber))
		_, headerErr := headerRepo.CreateOrUpdateHeader(header)
		Expect(headerErr).NotTo(HaveOccurred())

		cdpStorageValuesOne := test_helpers.GetCdpManagerStorageValues(1, test_helpers.FakeIlk.Hex,
			test_data.FakeUrn, fakeCdpiOne)
		cdpErr1 := test_helpers.CreateManagedCdp(db, header, cdpStorageValuesOne,
			test_helpers.GetCdpManagerMetadatas(strconv.Itoa(fakeCdpiOne)))
		Expect(cdpErr1).NotTo(HaveOccurred())

		cdpStorageValuesTwo := test_helpers.GetCdpManagerStorageValues(2, test_helpers.FakeIlk.Hex,
			fakes.FakeAddress.Hex(), fakeCdpiTwo)
		cdpErr2 := test_helpers.CreateManagedCdp(db, header, cdpStorageValuesTwo,
			test_helpers.GetCdpManagerMetadatas(strconv.Itoa(fakeCdpiTwo)))
		Expect(cdpErr2).NotTo(HaveOccurred())

		var actualCdps []test_helpers.ManagedCdp
		queryErr := db.Select(&actualCdps,
			`SELECT usr, id, urn_identifier, ilk_identifier, created FROM api.all_managed_cdps()`)
		Expect(queryErr).NotTo(HaveOccurred())

		expectedCdpOne := test_helpers.ManagedCdpFromValues(test_helpers.FakeIlk.Identifier,
			header.Timestamp, cdpStorageValuesOne)
		expectedCdpTwo := test_helpers.ManagedCdpFromValues(test_helpers.FakeIlk.Identifier,
			header.Timestamp, cdpStorageValuesTwo)
		Expect(actualCdps).To(Equal([]test_helpers.ManagedCdp{expectedCdpOne, expectedCdpTwo}))
	})

	It("returns a CDP for every unique CDPI (different blocks)", func() {
		fakeCdpiOne := rand.Int()
		fakeCdpiTwo := fakeCdpiOne + 1
		blockOne := rand.Int()
		timestampOne := int(rand.Int31())
		blockTwo := blockOne + 1
		timestampTwo := timestampOne + 1000

		blockOneHeader := fakes.GetFakeHeaderWithTimestamp(int64(timestampOne), int64(blockOne))
		_, headerOneErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
		Expect(headerOneErr).NotTo(HaveOccurred())
		cdpStorageValuesOne := test_helpers.GetCdpManagerStorageValues(1, test_helpers.FakeIlk.Hex,
			test_data.FakeUrn, fakeCdpiOne)
		cdpErr1 := test_helpers.CreateManagedCdp(db, blockOneHeader, cdpStorageValuesOne,
			test_helpers.GetCdpManagerMetadatas(strconv.Itoa(fakeCdpiOne)))
		Expect(cdpErr1).NotTo(HaveOccurred())

		blockTwoHeader := fakes.GetFakeHeaderWithTimestamp(int64(timestampTwo), int64(blockTwo))
		_, headerTwoErr := headerRepo.CreateOrUpdateHeader(blockTwoHeader)
		Expect(headerTwoErr).NotTo(HaveOccurred())
		cdpStorageValuesTwo := test_helpers.GetCdpManagerStorageValues(2, test_helpers.AnotherFakeIlk.Hex,
			test_data.FakeUrn, fakeCdpiTwo)
		cdpErr2 := test_helpers.CreateManagedCdp(db, blockTwoHeader, cdpStorageValuesTwo,
			test_helpers.GetCdpManagerMetadatas(strconv.Itoa(fakeCdpiTwo)))
		Expect(cdpErr2).NotTo(HaveOccurred())

		var actualCdps []test_helpers.ManagedCdp
		queryErr := db.Select(&actualCdps,
			`SELECT usr, id, urn_identifier, ilk_identifier, created FROM api.all_managed_cdps()`)
		Expect(queryErr).NotTo(HaveOccurred())

		expectedCdpOne := test_helpers.ManagedCdpFromValues(test_helpers.FakeIlk.Identifier,
			blockOneHeader.Timestamp, cdpStorageValuesOne)
		expectedCdpTwo := test_helpers.ManagedCdpFromValues(test_helpers.AnotherFakeIlk.Identifier,
			blockTwoHeader.Timestamp, cdpStorageValuesTwo)
		Expect(actualCdps).To(Equal([]test_helpers.ManagedCdp{expectedCdpOne, expectedCdpTwo}))
	})

	It("optionally accepts a usr argument", func() {
		fakeCdpiOne := rand.Int()
		fakeCdpiTwo := fakeCdpiOne + 1
		blockNumber := rand.Int()

		header := fakes.GetFakeHeader(int64(blockNumber))
		_, headerErr := headerRepo.CreateOrUpdateHeader(header)
		Expect(headerErr).NotTo(HaveOccurred())

		ownerOne := "fakeUsr1"
		cdpStorageValuesOne := test_helpers.GetCdpManagerStorageValues(1, test_helpers.FakeIlk.Hex,
			test_data.FakeUrn, fakeCdpiOne)
		cdpStorageValuesOne[cdp_manager.CdpManagerOwns] = ownerOne
		cdpErr1 := test_helpers.CreateManagedCdp(db, header, cdpStorageValuesOne,
			test_helpers.GetCdpManagerMetadatas(strconv.Itoa(fakeCdpiOne)))
		Expect(cdpErr1).NotTo(HaveOccurred())

		ownerTwo := "fakeUsr2"
		cdpStorageValuesTwo := test_helpers.GetCdpManagerStorageValues(1, test_helpers.FakeIlk.Hex,
			test_data.FakeUrn, fakeCdpiTwo)
		cdpStorageValuesTwo[cdp_manager.CdpManagerOwns] = ownerTwo
		cdpErr2 := test_helpers.CreateManagedCdp(db, header, cdpStorageValuesTwo,
			test_helpers.GetCdpManagerMetadatas(strconv.Itoa(fakeCdpiTwo)))
		Expect(cdpErr2).NotTo(HaveOccurred())

		var actualCdps []test_helpers.ManagedCdp
		queryErr := db.Select(&actualCdps,
			`SELECT usr, id, urn_identifier, ilk_identifier, created FROM api.all_managed_cdps($1)`, ownerOne)
		Expect(queryErr).NotTo(HaveOccurred())

		expectedCdp := test_helpers.ManagedCdpFromValues(test_helpers.FakeIlk.Identifier, header.Timestamp,
			cdpStorageValuesOne)
		Expect(actualCdps).To(Equal([]test_helpers.ManagedCdp{expectedCdp}))
	})
})
