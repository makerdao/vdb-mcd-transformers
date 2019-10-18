// VulcanizeDB
// Copyright © 2019 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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

var _ = Describe("Managed CDP trigger-populated table", func() {
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

	It("stores the state of each managed CDP, unique by cdpi", func() {
		fakeIlk := test_helpers.FakeIlk.Hex
		fakeUrn := test_data.FakeUrn
		headerBlock := rand.Int()

		header := fakes.GetFakeHeader(int64(headerBlock))
		_, headerErr := headerRepo.CreateOrUpdateHeader(header)
		Expect(headerErr).NotTo(HaveOccurred())

		cdpManagerStorageValues1 := test_helpers.GetCdpManagerStorageValues(1, fakeIlk, fakeUrn, fakeCdpi)
		cdpErr1 := test_helpers.CreateManagedCdp(db, header, cdpManagerStorageValues1,
			test_helpers.GetCdpManagerMetadatas(strconv.Itoa(fakeCdpi)))
		Expect(cdpErr1).NotTo(HaveOccurred())

		fakeCdpi2 := fakeCdpi + 1
		cdpManagerStorageValues2 := test_helpers.GetCdpManagerStorageValues(2, fakeIlk, fakeUrn, fakeCdpi2)
		cdpErr2 := test_helpers.CreateManagedCdp(db, header, cdpManagerStorageValues2,
			test_helpers.GetCdpManagerMetadatas(strconv.Itoa(fakeCdpi2)))
		Expect(cdpErr2).NotTo(HaveOccurred())

		expectedCdp1 := test_helpers.ManagedCdpFromValues(
			test_helpers.FakeIlk.Identifier, header.Timestamp, cdpManagerStorageValues1)
		expectedCdp2 := test_helpers.ManagedCdpFromValues(
			test_helpers.FakeIlk.Identifier, header.Timestamp, cdpManagerStorageValues2)

		var actualCdps []test_helpers.ManagedCdp
		queryErr := db.Select(&actualCdps,
			`SELECT usr, cdpi, urn_identifier, ilk_identifier, created FROM api.managed_cdp`)
		Expect(queryErr).NotTo(HaveOccurred())

		Expect(actualCdps).To(ConsistOf([]test_helpers.ManagedCdp{expectedCdp1, expectedCdp2}))
	})

	It("stores the latest owner of the CDP according to block number", func() {
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

		newOwner := "0x16Fb96a5fa0427Af0C8F7cF1eB4870231c8154B6"
		_, ownsErr := db.Exec(cdp_manager.InsertOwnsQuery, headerTwo.BlockNumber, headerTwo.Hash, fakeCdpi, newOwner)
		Expect(ownsErr).NotTo(HaveOccurred())

		cdpManagerStorageValues := test_helpers.GetCdpManagerStorageValues(1, fakeIlk, fakeUrn, fakeCdpi)
		cdpErr := test_helpers.CreateManagedCdp(db, headerOne, cdpManagerStorageValues,
			test_helpers.GetCdpManagerMetadatas(strconv.Itoa(fakeCdpi)))
		Expect(cdpErr).NotTo(HaveOccurred())

		cdpManagerStorageValues[cdp_manager.Owns] = newOwner
		expectedCdp := test_helpers.ManagedCdpFromValues(
			test_helpers.FakeIlk.Identifier, headerOne.Timestamp, cdpManagerStorageValues)

		var actualCdps []test_helpers.ManagedCdp
		queryErr := db.Select(&actualCdps,
			`SELECT usr, cdpi, urn_identifier, ilk_identifier, created FROM api.managed_cdp WHERE cdpi = $1`,
			fakeCdpi)
		Expect(queryErr).NotTo(HaveOccurred())

		Expect(actualCdps).To(Equal([]test_helpers.ManagedCdp{expectedCdp}))
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
		queryErr := db.Get(&actualCdp,
			`SELECT usr, cdpi, urn_identifier, ilk_identifier, created FROM api.managed_cdp WHERE cdpi = $1`,
			fakeCdpi)
		Expect(queryErr).NotTo(HaveOccurred())

		Expect(expectedCdp).To(Equal(actualCdp))
	})
})
